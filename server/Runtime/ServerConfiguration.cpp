#include "Runtime/ServerConfiguration.hpp"

#include <boost/json/parse.hpp>

#include <fstream>
#include <iterator>
#include <string_view>
#include <utility>

namespace Menu::Runtime {
namespace {

Foundation::Error InvalidConfiguration(std::string Message) {
    return Foundation::Error(
        Foundation::ErrorCode::InvalidArgument, std::move(Message));
}

Foundation::Result<std::string> ReadString(
    const boost::json::object& Object,
    std::string_view Key,
    std::string DefaultValue,
    std::size_t MaximumLength) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr) {
        return DefaultValue;
    }
    if (!Value->is_string() || Value->as_string().size() > MaximumLength) {
        return Foundation::Result<std::string>::FromError(
            InvalidConfiguration("配置字符串无效"));
    }
    return std::string(Value->as_string().c_str());
}

Foundation::Result<std::int64_t> ReadInteger(
    const boost::json::object& Object,
    std::string_view Key,
    std::int64_t DefaultValue,
    std::int64_t Minimum,
    std::int64_t Maximum) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr) {
        return DefaultValue;
    }
    if (!Value->is_int64()) {
        return Foundation::Result<std::int64_t>::FromError(
            InvalidConfiguration("配置数字无效"));
    }
    const std::int64_t Parsed = Value->as_int64();
    if (Parsed < Minimum || Parsed > Maximum) {
        return Foundation::Result<std::int64_t>::FromError(
            InvalidConfiguration("配置数字超出范围"));
    }
    return Parsed;
}

Foundation::Result<std::vector<std::string>> ReadStringArray(
    const boost::json::object& Object,
    std::string_view Key) {
    const boost::json::value* Value = Object.if_contains(Key);
    if (Value == nullptr) {
        return std::vector<std::string>();
    }
    if (!Value->is_array() || Value->as_array().size() > 32U) {
        return Foundation::Result<std::vector<std::string>>::FromError(
            InvalidConfiguration("配置数组无效"));
    }
    std::vector<std::string> Values;
    for (const boost::json::value& Item : Value->as_array()) {
        if (!Item.is_string() || Item.as_string().size() > 256U) {
            return Foundation::Result<std::vector<std::string>>::FromError(
                InvalidConfiguration("配置数组元素无效"));
        }
        Values.emplace_back(Item.as_string().c_str());
    }
    return Values;
}

}  // namespace

Foundation::Result<ServerConfiguration> LoadServerConfiguration(
    const std::filesystem::path& ConfigPath) {
    std::ifstream Input(ConfigPath, std::ios::binary);
    if (!Input) {
        return Foundation::Result<ServerConfiguration>::FromError(
            InvalidConfiguration("无法打开配置文件"));
    }
    const std::string Text{
        std::istreambuf_iterator<char>(Input), std::istreambuf_iterator<char>()};
    if (Text.size() > 128U * 1024U) {
        return Foundation::Result<ServerConfiguration>::FromError(
            InvalidConfiguration("配置文件过大"));
    }

    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Text, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Foundation::Result<ServerConfiguration>::FromError(
            InvalidConfiguration("配置 JSON 无效"));
    }

    const boost::json::object& Object = Parsed.as_object();
    ServerConfiguration Configuration;
    const auto Address = ReadString(Object, "address", Configuration.Http.Address, 64U);
    const auto Port = ReadInteger(Object, "port", Configuration.Http.Port, 1, 65535);
    const auto DatabasePath = ReadString(
        Object, "databasePath", Configuration.DatabasePath.string(), 512U);
    const auto MediaDirectory = ReadString(
        Object, "mediaDirectory", Configuration.MediaDirectory.string(), 512U);
    const auto MaxBodyBytes = ReadInteger(
        Object, "maxBodyBytes", Configuration.Http.MaxBodyBytes, 1024, 16 * 1024 * 1024);
    const auto MaxTargetBytes = ReadInteger(
        Object, "maxTargetBytes", Configuration.Http.MaxTargetBytes, 256, 1024 * 1024);
    const auto TimeoutSeconds = ReadInteger(
        Object, "timeoutSeconds", Configuration.Http.TimeoutSeconds, 1, 300);
    const auto CorsOrigins = ReadStringArray(Object, "corsOrigins");
    if (!Address.HasValue() || !Port.HasValue() || !DatabasePath.HasValue() ||
        !MediaDirectory.HasValue() || !MaxBodyBytes.HasValue() ||
        !MaxTargetBytes.HasValue() || !TimeoutSeconds.HasValue() ||
        !CorsOrigins.HasValue()) {
        return Foundation::Result<ServerConfiguration>::FromError(
            InvalidConfiguration("配置字段无效"));
    }

    Configuration.Http.Address = Address.Value();
    Configuration.Http.Port = static_cast<std::uint16_t>(Port.Value());
    Configuration.DatabasePath = DatabasePath.Value();
    Configuration.MediaDirectory = MediaDirectory.Value();
    Configuration.Http.MaxBodyBytes = static_cast<std::size_t>(MaxBodyBytes.Value());
    Configuration.Http.MaxTargetBytes = static_cast<std::size_t>(MaxTargetBytes.Value());
    Configuration.Http.TimeoutSeconds = static_cast<int>(TimeoutSeconds.Value());
    Configuration.CorsOrigins = CorsOrigins.Value();
    return Configuration;
}

}  // namespace Menu::Runtime
