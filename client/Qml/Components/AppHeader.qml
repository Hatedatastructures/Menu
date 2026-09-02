import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import MenuClient 1.0

Rectangle {
    id: root

    property string title: "今晚"
    property string kicker: "MENU · 厨房节奏"
    property string iconName: "tonight"
    property bool offline: false
    property real topInset: 0

    color: Theme.headerBackground
    implicitHeight: Theme.topBarHeight + root.topInset
    border.color: Theme.headerLine
    border.width: 1

    RowLayout {
        anchors.fill: parent
        anchors.topMargin: root.topInset
        anchors.leftMargin: 24
        anchors.rightMargin: 20
        spacing: 12

        Rectangle {
            Layout.preferredWidth: 34
            Layout.preferredHeight: 34
            radius: Theme.radiusSmall
            color: Theme.greenStrong

            AppIcon {
                anchors.centerIn: parent
                name: root.iconName
                color: Theme.onAccent
                width: 22
                height: 22
                strokeWidth: 1.7
            }
        }

        Column {
            Layout.fillWidth: true
            spacing: 1

            Label {
                text: root.kicker
                color: Theme.headerMuted
                font.pixelSize: 10
                font.weight: Font.DemiBold
                font.letterSpacing: 0.8
            }
            Label {
                text: root.title
                color: Theme.headerText
                font.pixelSize: 22
                font.weight: Font.DemiBold
            }
        }

        AppChip {
            visible: root.offline
            label: "离线"
            accent: Theme.coral
        }
    }
}
