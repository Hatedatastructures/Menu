import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "."
import MenuClient 1.0

Rectangle {
    id: root

    property bool active: true
    property string endpointDraft: MenuApi.BaseUrl
    width: parent ? parent.width : 360
    visible: active
    height: active ? connectionColumn.implicitHeight + 32 : 0
    radius: Theme.radius
    color: Theme.surface
    border.color: Theme.line
    border.width: 1

    Connections {
        target: MenuApi
        function onBaseUrlChanged() { root.endpointDraft = MenuApi.BaseUrl }
    }

    Column {
        id: connectionColumn
        anchors.left: parent.left
        anchors.right: parent.right
        anchors.top: parent.top
        anchors.margins: 16
        spacing: 12

        SectionHeading {
            width: parent.width
            title: "服务器地址"
            iconName: "settings"
        }

        Label {
            width: parent.width
            text: "在这里切换本地服务或局域网服务"
            color: Theme.mutedText
            font.pixelSize: 12
            wrapMode: Text.WordWrap
        }

        AppTextField {
            id: endpointField
            width: parent.width
            height: Theme.touchTarget
            text: root.endpointDraft
            placeholderText: "http://主机或 IP:端口"
            inputMethodHints: Qt.ImhUrlCharactersOnly
            Accessible.name: "服务器地址"
            onTextChanged: root.endpointDraft = text
        }

        Flow {
            width: parent.width
            spacing: 8

            AppButton {
                text: "测试连接"
                iconName: "refresh"
                variant: "secondary"
                enabled: !MenuApi.IsTestingConnection
                onClicked: MenuApi.TestBaseUrl(root.endpointDraft)
            }
            AppButton {
                text: "保存地址"
                variant: "primary"
                onClicked: {
                    if (MenuApi.ApplyBaseUrl(root.endpointDraft)) {
                        MenuApi.TestBaseUrl(MenuApi.BaseUrl)
                        MenuApi.LoadData()
                    }
                }
            }
        }

        RowLayout {
            width: parent.width
            height: 32
            Label {
                Layout.fillWidth: true
                text: MenuApi.ConnectionStatus
                color: MenuApi.ConnectionReachable ? Theme.green : Theme.mutedText
                font.pixelSize: 12
                elide: Text.ElideRight
            }
            AppChip {
                label: MenuApi.ConnectionReachable ? "在线" : "未连接"
                accent: MenuApi.ConnectionReachable ? Theme.green : Theme.mutedText
            }
        }

        Label {
            visible: root.endpointDraft.toLowerCase().indexOf("http://") === 0
            width: parent.width
            text: "HTTP 仅建议用于可信局域网；公网服务请使用 HTTPS"
            color: Theme.amber
            font.pixelSize: 12
            wrapMode: Text.WordWrap
        }
    }
}
