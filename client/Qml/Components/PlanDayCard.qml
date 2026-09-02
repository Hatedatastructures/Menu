import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "."
import MenuClient 1.0

Rectangle {
    id: root

    property string planDateValue: ""
    property var planItems: []
    property var combinedItems: []

    signal recipeRemoved(string recipeId)

    width: parent ? parent.width : 360
    height: dayContent.implicitHeight + 28
    color: Theme.surface
    border.color: Theme.line
    border.width: 1
    radius: Theme.radius

    Column {
        id: dayContent
        anchors.fill: parent
        anchors.margins: 14
        spacing: 8

        RowLayout {
            width: parent.width
            height: Theme.touchTarget
            Label {
                Layout.fillWidth: true
                text: root.planDateValue
                color: Theme.text
                font.pixelSize: 17
                font.weight: Font.DemiBold
                elide: Text.ElideRight
            }
            AppChip {
                label: root.planItems.length + " 道菜"
                accent: Theme.green
            }
        }

        Repeater {
            model: root.planItems
            delegate: RowLayout {
                width: parent.width
                height: Theme.touchTarget
                Label {
                    Layout.fillWidth: true
                    text: modelData.recipeName || modelData.recipeId
                    color: Theme.text
                    elide: Text.ElideRight
                }
                Label {
                    text: modelData.servings + " 人份"
                    color: Theme.mutedText
                    font.pixelSize: 12
                }
                AppButton {
                    text: "移除"
                    variant: "ghost"
                    onClicked: root.recipeRemoved(modelData.recipeId)
                }
            }
        }

        Label {
            text: "合并清单"
            color: Theme.mutedText
            font.pixelSize: 12
            font.weight: Font.DemiBold
        }

        ListView {
            width: parent.width
            height: Math.min(192, root.combinedItems.length * 48)
            visible: root.combinedItems.length > 0
            clip: true
            model: root.combinedItems
            delegate: AppCheckBox {
                width: ListView.view.width
                height: 48
                enabled: false
                text: (modelData.ingredientName || modelData.ingredientId) + " " +
                    modelData.quantity + modelData.unit +
                    (modelData.isPantryStaple ? " · 已有/常备" : " · 缺少/需准备")
                checked: modelData.isPantryStaple
            }
        }

        Label {
            visible: root.combinedItems.length === 0
            width: parent.width
            text: "暂无合并食材"
            color: Theme.text
            font.pixelSize: 13
        }
    }
}
