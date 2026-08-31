#pragma once

#include <Application/RecommendationProvider.hpp>

namespace Menu::Application {

class RuleBasedRecommendationProvider final : public RecommendationProvider {
public:
    Foundation::Result<std::vector<Domain::Recommendation>> Recommend(
        const Domain::RecommendationRequest& Request,
        const std::vector<Domain::Recipe>& Recipes) override;
};

}  // namespace Menu::Application
