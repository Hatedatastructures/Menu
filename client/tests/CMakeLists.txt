add_executable(MenuClientModelTests
    "${PROJECT_SOURCE_DIR}/client/tests/RecipeModelTest.cpp"
)
target_link_libraries(MenuClientModelTests PRIVATE
    Qt6::Test
    MenuClientCore
)
target_compile_features(MenuClientModelTests PRIVATE cxx_std_20)

add_executable(MenuClientApiTests
    "${PROJECT_SOURCE_DIR}/client/tests/ClientApiTest.cpp"
)
target_link_libraries(MenuClientApiTests PRIVATE
    Qt6::Test
    MenuClientCore
)
target_compile_features(MenuClientApiTests PRIVATE cxx_std_20)

add_custom_target(MenuClientTests)
add_dependencies(MenuClientTests MenuClientModelTests MenuClientApiTests)
add_test(NAME MenuClientModelTests COMMAND MenuClientModelTests)

add_executable(MenuClientPlanModelTests
    "${PROJECT_SOURCE_DIR}/client/tests/MealPlanModelTest.cpp"
)
target_link_libraries(MenuClientPlanModelTests PRIVATE
    Qt6::Test
    MenuClientCore
)
target_compile_features(MenuClientPlanModelTests PRIVATE cxx_std_20)
add_test(NAME MenuClientPlanModelTests COMMAND MenuClientPlanModelTests)
add_test(NAME MenuClientApiTests COMMAND MenuClientApiTests
    -o "${CMAKE_CURRENT_BINARY_DIR}/MenuClientApiTests.log,txt")
