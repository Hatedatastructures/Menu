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

add_library(MenuBoostJson STATIC "${PROJECT_SOURCE_DIR}/cmake/BoostJson.cpp")
target_link_libraries(MenuBoostJson PUBLIC MenuBoost)
set_target_properties(MenuBoostJson PROPERTIES POSITION_INDEPENDENT_CODE ON)

