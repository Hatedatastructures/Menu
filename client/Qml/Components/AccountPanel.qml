import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "."
import MenuClient 1.0

Rectangle {
    id: root

    property bool active: true
    width: parent ? parent.width : 360
    visible: active && MenuApi.IsAuthenticated
    height: visible ? accountColumn.implicitHeight + 32 : 0
    radius: Theme.radius
    color: Theme.surface
    border.color: Theme.line
    border.width: 1

    Column {
        id: accountColumn
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: parent.top
        anchors.margins: 16
        spacing: 14

        SectionHeading {
            width: parent.width
            title: "账户"
            iconName: "profile"
        }

        RowLayout {
            width: parent.width
            spacing: 12

            Rectangle {
                Layout.preferredWidth: 48
                Layout.preferredHeight: 48
                radius: 24
                color: Theme.tint(Theme.green, Theme.dark ? 0.20 : 0.10)

                AppIcon {
                    anchors.centerIn: parent
                    name: "profile"
                    width: 28
                    height: 28
                    color: Theme.green
                }
            }

            Column {
                Layout.fillWidth: true
                spacing: 3

                Label {
                    width: parent.width
                    text: MenuApi.DisplayName
                    color: Theme.ink
                    font.pixelSize: 18
                    font.weight: Font.DemiBold
                    elide: Text.ElideRight
                }
                Label {
                    text: "已登录 · 偏好和计划会同步"
                    color: Theme.mutedText
                    font.pixelSize: 13
                }
            }

            AppChip { label: "已连接"; accent: Theme.green }
        }

        Rectangle {
            width: parent.width
            height: 1
            color: Theme.line
            opacity: 0.75
        }

        RowLayout {
            width: parent.width
            height: Theme.touchTarget

            Column {
                Layout.fillWidth: true
                spacing: 2
                Label {
                    text: "本机数据"
                    color: Theme.ink
                    font.pixelSize: 14
                }
                Label {
                    text: "离线缓存和做饭进度保留在本机"
                    color: Theme.mutedText
                    font.pixelSize: 12
                }
            }
            AppButton {
                text: "退出登录"
                variant: "ghost"
                onClicked: MenuApi.Logout()
            }
        }
    }
}
