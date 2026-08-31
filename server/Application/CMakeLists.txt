add_library(MenuApplication INTERFACE)
target_include_directories(MenuApplication INTERFACE "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuApplication INTERFACE MenuDomain MenuFoundation)
