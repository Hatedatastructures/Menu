import QtQuick

QtObject {
    property int navigationIndex: 0
    property bool cookingVisible: false
    readonly property int pageIndex: cookingVisible ? 4 : navigationIndex
}
