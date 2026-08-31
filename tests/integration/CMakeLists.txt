add_executable(MenuTransportIntegrationTests
    "${PROJECT_SOURCE_DIR}/tests/TestMain.cpp"
    "${PROJECT_SOURCE_DIR}/tests/integration/HttpTransportTest.cpp"
)
target_include_directories(MenuTransportIntegrationTests PRIVATE
    "${PROJECT_SOURCE_DIR}/server"
)
target_link_libraries(MenuTransportIntegrationTests PRIVATE
    GTest::gtest
    MenuTransport
)
target_compile_features(MenuTransportIntegrationTests PRIVATE cxx_std_20)

include(GoogleTest)
gtest_discover_tests(MenuTransportIntegrationTests
    DISCOVERY_MODE PRE_TEST
    PROPERTIES TIMEOUT 30
)
