import QtQuick
import QtQuick.Controls
import QtQuick.Layouts

ApplicationWindow {
    id: root
    visible: true
    width: 430
    height: 860
    minimumWidth: 360
    minimumHeight: 640
    title: "Menu"
    color: Theme.paper
    property bool cookingVisible: Boolean(MenuApi.ActiveRecipe.id)
    property int pageIndex: pageState.pageIndex

    ColumnLayout {
        anchors.fill: parent
        spacing: 0

        Rectangle {
            Layout.fillWidth: true
            Layout.preferredHeight: 72
            color: Theme.ink

            Column {
                anchors.left: parent.left
                anchors.leftMargin: 24
                anchors.verticalCenter: parent.verticalCenter
                spacing: 4

                Label {
                    text: "Menu"
                    color: Theme.paper
                    font.pixelSize: 24
                    font.bold: true
                }
                Label {
                    text: root.pageIndex === 0 ? "今晚"
                        : root.pageIndex === 1 ? "本周"
                        : root.pageIndex === 2 ? "食材"
                        : root.pageIndex === 3 ? "我的" : "做饭"
                    color: Theme.mutedPaper
                    font.pixelSize: 13
                }
            }
        }

        StackLayout {
            id: pages
            Layout.fillWidth: true
            Layout.fillHeight: true
            currentIndex: root.pageIndex

            Loader {
                Layout.fillWidth: true
                Layout.fillHeight: true
                active: root.pageIndex === 0 || hasLoaded
                property bool hasLoaded: false
                sourceComponent: tonightPageComponent
                onLoaded: hasLoaded = true
            }
            Loader {
                Layout.fillWidth: true
                Layout.fillHeight: true
                active: root.pageIndex === 1 || hasLoaded
                property bool hasLoaded: false
                sourceComponent: weekPageComponent
                onLoaded: hasLoaded = true
            }
            Loader {
                Layout.fillWidth: true
                Layout.fillHeight: true
                active: root.pageIndex === 2 || hasLoaded
                property bool hasLoaded: false
                sourceComponent: ingredientsPageComponent
                onLoaded: hasLoaded = true
            }
            Loader {
                Layout.fillWidth: true
                Layout.fillHeight: true
                active: root.pageIndex === 3 || hasLoaded
                property bool hasLoaded: false
                sourceComponent: profilePageComponent
                onLoaded: hasLoaded = true
            }
            Loader {
                Layout.fillWidth: true
                Layout.fillHeight: true
                active: root.pageIndex === 4 || hasLoaded
                property bool hasLoaded: false
                sourceComponent: cookingPageComponent
                onLoaded: hasLoaded = true
            }
        }

        BottomNavigation {
            id: navigation
            Layout.fillWidth: true
            visible: !root.cookingVisible
        }

        NavigationState {
            id: pageState
            navigationIndex: navigation.currentIndex
            cookingVisible: root.cookingVisible
        }
    }

    Component {
        id: tonightPageComponent
        TonightPage {}
    }

    Component {
        id: weekPageComponent
        WeekPage {}
    }

    Component {
        id: ingredientsPageComponent
        IngredientsPage {}
    }

    Component {
        id: profilePageComponent
        ProfilePage {}
    }

    Component {
        id: cookingPageComponent
        CookingPage {}
    }

}
