import QtQuick
import QtQuick.Controls
import QtQuick.Layouts
import "Components"
import "Pages"

ApplicationWindow {
    id: root

    visible: true
    width: 430
    height: 860
    minimumWidth: 360
    minimumHeight: 640
    title: "Menu"
    color: Theme.background

    property bool cookingVisible: Boolean(MenuApi.ActiveRecipe.id)
    property int pageIndex: pageState.pageIndex
    property bool wideLayout: width >= 720

    onClosing: function(closeEvent) {
        if (root.cookingVisible) {
            MenuApi.CloseRecipe()
            closeEvent.accepted = false
        }
    }

    Item {
        id: backKeySink
        anchors.fill: parent
        focus: true
        z: -1

        Keys.onReleased: function(event) {
            if ((event.key === Qt.Key_Back || event.key === Qt.Key_Escape) &&
                root.cookingVisible) {
                MenuApi.CloseRecipe()
                event.accepted = true
            }
        }
    }

    onWideLayoutChanged: {
        if (root.wideLayout) {
            sideNavigation.currentIndex = navigation.currentIndex
        } else {
            navigation.currentIndex = sideNavigation.currentIndex
        }
    }

    function pageTitle() {
        return root.pageIndex === 0 ? "今晚"
            : root.pageIndex === 1 ? "本周"
            : root.pageIndex === 2 ? "食材"
            : root.pageIndex === 3 ? "我的" : "做饭"
    }

    function pageKicker() {
        return root.pageIndex === 4 ? "正在做饭" : "MENU · 厨房节奏"
    }

    ColumnLayout {
        anchors.fill: parent
        spacing: 0

        AppHeader {
            Layout.fillWidth: true
            Layout.preferredHeight: Theme.topBarHeight + MenuDisplay.TopInset
            visible: true
            title: root.pageTitle()
            kicker: root.pageKicker()
            iconName: root.pageIndex === 1 ? "week"
                : root.pageIndex === 2 ? "ingredients"
                : root.pageIndex === 3 ? "profile"
                : root.pageIndex === 4 ? "check" : "tonight"
            topInset: MenuDisplay.TopInset
            offline: MenuApi.IsOffline
        }

        RowLayout {
            Layout.fillWidth: true
            Layout.fillHeight: true
            spacing: 0

            NavigationRail {
                id: sideNavigation
                visible: root.wideLayout && !root.cookingVisible
                Layout.fillHeight: true
                Layout.preferredWidth: 216
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
        }

        BottomNavigation {
            id: navigation
            Layout.fillWidth: true
            Layout.preferredHeight: 84
            visible: !root.wideLayout && !root.cookingVisible
        }

        NavigationState {
            id: pageState
            navigationIndex: root.wideLayout ? sideNavigation.currentIndex : navigation.currentIndex
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
