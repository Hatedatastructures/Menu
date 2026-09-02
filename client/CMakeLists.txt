cmake_minimum_required(VERSION 3.23)

find_package(Qt6 6.8 REQUIRED COMPONENTS Core Gui Qml Quick QuickControls2 Network Concurrent)
if(MENU_ENABLE_TESTS)
    find_package(Qt6 6.8 REQUIRED COMPONENTS Test QuickTest)
endif()

set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_AUTOUIC ON)
qt_policy(SET QTP0001 NEW)
qt_policy(SET QTP0002 NEW)
qt_policy(SET QTP0004 NEW)

add_library(MenuClientCore STATIC
    Cpp/ClientApi.cpp
    Cpp/ClientApiSupport.cpp
    Cpp/ClientApiState.cpp
    Cpp/ClientApiCache.cpp
    Cpp/ClientApiCacheStorage.cpp
    Cpp/ClientApiPreferences.cpp
    Cpp/ClientApiRequests.cpp
    Cpp/ClientApiAuthentication.cpp
    Cpp/ClientApiWorkflow.cpp
    Cpp/ClientSessionState.cpp
    Cpp/ClientPreferenceState.cpp
    Cpp/ClientTransport.cpp
    Cpp/ClientApiConnection.cpp
    Cpp/ConnectionEndpoint.cpp
    Cpp/ConnectionSettingsStore.cpp
    Cpp/PreferencesStore.cpp
    Cpp/ClientApi.hpp
    Cpp/ClientSessionState.hpp
    Cpp/ClientPreferenceState.hpp
    Cpp/ClientTransport.hpp
    Cpp/ClientApiSupport.hpp
    Cpp/ConnectionEndpoint.hpp
    Cpp/ConnectionSettingsStore.hpp
    Cpp/PreferencesStore.hpp
    Cpp/RecipeModel.cpp
    Cpp/IngredientModel.cpp
    Cpp/RecommendationModel.cpp
    Cpp/MealPlanModel.cpp
    Cpp/IngredientModel.hpp
    Cpp/RecipeModel.hpp
    Cpp/RecommendationModel.hpp
    Cpp/MealPlanModel.hpp
)
target_include_directories(MenuClientCore PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/Cpp")
target_link_libraries(MenuClientCore PUBLIC
    Qt6::Core
    Qt6::Network
    Qt6::Concurrent
)
target_compile_features(MenuClientCore PUBLIC cxx_std_20)

add_library(MenuClientPlatform STATIC
    Cpp/DisplayMetrics.cpp
    Cpp/DisplayMetrics.hpp
    Cpp/NotificationController.cpp
    Cpp/NotificationController.hpp
    Cpp/NotificationBackend.cpp
    Cpp/NotificationBackend.hpp
    Cpp/CookingTimerCoordinator.cpp
    Cpp/CookingTimerCoordinator.hpp
    Cpp/ThemeController.cpp
    Cpp/ThemeController.hpp
    Cpp/SystemUiBackend.cpp
    Cpp/SystemUiBackend.hpp
    Cpp/ScreenAwakeBackend.cpp
    Cpp/ScreenAwakeBackend.hpp
    Cpp/ScreenAwakeController.cpp
    Cpp/ScreenAwakeController.hpp
)
target_include_directories(MenuClientPlatform PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/Cpp")
target_link_libraries(MenuClientPlatform PUBLIC
    Qt6::Core
    Qt6::Gui
)
target_compile_features(MenuClientPlatform PUBLIC cxx_std_20)

if(MENU_ENABLE_TESTS)
    add_subdirectory(Tests)
    add_subdirectory(QmlTests)
endif()

set(MenuMediaFiles
    "${PROJECT_SOURCE_DIR}/assets/media/tomato-egg.png"
    "${PROJECT_SOURCE_DIR}/assets/media/soy-chicken.png"
    "${PROJECT_SOURCE_DIR}/assets/media/beef-onion.png"
    "${PROJECT_SOURCE_DIR}/assets/media/tofu-spinach.png"
    "${PROJECT_SOURCE_DIR}/assets/media/pasta-tomato.png"
    "${PROJECT_SOURCE_DIR}/assets/media/lemon-salmon.png"
    "${PROJECT_SOURCE_DIR}/assets/media/cream-spinach.png"
    "${PROJECT_SOURCE_DIR}/assets/media/oyakodon.png"
    "${PROJECT_SOURCE_DIR}/assets/media/tofu-teriyaki.png"
    "${PROJECT_SOURCE_DIR}/assets/media/ginger-beef.png"
    "${PROJECT_SOURCE_DIR}/assets/media/tomato-rice.png"
    "${PROJECT_SOURCE_DIR}/assets/media/tofu-salad.png"
    "${PROJECT_SOURCE_DIR}/assets/media/menu-placeholder.png"
)
foreach(MenuMediaFile IN LISTS MenuMediaFiles)
    get_filename_component(MenuMediaName "${MenuMediaFile}" NAME)
    set_source_files_properties("${MenuMediaFile}" PROPERTIES
        QT_RESOURCE_ALIAS "assets/media/${MenuMediaName}")
endforeach()

set_source_files_properties(Qml/Theme.qml PROPERTIES QT_QML_SINGLETON_TYPE TRUE)

qt_add_executable(MenuClientApp MenuMain.cpp)
qt_add_qml_module(MenuClientApp
    URI MenuClient
    VERSION 1.0
    QML_FILES
        Qml/Main.qml
        Qml/Theme.qml
        Qml/Components/AppButton.qml
        Qml/Components/AppIcon.qml
        Qml/Components/AppSwitch.qml
        Qml/Components/AppSegmentedControl.qml
        Qml/Components/AppStepper.qml
        Qml/Components/CookingStepCard.qml
        Qml/Components/PlanDayCard.qml
        Qml/Components/AppTextField.qml
        Qml/Components/AppChip.qml
        Qml/Components/AppCheckBox.qml
        Qml/Components/AppHeader.qml
        Qml/Components/NavigationRail.qml
        Qml/Components/AuthPanel.qml
        Qml/Components/AccountPanel.qml
        Qml/Components/PreferencesPanel.qml
        Qml/Components/ConnectionPanel.qml
        Qml/Components/DeviceNotificationPanel.qml
        Qml/Components/SectionHeading.qml
        Qml/Components/RecipeCard.qml
        Qml/Components/BottomNavigation.qml
        Qml/Components/NavigationState.qml
        Qml/Components/DateFormatter.qml
        Qml/Pages/TonightPage.qml
        Qml/Pages/WeekPage.qml
        Qml/Pages/IngredientsPage.qml
        Qml/Pages/ProfilePage.qml
        Qml/Pages/CookingPage.qml
    RESOURCES
        ${MenuMediaFiles}
)
target_link_libraries(MenuClientApp PRIVATE
    MenuClientCore
    MenuClientPlatform
    Qt6::Gui
    Qt6::Qml
    Qt6::Quick
    Qt6::QuickControls2
)
set_target_properties(MenuClientApp PROPERTIES
    OUTPUT_NAME MenuClient
    RUNTIME_OUTPUT_DIRECTORY "${PROJECT_BINARY_DIR}/client"
    QT_ANDROID_PACKAGE_NAME "com.menu.cookflow"
    QT_ANDROID_VERSION_CODE 1
    QT_ANDROID_VERSION_NAME "0.1.0"
    QT_ANDROID_MIN_SDK_VERSION 28
    QT_ANDROID_TARGET_SDK_VERSION 35
)
if(ANDROID)
    set_property(TARGET MenuClientApp PROPERTY QT_ANDROID_PACKAGE_SOURCE_DIR
        "${CMAKE_CURRENT_SOURCE_DIR}/Android")
endif()
