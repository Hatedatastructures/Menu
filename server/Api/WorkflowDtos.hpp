#pragma once

#include <Domain/CookingSession.hpp>
#include <Domain/Feedback.hpp>
#include <Domain/MealPlan.hpp>

#include <boost/json/array.hpp>
#include <boost/json/object.hpp>

#include <vector>

namespace Menu::Api {

class WorkflowDtos final {
public:
    static boost::json::object ToObject(const Domain::MealPlan& Plan);
    static boost::json::object ToObject(const Domain::CookingSession& Session);
    static boost::json::object ToObject(const Domain::Feedback& FeedbackValue);
    static boost::json::array ToArray(const std::vector<Domain::MealPlan>& Plans);
};

}  // namespace Menu::Api
