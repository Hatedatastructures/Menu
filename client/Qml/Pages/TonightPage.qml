import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "../Components"
import MenuClient 1.0

Page {
    id: root
    title: "今晚"

    background: Rectangle { color: Theme.background }

    DateFormatter {
        id: dateFormatter
    }

    Component.onCompleted: MenuApi.LoadData()

    ScrollView {
        anchors.fill: parent
        anchors.margins: Theme.pageMargin
        clip: true

        Column {
            id: content
            width: Math.min(root.width - Theme.pageMargin * 2, Theme.contentMaxWidth)
            anchors.horizontalCenter: parent.horizontalCenter
            spacing: 20

            RowLayout {
                width: parent.width
                spacing: 16

                Column {
                    Layout.fillWidth: true
                    spacing: 4
                    Label {
                        text: "今晚 · " + dateFormatter.formatLocalDate(new Date())
                        color: Theme.mutedText
                        font.pixelSize: 12
                        font.weight: Font.DemiBold
                    }
                    Label {
                        text: MenuApi.IsLoading ? "正在为你挑选" : "选一顿刚刚好的饭"
                        color: Theme.ink
                        font.pixelSize: 25
                        font.weight: Font.DemiBold
                    }
                    Label {
                        text: "根据时间、偏好和常备食材整理"
                        color: Theme.mutedText
                        font.pixelSize: 13
                    }
                }
                AppButton {
                    text: "换一批"
                    iconName: "refresh"
                    variant: "secondary"
                    enabled: !MenuApi.IsLoading
                    onClicked: MenuApi.LoadData()
                }
            }

            Flow {
                width: parent.width
                spacing: 6
                AppChip { label: MenuApi.ServingCount + " 人"; accent: Theme.green }
                AppChip { label: MenuApi.AvailableMinutes + " 分钟"; accent: Theme.sky }
                AppChip {
                    label: MenuApi.PantryIngredientIds.length + " 项常备"
                    accent: Theme.amber
                }
            }

            Rectangle {
                width: parent.width
                height: MenuApi.IsOffline ? 40 : 0
                visible: MenuApi.IsOffline
                radius: Theme.radius
                color: Theme.tint(Theme.coral, Theme.dark ? 0.18 : 0.10)
                Label {
                    anchors.left: parent.left
                    anchors.leftMargin: 14
                    anchors.verticalCenter: parent.verticalCenter
                    text: "离线模式 · 显示最近一次可用方案"
                    color: Theme.coral
                    font.pixelSize: 13
                    font.weight: Font.DemiBold
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

            Column {
                visible: MenuApi.IsLoading
                height: visible ? implicitHeight : 0
                width: parent.width
                spacing: 12

                Repeater {
                    model: 2
                    delegate: Rectangle {
                        width: parent.width
                        height: index === 0 ? 230 : 188
                        radius: Theme.radius
                        color: Theme.surface
                        border.color: Theme.line
                        border.width: 1

                        Column {
                            anchors.fill: parent
                            anchors.margins: 16
                            spacing: 12
                            Rectangle {
                                width: parent.width
                                height: index === 0 ? 112 : 84
                                radius: Theme.radiusSmall
                                color: Theme.surfaceRaised
                            }
                            Rectangle {
                                width: parent.width * 0.64
                                height: 18
                                radius: 9
                                color: Theme.surfaceRaised
                            }
                            Rectangle {
                                width: parent.width * 0.42
                                height: 14
                                radius: 7
                                color: Theme.surfaceRaised
                            }
                        }
                    }
                }
            }

            SectionHeading { title: "为你筛过的方案"; width: parent.width }

            Grid {
                id: recipeGrid
                visible: !MenuApi.IsLoading
                height: visible ? implicitHeight : 0
                width: parent.width
                columns: root.width >= 980 ? 3 : root.width >= 620 ? 2 : 1
                columnSpacing: 12
                rowSpacing: 12
                Repeater {
                    model: MenuApi.Recommendations
                    delegate: RecipeCard {
                        width: (recipeGrid.width - (recipeGrid.columns - 1) * recipeGrid.columnSpacing) /
                            recipeGrid.columns
                        featured: index === 0
                        title: name
                        subtitle: cuisine + " · " + servings + " 人份"
                        meta: totalMinutes + " 分钟 · 难度 " + difficulty
                        imageSource: mediaUrl
                        availableCount: availableIngredientCount
                        missingCount: missingIngredientCount
                        onClicked: MenuApi.OpenRecipe(index)
                    }
                }
            }

            Rectangle {
                visible: !MenuApi.IsLoading && !MenuApi.ErrorMessage &&
                    MenuApi.Recommendations.Count === 0
                width: parent.width
                height: 120
                radius: Theme.radius
                color: Theme.surfaceRaised
                Label {
                    anchors.centerIn: parent
                    text: "还没有合适的方案"
                    color: Theme.mutedText
                    font.pixelSize: 15
                }
            }
        }
    }
}
