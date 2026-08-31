add_library(MenuApi INTERFACE)
target_include_directories(MenuApi INTERFACE "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuApi INTERFACE MenuApplication MenuTransportCore MenuBoost MenuBoostJson)
