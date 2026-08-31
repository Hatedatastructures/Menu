add_library(MenuTransportCore INTERFACE)
target_include_directories(MenuTransportCore INTERFACE "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuTransportCore INTERFACE MenuFoundation)

add_library(MenuTransport STATIC
    HttpServer.cpp
    HttpSession.cpp
)
target_include_directories(MenuTransport PUBLIC "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuTransport PUBLIC MenuTransportCore MenuBoost)
target_compile_features(MenuTransport PUBLIC cxx_std_20)
