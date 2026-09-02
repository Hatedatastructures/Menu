#include "Api/ApiParsing.hpp"

#include "Api/ApiParsingSupport.hpp"

#include <boost/json/parse.hpp>

#include <set>
#include <string>

namespace Menu::Api::Parsing {

Foundation::Result<Domain::RecommendationRequest> ReadRecommendationRequest(
    std::string_view Body) {
    if (Body.empty() || Body.size() > 1024U * 1024U) {
        return Foundation::Result<Domain::RecommendationRequest>::FromError(
            Support::InvalidRequest("推荐请求体为空或过大"));
    }
    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Body, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Foundation::Result<Domain::RecommendationRequest>::FromError(
            Support::InvalidRequest("推荐请求 JSON 无效"));
    }

    static const std::set<std::string> AllowedKeys = {
        "servings", "availableMinutes", "cuisines", "allergies",
        "pantryIngredientIds", "cookware", "recentRecipeIds"};
    for (const auto& Entry : Parsed.as_object()) {
        if (!AllowedKeys.contains(std::string(Entry.key()))) {
            return Foundation::Result<Domain::RecommendationRequest>::FromError(
                Support::InvalidRequest("推荐请求包含未知字段"));
        }
    }

    const auto Servings = ReadJsonInteger(Parsed.as_object(), "servings", 1, 24);
    const auto AvailableMinutes = ReadJsonInteger(
        Parsed.as_object(), "availableMinutes", 1, 24 * 60);
    if (!Servings.HasValue() || !AvailableMinutes.HasValue()) {
        return Foundation::Result<Domain::RecommendationRequest>::FromError(
            Support::InvalidRequest("推荐请求份量或时间无效"));
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
            Support::InvalidRequest("推荐请求数组无效"));
    }
    Request.Cuisines = Cuisines.Value();
    Request.Allergies = Allergies.Value();
    Request.PantryIngredientIds = PantryIds.Value();
    Request.Cookware = Cookware.Value();
    Request.RecentRecipeIds = RecentIds.Value();
    return Request;
}

}  // namespace Menu::Api::Parsing
