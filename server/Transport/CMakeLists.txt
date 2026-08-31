add_library(MenuTransportCore INTERFACE)
target_include_directories(MenuTransportCore INTERFACE "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuTransportCore INTERFACE MenuFoundation)

add_library(MenuTransport INTERFACE)
target_include_directories(MenuTransport INTERFACE "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuTransport INTERFACE MenuTransportCore MenuBoost)
