#include "ScreenAwakeController.hpp"

namespace Menu::Client {

ScreenAwakeController::ScreenAwakeController(
    QObject* Parent,
    BackendFactory Factory)
    : QObject(Parent), Backend(Factory ? Factory() : CreateScreenAwakeBackend()) {}

ScreenAwakeController::~ScreenAwakeController() {
    if (Backend && EnabledValue) {
        Backend->Apply(false);
    }
}

bool ScreenAwakeController::Enabled() const noexcept {
    return EnabledValue;
}

void ScreenAwakeController::SetEnabled(bool EnabledValueInput) {
    if (EnabledValue == EnabledValueInput) {
        return;
    }
    EnabledValue = EnabledValueInput;
    if (Backend) {
        Backend->Apply(EnabledValue);
    }
    emit EnabledChanged();
}

}  // namespace Menu::Client
