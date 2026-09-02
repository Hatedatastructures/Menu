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
    height: visible ? preferencesColumn.implicitHeight + 32 : 0
    radius: Theme.radius
    color: Theme.surface
    border.color: Theme.line
    border.width: 1

    Column {
        id: preferencesColumn
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: parent.top
        anchors.margins: 16
        spacing: 14

        SectionHeading {
            width: parent.width
            title: "烹饪偏好"
            iconName: "ingredients"
        }

        RowLayout {
            width: parent.width
            height: Theme.touchTarget
            Label {
                Layout.fillWidth: true
                text: "每餐人数"
                color: Theme.ink
                font.pixelSize: 14
            }
            AppStepper {
                value: MenuApi.ServingCount
                from: 1
                to: 24
                onValueChanged: MenuApi.ServingCount = value
            }
        }

        RowLayout {
            width: parent.width
            height: Theme.touchTarget
            Label {
                Layout.fillWidth: true
                text: "可用时间"
                color: Theme.ink
                font.pixelSize: 14
            }
            AppStepper {
                id: availableTimeStepper
                value: MenuApi.AvailableMinutes
                from: 10
                to: 240
                stepSize: 5
                onValueChanged: MenuApi.AvailableMinutes = value
            }
            Label {
                text: "分钟"
                color: Theme.mutedText
                font.pixelSize: 13
            }
        }

        Label {
            text: "忌口 / 过敏原"
            color: Theme.mutedText
            font.pixelSize: 12
            font.weight: Font.DemiBold
        }
        AppTextField {
            width: parent.width
            text: MenuApi.Allergies.join("、")
            placeholderText: "例如：花生、虾"
            onEditingFinished: MenuApi.Allergies = text.length === 0 ? [] : text.split("、")
        }

        Label {
            text: "偏好菜系"
            color: Theme.mutedText
            font.pixelSize: 12
            font.weight: Font.DemiBold
        }
        Flow {
            width: parent.width
            spacing: 4
            Repeater {
                model: ["中餐", "西餐", "日系"]
                delegate: AppCheckBox {
                    text: modelData
                    checked: MenuApi.PreferredCuisines.indexOf(modelData) >= 0
                    onToggled: {
                        var values = MenuApi.PreferredCuisines.slice()
                        if (checked && values.indexOf(modelData) < 0) values.push(modelData)
                        else if (!checked) values = values.filter(function(value) { return value !== modelData })
                        MenuApi.PreferredCuisines = values
                    }
                }
            }
        }

        Label {
            text: "家中常备食材"
            color: Theme.mutedText
            font.pixelSize: 12
            font.weight: Font.DemiBold
        }
        Flow {
            width: parent.width
            spacing: 4
            Repeater {
                model: MenuApi.Ingredients
                delegate: AppCheckBox {
                    text: name
                    checked: MenuApi.PantryIngredientIds.indexOf(id) >= 0
                    onToggled: {
                        var values = MenuApi.PantryIngredientIds.slice()
                        if (checked && values.indexOf(id) < 0) values.push(id)
                        else if (!checked) values = values.filter(function(value) { return value !== id })
                        MenuApi.PantryIngredientIds = values
                    }
                }
            }
        }

        RowLayout {
            width: parent.width
            height: Theme.touchTarget
            Label {
                Layout.fillWidth: true
                text: "备菜提醒"
                color: Theme.ink
                font.pixelSize: 14
            }
            AppSwitch {
                checked: MenuApi.RemindersEnabled
                accessibleName: "备菜提醒"
                onToggled: MenuApi.RemindersEnabled = checked
            }
        }
    }
}
