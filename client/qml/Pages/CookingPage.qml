import QtQuick
import QtQuick.Controls

Page {
    id: root
    title: "做饭"

    property var recipe: MenuApi.ActiveRecipe

    function startSessionIfReady() {
        if (root.recipe.id !== undefined && root.recipe.id !== "") {
            MenuApi.StartCookingSession()
        }
    }

    Component.onCompleted: root.startSessionIfReady()
    onRecipeChanged: root.startSessionIfReady()

    Column {
        anchors.fill: parent
        anchors.margins: 16
        spacing: 10

        Row {
            width: parent.width
            spacing: 12
            Button {
                text: "返回"
                onClicked: MenuApi.CloseRecipe()
            }
            Column {
                width: parent.width - 100
                spacing: 2
                Label {
                    text: root.recipe.name || "做饭"
                    color: Theme.text
                    font.pixelSize: 22
                    font.bold: true
                    elide: Text.ElideRight
                }
                Label {
                    text: "按步骤完成 · " + (root.recipe.totalMinutes || 0) + " 分钟"
                    color: Theme.mutedText
                    font.pixelSize: 13
                }
            }
        }

        Image {
            width: parent.width
            height: 156
            source: root.recipe.mediaUrl || "qrc:/qt/qml/MenuClient/assets/media/menu-placeholder.png"
            sourceSize.width: 720
            sourceSize.height: 312
            asynchronous: true
            cache: true
            fillMode: Image.PreserveAspectCrop
        }

        Label {
            width: parent.width
            visible: Boolean(root.recipe.description)
            text: root.recipe.description || ""
            color: Theme.mutedText
            wrapMode: Text.WordWrap
            font.pixelSize: 14
        }

        Label {
            width: parent.width
            text: MenuApi.IsAuthenticated
                ? "当前步骤已连接服务端；勾选和计时仅在本次打开期间保留"
                : "未登录：当前步骤、勾选和计时仅在本次打开期间保留"
            color: Theme.mutedText
            wrapMode: Text.WordWrap
            font.pixelSize: 12
        }

        Label {
            visible: MenuApi.RemindersEnabled
            width: parent.width
            text: "备菜提醒已开启：只在需要准备或计时的节点提示"
            color: Theme.green
            font.pixelSize: 12
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
            width: parent.width
            height: parent.height - y
            spacing: 8
            clip: true
            model: root.recipe.steps || []
            delegate: Rectangle {
                id: stepRow
                width: ListView.view.width
                height: 172
                color: Theme.surface
                border.color: completedCheck.checked ? Theme.green : Theme.line
                border.width: 1
                radius: 4

                property int remainingSeconds: Number(modelData.durationSeconds || 0)

                Timer {
                    id: stepTimer
                    interval: 1000
                    repeat: true
                    running: false
                    onTriggered: {
                        if (stepRow.remainingSeconds > 0) {
                            stepRow.remainingSeconds -= 1
                        } else {
                            stop()
                        }
                    }
                }

                Column {
                    anchors.fill: parent
                    anchors.margins: 14
                    spacing: 8

                    Row {
                        width: parent.width
                        spacing: 10
                        Label {
                            text: (index + 1).toString()
                            width: 28
                            height: 28
                            horizontalAlignment: Text.AlignHCenter
                            verticalAlignment: Text.AlignVCenter
                            color: Theme.paper
                            font.bold: true
                            background: Rectangle { color: Theme.green; radius: 14 }
                        }
                        Label {
                            width: parent.width - 42
                            text: modelData.title || "步骤 " + (index + 1)
                            color: Theme.text
                            font.pixelSize: 17
                            font.bold: true
                            elide: Text.ElideRight
                        }
                    }

                    Label {
                        width: parent.width
                        height: 42
                        text: modelData.instruction || ""
                        color: Theme.mutedText
                        wrapMode: Text.WordWrap
                        elide: Text.ElideRight
                        maximumLineCount: 2
                    }

                    Row {
                        width: parent.width
                        spacing: 12
                        CheckBox {
                            id: completedCheck
                            text: "完成"
                            onToggled: {
                                MenuApi.UpdateCookingSession(
                                    index + 1,
                                    completedCheck.checked ? "paused" : "active")
                            }
                        }
                        Button {
                            visible: Boolean(modelData.hasTimer) || stepRow.remainingSeconds > 0
                            text: stepTimer.running ? "暂停 " + stepRow.remainingSeconds + " 秒" : "计时 " + stepRow.remainingSeconds + " 秒"
                            onClicked: {
                                if (stepRow.remainingSeconds === 0) {
                                    stepRow.remainingSeconds = Number(modelData.durationSeconds || 0)
                                }
                                stepTimer.running = !stepTimer.running
                            }
                        }
                    }
                }
            }
        }

        Row {
            width: parent.width
            spacing: 10
            Button {
                text: "做了"
                onClicked: MenuApi.SubmitFeedback("made", ["下次还想做"], "")
            }
            Button {
                text: "这次没做"
                onClicked: MenuApi.SubmitFeedback("skipped", [], "")
            }
        }
    }
}
