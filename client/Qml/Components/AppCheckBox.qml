import QtQuick
import QtQuick.Controls
import MenuClient 1.0

CheckBox {
    id: root

    implicitHeight: Theme.touchTarget
    leftPadding: 34
    rightPadding: 4
    topPadding: 4
    bottomPadding: 4
    Accessible.name: root.text

    contentItem: Label {
        text: root.text
        color: root.enabled ? Theme.ink : Theme.mutedText
        font.pixelSize: 14
        verticalAlignment: Text.AlignVCenter
        elide: Text.ElideRight
    }

    indicator: Rectangle {
        x: root.leftPadding - width
        y: (root.height - height) / 2
        width: 24
        height: 24
        radius: 5
        color: root.checked ? Theme.greenStrong : "transparent"
        border.width: root.checked ? 0 : 2
        border.color: root.enabled ? Theme.line : Theme.mutedText

        Label {
            anchors.centerIn: parent
            visible: root.checked
            text: "✓"
            color: Theme.onColor(Theme.greenStrong)
            font.pixelSize: 17
            font.weight: Font.DemiBold
        }
    }
}
