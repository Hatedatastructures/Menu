import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "../Components"
import MenuClient 1.0

Page {
    id: root

    title: "我的"
    background: Rectangle { color: Theme.background }

    property string settingsSection: "account"

    Component.onCompleted: MenuNotifications.RefreshPermission()

    Connections {
        target: MenuApi
        function onAuthenticationChanged() {
            if (!MenuApi.IsAuthenticated) {
                root.settingsSection = "account"
            }
        }
    }

    ScrollView {
        anchors.fill: parent
        anchors.margins: Theme.pageMargin
        clip: true

        Column {
            id: content
            width: Math.min(root.width - Theme.pageMargin * 2, Theme.contentMaxWidth)
            anchors.horizontalCenter: parent.horizontalCenter
            spacing: 16

            RowLayout {
                width: parent.width
                spacing: 12

                Column {
                    Layout.fillWidth: true
                    spacing: 3
                    Label {
                        text: MenuApi.IsAuthenticated ? "你好，" + MenuApi.DisplayName : "你的厨房"
                        color: Theme.ink
                        font.pixelSize: 26
                        font.weight: Font.DemiBold
                    }
                    Label {
                        text: MenuApi.IsAuthenticated
                            ? "偏好、连接和提醒都在这里管理"
                            : "先连接服务，再登录同步你的厨房"
                        color: Theme.mutedText
                        font.pixelSize: 13
                        wrapMode: Text.WordWrap
                    }
                }
                AppIcon {
                    name: "settings"
                    width: 28
                    height: 28
                    color: Theme.green
                }
            }

            AppSegmentedControl {
                width: parent.width
                options: MenuApi.IsAuthenticated ? ["账户", "偏好", "设备"] : ["账户", "设备"]
                currentIndex: root.settingsSection === "account" ? 0
                    : root.settingsSection === "preferences" && MenuApi.IsAuthenticated ? 1
                    : MenuApi.IsAuthenticated ? 2 : 1
                onIndexActivated: {
                    if (index === 0) {
                        root.settingsSection = "account"
                    } else if (MenuApi.IsAuthenticated && index === 1) {
                        root.settingsSection = "preferences"
                    } else {
                        root.settingsSection = "device"
                    }
                }
            }

            Label {
                visible: Boolean(MenuApi.ErrorMessage)
                width: parent.width
                text: MenuApi.ErrorMessage
                color: Theme.coral
                wrapMode: Text.WordWrap
                font.pixelSize: 13
            }

            AuthPanel {
                width: parent.width
                active: root.settingsSection === "account"
            }
            AccountPanel {
                width: parent.width
                active: root.settingsSection === "account"
            }
            PreferencesPanel {
                width: parent.width
                active: root.settingsSection === "preferences"
            }
            Column {
                width: parent.width
                spacing: 16
                visible: root.settingsSection === "device"
                height: visible ? implicitHeight : 0
                ConnectionPanel { width: parent.width; active: root.settingsSection === "device" }
                DeviceNotificationPanel { width: parent.width; active: root.settingsSection === "device" }
            }

            Flow {
                width: parent.width
                spacing: 8
                AppButton {
                    text: "刷新数据"
                    iconName: "refresh"
                    variant: "secondary"
                    onClicked: MenuApi.LoadData()
                }
                AppButton {
                    text: "清除缓存"
                    variant: "ghost"
                    onClicked: MenuApi.ClearCache()
                }
            }
        }
    }
}
