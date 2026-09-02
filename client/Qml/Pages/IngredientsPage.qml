import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "../Components"
import MenuClient 1.0

Page {
    id: root
    title: "食材"
    background: Rectangle { color: Theme.background }

    property string filterText: ""

    function togglePantry(IngredientId) {
        var values = MenuApi.PantryIngredientIds.slice()
        var position = values.indexOf(IngredientId)
        if (position >= 0) {
            values.splice(position, 1)
        } else {
            values.push(IngredientId)
        }
        MenuApi.PantryIngredientIds = values
    }

    ListModel {
        id: filteredIngredients
    }

    function rebuildFilteredIngredients() {
        filteredIngredients.clear()
        var needle = root.filterText
        for (var index = 0; index < MenuApi.Ingredients.Count; ++index) {
            var ingredient = MenuApi.Ingredients.IngredientAt(index)
            if (needle.length > 0 &&
                ingredient.name.toLowerCase().indexOf(needle) < 0 &&
                ingredient.category.toLowerCase().indexOf(needle) < 0) {
                continue
            }
            filteredIngredients.append(ingredient)
        }
    }

    Connections {
        target: MenuApi.Ingredients
        function onCountChanged() { root.rebuildFilteredIngredients() }
    }

    Component.onCompleted: root.rebuildFilteredIngredients()

    ColumnLayout {
        anchors.top: parent.top
        anchors.bottom: parent.bottom
        anchors.margins: Theme.pageMargin
        width: Math.min(parent.width - Theme.pageMargin * 2, Theme.contentMaxWidth)
        anchors.horizontalCenter: parent.horizontalCenter
        spacing: 12

        RowLayout {
            Layout.fillWidth: true
            spacing: 10
            Column {
                Layout.fillWidth: true
                spacing: 3
                Label {
                    text: "厨房库存"
                    color: Theme.ink
                    font.pixelSize: 24
                    font.weight: Font.DemiBold
                }
                Label {
                    text: MenuApi.Ingredients.Count + " 项食材 · 随手掌握缺口"
                    color: Theme.mutedText
                    font.pixelSize: 13
                }
            }
            AppChip {
                label: MenuApi.PantryIngredientIds.length + " 常备"
                accent: Theme.green
            }
        }

        AppTextField {
            Layout.fillWidth: true
            placeholderText: "搜索食材或分类"
            onTextChanged: {
                root.filterText = text.trim().toLowerCase()
                root.rebuildFilteredIngredients()
            }
        }

        ListView {
            id: ingredientsList
            visible: filteredIngredients.count > 0
            Layout.fillWidth: true
            Layout.fillHeight: true
            spacing: 8
            clip: true
            model: filteredIngredients
            delegate: Rectangle {
                property bool inPantry: MenuApi.PantryIngredientIds.indexOf(id) >= 0
                width: ingredientsList.width
                height: 72
                radius: Theme.radius
                color: Theme.surface
                border.color: Theme.line
                border.width: 1

                RowLayout {
                    anchors.fill: parent
                    anchors.leftMargin: 14
                    anchors.rightMargin: 14
                    spacing: 12
                    Label {
                        Layout.preferredWidth: 96
                        text: name
                        color: Theme.ink
                        font.pixelSize: 16
                        font.weight: Font.DemiBold
                        elide: Text.ElideRight
                    }
                    Column {
                        Layout.fillWidth: true
                        spacing: 3
                        Label {
                            text: category + " · " + defaultUnit
                            color: Theme.mutedText
                            font.pixelSize: 13
                            elide: Text.ElideRight
                        }
                        Label {
                            text: inPantry ? "常备 · 不计入缺口" : "按需 · 加入准备清单"
                            color: inPantry ? Theme.green : Theme.mutedText
                            font.pixelSize: 12
                        }
                    }
                    AppButton {
                        Layout.preferredWidth: 104
                        Layout.preferredHeight: 40
                        text: inPantry ? "常备中" : "设为常备"
                        iconName: inPantry ? "check" : "plus"
                        variant: inPantry ? "secondary" : "ghost"
                        Accessible.name: inPantry ? name + " 已是常备" : "将 " + name + " 设为常备"
                        onClicked: root.togglePantry(id)
                    }
                }
            }
        }

        Label {
            Layout.fillWidth: true
            Layout.fillHeight: true
            visible: filteredIngredients.count === 0
            text: root.filterText.length > 0 ? "没有匹配的食材" : "还没有食材数据"
            color: Theme.mutedText
            horizontalAlignment: Text.AlignHCenter
            verticalAlignment: Text.AlignVCenter
        }
    }
}
