import QtQuick 2.15
import QtTest 1.3
import "../Qml/Components"

TestCase {
    name: "NavigationStateTest"

    function test_OpenRecipeSwitchesToCookingAndRestoresNavigation() {
        const State = navigationStateComponent.createObject(this)
        verify(State !== null)
        compare(State.pageIndex, 0)

        State.navigationIndex = 2
        compare(State.pageIndex, 2)

        State.cookingVisible = true
        compare(State.pageIndex, 4)

        State.cookingVisible = false
        compare(State.pageIndex, 2)
        State.destroy()
    }

    function test_LocalDateDoesNotUseUtcDay() {
        const Formatter = dateFormatterComponent.createObject(this)
        verify(Formatter !== null)
        const LocalMorning = new Date(2026, 8, 1, 0, 30, 0)
        compare(Formatter.formatLocalDate(LocalMorning), "2026-09-01")
        Formatter.destroy()
    }

    Component {
        id: navigationStateComponent
        NavigationState {}
    }

    Component {
        id: dateFormatterComponent
        DateFormatter {}
    }
}
