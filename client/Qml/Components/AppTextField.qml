import QtQuick
import QtQuick.Controls
import MenuClient 1.0

TextField {
    id: root

    implicitHeight: Theme.touchTarget
    leftPadding: 14
    rightPadding: 14
    topPadding: 8
    bottomPadding: 8
    color: Theme.ink
    font.pixelSize: 15
    verticalAlignment: TextInput.AlignVCenter
    placeholderTextColor: Theme.mutedText
    selectionColor: Theme.tint(Theme.green, 0.35)
    selectedTextColor: Theme.ink

    background: Rectangle {
        radius: Theme.radius
        color: Theme.surface
        border.width: root.activeFocus ? 2 : 1
        border.color: root.activeFocus ? Theme.green : Theme.line
    }
}
