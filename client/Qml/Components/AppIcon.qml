import QtQuick

Item {
    id: root

    property string name: "dot"
    property color color: "#17231c"
    property real strokeWidth: 1.8

    implicitWidth: 24
    implicitHeight: 24

    Canvas {
        id: canvas
        anchors.fill: parent
        antialiasing: true

        onPaint: {
            var ctx = getContext("2d")
            ctx.reset()
            ctx.strokeStyle = root.color
            ctx.fillStyle = root.color
            ctx.lineWidth = root.strokeWidth
            ctx.lineCap = "round"
            ctx.lineJoin = "round"

            var w = width
            var h = height
            var cx = w / 2
            var cy = h / 2

            if (root.name === "tonight") {
                ctx.beginPath()
                ctx.arc(cx, cy + 2, 7.5, Math.PI, 0, false)
                ctx.stroke()
                ctx.beginPath()
                ctx.moveTo(cx - 8, cy + 2)
                ctx.lineTo(cx + 8, cy + 2)
                ctx.stroke()
                ctx.beginPath()
                ctx.moveTo(cx - 4, cy - 10)
                ctx.bezierCurveTo(cx - 7, cy - 7, cx - 1, cy - 6, cx - 4, cy - 3)
                ctx.moveTo(cx + 3, cy - 9)
                ctx.bezierCurveTo(cx, cy - 6, cx + 6, cy - 5, cx + 3, cy - 2)
                ctx.stroke()
            } else if (root.name === "week") {
                ctx.strokeRect(4, 5, w - 8, h - 7)
                ctx.beginPath()
                ctx.moveTo(4, 9)
                ctx.lineTo(w - 4, 9)
                ctx.moveTo(8, 3)
                ctx.lineTo(8, 7)
                ctx.moveTo(w - 8, 3)
                ctx.lineTo(w - 8, 7)
                ctx.moveTo(8, 13)
                ctx.lineTo(w - 8, 13)
                ctx.moveTo(8, 17)
                ctx.lineTo(w - 8, 17)
                ctx.moveTo(13, 9)
                ctx.lineTo(13, h - 2)
                ctx.moveTo(18, 9)
                ctx.lineTo(18, h - 2)
                ctx.stroke()
            } else if (root.name === "ingredients") {
                ctx.beginPath()
                ctx.moveTo(cx, h - 4)
                ctx.bezierCurveTo(5, h - 8, 5, 7, cx, 4)
                ctx.bezierCurveTo(w - 5, 7, w - 5, h - 8, cx, h - 4)
                ctx.stroke()
                ctx.beginPath()
                ctx.moveTo(cx, h - 4)
                ctx.lineTo(cx, 5)
                ctx.moveTo(cx, 13)
                ctx.lineTo(8, 9)
                ctx.moveTo(cx, 17)
                ctx.lineTo(w - 8, 12)
                ctx.stroke()
            } else if (root.name === "profile") {
                ctx.beginPath()
                ctx.arc(cx, 8, 4, 0, Math.PI * 2)
                ctx.stroke()
                ctx.beginPath()
                ctx.arc(cx, 21, 8, Math.PI, 0, false)
                ctx.stroke()
            } else if (root.name === "refresh") {
                ctx.beginPath()
                ctx.arc(cx, cy, 7, -0.7, 3.9, false)
                ctx.stroke()
                ctx.beginPath()
                ctx.moveTo(18, 5)
                ctx.lineTo(18, 10)
                ctx.lineTo(13, 8)
                ctx.stroke()
            } else if (root.name === "settings") {
                ctx.beginPath()
                ctx.arc(cx, cy, 3.5, 0, Math.PI * 2)
                ctx.stroke()
                for (var i = 0; i < 8; ++i) {
                    var angle = i * Math.PI / 4
                    ctx.moveTo(cx + Math.cos(angle) * 6, cy + Math.sin(angle) * 6)
                    ctx.lineTo(cx + Math.cos(angle) * 9, cy + Math.sin(angle) * 9)
                }
                ctx.stroke()
            } else if (root.name === "back") {
                ctx.beginPath()
                ctx.moveTo(19, cy)
                ctx.lineTo(5, cy)
                ctx.moveTo(5, cy)
                ctx.lineTo(11, cy - 6)
                ctx.moveTo(5, cy)
                ctx.lineTo(11, cy + 6)
                ctx.stroke()
            } else if (root.name === "check") {
                ctx.beginPath()
                ctx.moveTo(4, 13)
                ctx.lineTo(10, 19)
                ctx.lineTo(20, 6)
                ctx.stroke()
            } else if (root.name === "plus" || root.name === "minus") {
                ctx.beginPath()
                ctx.moveTo(5, cy)
                ctx.lineTo(w - 5, cy)
                if (root.name === "plus") {
                    ctx.moveTo(cx, 5)
                    ctx.lineTo(cx, h - 5)
                }
                ctx.stroke()
            } else {
                ctx.beginPath()
                ctx.arc(cx, cy, 3, 0, Math.PI * 2)
                ctx.fill()
            }
        }

        Component.onCompleted: requestPaint()
    }

    onNameChanged: canvas.requestPaint()
    onColorChanged: canvas.requestPaint()
    onStrokeWidthChanged: canvas.requestPaint()
    onWidthChanged: canvas.requestPaint()
    onHeightChanged: canvas.requestPaint()
}
