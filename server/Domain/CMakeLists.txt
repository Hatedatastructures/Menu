add_library(MenuDomain INTERFACE)
target_include_directories(MenuDomain INTERFACE "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuDomain INTERFACE MenuFoundation)
