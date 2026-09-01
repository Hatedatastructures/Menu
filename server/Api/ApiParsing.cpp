#include "Api/ApiParsing.hpp"

#include "Api/ApiParsingSupport.hpp"

#include <charconv>
#include <cctype>
#include <map>
#include <ranges>

namespace Menu::Api::Parsing {
namespace {

int HexValue(char Character) {
    if (Character >= '0' && Character <= '9') {
        return Character - '0';
    }
    if (Character >= 'a' && Character <= 'f') {
        return Character - 'a' + 10;
    }
    if (Character >= 'A' && Character <= 'F') {
        return Character - 'A' + 10;
    }
    return -1;
}

Foundation::Result<std::string> DecodeComponent(std::string_view Value) {
    std::string Decoded;
    Decoded.reserve(Value.size());
    for (std::size_t Index = 0; Index < Value.size(); ++Index) {
        if (Value[Index] == '%') {
            if (Index + 2 >= Value.size()) {
                return Foundation::Result<std::string>::FromError(
                    Support::InvalidRequest("查询参数编码无效"));
            }
            const int High = HexValue(Value[Index + 1]);
            const int Low = HexValue(Value[Index + 2]);
            if (High < 0 || Low < 0) {
                return Foundation::Result<std::string>::FromError(
                    Support::InvalidRequest("查询参数编码无效"));
            }
            Decoded.push_back(static_cast<char>((High << 4) | Low));
            Index += 2;
        } else if (Value[Index] == '+') {
            Decoded.push_back(' ');
        } else {
            Decoded.push_back(Value[Index]);
        }
    }
    return Decoded;
}

}  // namespace

TargetParts SplitTarget(std::string_view Target) {
    const std::size_t QueryPosition = Target.find('?');
    if (QueryPosition == std::string_view::npos) {
        return {std::string(Target), {}};
    }
    return {std::string(Target.substr(0, QueryPosition)),
            std::string(Target.substr(QueryPosition + 1))};
}

Foundation::Result<std::map<std::string, std::string>> ParseQuery(
    std::string_view Query) {
    std::map<std::string, std::string> Values;
    if (Query.empty()) {
        return Values;
    }
    std::size_t Start = 0;
    while (Start <= Query.size()) {
        const std::size_t End = Query.find('&', Start);
        const std::string_view Pair = Query.substr(
            Start, End == std::string_view::npos ? Query.size() - Start : End - Start);
        const std::size_t Separator = Pair.find('=');
        if (Separator == std::string_view::npos || Separator == 0) {
            return Foundation::Result<std::map<std::string, std::string>>::FromError(
                Support::InvalidRequest("查询参数格式无效"));
        }
        const auto Key = DecodeComponent(Pair.substr(0, Separator));
        const auto Value = DecodeComponent(Pair.substr(Separator + 1));
        if (!Key.HasValue() || !Value.HasValue() || Key.Value().size() > 64U ||
            Value.Value().size() > 256U) {
            return Foundation::Result<std::map<std::string, std::string>>::FromError(
                Support::InvalidRequest("查询参数过长或编码无效"));
        }
        if (!Values.emplace(Key.Value(), Value.Value()).second) {
            return Foundation::Result<std::map<std::string, std::string>>::FromError(
                Support::InvalidRequest("查询参数重复"));
        }
        if (End == std::string_view::npos) {
            break;
        }
        Start = End + 1;
    }
    return Values;
}

Foundation::Result<int> ParseInteger(
    std::string_view Value,
    int Minimum,
    int Maximum) {
    int Parsed = 0;
    const auto ParseResult = std::from_chars(
        Value.data(), Value.data() + Value.size(), Parsed);
    if (ParseResult.ec != std::errc() || ParseResult.ptr != Value.data() + Value.size() ||
        Parsed < Minimum || Parsed > Maximum) {
        return Foundation::Result<int>::FromError(
            Support::InvalidRequest("数字参数无效"));
    }
    return Parsed;
}

Foundation::Result<RecipeListOptions> ReadRecipeListOptions(
    std::string_view Query) {
    const auto QueryResult = ParseQuery(Query);
    if (!QueryResult.HasValue()) {
        return Foundation::Result<RecipeListOptions>::FromError(QueryResult.ErrorValue());
    }
    RecipeListOptions Options;
    for (const auto& [Key, Value] : QueryResult.Value()) {
        if (Key == "cuisine") {
            if (Value.empty()) {
                return Foundation::Result<RecipeListOptions>::FromError(
                    Support::InvalidRequest("菜系不能为空"));
            }
            Options.Cuisine = Value;
        } else if (Key == "maxMinutes") {
            const auto Parsed = ParseInteger(Value, 0, 24 * 60);
            if (!Parsed.HasValue()) {
                return Foundation::Result<RecipeListOptions>::FromError(Parsed.ErrorValue());
            }
            Options.MaxMinutes = Parsed.Value();
            Options.HasMaxMinutes = true;
        } else if (Key == "difficulty") {
            const auto Parsed = ParseInteger(Value, 1, 5);
            if (!Parsed.HasValue()) {
                return Foundation::Result<RecipeListOptions>::FromError(Parsed.ErrorValue());
            }
            Options.Difficulty = Parsed.Value();
        } else if (Key == "limit") {
            const auto Parsed = ParseInteger(Value, 1, 100);
            if (!Parsed.HasValue()) {
                return Foundation::Result<RecipeListOptions>::FromError(Parsed.ErrorValue());
            }
            Options.Limit = static_cast<std::size_t>(Parsed.Value());
        } else {
            return Foundation::Result<RecipeListOptions>::FromError(
                Support::InvalidRequest("未知查询参数"));
        }
    }
    return Options;
}

Foundation::Result<std::pair<std::string, std::string>> ReadPlanQuery(
    std::string_view Query) {
    const auto QueryResult = ParseQuery(Query);
    if (!QueryResult.HasValue()) {
        return Foundation::Result<std::pair<std::string, std::string>>::FromError(
            QueryResult.ErrorValue());
    }
    const auto From = QueryResult.Value().find("from");
    const auto To = QueryResult.Value().find("to");
    if (From == QueryResult.Value().end() || To == QueryResult.Value().end() ||
        QueryResult.Value().size() != 2U) {
        return Foundation::Result<std::pair<std::string, std::string>>::FromError(
            Support::InvalidRequest("计划日期范围缺失"));
    }
    return std::pair<std::string, std::string>{From->second, To->second};
}

bool IsSafeIdentifier(std::string_view Value) {
    if (Value.empty() || Value.size() > 128U) {
        return false;
    }
    return std::ranges::all_of(Value, [](char Character) {
        return std::isalnum(static_cast<unsigned char>(Character)) != 0 ||
               Character == '.' || Character == '-' || Character == '_';
    });
}

}  // namespace Menu::Api::Parsing
