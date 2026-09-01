cmake_minimum_required(VERSION 3.23)

find_package(Qt6 6.8 REQUIRED COMPONENTS Core Gui Qml Quick QuickControls2 Network Sql Concurrent Test QuickTest)

set(CMAKE_AUTOMOC ON)
set(CMAKE_AUTORCC ON)
set(CMAKE_AUTOUIC ON)
qt_policy(SET QTP0001 NEW)
qt_policy(SET QTP0002 NEW)
qt_policy(SET QTP0004 NEW)

add_library(MenuClientCore STATIC
    cpp/ClientApi.cpp
    cpp/ClientApiSupport.cpp
    cpp/ClientApiState.cpp
    cpp/ClientApiCache.cpp
    cpp/ClientApiRequests.cpp
    cpp/ClientApiAuthentication.cpp
    cpp/ClientApiWorkflow.cpp
    cpp/ClientApi.hpp
    cpp/ClientApiSupport.hpp
    cpp/RecipeModel.cpp
    cpp/IngredientModel.cpp
    cpp/RecommendationModel.cpp
    cpp/MealPlanModel.cpp
    cpp/IngredientModel.hpp
    cpp/RecipeModel.hpp
    cpp/RecommendationModel.hpp
    cpp/MealPlanModel.hpp
)
target_include_directories(MenuClientCore PUBLIC "${CMAKE_CURRENT_SOURCE_DIR}/cpp")
target_link_libraries(MenuClientCore PUBLIC
    Qt6::Core
    Qt6::Network
    Qt6::Sql
    Qt6::Concurrent
)
target_compile_features(MenuClientCore PUBLIC cxx_std_20)

if(MENU_ENABLE_TESTS)
    add_subdirectory(tests)
    add_subdirectory(qmltests)
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

set_source_files_properties(qml/Theme.qml PROPERTIES QT_QML_SINGLETON_TYPE TRUE)

qt_add_executable(MenuAndroid MenuMain.cpp)
qt_add_qml_module(MenuAndroid
    URI MenuClient
    VERSION 1.0
    QML_FILES
        qml/Main.qml
        qml/Theme.qml
        qml/Components/RecipeCard.qml
        qml/Components/BottomNavigation.qml
        qml/Components/NavigationState.qml
        qml/Components/DateFormatter.qml
        qml/Pages/TonightPage.qml
        qml/Pages/WeekPage.qml
        qml/Pages/IngredientsPage.qml
        qml/Pages/ProfilePage.qml
        qml/Pages/CookingPage.qml
    RESOURCES
        ${MenuMediaFiles}
)
target_link_libraries(MenuAndroid PRIVATE
    MenuClientCore
    Qt6::Gui
    Qt6::Qml
    Qt6::Quick
    Qt6::QuickControls2
)
set_target_properties(MenuAndroid PROPERTIES
    OUTPUT_NAME MenuClient
    RUNTIME_OUTPUT_DIRECTORY "${PROJECT_BINARY_DIR}/client"
    QT_ANDROID_PACKAGE_NAME "com.menu.cookflow"
    QT_ANDROID_VERSION_CODE 1
    QT_ANDROID_VERSION_NAME "0.1.0"
)
if(ANDROID)
    set_property(TARGET MenuAndroid PROPERTY QT_ANDROID_PACKAGE_SOURCE_DIR
        "${CMAKE_CURRENT_SOURCE_DIR}/android")
endif()
