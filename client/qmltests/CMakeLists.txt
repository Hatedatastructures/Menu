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
