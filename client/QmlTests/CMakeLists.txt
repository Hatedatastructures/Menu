add_executable(MenuClientQmlTests
    QmlTestMain.cpp
)
target_link_libraries(MenuClientQmlTests PRIVATE
    Qt6::QuickTest
    Qt6::Qml
    Qt6::Quick
)
target_compile_features(MenuClientQmlTests PRIVATE cxx_std_20)
target_compile_definitions(MenuClientQmlTests PRIVATE
    QUICK_TEST_SOURCE_DIR="${CMAKE_CURRENT_SOURCE_DIR}"
)
add_test(NAME MenuClientQmlTests COMMAND MenuClientQmlTests)

if(WIN32)
    get_filename_component(MenuQtCMakeDirectory "${Qt6Core_DIR}" DIRECTORY)
    get_filename_component(MenuQtCMakeRoot "${MenuQtCMakeDirectory}" DIRECTORY)
    get_filename_component(MenuQtInstallRoot "${MenuQtCMakeRoot}" DIRECTORY)
    set_tests_properties(MenuClientQmlTests PROPERTIES
        ENVIRONMENT "PATH=${MenuQtInstallRoot}/bin"
    )
endif()
