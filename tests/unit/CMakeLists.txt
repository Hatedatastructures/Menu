add_executable(MenuUnitTests
    "${PROJECT_SOURCE_DIR}/tests/TestMain.cpp"
    "${PROJECT_SOURCE_DIR}/tests/unit/BootstrapContractTest.cpp"
)
target_include_directories(MenuUnitTests PRIVATE "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuUnitTests PRIVATE GTest::gtest MenuFoundation)
target_compile_features(MenuUnitTests PRIVATE cxx_std_20)

include(GoogleTest)
gtest_discover_tests(MenuUnitTests
    DISCOVERY_MODE PRE_TEST
    PROPERTIES TIMEOUT 30
)
