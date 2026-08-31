file(READ "${CMAKE_CURRENT_SOURCE_DIR}/Migrations/001Initial.sql" MenuInitialMigrationSql)
file(GENERATE
    OUTPUT "${CMAKE_CURRENT_BINARY_DIR}/MigrationSql.hpp"
    CONTENT "#pragma once\n\nnamespace Menu::Infrastructure::Generated {\ninline constexpr char InitialMigrationSql[] = R\"MenuSql(\n${MenuInitialMigrationSql}\n)MenuSql\";\n}\n"
)

add_library(MenuInfrastructure STATIC
    SqliteDatabase.cpp
    MigrationRunner.cpp
    SqliteRecipeRepository.cpp
    SqliteIngredientRepository.cpp
    SeedData.cpp
)
target_include_directories(MenuInfrastructure PUBLIC
    "${PROJECT_SOURCE_DIR}/server"
    "${CMAKE_CURRENT_BINARY_DIR}"
)
target_link_libraries(MenuInfrastructure PUBLIC
    MenuApplication MenuDomain MenuFoundation MenuSqlite MenuBoostJson)
target_compile_features(MenuInfrastructure PUBLIC cxx_std_20)
