#include "Application/RuleBasedRecommendationProvider.hpp"

namespace Menu::Application {

Foundation::Result<std::vector<Domain::Recommendation>>
RuleBasedRecommendationProvider::Recommend(
    const Domain::RecommendationRequest& Request,
    const std::vector<Domain::Recipe>& Recipes) {
    return Domain::RuleBasedRecommendation::Rank(Request, Recipes);
}

}  // namespace Menu::Application
