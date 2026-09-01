#pragma once

#include <Foundation/Result.hpp>

#include <boost/json/object.hpp>

#include <string>
#include <string_view>

namespace Menu::Api::Parsing::Support {

Foundation::Error InvalidRequest(std::string Message);

Foundation::Result<double> ReadJsonDouble(
    const boost::json::object& Object,
    std::string_view Key);

Foundation::Result<bool> ReadJsonBool(
    const boost::json::object& Object,
    std::string_view Key);

}  // namespace Menu::Api::Parsing::Support
