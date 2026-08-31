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

add_executable(MenuMigrationTests
    "${PROJECT_SOURCE_DIR}/tests/TestMain.cpp"
    "${PROJECT_SOURCE_DIR}/tests/unit/SqliteMigrationTest.cpp"
)
target_include_directories(MenuMigrationTests PRIVATE
    "${PROJECT_SOURCE_DIR}/server"
    "${PROJECT_SOURCE_DIR}/tests"
)
target_link_libraries(MenuMigrationTests PRIVATE GTest::gtest MenuInfrastructure)
target_compile_features(MenuMigrationTests PRIVATE cxx_std_20)
gtest_discover_tests(MenuMigrationTests
    DISCOVERY_MODE PRE_TEST
    PROPERTIES TIMEOUT 30
)

add_executable(MenuRecipeRepositoryTests
    "${PROJECT_SOURCE_DIR}/tests/TestMain.cpp"
    "${PROJECT_SOURCE_DIR}/tests/unit/SqliteRecipeRepositoryTest.cpp"
)
target_include_directories(MenuRecipeRepositoryTests PRIVATE
    "${PROJECT_SOURCE_DIR}/server"
    "${PROJECT_SOURCE_DIR}/tests"
)
target_link_libraries(MenuRecipeRepositoryTests PRIVATE GTest::gtest MenuInfrastructure)
target_compile_features(MenuRecipeRepositoryTests PRIVATE cxx_std_20)
gtest_discover_tests(MenuRecipeRepositoryTests
    DISCOVERY_MODE PRE_TEST
    PROPERTIES TIMEOUT 30
)

add_executable(MenuApplicationTests
    "${PROJECT_SOURCE_DIR}/tests/TestMain.cpp"
    "${PROJECT_SOURCE_DIR}/tests/unit/RecipeApplicationServiceTest.cpp"
    "${PROJECT_SOURCE_DIR}/tests/unit/StorageExecutorTest.cpp"
)
target_include_directories(MenuApplicationTests PRIVATE
    "${PROJECT_SOURCE_DIR}/server"
    "${PROJECT_SOURCE_DIR}/tests"
)
target_link_libraries(MenuApplicationTests PRIVATE GTest::gtest MenuApplication MenuDomain)
target_compile_features(MenuApplicationTests PRIVATE cxx_std_20)
gtest_discover_tests(MenuApplicationTests
    DISCOVERY_MODE PRE_TEST
    PROPERTIES TIMEOUT 30
)

add_executable(MenuAuthTests
    "${PROJECT_SOURCE_DIR}/tests/TestMain.cpp"
    "${PROJECT_SOURCE_DIR}/tests/unit/PasswordHasherTest.cpp"
    "${PROJECT_SOURCE_DIR}/tests/unit/TokenServiceTest.cpp"
)
target_include_directories(MenuAuthTests PRIVATE "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuAuthTests PRIVATE GTest::gtest MenuInfrastructure)
target_compile_features(MenuAuthTests PRIVATE cxx_std_20)
gtest_discover_tests(MenuAuthTests
    DISCOVERY_MODE PRE_TEST
    PROPERTIES TIMEOUT 30
)
