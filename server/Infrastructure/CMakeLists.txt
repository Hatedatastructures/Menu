file(READ "${CMAKE_CURRENT_SOURCE_DIR}/Migrations/001Initial.sql" MenuInitialMigrationSql)
file(READ "${CMAKE_CURRENT_SOURCE_DIR}/Migrations/002AuthTokens.sql" MenuAuthMigrationSql)
file(GENERATE
    OUTPUT "${CMAKE_CURRENT_BINARY_DIR}/MigrationSql.hpp"
    CONTENT "#pragma once\n\nnamespace Menu::Infrastructure::Generated {\ninline constexpr char InitialMigrationSql[] = R\"MenuSqlInitial(\n${MenuInitialMigrationSql}\n)MenuSqlInitial\";\ninline constexpr char AuthMigrationSql[] = R\"MenuSqlAuth(\n${MenuAuthMigrationSql}\n)MenuSqlAuth\";\n}\n"
)

add_library(MenuInfrastructure STATIC
    SqliteDatabase.cpp
    SqliteCookingSessionRepository.cpp
    SqliteAdminRecipeRepository.cpp
    SqliteAdminIngredientRepository.cpp
    SqliteFeedbackRepository.cpp
    SqliteMealPlanRepositoryRead.cpp
    SqliteMealPlanRepositorySupport.cpp
    SqliteMealPlanRepositoryWrite.cpp
    MigrationRunner.cpp
    SqliteRecipeRepositoryRead.cpp
    SqliteRecipeRepositoryReadSupport.cpp
    SqliteRecipeRepositoryWriteSupport.cpp
    SqliteRecipeRepositoryWrite.cpp
    SqliteIngredientRepositoryRead.cpp
    SqliteIngredientRepositorySupport.cpp
    SqliteIngredientRepositoryWrite.cpp
    PasswordHasher.cpp
    TokenService.cpp
    SqliteAuthRepository.cpp
    SqliteAuthRepositorySupport.cpp
    SqliteAuthRepositoryTokens.cpp
    SqliteAuthService.cpp
    SqliteAuthServiceTokenFlow.cpp
    SqliteAuthServiceValidation.cpp
    SqliteWorkflowSupport.cpp
    SeedCatalog.cpp
    SeedCatalog.hpp
    SeedData.cpp
    SeedDataIngredient.cpp
    SeedDataRecipe.cpp
    SeedDataRows.hpp
)
target_include_directories(MenuInfrastructure PUBLIC
    "${PROJECT_SOURCE_DIR}/server"
    "${CMAKE_CURRENT_BINARY_DIR}"
)
target_link_libraries(MenuInfrastructure PUBLIC
    MenuApplication MenuDomain MenuFoundation MenuSqlite MenuBoostJson OpenSSL::Crypto)
target_compile_features(MenuInfrastructure PUBLIC cxx_std_20)
