cmake_minimum_required(VERSION 3.23)

project(Menu VERSION 0.1.0 LANGUAGES C CXX)

if(POLICY CMP0167)
    cmake_policy(SET CMP0167 NEW)
endif()

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

option(MENU_ENABLE_TESTS "Build Menu tests" ON)
option(MENU_ENABLE_CLIENT "Build the Qt Quick client" OFF)
option(MENU_ENABLE_SERVER "Build the C++ server" ON)

if(MENU_ENABLE_SERVER)
    include(cmake/Dependencies.cmake)
    add_subdirectory(server)
endif()

if(MENU_ENABLE_TESTS)
    include(CTest)
    enable_testing()
    if(MENU_ENABLE_SERVER)
        add_subdirectory(tests)
    endif()
endif()

if(MENU_ENABLE_CLIENT)
    add_subdirectory(client)
endif()

if(MENU_ENABLE_SERVER)
    include(cmake/TargetBoundaries.cmake)
    AssertMenuTargetBoundaries()
endif()
