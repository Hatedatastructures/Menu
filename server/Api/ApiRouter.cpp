#include "Api/ApiRouter.hpp"

#include "Api/ApiErrors.hpp"
#include "Api/AuthDtos.hpp"
#include "Api/RecipeDtos.hpp"

#include <boost/asio/system_executor.hpp>
#include <boost/json/object.hpp>
#include <boost/json/parse.hpp>
#include <boost/json/serialize.hpp>

#include <algorithm>
#include <atomic>
#include <charconv>
#include <cctype>
#include <cstdint>
#include <map>
#include <memory>
#include <optional>
#include <ranges>
#include <set>
#include <string>
#include <string_view>
#include <utility>

namespace Menu::Api {
namespace {

struct TargetParts {
    std::string Path;
    std::string Query;
};

struct RouteState {
    std::optional<Transport::HttpResponse> Response;
};

struct RecipeListOptions {
    std::string Cuisine;
    int MaxMinutes = 0;
    int Difficulty = 0;
    std::size_t Limit = 20;
    bool HasMaxMinutes = false;
};

struct AuthPayload {
    std::string Email;
    std::string Password;
    std::string DisplayName;
    std::string RefreshToken;
};

Foundation::Error InvalidRequest(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

TargetParts SplitTarget(std::string_view Target) {
    const std::size_t QueryPosition = Target.find('?');
    if (QueryPosition == std::string_view::npos) {
        return {std::string(Target), {}};
    }
    return {std::string(Target.substr(0, QueryPosition)),
            std::string(Target.substr(QueryPosition + 1))};
}

int HexValue(char Character) {
    if (Character >= '0' && Character <= '9') {
        return Character - '0';
    }
    if (Character >= 'a' && Character <= 'f') {
        return Character - 'a' + 10;
    }
    if (Character >= 'A' && Character <= 'F') {
        return Character - 'A' + 10;
    }
    return -1;
}

Foundation::Result<std::string> DecodeComponent(std::string_view Value) {
    std::string Decoded;
    Decoded.reserve(Value.size());
    for (std::size_t Index = 0; Index < Value.size(); ++Index) {
        if (Value[Index] == '%') {
            if (Index + 2 >= Value.size()) {
                return Foundation::Result<std::string>::FromError(
                    InvalidRequest("查询参数编码无效"));
            }
            const int High = HexValue(Value[Index + 1]);
            const int Low = HexValue(Value[Index + 2]);
            if (High < 0 || Low < 0) {
                return Foundation::Result<std::string>::FromError(
                    InvalidRequest("查询参数编码无效"));
            }
            Decoded.push_back(static_cast<char>((High << 4) | Low));
            Index += 2;
        } else if (Value[Index] == '+') {
            Decoded.push_back(' ');
        } else {
            Decoded.push_back(Value[Index]);
        }
    }
    return Decoded;
}

Foundation::Result<std::map<std::string, std::string>> ParseQuery(std::string_view Query) {
    std::map<std::string, std::string> Values;
    if (Query.empty()) {
        return Values;
    }
    std::size_t Start = 0;
    while (Start <= Query.size()) {
        const std::size_t End = Query.find('&', Start);
        const std::string_view Pair = Query.substr(
            Start, End == std::string_view::npos ? Query.size() - Start : End - Start);
        const std::size_t Separator = Pair.find('=');
        if (Separator == std::string_view::npos || Separator == 0) {
            return Foundation::Result<std::map<std::string, std::string>>::FromError(
                InvalidRequest("查询参数格式无效"));
        }
        const auto Key = DecodeComponent(Pair.substr(0, Separator));
        const auto Value = DecodeComponent(Pair.substr(Separator + 1));
        if (!Key.HasValue() || !Value.HasValue() || Key.Value().size() > 64U ||
            Value.Value().size() > 256U) {
            return Foundation::Result<std::map<std::string, std::string>>::FromError(
                InvalidRequest("查询参数过长或编码无效"));
        }
        if (!Values.emplace(Key.Value(), Value.Value()).second) {
            return Foundation::Result<std::map<std::string, std::string>>::FromError(
                InvalidRequest("查询参数重复"));
        }
        if (End == std::string_view::npos) {
            break;
        }
        Start = End + 1;
    }
    return Values;
}

Foundation::Result<int> ParseInteger(std::string_view Value, int Minimum, int Maximum) {
    int Parsed = 0;
    const auto ParseResult = std::from_chars(
        Value.data(), Value.data() + Value.size(), Parsed);
    if (ParseResult.ec != std::errc() || ParseResult.ptr != Value.data() + Value.size() ||
        Parsed < Minimum || Parsed > Maximum) {
        return Foundation::Result<int>::FromError(InvalidRequest("数字参数无效"));
    }
    return Parsed;
}

Foundation::Result<std::string> ReadRequiredString(
    const boost::json::object& Object,
    std::string_view Key,
    std::size_t MaximumLength) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr || !Value->is_string() || Value->as_string().empty() ||
        Value->as_string().size() > MaximumLength) {
        return Foundation::Result<std::string>::FromError(
            InvalidRequest("认证请求字段无效"));
    }
    return std::string(Value->as_string().c_str());
}

Foundation::Result<int> ReadJsonInteger(
    const boost::json::object& Object,
    std::string_view Key,
    int Minimum,
    int Maximum);

Foundation::Result<std::vector<std::string>> ReadJsonStringArray(
    const boost::json::object& Object,
    std::string_view Key);

Foundation::Result<double> ReadJsonDouble(
    const boost::json::object& Object,
    std::string_view Key) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr) {
        return Foundation::Result<double>::FromError(
            InvalidRequest("菜谱数字字段缺失"));
    }
    if (Value->is_double()) {
        return Value->as_double();
    }
    if (Value->is_int64()) {
        return static_cast<double>(Value->as_int64());
    }
    return Foundation::Result<double>::FromError(
        InvalidRequest("菜谱数字字段无效"));
}

Foundation::Result<bool> ReadJsonBool(
    const boost::json::object& Object,
    std::string_view Key) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr || !Value->is_bool()) {
        return Foundation::Result<bool>::FromError(
            InvalidRequest("菜谱布尔字段无效"));
    }
    return Value->as_bool();
}

Foundation::Result<Domain::Recipe> ReadRecipePayload(std::string_view Body) {
    if (Body.empty() || Body.size() > 1024U * 1024U) {
        return Foundation::Result<Domain::Recipe>::FromError(
            InvalidRequest("菜谱请求体为空或过大"));
    }
    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Body, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Foundation::Result<Domain::Recipe>::FromError(
            InvalidRequest("菜谱 JSON 无效"));
    }
    static const std::set<std::string> AllowedKeys = {
        "id", "slug", "name", "cuisine", "description", "prepMinutes",
        "cookMinutes", "servings", "difficulty", "imagePath", "status",
        "allergens", "cookware", "ingredients", "steps"};
    for (const auto& Entry : Parsed.as_object()) {
        if (!AllowedKeys.contains(std::string(Entry.key().data(), Entry.key().size()))) {
            return Foundation::Result<Domain::Recipe>::FromError(
                InvalidRequest("菜谱包含未知字段"));
        }
    }

    const boost::json::object& Object = Parsed.as_object();
    const auto Id = ReadRequiredString(Object, "id", 128U);
    const auto Slug = ReadRequiredString(Object, "slug", 160U);
    const auto Name = ReadRequiredString(Object, "name", 160U);
    const auto Cuisine = ReadRequiredString(Object, "cuisine", 32U);
    const auto Description = ReadRequiredString(Object, "description", 4096U);
    const auto ImagePath = ReadRequiredString(Object, "imagePath", 512U);
    const auto Status = ReadRequiredString(Object, "status", 16U);
    const auto PrepMinutes = ReadJsonInteger(Object, "prepMinutes", 0, 24 * 60);
    const auto CookMinutes = ReadJsonInteger(Object, "cookMinutes", 0, 24 * 60);
    const auto Servings = ReadJsonInteger(Object, "servings", 1, 24);
    const auto Difficulty = ReadJsonInteger(Object, "difficulty", 1, 5);
    const auto Allergens = ReadJsonStringArray(Object, "allergens");
    const auto Cookware = ReadJsonStringArray(Object, "cookware");
    if (!Id.HasValue() || !Slug.HasValue() || !Name.HasValue() || !Cuisine.HasValue() ||
        !Description.HasValue() || !ImagePath.HasValue() || !Status.HasValue() ||
        !PrepMinutes.HasValue() || !CookMinutes.HasValue() || !Servings.HasValue() ||
        !Difficulty.HasValue() || !Allergens.HasValue() || !Cookware.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(
            InvalidRequest("菜谱基础字段无效"));
    }

    const boost::json::value* IngredientsValue = Object.if_contains("ingredients");
    const boost::json::value* StepsValue = Object.if_contains("steps");
    if (IngredientsValue == nullptr || !IngredientsValue->is_array() ||
        IngredientsValue->as_array().size() > 128U || StepsValue == nullptr ||
        !StepsValue->is_array() || StepsValue->as_array().size() > 128U) {
        return Foundation::Result<Domain::Recipe>::FromError(
            InvalidRequest("菜谱食材或步骤无效"));
    }

    Domain::Recipe RecipeValue;
    RecipeValue.Id = Id.Value();
    RecipeValue.Slug = Slug.Value();
    RecipeValue.Name = Name.Value();
    RecipeValue.Cuisine = Cuisine.Value();
    RecipeValue.Description = Description.Value();
    RecipeValue.PrepMinutes = PrepMinutes.Value();
    RecipeValue.CookMinutes = CookMinutes.Value();
    RecipeValue.Servings = Servings.Value();
    RecipeValue.Difficulty = Difficulty.Value();
    RecipeValue.ImagePath = ImagePath.Value();
    RecipeValue.Status = Status.Value();
    RecipeValue.Allergens = Allergens.Value();
    RecipeValue.Cookware = Cookware.Value();

    for (const boost::json::value& Value : IngredientsValue->as_array()) {
        if (!Value.is_object()) {
            return Foundation::Result<Domain::Recipe>::FromError(
                InvalidRequest("菜谱食材项无效"));
        }
        const boost::json::object& Item = Value.as_object();
        static const std::set<std::string> IngredientKeys = {
            "ingredientId", "quantity", "unit", "required", "servingFactor", "preparation"};
        for (const auto& Entry : Item) {
            if (!IngredientKeys.contains(
                    std::string(Entry.key().data(), Entry.key().size()))) {
                return Foundation::Result<Domain::Recipe>::FromError(
                    InvalidRequest("菜谱食材包含未知字段"));
            }
        }
        const auto IngredientId = ReadRequiredString(Item, "ingredientId", 128U);
        const auto Quantity = ReadJsonDouble(Item, "quantity");
        const auto Unit = ReadRequiredString(Item, "unit", 16U);
        const auto Required = ReadJsonBool(Item, "required");
        const auto ServingFactor = ReadJsonDouble(Item, "servingFactor");
        std::string Preparation;
        if (const boost::json::value* PreparationValue = Item.if_contains("preparation");
            PreparationValue != nullptr) {
            if (!PreparationValue->is_string() || PreparationValue->as_string().size() > 256U) {
                return Foundation::Result<Domain::Recipe>::FromError(
                    InvalidRequest("菜谱处理方式无效"));
            }
            Preparation = std::string(PreparationValue->as_string().c_str());
        }
        if (!IngredientId.HasValue() || !Quantity.HasValue() || !Unit.HasValue() ||
            !Required.HasValue() || !ServingFactor.HasValue()) {
            return Foundation::Result<Domain::Recipe>::FromError(
                InvalidRequest("菜谱食材字段无效"));
        }
        RecipeValue.Ingredients.emplace_back(
            IngredientId.Value(), Quantity.Value(), Unit.Value(), Required.Value(),
            ServingFactor.Value(), std::move(Preparation));
    }

    for (const boost::json::value& Value : StepsValue->as_array()) {
        if (!Value.is_object()) {
            return Foundation::Result<Domain::Recipe>::FromError(
                InvalidRequest("菜谱步骤项无效"));
        }
        const boost::json::object& Item = Value.as_object();
        const auto StepOrder = ReadJsonInteger(Item, "stepOrder", 1, 128);
        const auto Title = ReadRequiredString(Item, "title", 160U);
        const auto Instruction = ReadRequiredString(Item, "instruction", 4096U);
        const auto Duration = ReadJsonInteger(Item, "durationSeconds", 0, 24 * 60 * 60);
        const auto HasTimer = ReadJsonBool(Item, "hasTimer");
        if (!StepOrder.HasValue() || !Title.HasValue() || !Instruction.HasValue() ||
            !Duration.HasValue() || !HasTimer.HasValue()) {
            return Foundation::Result<Domain::Recipe>::FromError(
                InvalidRequest("菜谱步骤字段无效"));
        }
        RecipeValue.Steps.push_back(Domain::RecipeStep{
            StepOrder.Value(), Title.Value(), Instruction.Value(), Duration.Value(),
            HasTimer.Value()});
    }
    return RecipeValue;
}

Foundation::Result<AuthPayload> ReadAuthPayload(
    std::string_view Body,
    std::string_view Kind) {
    if (Body.empty() || Body.size() > 64U * 1024U) {
        return Foundation::Result<AuthPayload>::FromError(
            InvalidRequest("认证请求体无效"));
    }
    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Body, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Foundation::Result<AuthPayload>::FromError(
            InvalidRequest("认证请求 JSON 无效"));
    }
    const std::set<std::string> AllowedKeys = Kind == "register"
                                                  ? std::set<std::string>{"email", "password", "displayName"}
                                                  : Kind == "login"
                                                        ? std::set<std::string>{"email", "password"}
                                                        : std::set<std::string>{"refreshToken"};
    for (const auto& Entry : Parsed.as_object()) {
        if (!AllowedKeys.contains(std::string(Entry.key().data(), Entry.key().size()))) {
            return Foundation::Result<AuthPayload>::FromError(
                InvalidRequest("认证请求包含未知字段"));
        }
    }

    AuthPayload Payload;
    if (Kind == "refresh") {
        const auto RefreshToken = ReadRequiredString(
            Parsed.as_object(), "refreshToken", 256U);
        if (!RefreshToken.HasValue()) {
            return Foundation::Result<AuthPayload>::FromError(RefreshToken.ErrorValue());
        }
        Payload.RefreshToken = RefreshToken.Value();
        return Payload;
    }
    const auto Email = ReadRequiredString(Parsed.as_object(), "email", 320U);
    const std::size_t PasswordMaximumLength = Kind == "login" ? 64U * 1024U : 128U;
    const auto Password = ReadRequiredString(
        Parsed.as_object(), "password", PasswordMaximumLength);
    if (!Email.HasValue() || !Password.HasValue()) {
        return Foundation::Result<AuthPayload>::FromError(
            InvalidRequest("认证请求字段无效"));
    }
    Payload.Email = Email.Value();
    Payload.Password = Password.Value();
    if (Kind == "register") {
        const auto DisplayName = ReadRequiredString(
            Parsed.as_object(), "displayName", 80U);
        if (!DisplayName.HasValue()) {
            return Foundation::Result<AuthPayload>::FromError(DisplayName.ErrorValue());
        }
        Payload.DisplayName = DisplayName.Value();
    }
    return Payload;
}

Transport::HttpResponse AuthError(
    const Foundation::Error& ErrorValue,
    std::string_view RequestId) {
    int Status = 503;
    std::string_view Code = "auth_unavailable";
    std::string_view Message = "认证服务暂时不可用";
    if (ErrorValue.CodeValue() == Foundation::ErrorCode::InvalidArgument) {
        Status = 400;
        Code = "invalid_body";
        Message = "请求内容无效";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::Conflict) {
        Status = 409;
        Code = "email_already_registered";
        Message = "邮箱已注册";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::AuthenticationFailed) {
        Status = 401;
        Code = "authentication_failed";
        Message = "认证信息无效";
    }
    return ApiErrors::Create(Status, Code, Message, RequestId);
}

Foundation::Result<Application::AuthUser> RequireAdmin(
    Application::AuthService* Authentication,
    const Transport::HttpRequest& Request) {
    if (Authentication == nullptr) {
        return Foundation::Result<Application::AuthUser>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "认证服务不可用"));
    }
    const std::string Authorization = Request.HeaderValue("Authorization");
    constexpr std::string_view Prefix = "Bearer ";
    if (!Authorization.starts_with(Prefix) || Authorization.size() == Prefix.size()) {
        return Foundation::Result<Application::AuthUser>::FromError(
            Foundation::Error(Foundation::ErrorCode::AuthenticationFailed, "需要认证"));
    }
    const auto UserResult = Authentication->Authenticate(
        std::string_view(Authorization).substr(Prefix.size()));
    if (!UserResult.HasValue()) {
        return UserResult;
    }
    if (!UserResult.Value().IsAdmin) {
        return Foundation::Result<Application::AuthUser>::FromError(
            Foundation::Error(Foundation::ErrorCode::Forbidden, "需要管理员权限"));
    }
    return UserResult.Value();
}

Transport::HttpResponse AdminError(
    const Foundation::Error& ErrorValue,
    std::string_view RequestIdValue) {
    int Status = 503;
    std::string_view Code = "admin_unavailable";
    std::string_view Message = "管理服务暂时不可用";
    if (ErrorValue.CodeValue() == Foundation::ErrorCode::AuthenticationFailed) {
        Status = 401;
        Code = "authentication_required";
        Message = "请先登录";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::Forbidden) {
        Status = 403;
        Code = "admin_forbidden";
        Message = "需要管理员权限";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::InvalidArgument) {
        Status = 400;
        Code = "invalid_body";
        Message = "请求内容无效";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::Conflict) {
        Status = 409;
        Code = "recipe_conflict";
        Message = "菜谱已存在";
    } else if (ErrorValue.CodeValue() == Foundation::ErrorCode::NotFound) {
        Status = 404;
        Code = "recipe_not_found";
        Message = "菜谱不存在";
    }
    return ApiErrors::Create(Status, Code, Message, RequestIdValue);
}

bool IsSafeIdentifier(std::string_view Value) {
    if (Value.empty() || Value.size() > 128U) {
        return false;
    }
    return std::ranges::all_of(Value, [](char Character) {
        return std::isalnum(static_cast<unsigned char>(Character)) != 0 ||
               Character == '.' || Character == '-' || Character == '_';
    });
}

bool IsAllowedOrigin(
    const std::vector<std::string>& AllowedOrigins,
    std::string_view Origin) {
    return std::ranges::any_of(AllowedOrigins, [Origin](const std::string& AllowedOrigin) {
        return AllowedOrigin == Origin;
    });
}

std::string RequestId(const Transport::HttpRequest& Request) {
    const std::string Provided = Request.HeaderValue("X-Request-Id");
    if (!Provided.empty() && Provided.size() <= 64U &&
        std::ranges::all_of(Provided, [](char Character) {
            return Character >= 33 && Character <= 126;
        })) {
        return Provided;
    }
    static std::atomic<std::uint64_t> NextId = 1;
    return "menu-" + std::to_string(NextId.fetch_add(1, std::memory_order_relaxed));
}

Transport::HttpResponse JsonResponse(
    int Status,
    boost::json::value Body,
    std::string_view Id) {
    Transport::HttpResponse Response;
    Response.Status = Status;
    Response.ContentType = "application/json; charset=utf-8";
    Response.Body = boost::json::serialize(std::move(Body));
    Response.Headers.push_back(Transport::HttpHeader{"X-Request-Id", std::string(Id)});
    return Response;
}

Transport::HttpResponse InternalError(std::string_view Id) {
    return ApiErrors::Create(503, "storage_unavailable", "服务暂时不可用", Id);
}

Foundation::Result<RecipeListOptions> ReadRecipeListOptions(std::string_view Query) {
    const auto QueryResult = ParseQuery(Query);
    if (!QueryResult.HasValue()) {
        return Foundation::Result<RecipeListOptions>::FromError(QueryResult.ErrorValue());
    }
    RecipeListOptions Options;
    for (const auto& [Key, Value] : QueryResult.Value()) {
        if (Key == "cuisine") {
            if (Value.empty()) {
                return Foundation::Result<RecipeListOptions>::FromError(
                    InvalidRequest("菜系不能为空"));
            }
            Options.Cuisine = Value;
        } else if (Key == "maxMinutes") {
            const auto Parsed = ParseInteger(Value, 0, 24 * 60);
            if (!Parsed.HasValue()) {
                return Foundation::Result<RecipeListOptions>::FromError(Parsed.ErrorValue());
            }
            Options.MaxMinutes = Parsed.Value();
            Options.HasMaxMinutes = true;
        } else if (Key == "difficulty") {
            const auto Parsed = ParseInteger(Value, 1, 5);
            if (!Parsed.HasValue()) {
                return Foundation::Result<RecipeListOptions>::FromError(Parsed.ErrorValue());
            }
            Options.Difficulty = Parsed.Value();
        } else if (Key == "limit") {
            const auto Parsed = ParseInteger(Value, 1, 100);
            if (!Parsed.HasValue()) {
                return Foundation::Result<RecipeListOptions>::FromError(Parsed.ErrorValue());
            }
            Options.Limit = static_cast<std::size_t>(Parsed.Value());
        } else {
            return Foundation::Result<RecipeListOptions>::FromError(
                InvalidRequest("未知查询参数"));
        }
    }
    return Options;
}

Foundation::Result<int> ReadJsonInteger(
    const boost::json::object& Object,
    std::string_view Key,
    int Minimum,
    int Maximum) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr || !Value->is_int64()) {
        return Foundation::Result<int>::FromError(
            InvalidRequest("推荐请求数字字段无效"));
    }
    const std::int64_t Parsed = Value->as_int64();
    if (Parsed < Minimum || Parsed > Maximum) {
        return Foundation::Result<int>::FromError(
            InvalidRequest("推荐请求数字字段超出范围"));
    }
    return static_cast<int>(Parsed);
}

Foundation::Result<std::vector<std::string>> ReadJsonStringArray(
    const boost::json::object& Object,
    std::string_view Key) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr) {
        return std::vector<std::string>();
    }
    if (!Value->is_array() || Value->as_array().size() > 64U) {
        return Foundation::Result<std::vector<std::string>>::FromError(
            InvalidRequest("推荐请求数组字段无效"));
    }
    std::vector<std::string> Values;
    for (const boost::json::value& Item : Value->as_array()) {
        if (!Item.is_string() || Item.as_string().size() > 128U) {
            return Foundation::Result<std::vector<std::string>>::FromError(
                InvalidRequest("推荐请求数组元素无效"));
        }
        Values.emplace_back(Item.as_string().c_str());
    }
    return Values;
}

Foundation::Result<Domain::RecommendationRequest> ReadRecommendationRequest(
    std::string_view Body) {
    if (Body.empty() || Body.size() > 1024U * 1024U) {
        return Foundation::Result<Domain::RecommendationRequest>::FromError(
            InvalidRequest("推荐请求体为空或过大"));
    }
    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Body, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Foundation::Result<Domain::RecommendationRequest>::FromError(
            InvalidRequest("推荐请求 JSON 无效"));
    }

    static const std::set<std::string> AllowedKeys = {
        "servings", "availableMinutes", "cuisines", "allergies",
        "pantryIngredientIds", "cookware", "recentRecipeIds"};
    for (const auto& Entry : Parsed.as_object()) {
        if (!AllowedKeys.contains(std::string(Entry.key()))) {
            return Foundation::Result<Domain::RecommendationRequest>::FromError(
                InvalidRequest("推荐请求包含未知字段"));
        }
    }

    const auto Servings = ReadJsonInteger(Parsed.as_object(), "servings", 1, 24);
    const auto AvailableMinutes = ReadJsonInteger(
        Parsed.as_object(), "availableMinutes", 0, 24 * 60);
    if (!Servings.HasValue() || !AvailableMinutes.HasValue()) {
        return Foundation::Result<Domain::RecommendationRequest>::FromError(
            InvalidRequest("推荐请求份量或时间无效"));
    }
    Domain::RecommendationRequest Request;
    Request.Servings = Servings.Value();
    Request.AvailableMinutes = AvailableMinutes.Value();
    const auto Cuisines = ReadJsonStringArray(Parsed.as_object(), "cuisines");
    const auto Allergies = ReadJsonStringArray(Parsed.as_object(), "allergies");
    const auto PantryIds = ReadJsonStringArray(
        Parsed.as_object(), "pantryIngredientIds");
    const auto Cookware = ReadJsonStringArray(Parsed.as_object(), "cookware");
    const auto RecentIds = ReadJsonStringArray(
        Parsed.as_object(), "recentRecipeIds");
    if (!Cuisines.HasValue() || !Allergies.HasValue() || !PantryIds.HasValue() ||
        !Cookware.HasValue() || !RecentIds.HasValue()) {
        return Foundation::Result<Domain::RecommendationRequest>::FromError(
            InvalidRequest("推荐请求数组无效"));
    }
    Request.Cuisines = Cuisines.Value();
    Request.Allergies = Allergies.Value();
    Request.PantryIngredientIds = PantryIds.Value();
    Request.Cookware = Cookware.Value();
    Request.RecentRecipeIds = RecentIds.Value();
    return Request;
}

}  // namespace

ApiRouter::ApiRouter(
    Application::RecipeApplicationService& ServiceValue,
    Application::StorageExecutor& StorageValue,
    Application::AuthService* AuthenticationValue,
    Application::AdminRecipeApplicationService* AdminServiceValue,
    std::vector<std::string> CorsOriginsValue)
    : Service(ServiceValue),
      Storage(StorageValue),
      Authentication(AuthenticationValue),
      AdminService(AdminServiceValue),
      CorsOrigins(std::move(CorsOriginsValue)) {}

void ApiRouter::Handle(
    Transport::HttpRequest RequestValue,
    Transport::HttpResponseCallback Complete) {
    if (!Complete) {
        return;
    }
    const std::string Id = RequestId(RequestValue);
    const std::string Origin = RequestValue.HeaderValue("Origin");
    if (!Origin.empty() && !IsAllowedOrigin(CorsOrigins, Origin)) {
        Complete(ApiErrors::Create(403, "cors_denied", "跨域来源未被允许", Id));
        return;
    }
    const Transport::HttpResponseCallback Finish =
        [Complete = std::move(Complete), Origin](Transport::HttpResponse Response) mutable {
            if (!Origin.empty()) {
                Response.Headers.push_back(
                    Transport::HttpHeader{"Access-Control-Allow-Origin", Origin});
                Response.Headers.push_back(Transport::HttpHeader{"Vary", "Origin"});
            }
            Complete(std::move(Response));
        };
    const TargetParts Target = SplitTarget(RequestValue.Target);
    if (RequestValue.Method == "OPTIONS") {
        Transport::HttpResponse Response;
        Response.Status = 204;
        Response.ContentType = "text/plain; charset=utf-8";
        Response.Headers = {
            {"Access-Control-Allow-Methods", "GET, POST, OPTIONS"},
            {"Access-Control-Allow-Headers", "Content-Type, X-Request-Id"}};
        Finish(std::move(Response));
        return;
    }
    if (RequestValue.Method == "GET" && Target.Path == "/healthz") {
        boost::json::object Body;
        Body["status"] = "ok";
        Finish(JsonResponse(200, std::move(Body), Id));
        return;
    }
    if (RequestValue.Method == "GET" && Target.Path == "/readyz") {
        boost::json::object Body;
        Body["status"] = Ready.load(std::memory_order_acquire) ? "ready" : "starting";
        Finish(JsonResponse(
            Ready.load(std::memory_order_acquire) ? 200 : 503,
            std::move(Body), Id));
        return;
    }

    auto State = std::make_shared<RouteState>();
    const bool Accepted = Storage.Submit(
        [this,
         RequestValue = std::move(RequestValue),
         Id,
         State]() mutable {
            State->Response = Route(RequestValue, Id);
        },
        [State, Id, Finish](
            Application::StorageExecutionResult Execution) mutable {
            if (!Execution.Succeeded() || !State->Response.has_value()) {
                Finish(InternalError(Id));
                return;
            }
            Finish(std::move(State->Response).value());
        },
        boost::asio::system_executor());
    if (!Accepted) {
        Finish(InternalError(Id));
    }
}

void ApiRouter::SetReady(bool ReadyValue) noexcept {
    Ready.store(ReadyValue, std::memory_order_release);
}

Transport::HttpResponse ApiRouter::Route(
    const Transport::HttpRequest& Request,
    const std::string& Id) {
    const TargetParts Target = SplitTarget(Request.Target);
    if (Request.Method == "POST" &&
        (Target.Path == "/api/v1/auth/register" ||
         Target.Path == "/api/v1/auth/login" ||
         Target.Path == "/api/v1/auth/refresh")) {
        if (Authentication == nullptr) {
            return ApiErrors::Create(503, "auth_unavailable", "认证服务暂时不可用", Id);
        }
        const std::string Kind = Target.Path.ends_with("/register")
                                     ? "register"
                                     : Target.Path.ends_with("/login") ? "login" : "refresh";
        const auto Payload = ReadAuthPayload(Request.Body, Kind);
        if (!Payload.HasValue()) {
            return ApiErrors::Create(400, "invalid_body", "请求内容无效", Id);
        }
        Foundation::Result<Application::AuthResponse> AuthResult =
            Foundation::Result<Application::AuthResponse>::FromError(
                Foundation::Error(Foundation::ErrorCode::AuthenticationFailed, "认证失败"));
        if (Kind == "register") {
            AuthResult = Authentication->Register(
                Payload.Value().Email,
                Payload.Value().Password,
                Payload.Value().DisplayName);
        } else if (Kind == "login") {
            AuthResult = Authentication->Login(
                Payload.Value().Email, Payload.Value().Password);
        } else {
            AuthResult = Authentication->Refresh(Payload.Value().RefreshToken);
        }
        if (!AuthResult.HasValue()) {
            return AuthError(AuthResult.ErrorValue(), Id);
        }
        return JsonResponse(
            Kind == "register" ? 201 : 200,
            AuthDtos::ToObject(AuthResult.Value()),
            Id);
    }
    constexpr std::string_view AdminRecipePrefix = "/api/v1/admin/recipes";
    if (Request.Method == "GET" && Target.Path == AdminRecipePrefix) {
        if (AdminService == nullptr) {
            return ApiErrors::Create(503, "admin_unavailable", "管理服务暂时不可用", Id);
        }
        const auto UserResult = RequireAdmin(Authentication, Request);
        if (!UserResult.HasValue()) {
            return AdminError(UserResult.ErrorValue(), Id);
        }
        const auto RecipesResult = AdminService->ListAllRecipes();
        if (!RecipesResult.HasValue()) {
            return AdminError(RecipesResult.ErrorValue(), Id);
        }
        return JsonResponse(200, RecipeDtos::ToArray(RecipesResult.Value()), Id);
    }
    if (Target.Path == AdminRecipePrefix ||
        Target.Path.starts_with(std::string(AdminRecipePrefix) + "/")) {
        if (AdminService == nullptr) {
            return ApiErrors::Create(503, "admin_unavailable", "管理服务暂时不可用", Id);
        }
        const auto UserResult = RequireAdmin(Authentication, Request);
        if (!UserResult.HasValue()) {
            return AdminError(UserResult.ErrorValue(), Id);
        }
        if (Request.Method == "POST" && Target.Path == AdminRecipePrefix) {
            const auto RecipeResult = ReadRecipePayload(Request.Body);
            if (!RecipeResult.HasValue()) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", Id);
            }
            const auto CreatedResult = AdminService->CreateRecipe(RecipeResult.Value());
            if (!CreatedResult.HasValue()) {
                return AdminError(CreatedResult.ErrorValue(), Id);
            }
            return JsonResponse(201, RecipeDtos::ToObject(CreatedResult.Value()), Id);
        }

        const std::string RecipeId = Target.Path.substr(AdminRecipePrefix.size() + 1U);
        if (!IsSafeIdentifier(RecipeId)) {
            return ApiErrors::Create(400, "invalid_id", "菜谱 ID 无效", Id);
        }
        if (Request.Method == "PATCH") {
            const auto RecipeResult = ReadRecipePayload(Request.Body);
            if (!RecipeResult.HasValue()) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", Id);
            }
            const auto UpdatedResult = AdminService->UpdateRecipe(
                RecipeId, RecipeResult.Value());
            if (!UpdatedResult.HasValue()) {
                return AdminError(UpdatedResult.ErrorValue(), Id);
            }
            return JsonResponse(200, RecipeDtos::ToObject(UpdatedResult.Value()), Id);
        }
        if (Request.Method == "DELETE") {
            const auto DeletedResult = AdminService->DeleteRecipe(RecipeId);
            if (!DeletedResult.HasValue()) {
                return AdminError(DeletedResult.ErrorValue(), Id);
            }
            Transport::HttpResponse Response;
            Response.Status = 204;
            Response.ContentType = "application/json; charset=utf-8";
            Response.Headers.push_back(Transport::HttpHeader{
                "X-Request-Id", Id});
            return Response;
        }
    }
    if (Request.Method == "GET" && Target.Path == "/api/v1/recipes") {
        const auto OptionsResult = ReadRecipeListOptions(Target.Query);
        if (!OptionsResult.HasValue()) {
            return ApiErrors::Create(400, "invalid_query", "查询参数无效", Id);
        }
        const auto RecipesResult = Service.ListPublishedRecipes();
        if (!RecipesResult.HasValue()) {
            return InternalError(Id);
        }
        std::vector<Domain::Recipe> Recipes;
        for (const Domain::Recipe& RecipeValue : RecipesResult.Value()) {
            const RecipeListOptions& Options = OptionsResult.Value();
            if (!Options.Cuisine.empty() && RecipeValue.Cuisine != Options.Cuisine) {
                continue;
            }
            if (Options.HasMaxMinutes &&
                RecipeValue.PrepMinutes + RecipeValue.CookMinutes > Options.MaxMinutes) {
                continue;
            }
            if (Options.Difficulty != 0 && RecipeValue.Difficulty != Options.Difficulty) {
                continue;
            }
            Recipes.push_back(RecipeValue);
            if (Recipes.size() == Options.Limit) {
                break;
            }
        }
        return JsonResponse(200, RecipeDtos::ToArray(Recipes), Id);
    }

    if (Request.Method == "GET" && Target.Path == "/api/v1/ingredients") {
        if (!Target.Query.empty()) {
            return ApiErrors::Create(400, "invalid_query", "查询参数无效", Id);
        }
        const auto IngredientsResult = Service.ListIngredients();
        if (!IngredientsResult.HasValue()) {
            return InternalError(Id);
        }
        return JsonResponse(
            200, RecipeDtos::ToArray(IngredientsResult.Value()), Id);
    }

    constexpr std::string_view RecipePrefix = "/api/v1/recipes/";
    if (Request.Method == "GET" && Target.Path.starts_with(RecipePrefix)) {
        const std::string RecipeId = Target.Path.substr(RecipePrefix.size());
        if (!IsSafeIdentifier(RecipeId)) {
            return ApiErrors::Create(400, "invalid_id", "菜谱 ID 无效", Id);
        }
        const auto RecipeResult = Service.FindPublishedRecipe(RecipeId);
        if (!RecipeResult.HasValue()) {
            return InternalError(Id);
        }
        if (!RecipeResult.Value().has_value()) {
            return ApiErrors::Create(404, "recipe_not_found", "菜谱不存在", Id);
        }
        return JsonResponse(
            200, RecipeDtos::ToObject(RecipeResult.Value().value()), Id);
    }

    if (Request.Method == "POST" &&
        Target.Path == "/api/v1/recommendations/tonight") {
        const auto RequestResult = ReadRecommendationRequest(Request.Body);
        if (!RequestResult.HasValue()) {
            return ApiErrors::Create(400, "invalid_body", "请求内容无效", Id);
        }
        const auto Recommendations = Service.RecommendTonight(RequestResult.Value());
        if (!Recommendations.HasValue()) {
            if (Recommendations.ErrorValue().CodeValue() ==
                Foundation::ErrorCode::InvalidArgument) {
                return ApiErrors::Create(400, "invalid_body", "请求内容无效", Id);
            }
            return InternalError(Id);
        }
        return JsonResponse(200, RecipeDtos::ToArray(Recommendations.Value()), Id);
    }

    return ApiErrors::Create(404, "route_not_found", "接口不存在", Id);
}

}  // namespace Menu::Api
