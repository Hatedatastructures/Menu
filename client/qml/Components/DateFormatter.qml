import QtQuick

QtObject {
    function formatLocalDate(dateValue) {
        function twoDigits(value) {
            return value < 10 ? "0" + value : value.toString()
        }
        return dateValue.getFullYear().toString() + "-" +
            twoDigits(dateValue.getMonth() + 1) + "-" +
            twoDigits(dateValue.getDate())
    }
}
