import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import MenuClient 1.0

RowLayout {
    id: root

    property string title: ""
    property string actionText: ""
    property string iconName: ""
    signal actionClicked()

    width: parent ? parent.width : 360
    height: Theme.touchTarget
    spacing: 10

    AppIcon {
        visible: root.iconName.length > 0
        Layout.preferredWidth: 20
        Layout.preferredHeight: 20
        name: root.iconName
        color: Theme.green
    }

    Label {
        Layout.fillWidth: true
        text: root.title
        color: Theme.ink
        font.pixelSize: 18
        font.weight: Font.DemiBold
        elide: Text.ElideRight
        verticalAlignment: Text.AlignVCenter
    }

    AppButton {
        visible: root.actionText.length > 0
        text: root.actionText
        variant: "ghost"
        Layout.minimumWidth: 48
        Layout.preferredHeight: 40
        onClicked: root.actionClicked()
    }
}
