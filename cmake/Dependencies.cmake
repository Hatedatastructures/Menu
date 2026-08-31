include(FetchContent)

if(NOT DEFINED FETCHCONTENT_BASE_DIR OR FETCHCONTENT_BASE_DIR STREQUAL "")
    set(FETCHCONTENT_BASE_DIR "${PROJECT_SOURCE_DIR}/.cache/cmake-fetch" CACHE PATH
        "Menu dependency source and download cache")
endif()

set(INSTALL_GTEST OFF CACHE BOOL "" FORCE)
set(BUILD_GMOCK OFF CACHE BOOL "" FORCE)
set(gtest_force_shared_crt ON CACHE BOOL "" FORCE)

FetchContent_Declare(
    GoogleTest
    URL https://github.com/google/googletest/archive/refs/tags/v1.16.0.tar.gz
    URL_HASH SHA256=78c676fc63881529bf97bf9d45948d905a66833fbfa5318ea2cd7478cb98f399
    DOWNLOAD_EXTRACT_TIMESTAMP TRUE
)
FetchContent_MakeAvailable(GoogleTest)

if(POLICY CMP0169)
    cmake_policy(SET CMP0169 OLD)
endif()

FetchContent_Declare(
    BoostSource
    URL https://archives.boost.io/release/1.89.0/source/boost_1_89_0.tar.bz2
    URL_HASH SHA256=85a33fa22621b4f314f8e85e1a5e2a9363d22e4f4992925d4bb3bc631b5a0c7a
    DOWNLOAD_EXTRACT_TIMESTAMP TRUE
)
FetchContent_GetProperties(BoostSource)
if(NOT boostsource_POPULATED)
    FetchContent_Populate(BoostSource)
endif()

add_library(MenuBoost INTERFACE)
target_include_directories(MenuBoost INTERFACE "${boostsource_SOURCE_DIR}")
target_compile_definitions(MenuBoost INTERFACE
    BOOST_ERROR_CODE_HEADER_ONLY
    BOOST_SYSTEM_NO_DEPRECATED
    BOOST_SYSTEM_NO_LIB
)
target_compile_features(MenuBoost INTERFACE cxx_std_20)
if(WIN32)
    target_link_libraries(MenuBoost INTERFACE ws2_32 mswsock)
endif()

add_library(MenuBoostJson STATIC "${PROJECT_SOURCE_DIR}/cmake/BoostJson.cpp")
target_link_libraries(MenuBoostJson PUBLIC MenuBoost)
set_target_properties(MenuBoostJson PROPERTIES POSITION_INDEPENDENT_CODE ON)

set(MenuSqliteCachedSourceDir
    "${PROJECT_SOURCE_DIR}/.cache/sqlite/sqlite-amalgamation-3530400"
)
if(EXISTS "${MenuSqliteCachedSourceDir}/sqlite3.c")
    set(FETCHCONTENT_SOURCE_DIR_SQLITESOURCE "${MenuSqliteCachedSourceDir}" CACHE PATH
        "Use the existing Menu SQLite source cache")
endif()

FetchContent_Declare(
    SQLiteSource
    URL https://www.sqlite.org/2026/sqlite-amalgamation-3530400.zip
    URL_HASH SHA256=1E71DDF93849C6A6ECF58B827C0692073D2DD7EE40196158068F7B29F422E87D
    DOWNLOAD_EXTRACT_TIMESTAMP TRUE
)
FetchContent_GetProperties(SQLiteSource)
if(NOT sqlitesource_POPULATED)
    FetchContent_Populate(SQLiteSource)
endif()

set(MenuSqliteSourceDir "${sqlitesource_SOURCE_DIR}")
if(EXISTS "${MenuSqliteSourceDir}/sqlite-amalgamation-3530400/sqlite3.c")
    set(MenuSqliteSourceDir
        "${MenuSqliteSourceDir}/sqlite-amalgamation-3530400")
endif()
if(NOT EXISTS "${MenuSqliteSourceDir}/sqlite3.c")
    message(FATAL_ERROR "SQLite amalgamation source was not found in ${MenuSqliteSourceDir}")
endif()

add_library(MenuSqlite STATIC "${MenuSqliteSourceDir}/sqlite3.c")
target_include_directories(MenuSqlite PUBLIC "${MenuSqliteSourceDir}")
target_compile_definitions(MenuSqlite PUBLIC
    SQLITE_THREADSAFE=1
    SQLITE_DQS=0
)
set_target_properties(MenuSqlite PROPERTIES POSITION_INDEPENDENT_CODE ON)
