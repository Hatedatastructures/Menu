#pragma once

#include <string>
#include <vector>

namespace Menu::Domain {

struct Ingredient {
    std::string Id;
    std::string Name;
    std::vector<std::string> Aliases;
    std::string Category;
    std::string DefaultUnit;
    bool IsPantryStaple = false;
    std::string SubstituteGroup;
    std::string StoreSkuMapping;
};

}  // namespace Menu::Domain
