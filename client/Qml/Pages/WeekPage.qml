import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "../Components"
import MenuClient 1.0

Page {
    id: root
    title: "本周"

    background: Rectangle { color: Theme.background }

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
        anchors.top: parent.top
        anchors.bottom: parent.bottom
        anchors.margins: Theme.pageMargin
        width: Math.min(parent.width - Theme.pageMargin * 2, Theme.contentMaxWidth)
        anchors.horizontalCenter: parent.horizontalCenter
        spacing: 14

        Row {
            width: parent.width
            spacing: 10
            Column {
                width: parent.width - dateField.width - 10
                spacing: 3
                Label {
                    text: "一周安排"
                    color: Theme.text
                    font.pixelSize: 24
                    font.weight: Font.DemiBold
                }
                Label {
                    text: MenuApi.IsAuthenticated ? "选好菜，食材自动合并" : "登录后保存跨设备计划"
                    color: Theme.mutedText
                    font.pixelSize: 13
                }
            }
            AppTextField {
                id: dateField
                width: 132
                enabled: root.authenticated
                opacity: root.authenticated ? 1 : 0.55
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
            text: "还没有安排，从今晚的方案开始添加"
            color: Theme.mutedText
        }

        Rectangle {
            visible: !root.authenticated
            width: parent.width
            height: 72
            radius: Theme.radius
            color: Theme.tint(Theme.amber, Theme.dark ? 0.14 : 0.10)
            Label {
                anchors.fill: parent
                anchors.margins: 14
                text: "登录后可保存一周计划，并在手机与电脑间同步。"
                color: Theme.amber
                font.pixelSize: 13
                wrapMode: Text.WordWrap
                verticalAlignment: Text.AlignVCenter
            }
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
                AppButton {
                    text: "加入"
                    iconText: "+"
                    variant: "ghost"
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
            delegate: PlanDayCard {
                width: ListView.view.width
                planDateValue: planDate
                planItems: items || []
                combinedItems: combinedIngredients || []
                onRecipeRemoved: root.removeRecipe(planDateValue, recipeId)
            }
        }
    }
}
