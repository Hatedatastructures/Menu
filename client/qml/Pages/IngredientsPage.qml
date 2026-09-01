import QtQuick
import QtQuick.Controls

Page {
    title: "食材"

    ListView {
        anchors.fill: parent
        anchors.margins: 16
        spacing: 1
        clip: true
        model: MenuApi.Ingredients
        delegate: ItemDelegate {
            width: ListView.view.width
            height: 72
            contentItem: Row {
                spacing: 12
                Label {
                    width: 88
                    anchors.verticalCenter: parent.verticalCenter
                    text: name
                    color: Theme.text
                    font.pixelSize: 16
                    font.bold: true
                }
                Column {
                    width: parent.width - 120
                    anchors.verticalCenter: parent.verticalCenter
                    spacing: 3
                    Label {
                        text: category + " · 默认单位 " + defaultUnit
                        color: Theme.mutedText
                        font.pixelSize: 13
                    }
                    Label {
                        text: isPantryStaple ? "常备食材" : "按需准备"
                        color: isPantryStaple ? Theme.green : Theme.mutedText
                        font.pixelSize: 12
                    }
                }
            }
        }
    }
}
