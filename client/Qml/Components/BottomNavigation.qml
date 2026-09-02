import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import MenuClient 1.0

Rectangle {
    id: root

    property int currentIndex: 0

    implicitHeight: 84
    color: Theme.surface
    border.color: Theme.line
    border.width: 1

    RowLayout {
        anchors.fill: parent
        anchors.leftMargin: 6
        anchors.rightMargin: 6
        spacing: 4

        Repeater {
            model: [
                { label: "今晚", icon: "tonight" },
                { label: "本周", icon: "week" },
                { label: "食材", icon: "ingredients" },
                { label: "我的", icon: "profile" }
            ]

            delegate: Button {
                Layout.fillWidth: true
                Layout.fillHeight: true
                Layout.minimumWidth: 0
                flat: true
                padding: 4
                text: modelData.label
                Accessible.name: modelData.label
                Accessible.role: Accessible.Button

                contentItem: Column {
                    spacing: 4
                    anchors.centerIn: parent

                    AppIcon {
                        anchors.horizontalCenter: parent.horizontalCenter
                        width: 22
                        height: 22
                        name: modelData.icon
                        color: index === root.currentIndex ? Theme.green : Theme.mutedText
                    }
                    Label {
                        text: modelData.label
                        color: index === root.currentIndex ? Theme.green : Theme.mutedText
                        font.pixelSize: 12
                        font.weight: index === root.currentIndex ? Font.DemiBold : Font.Normal
                        horizontalAlignment: Text.AlignHCenter
                    }
                }

                background: Rectangle {
                    radius: Theme.radius
                    color: index === root.currentIndex
                        ? Theme.tint(Theme.green, Theme.dark ? 0.22 : 0.10)
                        : "transparent"
                    border.width: root.activeFocus && index === root.currentIndex ? 2 : 0
                    border.color: Theme.focus
                }

                onClicked: root.currentIndex = index
            }
        }
    }
}
