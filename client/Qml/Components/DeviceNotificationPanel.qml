import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "."
import MenuClient 1.0

Rectangle {
    id: root

    property bool active: true
    width: parent ? parent.width : 360
    visible: active
    height: active ? systemColumn.implicitHeight + 32 : 0
    radius: Theme.radius
    color: Theme.surface
    border.color: Theme.line
    border.width: 1

    Column {
        id: systemColumn
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: parent.top
        anchors.margins: 16
        spacing: 14

        SectionHeading {
            width: parent.width
            title: "设备与通知"
            iconName: "settings"
        }

        RowLayout {
            width: parent.width
            height: 44
            Label {
                Layout.fillWidth: true
                text: "屏幕刷新率"
                color: Theme.ink
                font.pixelSize: 14
                font.weight: Font.DemiBold
            }
            Label {
                text: MenuDisplay.UsingFallback
                    ? "120 Hz 预算" : Math.round(MenuDisplay.ReportedRefreshRateHz) + " Hz"
                color: Theme.green
                font.pixelSize: 18
                font.weight: Font.DemiBold
            }
        }
        Label {
            width: parent.width
            text: MenuDisplay.UsingFallback
                ? "设备未提供屏幕指标；渲染预算按 120 Hz · " +
                  MenuDisplay.FrameBudgetMilliseconds.toFixed(2) + " ms/帧"
                : "跟随本机屏幕刷新 · " +
                  MenuDisplay.FrameBudgetMilliseconds.toFixed(2) + " ms/帧"
            color: Theme.mutedText
            font.pixelSize: 12
            wrapMode: Text.WordWrap
        }

        RowLayout {
            width: parent.width
            height: 44
            Label {
                Layout.fillWidth: true
                text: "界面主题"
                color: Theme.ink
                font.pixelSize: 14
            }
            AppSegmentedControl {
                Layout.preferredWidth: Math.min(240, parent.width * 0.62)
                options: ["跟随系统", "浅色", "深色"]
                currentIndex: MenuTheme.Mode === "light" ? 1
                    : MenuTheme.Mode === "dark" ? 2 : 0
                onIndexActivated: MenuTheme.Mode = index === 1 ? "light"
                    : index === 2 ? "dark" : "system"
            }
        }

        Rectangle {
            width: parent.width
            height: 1
            color: Theme.line
            opacity: 0.75
        }

        RowLayout {
            width: parent.width
            height: 44
            Label {
                Layout.fillWidth: true
                text: "系统通知"
                color: Theme.ink
                font.pixelSize: 14
            }
            Label {
                text: MenuNotifications.Status
                color: MenuNotifications.PermissionGranted ? Theme.green : Theme.mutedText
                font.pixelSize: 12
                elide: Text.ElideRight
            }
        }

        Flow {
            width: parent.width
            spacing: 8
            AppButton {
                visible: MenuNotifications.Supported && !MenuNotifications.PermissionGranted
                text: "允许通知"
                iconName: "check"
                variant: "accent"
                onClicked: MenuNotifications.RequestPermission()
            }
            AppButton {
                visible: MenuNotifications.Supported
                text: "通知设置"
                iconName: "settings"
                variant: "secondary"
                onClicked: MenuNotifications.OpenSettings()
            }
            AppButton {
                visible: MenuNotifications.Supported && MenuNotifications.PermissionGranted
                text: "测试通知"
                variant: "ghost"
                onClicked: MenuNotifications.ShowTestNotification()
            }
        }

        Label {
            width: parent.width
            visible: MenuNotifications.Supported
            text: "计时完成后由系统通知提醒；声音和震动可在系统渠道中调整。"
            color: Theme.mutedText
            font.pixelSize: 12
            wrapMode: Text.WordWrap
        }
    }
}
