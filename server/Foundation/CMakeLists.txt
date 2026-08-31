add_library(MenuFoundation INTERFACE)
target_include_directories(MenuFoundation INTERFACE "${PROJECT_SOURCE_DIR}/server")
target_compile_features(MenuFoundation INTERFACE cxx_std_20)
