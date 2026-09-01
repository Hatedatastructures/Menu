import QtQuick
import QtQuick.Controls

Page {
    id: root
    title: "今晚"

    Component.onCompleted: MenuApi.LoadData()

    Column {
        anchors.fill: parent
        anchors.margins: 16
        spacing: 12

        Row {
            width: parent.width
            spacing: 12

            Column {
                width: parent.width - refreshButton.width - 12
                spacing: 3

                Label {
                    text: "今晚吃什么"
                    color: Theme.text
                    font.pixelSize: 25
                    font.bold: true
                }
                Label {
                    text: MenuApi.IsLoading ? "正在寻找合适的方案" : "按时间、偏好和家中食材筛选"
                    color: Theme.mutedText
                    font.pixelSize: 13
                }
            }

            Button {
                id: refreshButton
                text: "重新推荐"
                enabled: !MenuApi.IsLoading
                onClicked: MenuApi.LoadData()
            }
        }

        Rectangle {
            width: parent.width
            height: MenuApi.IsOffline ? 38 : 0
            visible: MenuApi.IsOffline
            color: "#fff0e8"
            radius: 4

            Label {
                anchors.centerIn: parent
                text: "当前离线，显示上次可用方案"
                color: Theme.coral
                font.pixelSize: 13
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

        ListView {
            id: recommendationList
            width: parent.width
            height: parent.height - y
            spacing: 8
            clip: true
            model: MenuApi.Recommendations
            delegate: RecipeCard {
                title: name
                subtitle: cuisine + " · " + servings + " 人份"
                meta: totalMinutes + " 分钟 · 难度 " + difficulty
                imageSource: mediaUrl
                availableCount: availableIngredientCount
                missingCount: missingIngredientCount
                onClicked: MenuApi.OpenRecipe(index)
            }
        }

        Label {
            visible: !MenuApi.IsLoading && !MenuApi.ErrorMessage && MenuApi.Recommendations.Count === 0
            width: parent.width
            text: "暂时没有符合条件的方案"
            color: Theme.mutedText
            horizontalAlignment: Text.AlignHCenter
            font.pixelSize: 15
        }
    }
}
