import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "."
import MenuClient 1.0

ItemDelegate {
    id: root

    width: ListView.view ? ListView.view.width : 360
    height: 246
    padding: 0
    text: root.title
    Accessible.name: root.title
    Accessible.role: Accessible.Button
    hoverEnabled: true

    property bool featured: false
    property string title: ""
    property string subtitle: ""
    property string meta: ""
    property string imageSource: "qrc:/qt/qml/MenuClient/assets/media/menu-placeholder.png"
    property int availableCount: 0
    property int missingCount: 0

    background: Rectangle {
        color: root.highlighted || root.hovered
            ? Theme.tint(Theme.green, Theme.dark ? 0.16 : 0.06)
            : Theme.surface
        border.color: root.activeFocus ? Theme.focus : Theme.line
        border.width: root.activeFocus ? 2 : 1
        radius: Theme.radius
    }

    contentItem: ColumnLayout {
        anchors.fill: parent
        anchors.margins: 14
        spacing: 9

        Item {
            Layout.fillWidth: true
            Layout.preferredHeight: 116
            clip: true

            Image {
                anchors.fill: parent
                source: root.imageSource
                sourceSize.width: 720
                sourceSize.height: 360
                asynchronous: true
                cache: true
                fillMode: Image.PreserveAspectCrop
            }

            Rectangle {
                anchors.left: parent.left
                anchors.top: parent.top
                anchors.margins: 8
                width: timeLabel.implicitWidth + 18
                height: 28
                radius: Theme.radiusSmall
                color: Theme.headerBackground

                Label {
                    id: timeLabel
                    anchors.centerIn: parent
                    text: root.meta.split(" · ")[0] || ""
                    color: Theme.headerText
                    font.pixelSize: 11
                    font.weight: Font.DemiBold
                }
            }
        }

        RowLayout {
            Layout.fillWidth: true
            spacing: 12

            Column {
                Layout.fillWidth: true
                spacing: 4

                Label {
                    width: parent.width
                    text: root.title
                    color: Theme.text
                    font.pixelSize: 18
                    font.weight: Font.DemiBold
                    elide: Text.ElideRight
                }
                Label {
                    width: parent.width
                    text: root.subtitle
                    color: Theme.mutedText
                    font.pixelSize: 13
                    elide: Text.ElideRight
                }
            }

            AppIcon {
                name: "back"
                rotation: 180
                width: 20
                height: 20
                color: Theme.mutedText
            }
        }

        RowLayout {
            Layout.fillWidth: true
            spacing: 8

            AppChip {
                label: root.meta.split(" · ")[1] || "轻松"
                accent: Theme.green
            }
            Label {
                Layout.fillWidth: true
                text: root.missingCount > 0
                    ? "缺 " + root.missingCount + " 种食材"
                    : "食材齐全"
                color: root.missingCount > 0 ? Theme.coral : Theme.green
                font.pixelSize: 12
                font.weight: Font.DemiBold
                elide: Text.ElideRight
                horizontalAlignment: Text.AlignRight
            }
        }
    }
}
