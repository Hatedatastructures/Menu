add_library(MenuRuntime STATIC
    ServerConfiguration.cpp
    ServerConfiguration.hpp
)
target_include_directories(MenuRuntime PUBLIC "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuRuntime PUBLIC MenuTransportCore MenuBoostJson)
target_compile_features(MenuRuntime PUBLIC cxx_std_20)
