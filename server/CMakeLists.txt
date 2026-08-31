add_subdirectory(Foundation)
add_subdirectory(Domain)
add_subdirectory(Application)
add_subdirectory(Infrastructure)
add_subdirectory(Transport)
add_subdirectory(Api)

add_executable(MenuServer Main.cpp)
target_link_libraries(MenuServer PRIVATE
    MenuApi
    MenuTransport
    MenuInfrastructure
    MenuApplication
)
set_target_properties(MenuServer PROPERTIES
    RUNTIME_OUTPUT_DIRECTORY "${PROJECT_BINARY_DIR}/bin"
)
