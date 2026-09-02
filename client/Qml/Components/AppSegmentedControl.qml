import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import MenuClient 1.0

Item {
    id: root

    property var options: []
    property int currentIndex: 0
    signal indexActivated(int index)

    implicitHeight: Theme.touchTarget
    implicitWidth: 260

    RowLayout {
        anchors.fill: parent
        spacing: 4

        Repeater {
            model: root.options
            delegate: Button {
                Layout.fillWidth: true
                Layout.fillHeight: true
                Layout.minimumWidth: 0
                flat: true
                padding: 0
                text: modelData
                Accessible.name: modelData
                Accessible.role: Accessible.RadioButton

                contentItem: Label {
                    text: modelData
                    color: index === root.currentIndex ? Theme.onColor(Theme.greenStrong) : Theme.ink
                    font.pixelSize: 13
                    font.weight: index === root.currentIndex ? Font.DemiBold : Font.Normal
                    horizontalAlignment: Text.AlignHCenter
                    verticalAlignment: Text.AlignVCenter
                    elide: Text.ElideRight
                }

                background: Rectangle {
                    radius: Theme.radiusSmall
                    color: index === root.currentIndex
                        ? Theme.greenStrong : Theme.surface
                    border.width: index === root.currentIndex ? 0 : 1
                    border.color: Theme.line
                }

                onClicked: {
                    root.currentIndex = index
                    root.indexActivated(index)
                }
            }
        }
    }
}
