#include "Api/ApiParsing.hpp"

#include "Api/ApiParsingSupport.hpp"

#include <boost/json/object.hpp>

#include <cstdint>
#include <utility>

namespace Menu::Api::Parsing::Support {

Foundation::Error InvalidRequest(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

Foundation::Result<double> ReadJsonDouble(
    const boost::json::object& Object,
    std::string_view Key) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr) {
        return Foundation::Result<double>::FromError(
            InvalidRequest("菜谱数字字段缺失"));
    }
    if (Value->is_double()) {
        return Value->as_double();
    }
    if (Value->is_int64()) {
        return static_cast<double>(Value->as_int64());
    }
    return Foundation::Result<double>::FromError(
        InvalidRequest("菜谱数字字段无效"));
}

Foundation::Result<bool> ReadJsonBool(
    const boost::json::object& Object,
    std::string_view Key) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr || !Value->is_bool()) {
        return Foundation::Result<bool>::FromError(
            InvalidRequest("菜谱布尔字段无效"));
    }
    return Value->as_bool();
}

}  // namespace Menu::Api::Parsing::Support

namespace Menu::Api::Parsing {

Foundation::Result<std::string> ReadRequiredString(
    const boost::json::object& Object,
    std::string_view Key,
    std::size_t MaximumLength) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr || !Value->is_string() || Value->as_string().empty() ||
        Value->as_string().size() > MaximumLength) {
        return Foundation::Result<std::string>::FromError(
            Support::InvalidRequest("请求字段无效"));
    }
    return std::string(Value->as_string().c_str());
}

Foundation::Result<int> ReadJsonInteger(
    const boost::json::object& Object,
    std::string_view Key,
    int Minimum,
    int Maximum) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr || !Value->is_int64()) {
        return Foundation::Result<int>::FromError(
            Support::InvalidRequest("请求数字字段无效"));
    }
    const std::int64_t Parsed = Value->as_int64();
    if (Parsed < Minimum || Parsed > Maximum) {
        return Foundation::Result<int>::FromError(
            Support::InvalidRequest("请求数字字段超出范围"));
    }
    return static_cast<int>(Parsed);
}

Foundation::Result<std::vector<std::string>> ReadJsonStringArray(
    const boost::json::object& Object,
    std::string_view Key) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr) {
        return std::vector<std::string>();
    }
    if (!Value->is_array() || Value->as_array().size() > 64U) {
        return Foundation::Result<std::vector<std::string>>::FromError(
            Support::InvalidRequest("请求数组字段无效"));
    }
    std::vector<std::string> Values;
    for (const boost::json::value& Item : Value->as_array()) {
        if (!Item.is_string() || Item.as_string().size() > 128U) {
            return Foundation::Result<std::vector<std::string>>::FromError(
                Support::InvalidRequest("请求数组元素无效"));
        }
        Values.emplace_back(Item.as_string().c_str());
    }
    return Values;
}

}  // namespace Menu::Api::Parsing
