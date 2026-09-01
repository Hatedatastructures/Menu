import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Page {
    id: root
    title: "我的"

    property bool registering: false

    Flickable {
        anchors.fill: parent
        anchors.margins: 16
        contentWidth: width
        contentHeight: settingsColumn.height
        clip: true

        Column {
            id: settingsColumn
            width: parent.width
            spacing: 14

            Label {
                text: "我的设置"
                color: Theme.text
                font.pixelSize: 25
                font.bold: true
            }

            Label {
                visible: Boolean(MenuApi.ErrorMessage)
                width: parent.width
                text: MenuApi.ErrorMessage
                color: Theme.coral
                wrapMode: Text.WordWrap
            }

            Column {
                visible: !MenuApi.IsAuthenticated
                width: parent.width
                spacing: 10
                Label {
                    text: root.registering ? "创建账户" : "登录 Menu"
                    color: Theme.text
                    font.pixelSize: 18
                    font.bold: true
                }
                TextField {
                    id: emailField
                    width: parent.width
                    placeholderText: "邮箱"
                    inputMethodHints: Qt.ImhEmailCharactersOnly
                }
                TextField {
                    id: passwordField
                    width: parent.width
                    placeholderText: "密码（至少 12 个字符）"
                    echoMode: TextInput.Password
                }
                TextField {
                    id: displayNameField
                    visible: root.registering
                    width: parent.width
                    placeholderText: "称呼"
                }
                Row {
                    spacing: 10
                    Button {
                        text: root.registering ? "创建账户" : "登录"
                        onClicked: {
                            if (root.registering) {
                                MenuApi.Register(emailField.text, passwordField.text, displayNameField.text)
                            } else {
                                MenuApi.Login(emailField.text, passwordField.text)
                            }
                        }
                    }
                    Button {
                        text: root.registering ? "已有账户，去登录" : "首次使用，创建账户"
                        onClicked: root.registering = !root.registering
                    }
                }
            }

            Column {
                visible: MenuApi.IsAuthenticated
                width: parent.width
                spacing: 12
                Label {
                    text: "你好，" + MenuApi.DisplayName
                    color: Theme.text
                    font.pixelSize: 18
                    font.bold: true
                }
                Row {
                    width: parent.width
                    spacing: 12
                    Label {
                        text: "每餐人数"
                        width: 100
                        anchors.verticalCenter: parent.verticalCenter
                        color: Theme.mutedText
                    }
                    SpinBox {
                        value: MenuApi.ServingCount
                        from: 1
                        to: 24
                        onValueChanged: MenuApi.ServingCount = value
                    }
                }
                Row {
                    width: parent.width
                    spacing: 12
                    Label {
                        text: "可用时间"
                        width: 100
                        anchors.verticalCenter: parent.verticalCenter
                        color: Theme.mutedText
                    }
                    SpinBox {
                        value: MenuApi.AvailableMinutes
                        from: 10
                        to: 240
                        stepSize: 5
                        onValueChanged: MenuApi.AvailableMinutes = value
                    }
                    Label {
                        text: "分钟"
                        anchors.verticalCenter: parent.verticalCenter
                        color: Theme.mutedText
                    }
                }
                Label {
                    text: "忌口 / 过敏原"
                    color: Theme.mutedText
                }
                TextField {
                    id: allergiesField
                    width: parent.width
                    text: MenuApi.Allergies.join("、")
                    placeholderText: "例如：花生、虾"
                    onEditingFinished: {
                        var Values = text.length === 0 ? [] : text.split("、")
                        MenuApi.Allergies = Values
                    }
                }
                Label {
                    text: "家中常备食材"
                    color: Theme.mutedText
                }
                Flow {
                    width: parent.width
                    spacing: 4
                    Repeater {
                        model: MenuApi.Ingredients
                        delegate: CheckBox {
                            text: name
                            checked: MenuApi.PantryIngredientIds.indexOf(id) >= 0
                            onToggled: {
                                var Values = MenuApi.PantryIngredientIds.slice()
                                if (checked && Values.indexOf(id) < 0) {
                                    Values.push(id)
                                } else if (!checked) {
                                    Values = Values.filter(function(Value) { return Value !== id })
                                }
                                MenuApi.PantryIngredientIds = Values
                            }
                        }
                    }
                }
                Label {
                    text: "偏好菜系"
                    color: Theme.mutedText
                }
                Row {
                    spacing: 8
                    Repeater {
                        model: ["中餐", "西餐", "日系"]
                        delegate: CheckBox {
                            text: modelData
                            checked: MenuApi.PreferredCuisines.indexOf(modelData) >= 0
                            onToggled: {
                                var Values = MenuApi.PreferredCuisines
                                if (checked && Values.indexOf(modelData) < 0) {
                                    Values.push(modelData)
                                } else if (!checked) {
                                    Values = Values.filter(function(Value) { return Value !== modelData })
                                }
                                MenuApi.PreferredCuisines = Values
                            }
                        }
                    }
                }
                Label {
                    text: "可用厨具"
                    color: Theme.mutedText
                }
                Row {
                    spacing: 8
                    Repeater {
                        model: ["炒锅", "平底锅", "烤箱", "电饭锅"]
                        delegate: CheckBox {
                            text: modelData
                            checked: MenuApi.Cookware.indexOf(modelData) >= 0
                            onToggled: {
                                var Values = MenuApi.Cookware
                                if (checked && Values.indexOf(modelData) < 0) {
                                    Values.push(modelData)
                                } else if (!checked) {
                                    Values = Values.filter(function(Value) { return Value !== modelData })
                                }
                                MenuApi.Cookware = Values
                            }
                        }
                    }
                }
                Row {
                    width: parent.width
                    Label {
                        text: "有意义的备菜提醒"
                        color: Theme.text
                        width: parent.width - reminderSwitch.width
                        anchors.verticalCenter: parent.verticalCenter
                    }
                    Switch {
                        id: reminderSwitch
                        checked: MenuApi.RemindersEnabled
                        onToggled: MenuApi.RemindersEnabled = checked
                    }
                }
                Label {
                    text: "服务地址"
                    color: Theme.mutedText
                }
                TextField {
                    width: parent.width
                    text: MenuApi.BaseUrl
                    onEditingFinished: MenuApi.BaseUrl = text
                }
                Row {
                    spacing: 10
                    Button {
                        text: "刷新计划"
                        onClicked: MenuApi.LoadPlans()
                    }
                    Button {
                        text: "清除离线缓存"
                        onClicked: MenuApi.ClearCache()
                    }
                    Button {
                        text: "退出登录"
                        onClicked: MenuApi.Logout()
                    }
                }
            }
        }
    }
}
