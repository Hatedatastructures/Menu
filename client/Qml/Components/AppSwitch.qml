import QtQuick
import QtQuick.Controls
import MenuClient 1.0

Switch {
    id: root

    implicitWidth: 56
    implicitHeight: Theme.touchTarget
    width: implicitWidth
    height: implicitHeight
    text: ""
    Accessible.name: root.accessibleName.length > 0 ? root.accessibleName : "开关"
    Accessible.role: Accessible.CheckBox

    property string accessibleName: ""

    contentItem: Item {}

    indicator: Item {
        width: 52
        height: 32
        x: (root.width - width) / 2
        y: (root.height - height) / 2

        Rectangle {
            anchors.fill: parent
            radius: height / 2
            color: root.checked ? Theme.greenStrong : Theme.surfaceRaised
            border.width: root.checked ? 0 : 1
            border.color: Theme.line
        }

        Rectangle {
            width: 24
            height: 24
            radius: 12
            y: 4
            x: root.checked ? parent.width - width - 4 : 4
            color: root.checked ? Theme.onAccent : Theme.mutedText
            Behavior on x { NumberAnimation { duration: 120; easing.type: Easing.OutCubic } }
        }

        Rectangle {
            anchors.fill: parent
            radius: height / 2
            color: "transparent"
            border.width: root.activeFocus ? 2 : 0
            border.color: Theme.focus
        }
    }
}
