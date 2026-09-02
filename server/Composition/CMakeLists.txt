add_library(MenuComposition STATIC
    ServerApplication.cpp
    ServerApplication.hpp
)
target_include_directories(MenuComposition PUBLIC "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuComposition PUBLIC
    MenuApi
    MenuTransport
    MenuInfrastructure
    MenuApplication
    MenuRuntime)
target_compile_features(MenuComposition PUBLIC cxx_std_20)
