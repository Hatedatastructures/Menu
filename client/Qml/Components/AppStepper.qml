import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "."
import MenuClient 1.0

Item {
    id: root

    property int value: 0
    property int from: 0
    property int to: 100
    property int stepSize: 1

    implicitWidth: 156
    implicitHeight: Theme.touchTarget

    RowLayout {
        anchors.fill: parent
        spacing: 4

        AppButton {
            Layout.preferredWidth: 44
            Layout.fillHeight: true
            text: ""
            iconName: "minus"
            variant: "secondary"
            enabled: root.value - root.stepSize >= root.from
            Accessible.name: "减少"
            onClicked: root.value = Math.max(root.from, root.value - root.stepSize)
        }

        Label {
            Layout.fillWidth: true
            Layout.fillHeight: true
            text: root.value
            color: Theme.ink
            font.pixelSize: 16
            font.weight: Font.DemiBold
            horizontalAlignment: Text.AlignHCenter
            verticalAlignment: Text.AlignVCenter
            background: Rectangle {
                radius: Theme.radiusSmall
                color: Theme.surface
                border.color: Theme.line
                border.width: 1
            }
        }

        AppButton {
            Layout.preferredWidth: 44
            Layout.fillHeight: true
            text: ""
            iconName: "plus"
            variant: "secondary"
            enabled: root.value + root.stepSize <= root.to
            Accessible.name: "增加"
            onClicked: root.value = Math.min(root.to, root.value + root.stepSize)
        }
    }
}
