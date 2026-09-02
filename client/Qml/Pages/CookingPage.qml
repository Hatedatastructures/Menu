import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "../Components"
import MenuClient 1.0

Page {
    id: root
    title: "做饭"
    background: Rectangle { color: Theme.background }

    property var recipe: MenuApi.ActiveRecipe
    property var completedSteps: ({})
    property var remainingByStep: ({})
    property var timerDeadlines: ({})

    function stepOrderFor(stepValue, fallbackIndex) {
        return Number(stepValue.stepOrder || (fallbackIndex + 1))
    }

    function remainingForStep(stepOrder, fallback) {
        var value = root.remainingByStep[stepOrder]
        return value === undefined ? fallback : Number(value)
    }

    function isTimerRunning(stepOrder) {
        return root.timerDeadlines[stepOrder] !== undefined
    }

    function setCompleted(stepOrder, completed) {
        var next = Object.assign({}, root.completedSteps)
        next[stepOrder] = completed
        root.completedSteps = next
        if (completed) {
            root.stopTimer(stepOrder)
        }
        var steps = root.recipe.steps || []
        var allCompleted = steps.length > 0
        for (var index = 0; index < steps.length; ++index) {
            var order = root.stepOrderFor(steps[index], index)
            if (!root.completedSteps[order]) {
                allCompleted = false
                break
            }
        }
        MenuApi.UpdateCookingSession(stepOrder, completed
            ? (allCompleted ? "completed" : "paused") : "active")
    }

    function hydrateCompletedSteps() {
        var session = MenuApi.CurrentCookingSession || {}
        var currentOrder = Number(session.currentStepOrder || 0)
        var state = session.state || ""
        if (currentOrder <= 0 && state !== "completed") {
            return
        }
        var steps = root.recipe.steps || []
        var next = Object.assign({}, root.completedSteps)
        for (var index = 0; index < steps.length; ++index) {
            var order = root.stepOrderFor(steps[index], index)
            if (state === "completed" || order < currentOrder) {
                next[order] = true
            }
        }
        root.completedSteps = next
    }

    function startTimer(stepOrder, seconds, title, body) {
        var duration = Math.max(1, Number(seconds || 0))
        var deadlines = Object.assign({}, root.timerDeadlines)
        deadlines[stepOrder] = Date.now() + duration * 1000
        root.timerDeadlines = deadlines
        var remaining = Object.assign({}, root.remainingByStep)
        remaining[stepOrder] = duration
        root.remainingByStep = remaining
                MenuTimers.StartTimer(root.recipe.id, stepOrder, duration, title, body)
    }

    function stopTimer(stepOrder) {
        var deadlines = Object.assign({}, root.timerDeadlines)
        delete deadlines[stepOrder]
        root.timerDeadlines = deadlines
        MenuTimers.StopTimer(root.recipe.id, stepOrder)
    }

    function toggleTimer(stepOrder, seconds, title, body) {
        if (root.isTimerRunning(stepOrder)) {
            root.stopTimer(stepOrder)
        } else {
            var remaining = root.remainingForStep(stepOrder, 0)
            if (remaining <= 0) {
                remaining = Number(seconds || 0)
            }
            root.startTimer(stepOrder, remaining, title, body)
        }
    }

    function reconcileTimers() {
        var now = Date.now()
        var deadlines = Object.assign({}, root.timerDeadlines)
        var remaining = Object.assign({}, root.remainingByStep)
        var keys = Object.keys(deadlines)
        for (var index = 0; index < keys.length; ++index) {
            var order = keys[index]
            var seconds = Math.max(0, Math.ceil((Number(deadlines[order]) - now) / 1000))
            remaining[order] = seconds
            if (seconds <= 0) {
                delete deadlines[order]
                MenuTimers.StopTimer(root.recipe.id, Number(order))
            }
        }
        root.remainingByStep = remaining
        root.timerDeadlines = deadlines
    }

    Timer {
        interval: 250
        repeat: true
        running: Object.keys(root.timerDeadlines).length > 0
        onTriggered: root.reconcileTimers()
    }

    function startSessionIfReady() {
        if (root.recipe.id !== undefined && root.recipe.id !== "") {
            MenuApi.StartCookingSession()
        }
    }

    function updateScreenAwake() {
        MenuScreenAwake.SetEnabled(root.recipe.id !== undefined && root.recipe.id !== "")
    }

    Component.onCompleted: {
        root.startSessionIfReady()
        root.updateScreenAwake()
    }
    onRecipeChanged: {
        root.completedSteps = ({})
        root.remainingByStep = ({})
        root.timerDeadlines = ({})
        MenuTimers.StopAllTimers()
        root.startSessionIfReady()
        root.updateScreenAwake()
    }

    Component.onDestruction: {
        MenuTimers.StopAllTimers()
        MenuScreenAwake.SetEnabled(false)
    }

    Connections {
        target: MenuApi
        function onCookingSessionChanged() { root.hydrateCompletedSteps() }
    }

    Column {
        anchors.top: parent.top
        anchors.bottom: parent.bottom
        anchors.margins: Theme.pageMargin
        width: Math.min(parent.width - Theme.pageMargin * 2, Theme.contentMaxWidth)
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.bottomMargin: 30
        spacing: 12

        RowLayout {
            width: parent.width
            spacing: 10
            AppButton {
                text: "返回"
                variant: "ghost"
                onClicked: MenuApi.CloseRecipe()
            }
            Column {
                Layout.fillWidth: true
                spacing: 2
                Label {
                    width: parent.width
                    text: root.recipe.name || "做饭"
                    color: Theme.ink
                    font.pixelSize: 22
                    font.weight: Font.DemiBold
                    elide: Text.ElideRight
                }
                Label {
                    text: "跟着步骤走 · " + (root.recipe.totalMinutes || 0) + " 分钟"
                    color: Theme.mutedText
                    font.pixelSize: 13
                }
            }
            AppChip { label: (root.recipe.steps || []).length + " 步"; accent: Theme.green }
        }

        Rectangle {
            width: parent.width
            height: 156
            radius: Theme.radius
            clip: true
            color: Theme.surfaceRaised
            Image {
                anchors.fill: parent
                source: root.recipe.mediaUrl || "qrc:/qt/qml/MenuClient/assets/media/menu-placeholder.png"
                sourceSize.width: 720
                sourceSize.height: 312
                asynchronous: true
                cache: true
                fillMode: Image.PreserveAspectCrop
            }
        }

        Label {
            width: parent.width
            visible: Boolean(root.recipe.description)
            text: root.recipe.description || ""
            color: Theme.mutedText
            wrapMode: Text.WordWrap
            font.pixelSize: 14
            maximumLineCount: 2
            elide: Text.ElideRight
        }

        RowLayout {
            width: parent.width
            Label {
                Layout.fillWidth: true
                text: MenuApi.RemindersEnabled ? "备菜提醒已开启" : "提醒已关闭"
                color: MenuApi.RemindersEnabled ? Theme.green : Theme.mutedText
                font.pixelSize: 12
                font.weight: Font.DemiBold
            }
            Label {
                text: MenuDisplay.UsingFallback
                    ? "120 Hz 预算"
                    : Math.round(MenuDisplay.ReportedRefreshRateHz) + " Hz"
                color: Theme.mutedText
                font.pixelSize: 12
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
            id: stepList
            width: parent.width
            height: parent.height - y - feedbackRow.height
            spacing: 8
            clip: true
            model: root.recipe.steps || []
            delegate: CookingStepCard {
                step: modelData
                stepIndex: index
                stepOrder: root.stepOrderFor(modelData, index)
                completed: Boolean(root.completedSteps[stepOrder])
                remainingSeconds: root.remainingForStep(
                    stepOrder, Number(modelData.durationSeconds || 0))
                timerRunning: root.isTimerRunning(stepOrder)
                onCompletedToggled: root.setCompleted(stepOrder, checked)
                onTimerToggled: function(order, seconds, title, body) {
                    root.toggleTimer(order, seconds, title, body)
                }
            }
        }

        Row {
            id: feedbackRow
            width: parent.width
            spacing: 8
            AppButton {
                text: "做了"
                variant: "primary"
                onClicked: MenuApi.SubmitFeedback("made", ["下次还想做"], "")
            }
            AppButton {
                text: "先跳过"
                variant: "ghost"
                onClicked: MenuApi.SubmitFeedback("skipped", [], "")
            }
        }
    }
}
