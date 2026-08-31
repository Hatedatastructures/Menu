#pragma once

#include <Foundation/Result.hpp>
#include <Domain/Recommendation.hpp>

#include <vector>

namespace Menu::Application {

class RecommendationProvider {
public:
    virtual ~RecommendationProvider() = default;

    virtual Foundation::Result<std::vector<Domain::Recommendation>> Recommend(
        const Domain::RecommendationRequest& Request,
        const std::vector<Domain::Recipe>& Recipes) = 0;
};

}  // namespace Menu::Application
