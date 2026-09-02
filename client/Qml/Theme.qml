pragma Singleton
import QtQuick

QtObject {
    readonly property bool dark: typeof MenuTheme !== "undefined"
        ? MenuTheme.Dark : Qt.application.styleHints.colorScheme === Qt.Dark
    // A quiet neutral canvas keeps the food imagery and the action color in focus.
    readonly property color background: dark ? "#101512" : "#f6f7f3"
    readonly property color paper: background
    readonly property color surface: dark ? "#171e1a" : "#ffffff"
    readonly property color surfaceRaised: dark ? "#202a24" : "#edf1ec"
    readonly property color ink: dark ? "#f3f5ef" : "#17231c"
    readonly property color text: ink
    readonly property color mutedText: dark ? "#aab7ae" : "#65736a"
    readonly property color line: dark ? "#35423a" : "#d9e1da"
    readonly property color green: dark ? "#7bd19d" : "#287951"
    readonly property color greenStrong: dark ? "#327e59" : "#236b49"
    readonly property color coral: dark ? "#ff9c80" : "#b25440"
    readonly property color coralStrong: dark ? "#b65440" : "#a84938"
    readonly property color amber: dark ? "#f3c76a" : "#a86d15"
    readonly property color sky: dark ? "#91cbd2" : "#477c83"
    readonly property color mutedPaper: dark ? "#cbd7ce" : "#cbd7ce"
    readonly property color headerBackground: dark ? "#0b100d" : "#e7f0ea"
    readonly property color headerText: dark ? "#f3f5ef" : "#17231c"
    readonly property color headerMuted: dark ? "#a9b9af" : "#5e7065"
    readonly property color headerLine: dark ? "#1f2c24" : "#d4e0d7"
    readonly property color onAccent: "#ffffff"
    readonly property color focus: dark ? "#9ae4b8" : "#1e6645"
    readonly property int spacingXs: 4
    readonly property int spacingSm: 8
    readonly property int spacingMd: 16
    readonly property int spacingLg: 24
    readonly property int spacingXl: 32
    readonly property int radius: 8
    readonly property int radiusSmall: 5
    readonly property int touchTarget: 48
    readonly property int topBarHeight: 72
    readonly property int pageMargin: 20
    readonly property int contentMaxWidth: 1120

    function tint(colorValue, opacityValue) {
        return Qt.rgba(colorValue.r, colorValue.g, colorValue.b, opacityValue)
    }

    function onColor(colorValue) {
        var luminance = 0.2126 * colorValue.r +
            0.7152 * colorValue.g + 0.0722 * colorValue.b
        return luminance > 0.58 ? ink : onAccent
    }
}
