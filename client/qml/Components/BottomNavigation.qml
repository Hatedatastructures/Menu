import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

Rectangle {
    id: root
    height: 76
    color: Theme.surface
    border.color: Theme.line
    border.width: 1

    property int currentIndex: 0

    RowLayout {
        anchors.fill: parent
        anchors.leftMargin: 8
        anchors.rightMargin: 8
        spacing: 4

        Repeater {
            model: ["今晚", "本周", "食材", "我的"]

            delegate: Button {
                Layout.fillWidth: true
                Layout.fillHeight: true
                flat: true
                text: modelData
                highlighted: index === root.currentIndex
                onClicked: root.currentIndex = index
            }
        }
    }
}
