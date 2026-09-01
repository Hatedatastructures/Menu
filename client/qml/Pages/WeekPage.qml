import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Page {
    id: root
    title: "本周"

    DateFormatter {
        id: dateFormatter
    }

    property string selectedDate: dateFormatter.formatLocalDate(new Date())
    property bool authenticated: MenuApi.IsAuthenticated

    function recipeIdsForDate(DateValue) {
        var Result = []
        for (var PlanIndex = 0; PlanIndex < MenuApi.Plans.Count; ++PlanIndex) {
            var Plan = MenuApi.Plans.PlanAt(PlanIndex)
            if (Plan.planDate !== DateValue) {
                continue
            }
            var Items = Plan.items || []
            for (var ItemIndex = 0; ItemIndex < Items.length; ++ItemIndex) {
                Result.push(Items[ItemIndex].recipeId)
            }
        }
        return Result
    }

    function addRecipe(RecipeId) {
        var RecipeIds = recipeIdsForDate(root.selectedDate)
        if (RecipeIds.indexOf(RecipeId) < 0) {
            RecipeIds.push(RecipeId)
            MenuApi.SavePlan(root.selectedDate, RecipeIds)
        }
    }

    function removeRecipe(DateValue, RecipeId) {
        var RecipeIds = recipeIdsForDate(DateValue)
        var Remaining = []
        for (var Index = 0; Index < RecipeIds.length; ++Index) {
            if (RecipeIds[Index] !== RecipeId) {
                Remaining.push(RecipeIds[Index])
            }
        }
        MenuApi.SavePlan(DateValue, Remaining)
    }

    Component.onCompleted: {
        if (root.authenticated) {
            MenuApi.LoadPlans()
        }
    }

    onAuthenticatedChanged: {
        if (root.authenticated) {
            MenuApi.LoadPlans()
        }
    }

    Column {
        anchors.fill: parent
        anchors.margins: 16
        spacing: 10

        Row {
            width: parent.width
            spacing: 10
            Column {
                width: parent.width - dateField.width - 10
                spacing: 3
                Label {
                    text: "本周计划"
                    color: Theme.text
                    font.pixelSize: 25
                    font.bold: true
                }
                Label {
                    text: MenuApi.IsAuthenticated ? "选好菜，清单会自动合并" : "登录后可跨设备保存计划"
                    color: Theme.mutedText
                    font.pixelSize: 13
                }
            }
            TextField {
                id: dateField
                width: 132
                text: root.selectedDate
                placeholderText: "YYYY-MM-DD"
                onEditingFinished: root.selectedDate = text
            }
        }

        Label {
            visible: Boolean(MenuApi.ErrorMessage)
            width: parent.width
            text: MenuApi.ErrorMessage
            color: Theme.coral
            wrapMode: Text.WordWrap
        }

        Label {
            visible: MenuApi.IsAuthenticated && MenuApi.Plans.Count === 0
            width: parent.width
            text: "还没有计划，从今晚的方案开始添加"
            color: Theme.mutedText
        }

        Label {
            visible: MenuApi.IsAuthenticated && MenuApi.Recommendations.Count > 0
            text: "今晚方案"
            color: Theme.mutedText
            font.pixelSize: 12
        }

        ListView {
            visible: MenuApi.IsAuthenticated && MenuApi.Recommendations.Count > 0
            width: parent.width
            height: Math.min(132, MenuApi.Recommendations.Count * 44)
            spacing: 1
            clip: true
            model: MenuApi.Recommendations
            delegate: RowLayout {
                width: ListView.view.width
                height: 42
                Label {
                    Layout.fillWidth: true
                    text: name + " · " + totalMinutes + " 分钟"
                    color: Theme.text
                    elide: Text.ElideRight
                }
                Button {
                    text: "加入 " + root.selectedDate
                    onClicked: root.addRecipe(id)
                }
            }
        }

        ListView {
            width: parent.width
            height: parent.height - y
            spacing: 8
            clip: true
            model: MenuApi.Plans
            delegate: Rectangle {
                id: planDay
                width: ListView.view.width
                height: 168 + (planItems.length * 42) + Math.min(152, combinedItems.length * 38)
                color: Theme.surface
                border.color: Theme.line
                border.width: 1
                radius: 4
                property string planDateValue: planDate
                property var planItems: items || []
                property var combinedItems: combinedIngredients || []

                Column {
                    anchors.fill: parent
                    anchors.margins: 14
                    spacing: 8
                    RowLayout {
                        width: parent.width
                        Label {
                            Layout.fillWidth: true
                            text: planDay.planDateValue
                            color: Theme.text
                            font.pixelSize: 17
                            font.bold: true
                        }
                        Label {
                            text: planDay.planItems.length + " 道菜"
                            color: Theme.green
                            font.pixelSize: 13
                        }
                    }
                    Repeater {
                        model: planDay.planItems
                        delegate: RowLayout {
                            width: parent.width
                            height: 34
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
                            Button {
                                text: "移除"
                                onClicked: root.removeRecipe(planDay.planDateValue, modelData.recipeId)
                            }
                        }
                    }
                    Label {
                        text: "合并清单"
                        color: Theme.mutedText
                        font.pixelSize: 12
                    }
                    ListView {
                        width: parent.width
                        height: Math.min(152, planDay.combinedItems.length * 38)
                        visible: planDay.combinedItems.length > 0
                        clip: true
                        model: planDay.combinedItems
                        delegate: CheckBox {
                            width: ListView.view.width
                            height: 38
                            text: (modelData.ingredientName || modelData.ingredientId) + " " +
                                modelData.quantity + modelData.unit +
                                (modelData.isPantryStaple ? " · 已有/常备" : " · 缺少/需准备")
                            checked: modelData.isPantryStaple
                        }
                    }
                    Label {
                        visible: planDay.combinedItems.length === 0
                        width: parent.width
                        text: "暂无合并食材"
                        color: Theme.text
                        font.pixelSize: 13
                    }
                }
            }
        }
    }
}
