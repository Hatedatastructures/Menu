import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import MenuClient 1.0

Rectangle {
    id: root

    property int currentIndex: 0
    signal indexActivated(int index)

    implicitWidth: 216
    color: Theme.surface
    border.color: Theme.line
    border.width: 1

    Column {
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: parent.top
        anchors.margins: 18
        spacing: 8

        Label {
            text: "浏览"
            color: Theme.mutedText
            font.pixelSize: 11
            font.weight: Font.DemiBold
            font.letterSpacing: 0.8
            leftPadding: 12
            bottomPadding: 6
        }

        Repeater {
            model: [
                { label: "今晚", icon: "tonight" },
                { label: "本周", icon: "week" },
                { label: "食材", icon: "ingredients" },
                { label: "我的", icon: "profile" }
            ]

            delegate: Button {
                width: parent.width
                height: 52
                flat: true
                padding: 0
                text: modelData.label
                Accessible.name: modelData.label
                Accessible.role: Accessible.Button

                contentItem: RowLayout {
                    anchors.fill: parent
                    anchors.leftMargin: 14
                    anchors.rightMargin: 14
                    spacing: 12

                    AppIcon {
                        Layout.preferredWidth: 22
                        Layout.preferredHeight: 22
                        name: modelData.icon
                        color: index === root.currentIndex ? Theme.green : Theme.mutedText
                    }
                    Label {
                        Layout.fillWidth: true
                        text: modelData.label
                        color: index === root.currentIndex ? Theme.green : Theme.ink
                        font.pixelSize: 14
                        font.weight: index === root.currentIndex ? Font.DemiBold : Font.Normal
                        verticalAlignment: Text.AlignVCenter
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

                onClicked: {
                    root.currentIndex = index
                    root.indexActivated(index)
                }
            }
        }

        Rectangle {
            width: parent.width
            height: 1
            color: Theme.line
            opacity: 0.7
        }

        Label {
            width: parent.width
            text: "做饭时会自动隐藏导航"
            color: Theme.mutedText
            font.pixelSize: 11
            wrapMode: Text.WordWrap
            leftPadding: 12
            rightPadding: 12
        }
    }
}
