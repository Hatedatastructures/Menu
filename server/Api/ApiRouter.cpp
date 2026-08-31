#include "Api/ApiRouter.hpp"

#include "Api/ApiErrors.hpp"
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

bool IsSafeIdentifier(std::string_view Value) {
    if (Value.empty() || Value.size() > 128U) {
        return false;
    }
    return std::ranges::all_of(Value, [](char Character) {
        return std::isalnum(static_cast<unsigned char>(Character)) != 0 ||
               Character == '.' || Character == '-' || Character == '_';
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
    Application::StorageExecutor& StorageValue)
    : Service(ServiceValue), Storage(StorageValue) {}

void ApiRouter::Handle(
    Transport::HttpRequest RequestValue,
    Transport::HttpResponseCallback Complete) {
    if (!Complete) {
        return;
    }
    const std::string Id = RequestId(RequestValue);
    const TargetParts Target = SplitTarget(RequestValue.Target);
    if (RequestValue.Method == "GET" && Target.Path == "/healthz") {
        boost::json::object Body;
        Body["status"] = "ok";
        Complete(JsonResponse(200, std::move(Body), Id));
        return;
    }
    if (RequestValue.Method == "GET" && Target.Path == "/readyz") {
        boost::json::object Body;
        Body["status"] = Ready.load(std::memory_order_acquire) ? "ready" : "starting";
        Complete(JsonResponse(
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
        [State, Id, Complete](
            Application::StorageExecutionResult Execution) mutable {
            if (!Execution.Succeeded() || !State->Response.has_value()) {
                Complete(InternalError(Id));
                return;
            }
            Complete(std::move(State->Response).value());
        },
        boost::asio::system_executor());
    if (!Accepted) {
        Complete(InternalError(Id));
    }
}

void ApiRouter::SetReady(bool ReadyValue) noexcept {
    Ready.store(ReadyValue, std::memory_order_release);
}

Transport::HttpResponse ApiRouter::Route(
    const Transport::HttpRequest& Request,
    const std::string& Id) {
    const TargetParts Target = SplitTarget(Request.Target);
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
