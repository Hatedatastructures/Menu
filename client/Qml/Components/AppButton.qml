import QtQuick
import QtQuick.Controls
import MenuClient 1.0

Button {
    id: root

    property string variant: "primary"
    property string iconText: ""
    property string iconName: ""

    implicitHeight: Theme.touchTarget
    implicitWidth: root.text.length === 0 &&
        (root.iconName.length > 0 || root.iconText.length > 0)
        ? Theme.touchTarget
        : Math.max(96, contentItem.implicitWidth + leftPadding + rightPadding)
    leftPadding: 16
    rightPadding: 16
    topPadding: 8
    bottomPadding: 8
    opacity: enabled ? 1 : 0.48
    Accessible.name: root.text
    Accessible.role: Accessible.Button
    hoverEnabled: true

    contentItem: Row {
        spacing: root.iconName.length > 0 || root.iconText.length > 0 ? 8 : 0
        anchors.centerIn: parent

        AppIcon {
            visible: root.iconName.length > 0
            name: root.iconName
            color: root.variant === "primary" ? Theme.onColor(Theme.greenStrong)
                : root.variant === "accent" ? Theme.onColor(Theme.coralStrong) : Theme.ink
            width: root.iconName.length > 0 ? 18 : 0
            height: root.iconName.length > 0 ? 18 : 0
            strokeWidth: 1.7
            anchors.verticalCenter: parent.verticalCenter
        }

        Label {
            visible: root.iconName.length === 0 && root.iconText.length > 0
            text: root.iconText
            color: root.variant === "primary" ? Theme.onColor(Theme.greenStrong)
                : root.variant === "accent" ? Theme.onColor(Theme.coralStrong) : Theme.ink
            font.pixelSize: 18
            horizontalAlignment: Text.AlignHCenter
            verticalAlignment: Text.AlignVCenter
        }
        Label {
            text: root.text
            color: root.variant === "primary" ? Theme.onColor(Theme.greenStrong)
                : root.variant === "accent" ? Theme.onColor(Theme.coralStrong) : Theme.ink
            font.pixelSize: 14
            font.weight: Font.DemiBold
            elide: Text.ElideRight
            fontSizeMode: Text.HorizontalFit
            minimumPixelSize: 12
            verticalAlignment: Text.AlignVCenter
            horizontalAlignment: Text.AlignHCenter
        }
    }

    background: Rectangle {
        radius: Theme.radius
            color: root.variant === "primary" ? Theme.greenStrong
                : root.variant === "accent" ? Theme.coralStrong
                : root.variant === "ghost" ? "transparent" : Theme.surface
        border.width: root.variant === "secondary" || root.variant === "ghost" ? 1 : 0
        border.color: root.variant === "ghost" ? Theme.line : Theme.line

        Rectangle {
            anchors.fill: parent
            anchors.margins: 1
            radius: parent.radius - 1
            color: "transparent"
            border.width: root.down || root.activeFocus || root.hovered ? 2 : 0
            border.color: root.variant === "accent" ? Theme.coralStrong : Theme.focus
        }
    }
}
