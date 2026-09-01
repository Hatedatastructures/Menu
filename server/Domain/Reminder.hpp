#pragma once

#include <string>

namespace Menu::Domain {

struct Reminder {
    std::string Id;
    std::string UserId;
    std::string MealPlanItemId;
    std::string ReminderType;
    std::string ScheduledAt;
    bool IsEnabled = true;
    std::string CompletedAt;
};

}  // namespace Menu::Domain
