add_executable(MenuUnitTests
    "${PROJECT_SOURCE_DIR}/tests/TestMain.cpp"
    "${PROJECT_SOURCE_DIR}/tests/unit/BootstrapContractTest.cpp"
    "${PROJECT_SOURCE_DIR}/tests/unit/RecipeRulesTest.cpp"
    "${PROJECT_SOURCE_DIR}/tests/unit/RecommendationRulesTest.cpp"
)
target_include_directories(MenuUnitTests PRIVATE
    "${PROJECT_SOURCE_DIR}/server"
    "${PROJECT_SOURCE_DIR}/tests"
)
target_link_libraries(MenuUnitTests PRIVATE GTest::gtest MenuDomain)
target_compile_features(MenuUnitTests PRIVATE cxx_std_20)

include(GoogleTest)
gtest_discover_tests(MenuUnitTests
    DISCOVERY_MODE PRE_TEST
    PROPERTIES TIMEOUT 30
)
