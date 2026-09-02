add_executable(MenuClientModelTests
    "${PROJECT_SOURCE_DIR}/client/Tests/RecipeModelTest.cpp"
)
target_link_libraries(MenuClientModelTests PRIVATE
    Qt6::Test
    MenuClientCore
)
target_compile_features(MenuClientModelTests PRIVATE cxx_std_20)

add_executable(MenuClientApiTests
    "${PROJECT_SOURCE_DIR}/client/Tests/ClientApiTest.cpp"
)
target_link_libraries(MenuClientApiTests PRIVATE
    Qt6::Test
    MenuClientCore
)
target_compile_features(MenuClientApiTests PRIVATE cxx_std_20)

add_custom_target(MenuClientTests)
add_dependencies(MenuClientTests
    MenuClientModelTests
    MenuClientApiTests
    MenuClientPlanModelTests
    MenuClientConnectionEndpointTests
    MenuClientDisplayMetricsTests
    MenuClientNotificationControllerTests
)
add_test(NAME MenuClientModelTests COMMAND MenuClientModelTests)

add_executable(MenuClientIngredientModelTests
    "${PROJECT_SOURCE_DIR}/client/Tests/IngredientModelTest.cpp"
)
target_link_libraries(MenuClientIngredientModelTests PRIVATE
    Qt6::Test
    MenuClientCore
)
target_compile_features(MenuClientIngredientModelTests PRIVATE cxx_std_20)
add_dependencies(MenuClientTests MenuClientIngredientModelTests)
add_test(NAME MenuClientIngredientModelTests COMMAND MenuClientIngredientModelTests)

add_executable(MenuClientPlanModelTests
    "${PROJECT_SOURCE_DIR}/client/Tests/MealPlanModelTest.cpp"
)
target_link_libraries(MenuClientPlanModelTests PRIVATE
    Qt6::Test
    MenuClientCore
)
target_compile_features(MenuClientPlanModelTests PRIVATE cxx_std_20)
add_test(NAME MenuClientPlanModelTests COMMAND MenuClientPlanModelTests)

add_executable(MenuClientConnectionEndpointTests
    "${PROJECT_SOURCE_DIR}/client/Tests/ConnectionEndpointTest.cpp"
)
target_link_libraries(MenuClientConnectionEndpointTests PRIVATE
    Qt6::Test
    MenuClientCore
)
target_compile_features(MenuClientConnectionEndpointTests PRIVATE cxx_std_20)
add_test(NAME MenuClientConnectionEndpointTests COMMAND MenuClientConnectionEndpointTests)

add_executable(MenuClientDisplayMetricsTests
    "${PROJECT_SOURCE_DIR}/client/Tests/DisplayMetricsTest.cpp"
)
target_link_libraries(MenuClientDisplayMetricsTests PRIVATE
    Qt6::Test
    MenuClientPlatform
)
target_compile_features(MenuClientDisplayMetricsTests PRIVATE cxx_std_20)
add_test(NAME MenuClientDisplayMetricsTests COMMAND MenuClientDisplayMetricsTests)

add_executable(MenuClientNotificationControllerTests
    "${PROJECT_SOURCE_DIR}/client/Tests/NotificationControllerTest.cpp"
)
target_link_libraries(MenuClientNotificationControllerTests PRIVATE
    Qt6::Test
    MenuClientPlatform
)
target_compile_features(MenuClientNotificationControllerTests PRIVATE cxx_std_20)
add_test(NAME MenuClientNotificationControllerTests COMMAND MenuClientNotificationControllerTests)

add_executable(MenuClientCookingTimerCoordinatorTests
    "${PROJECT_SOURCE_DIR}/client/Tests/CookingTimerCoordinatorTest.cpp"
)
target_link_libraries(MenuClientCookingTimerCoordinatorTests PRIVATE
    Qt6::Test
    MenuClientPlatform
)
target_compile_features(MenuClientCookingTimerCoordinatorTests PRIVATE cxx_std_20)
add_dependencies(MenuClientTests MenuClientCookingTimerCoordinatorTests)
add_test(NAME MenuClientCookingTimerCoordinatorTests COMMAND MenuClientCookingTimerCoordinatorTests)

add_executable(MenuClientThemeControllerTests
    "${PROJECT_SOURCE_DIR}/client/Tests/ThemeControllerTest.cpp"
)
target_link_libraries(MenuClientThemeControllerTests PRIVATE
    Qt6::Test
    MenuClientPlatform
)
target_compile_features(MenuClientThemeControllerTests PRIVATE cxx_std_20)
add_dependencies(MenuClientTests MenuClientThemeControllerTests)
add_test(NAME MenuClientThemeControllerTests COMMAND MenuClientThemeControllerTests)

add_executable(MenuClientPreferencesStoreTests
    "${PROJECT_SOURCE_DIR}/client/Tests/PreferencesStoreTest.cpp"
)
target_link_libraries(MenuClientPreferencesStoreTests PRIVATE
    Qt6::Test
    MenuClientCore
)
target_compile_features(MenuClientPreferencesStoreTests PRIVATE cxx_std_20)
add_dependencies(MenuClientTests MenuClientPreferencesStoreTests)
add_test(NAME MenuClientPreferencesStoreTests COMMAND MenuClientPreferencesStoreTests)

add_executable(MenuClientPreferenceStateTests
    "${PROJECT_SOURCE_DIR}/client/Tests/ClientPreferenceStateTest.cpp"
)
target_link_libraries(MenuClientPreferenceStateTests PRIVATE
    Qt6::Test
    MenuClientCore
)
target_compile_features(MenuClientPreferenceStateTests PRIVATE cxx_std_20)
add_dependencies(MenuClientTests MenuClientPreferenceStateTests)
add_test(NAME MenuClientPreferenceStateTests COMMAND MenuClientPreferenceStateTests)

add_executable(MenuClientScreenAwakeControllerTests
    "${PROJECT_SOURCE_DIR}/client/Tests/ScreenAwakeControllerTest.cpp"
)
target_link_libraries(MenuClientScreenAwakeControllerTests PRIVATE
    Qt6::Test
    MenuClientPlatform
)
target_compile_features(MenuClientScreenAwakeControllerTests PRIVATE cxx_std_20)
add_dependencies(MenuClientTests MenuClientScreenAwakeControllerTests)
add_test(NAME MenuClientScreenAwakeControllerTests COMMAND MenuClientScreenAwakeControllerTests)

add_test(NAME MenuClientApiTests COMMAND MenuClientApiTests
    -o "${CMAKE_CURRENT_BINARY_DIR}/MenuClientApiTests.log,txt")

if(WIN32)
    get_filename_component(MenuQtCMakeDirectory "${Qt6Core_DIR}" DIRECTORY)
    get_filename_component(MenuQtCMakeRoot "${MenuQtCMakeDirectory}" DIRECTORY)
    get_filename_component(MenuQtInstallRoot "${MenuQtCMakeRoot}" DIRECTORY)
    set(MenuQtRuntimeDirectory "${MenuQtInstallRoot}/bin")
    set_tests_properties(
        MenuClientModelTests
        MenuClientApiTests
        MenuClientIngredientModelTests
        MenuClientPlanModelTests
        MenuClientConnectionEndpointTests
        MenuClientDisplayMetricsTests
        MenuClientNotificationControllerTests
        MenuClientCookingTimerCoordinatorTests
        MenuClientThemeControllerTests
        MenuClientPreferencesStoreTests
        MenuClientPreferenceStateTests
        MenuClientScreenAwakeControllerTests
        PROPERTIES ENVIRONMENT "PATH=${MenuQtRuntimeDirectory}"
    )
endif()
