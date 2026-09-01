#include "Api/ApiParsing.hpp"

#include "Api/ApiParsingSupport.hpp"

#include <boost/json/parse.hpp>

#include <set>
#include <string>
#include <utility>

namespace Menu::Api::Parsing {

Foundation::Result<MealPlanPayload> ReadMealPlanPayload(std::string_view Body) {
    if (Body.empty() || Body.size() > 256U * 1024U) {
        return Foundation::Result<MealPlanPayload>::FromError(
            Support::InvalidRequest("计划请求体无效"));
    }
    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Body, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Foundation::Result<MealPlanPayload>::FromError(
            Support::InvalidRequest("计划 JSON 无效"));
    }
    static const std::set<std::string> AllowedKeys = {"planDate", "items"};
    for (const auto& Entry : Parsed.as_object()) {
        if (!AllowedKeys.contains(std::string(Entry.key()))) {
            return Foundation::Result<MealPlanPayload>::FromError(
                Support::InvalidRequest("计划包含未知字段"));
        }
    }
    const auto PlanDate = ReadRequiredString(Parsed.as_object(), "planDate", 10U);
    const boost::json::value* ItemsValue = Parsed.as_object().if_contains("items");
    if (!PlanDate.HasValue() || ItemsValue == nullptr || !ItemsValue->is_array() ||
        ItemsValue->as_array().size() > 14U) {
        return Foundation::Result<MealPlanPayload>::FromError(
            Support::InvalidRequest("计划日期或菜谱项无效"));
    }
    MealPlanPayload Payload;
    Payload.PlanDate = PlanDate.Value();
    for (std::size_t Index = 0; Index < ItemsValue->as_array().size(); ++Index) {
        const boost::json::value& ItemValue = ItemsValue->as_array().at(Index);
        if (!ItemValue.is_object()) {
            return Foundation::Result<MealPlanPayload>::FromError(
                Support::InvalidRequest("计划菜谱项无效"));
        }
        const boost::json::object& ItemObject = ItemValue.as_object();
        static const std::set<std::string> ItemKeys = {
            "recipeId", "servings", "sortOrder"};
        for (const auto& Entry : ItemObject) {
            if (!ItemKeys.contains(std::string(Entry.key()))) {
                return Foundation::Result<MealPlanPayload>::FromError(
                    Support::InvalidRequest("计划菜谱项包含未知字段"));
            }
        }
        const auto RecipeId = ReadRequiredString(ItemObject, "recipeId", 128U);
        const auto Servings = ReadJsonInteger(ItemObject, "servings", 1, 24);
        int SortOrder = static_cast<int>(Index);
        if (ItemObject.if_contains("sortOrder") != nullptr) {
            const auto ParsedSortOrder = ReadJsonInteger(ItemObject, "sortOrder", 0, 128);
            if (!ParsedSortOrder.HasValue()) {
                return Foundation::Result<MealPlanPayload>::FromError(
                    ParsedSortOrder.ErrorValue());
            }
            SortOrder = ParsedSortOrder.Value();
        }
        if (!RecipeId.HasValue() || !Servings.HasValue()) {
            return Foundation::Result<MealPlanPayload>::FromError(
                Support::InvalidRequest("计划菜谱项字段无效"));
        }
        Domain::MealPlanItem Item;
        Item.RecipeId = RecipeId.Value();
        Item.Servings = Servings.Value();
        Item.SortOrder = SortOrder;
        Payload.Items.push_back(std::move(Item));
    }
    return Payload;
}

Foundation::Result<CookingSessionUpdatePayload> ReadCookingSessionUpdatePayload(
    std::string_view Body) {
    if (Body.empty() || Body.size() > 16U * 1024U) {
        return Foundation::Result<CookingSessionUpdatePayload>::FromError(
            Support::InvalidRequest("做饭会话请求体无效"));
    }
    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Body, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Foundation::Result<CookingSessionUpdatePayload>::FromError(
            Support::InvalidRequest("做饭会话 JSON 无效"));
    }
    static const std::set<std::string> AllowedKeys = {"currentStepOrder", "state"};
    for (const auto& Entry : Parsed.as_object()) {
        if (!AllowedKeys.contains(std::string(Entry.key()))) {
            return Foundation::Result<CookingSessionUpdatePayload>::FromError(
                Support::InvalidRequest("做饭会话包含未知字段"));
        }
    }
    const auto StepOrder = ReadJsonInteger(
        Parsed.as_object(), "currentStepOrder", 1, 128);
    const auto State = ReadRequiredString(Parsed.as_object(), "state", 16U);
    if (!StepOrder.HasValue() || !State.HasValue()) {
        return Foundation::Result<CookingSessionUpdatePayload>::FromError(
            Support::InvalidRequest("做饭会话字段无效"));
    }
    return CookingSessionUpdatePayload{StepOrder.Value(), State.Value()};
}

Foundation::Result<FeedbackPayload> ReadFeedbackPayload(std::string_view Body) {
    if (Body.empty() || Body.size() > 64U * 1024U) {
        return Foundation::Result<FeedbackPayload>::FromError(
            Support::InvalidRequest("反馈请求体无效"));
    }
    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Body, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Foundation::Result<FeedbackPayload>::FromError(
            Support::InvalidRequest("反馈 JSON 无效"));
    }
    static const std::set<std::string> AllowedKeys = {
        "recipeId", "outcome", "tags", "comment"};
    for (const auto& Entry : Parsed.as_object()) {
        if (!AllowedKeys.contains(std::string(Entry.key()))) {
            return Foundation::Result<FeedbackPayload>::FromError(
                Support::InvalidRequest("反馈包含未知字段"));
        }
    }
    const auto RecipeId = ReadRequiredString(Parsed.as_object(), "recipeId", 128U);
    const auto Outcome = ReadRequiredString(Parsed.as_object(), "outcome", 16U);
    const auto Tags = ReadJsonStringArray(Parsed.as_object(), "tags");
    if (!RecipeId.HasValue() || !Outcome.HasValue() || !Tags.HasValue()) {
        return Foundation::Result<FeedbackPayload>::FromError(
            Support::InvalidRequest("反馈字段无效"));
    }
    FeedbackPayload Payload;
    Payload.RecipeId = RecipeId.Value();
    Payload.Outcome = Outcome.Value();
    Payload.Tags = Tags.Value();
    if (const boost::json::value* Comment = Parsed.as_object().if_contains("comment");
        Comment != nullptr) {
        if (!Comment->is_string() || Comment->as_string().size() > 2000U) {
            return Foundation::Result<FeedbackPayload>::FromError(
                Support::InvalidRequest("反馈评论无效"));
        }
        Payload.Comment = std::string(Comment->as_string().c_str());
    }
    return Payload;
}

}  // namespace Menu::Api::Parsing
