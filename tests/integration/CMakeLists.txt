add_executable(MenuTransportIntegrationTests
    "${PROJECT_SOURCE_DIR}/tests/TestMain.cpp"
    "${PROJECT_SOURCE_DIR}/tests/Integration/HttpTransportTest.cpp"
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

add_executable(MenuApiIntegrationTests
    "${PROJECT_SOURCE_DIR}/tests/TestMain.cpp"
    "${PROJECT_SOURCE_DIR}/tests/Integration/MenuApiTest.cpp"
)
target_include_directories(MenuApiIntegrationTests PRIVATE
    "${PROJECT_SOURCE_DIR}/server"
    "${PROJECT_SOURCE_DIR}/tests"
)
target_link_libraries(MenuApiIntegrationTests PRIVATE
    GTest::gtest
    MenuComposition
)
target_compile_features(MenuApiIntegrationTests PRIVATE cxx_std_20)
gtest_discover_tests(MenuApiIntegrationTests
    DISCOVERY_MODE PRE_TEST
    PROPERTIES TIMEOUT 30
)
