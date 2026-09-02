import QtQuick
import QtQuick.Controls
import "."
import MenuClient 1.0

Rectangle {
    id: root

    property bool active: true
    property bool registering: false
    width: parent ? parent.width : 360
    visible: active && !MenuApi.IsAuthenticated
    height: visible ? loginColumn.implicitHeight + 28 : 0
    radius: Theme.radius
    color: Theme.surface
    border.color: Theme.line
    border.width: 1

    Column {
        id: loginColumn
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: parent.top
        anchors.margins: 14
        spacing: 10

        Row {
            width: parent.width
            spacing: 8
            Label {
                text: root.registering ? "创建账户" : "登录 Menu"
                color: Theme.ink
                font.pixelSize: 18
                font.weight: Font.DemiBold
            }
            AppChip { label: "本地优先"; accent: Theme.green }
        }

        AppTextField {
            id: emailField
            width: parent.width
            height: Theme.touchTarget
            placeholderText: "邮箱"
            inputMethodHints: Qt.ImhEmailCharactersOnly
            Accessible.name: "邮箱"
        }
        AppTextField {
            id: passwordField
            width: parent.width
            height: Theme.touchTarget
            placeholderText: "密码（至少 12 个字符）"
            echoMode: TextInput.Password
            Accessible.name: "密码"
        }
        AppTextField {
            id: displayNameField
            visible: root.registering
            width: parent.width
            height: Theme.touchTarget
            placeholderText: "称呼"
            Accessible.name: "称呼"
        }
        Row {
            spacing: 8
            AppButton {
                text: root.registering ? "创建账户" : "登录"
                iconName: "check"
                variant: "primary"
                onClicked: root.registering
                    ? MenuApi.Register(emailField.text, passwordField.text, displayNameField.text)
                    : MenuApi.Login(emailField.text, passwordField.text)
            }
            AppButton {
                text: root.registering ? "返回登录" : "创建账户"
                iconName: "back"
                variant: "ghost"
                onClicked: root.registering = !root.registering
            }
        }
    }
}
