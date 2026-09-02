import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "."
import MenuClient 1.0

Rectangle {
    id: root

    property var step: ({})
    property int stepIndex: 0
    property int stepOrder: stepIndex + 1
    property bool completed: false
    property int remainingSeconds: 0
    property bool timerRunning: false

    signal completedToggled(bool checked)
    signal timerToggled(int stepOrder, int seconds, string title, string body)

    function durationLabel(seconds) {
        var value = Number(seconds || 0)
        if (value >= 60 && value % 60 === 0) {
            return (value / 60) + " 分钟"
        }
        return value + " 秒"
    }

    width: parent ? parent.width : 360
    height: Math.max(176, 122 + instructionLabel.implicitHeight)
    color: root.completed ? Theme.tint(Theme.green, Theme.dark ? 0.10 : 0.04) : Theme.surface
    border.color: root.completed ? Theme.green : Theme.line
    border.width: root.completed ? 2 : 1
    radius: Theme.radius

    Column {
        anchors.fill: parent
        anchors.margins: 14
        spacing: 8

        Row {
            width: parent.width
            spacing: 10

            Label {
                text: (root.stepIndex + 1).toString()
                width: 30
                height: 30
                horizontalAlignment: Text.AlignHCenter
                verticalAlignment: Text.AlignVCenter
                color: Theme.onAccent
                font.weight: Font.DemiBold
                background: Rectangle {
                    color: root.completed ? Theme.greenStrong : Theme.green
                    radius: 15
                }
            }
            Label {
                width: parent.width - 42
                text: root.step.title || "步骤 " + (root.stepIndex + 1)
                color: Theme.ink
                font.pixelSize: 17
                font.weight: Font.DemiBold
                elide: Text.ElideRight
            }
        }

        Label {
            id: instructionLabel
            width: parent.width
            text: root.step.instruction || ""
            color: Theme.mutedText
            wrapMode: Text.WordWrap
            maximumLineCount: 4
            elide: Text.ElideRight
        }

        Row {
            width: parent.width
            spacing: 10

            AppCheckBox {
                text: "完成"
                checked: root.completed
                onToggled: root.completedToggled(checked)
            }
            AppButton {
                visible: Number(root.step.durationSeconds || 0) > 0 &&
                    (Boolean(root.step.hasTimer) || root.remainingSeconds > 0)
                text: root.timerRunning
                    ? "暂停 " + root.remainingSeconds + " 秒"
                    : "计时 " + root.durationLabel(root.remainingSeconds > 0
                        ? root.remainingSeconds : root.step.durationSeconds)
                variant: root.timerRunning ? "accent" : "secondary"
                onClicked: root.timerToggled(
                    root.stepOrder,
                    Math.max(1, Number(root.step.durationSeconds || 0)),
                    root.step.title || "步骤完成",
                    root.step.instruction || "可以继续下一步了")
            }
        }
    }
}
