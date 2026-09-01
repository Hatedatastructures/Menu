import QtQuick
import QtQuick.Controls

ItemDelegate {
    id: root
    width: ListView.view ? ListView.view.width : 360
    height: 132
    padding: 16

    background: Rectangle {
        color: root.highlighted ? "#edf5ef" : Theme.surface
        border.color: root.highlighted ? Theme.green : Theme.line
        border.width: 1
        radius: 4
    }

    contentItem: Row {
        spacing: 12

        Image {
            width: 104
            height: 96
            source: root.imageSource
            sourceSize.width: 208
            sourceSize.height: 192
            asynchronous: true
            cache: true
            fillMode: Image.PreserveAspectCrop
        }

        Column {
            width: parent.width - 116
            spacing: 6

            Label {
                width: parent.width
                text: root.title
                color: Theme.text
                font.pixelSize: 18
                font.bold: true
                elide: Text.ElideRight
            }
            Label {
                width: parent.width
                text: root.subtitle
                color: Theme.mutedText
                font.pixelSize: 13
                elide: Text.ElideRight
            }
            Label {
                width: parent.width
                text: root.meta
                color: Theme.green
                font.pixelSize: 13
                elide: Text.ElideRight
            }
            Label {
                width: parent.width
                text: root.missingCount > 0
                    ? "家里有 " + root.availableCount + " 种 · 还缺 " + root.missingCount + " 种"
                    : "家里已有食材"
                color: root.missingCount > 0 ? Theme.coral : Theme.green
                font.pixelSize: 12
                elide: Text.ElideRight
            }
        }
    }

    property string title: ""
    property string subtitle: ""
    property string meta: ""
    property string imageSource: "qrc:/qt/qml/MenuClient/assets/media/menu-placeholder.png"
    property int availableCount: 0
    property int missingCount: 0
}
