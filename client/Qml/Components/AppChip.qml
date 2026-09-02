import QtQuick
import QtQuick.Controls
import MenuClient 1.0

Item {
    id: root

    property string label: ""
    property color accent: Theme.green
    property bool selected: false
    implicitWidth: chipLabel.implicitWidth + 24
    implicitHeight: 30

    Rectangle {
        anchors.fill: parent
        radius: height / 2
        color: root.selected ? root.accent : Theme.tint(root.accent, Theme.dark ? 0.18 : 0.12)
        border.width: root.selected ? 0 : 1
        border.color: Theme.tint(root.accent, 0.36)
    }

    Label {
        id: chipLabel
        anchors.centerIn: parent
        text: root.label
        color: root.selected ? Theme.onColor(root.accent) : root.accent
        font.pixelSize: 12
        font.weight: Font.DemiBold
    }
}
