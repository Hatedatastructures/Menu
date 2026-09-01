#include "Api/ApiParsing.hpp"

#include "Api/ApiParsingSupport.hpp"

#include <boost/json/parse.hpp>

#include <set>
#include <string>
#include <utility>

namespace Menu::Api::Parsing {

Foundation::Result<Domain::Recipe> ReadRecipePayload(std::string_view Body) {
    if (Body.empty() || Body.size() > 1024U * 1024U) {
        return Foundation::Result<Domain::Recipe>::FromError(
            Support::InvalidRequest("菜谱请求体为空或过大"));
    }
    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Body, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Foundation::Result<Domain::Recipe>::FromError(
            Support::InvalidRequest("菜谱 JSON 无效"));
    }
    static const std::set<std::string> AllowedKeys = {
        "id", "slug", "name", "cuisine", "description", "prepMinutes",
        "cookMinutes", "servings", "difficulty", "imagePath", "status",
        "allergens", "cookware", "ingredients", "steps"};
    for (const auto& Entry : Parsed.as_object()) {
        if (!AllowedKeys.contains(std::string(Entry.key().data(), Entry.key().size()))) {
            return Foundation::Result<Domain::Recipe>::FromError(
                Support::InvalidRequest("菜谱包含未知字段"));
        }
    }

    const boost::json::object& Object = Parsed.as_object();
    const auto Id = ReadRequiredString(Object, "id", 128U);
    const auto Slug = ReadRequiredString(Object, "slug", 160U);
    const auto Name = ReadRequiredString(Object, "name", 160U);
    const auto Cuisine = ReadRequiredString(Object, "cuisine", 32U);
    const auto Description = ReadRequiredString(Object, "description", 4096U);
    const auto ImagePath = ReadRequiredString(Object, "imagePath", 512U);
    const auto Status = ReadRequiredString(Object, "status", 16U);
    const auto PrepMinutes = ReadJsonInteger(Object, "prepMinutes", 0, 24 * 60);
    const auto CookMinutes = ReadJsonInteger(Object, "cookMinutes", 0, 24 * 60);
    const auto Servings = ReadJsonInteger(Object, "servings", 1, 24);
    const auto Difficulty = ReadJsonInteger(Object, "difficulty", 1, 5);
    const auto Allergens = ReadJsonStringArray(Object, "allergens");
    const auto Cookware = ReadJsonStringArray(Object, "cookware");
    if (!Id.HasValue() || !Slug.HasValue() || !Name.HasValue() || !Cuisine.HasValue() ||
        !Description.HasValue() || !ImagePath.HasValue() || !Status.HasValue() ||
        !PrepMinutes.HasValue() || !CookMinutes.HasValue() || !Servings.HasValue() ||
        !Difficulty.HasValue() || !Allergens.HasValue() || !Cookware.HasValue()) {
        return Foundation::Result<Domain::Recipe>::FromError(
            Support::InvalidRequest("菜谱基础字段无效"));
    }

    const boost::json::value* IngredientsValue = Object.if_contains("ingredients");
    const boost::json::value* StepsValue = Object.if_contains("steps");
    if (IngredientsValue == nullptr || !IngredientsValue->is_array() ||
        IngredientsValue->as_array().size() > 128U || StepsValue == nullptr ||
        !StepsValue->is_array() || StepsValue->as_array().size() > 128U) {
        return Foundation::Result<Domain::Recipe>::FromError(
            Support::InvalidRequest("菜谱食材或步骤无效"));
    }

    Domain::Recipe RecipeValue;
    RecipeValue.Id = Id.Value();
    RecipeValue.Slug = Slug.Value();
    RecipeValue.Name = Name.Value();
    RecipeValue.Cuisine = Cuisine.Value();
    RecipeValue.Description = Description.Value();
    RecipeValue.PrepMinutes = PrepMinutes.Value();
    RecipeValue.CookMinutes = CookMinutes.Value();
    RecipeValue.Servings = Servings.Value();
    RecipeValue.Difficulty = Difficulty.Value();
    RecipeValue.ImagePath = ImagePath.Value();
    RecipeValue.Status = Status.Value();
    RecipeValue.Allergens = Allergens.Value();
    RecipeValue.Cookware = Cookware.Value();

    for (const boost::json::value& Value : IngredientsValue->as_array()) {
        if (!Value.is_object()) {
            return Foundation::Result<Domain::Recipe>::FromError(
                Support::InvalidRequest("菜谱食材项无效"));
        }
        const boost::json::object& Item = Value.as_object();
        static const std::set<std::string> IngredientKeys = {
            "ingredientId", "quantity", "unit", "required", "servingFactor", "preparation"};
        for (const auto& Entry : Item) {
            if (!IngredientKeys.contains(std::string(Entry.key().data(), Entry.key().size()))) {
                return Foundation::Result<Domain::Recipe>::FromError(
                    Support::InvalidRequest("菜谱食材包含未知字段"));
            }
        }
        const auto IngredientId = ReadRequiredString(Item, "ingredientId", 128U);
        const auto Quantity = Support::ReadJsonDouble(Item, "quantity");
        const auto Unit = ReadRequiredString(Item, "unit", 16U);
        const auto Required = Support::ReadJsonBool(Item, "required");
        const auto ServingFactor = Support::ReadJsonDouble(Item, "servingFactor");
        std::string Preparation;
        if (const boost::json::value* PreparationValue = Item.if_contains("preparation");
            PreparationValue != nullptr) {
            if (!PreparationValue->is_string() || PreparationValue->as_string().size() > 256U) {
                return Foundation::Result<Domain::Recipe>::FromError(
                    Support::InvalidRequest("菜谱处理方式无效"));
            }
            Preparation = std::string(PreparationValue->as_string().c_str());
        }
        if (!IngredientId.HasValue() || !Quantity.HasValue() || !Unit.HasValue() ||
            !Required.HasValue() || !ServingFactor.HasValue()) {
            return Foundation::Result<Domain::Recipe>::FromError(
                Support::InvalidRequest("菜谱食材字段无效"));
        }
        RecipeValue.Ingredients.emplace_back(
            IngredientId.Value(), Quantity.Value(), Unit.Value(), Required.Value(),
            ServingFactor.Value(), std::move(Preparation));
    }

    for (const boost::json::value& Value : StepsValue->as_array()) {
        if (!Value.is_object()) {
            return Foundation::Result<Domain::Recipe>::FromError(
                Support::InvalidRequest("菜谱步骤项无效"));
        }
        const boost::json::object& Item = Value.as_object();
        const auto StepOrder = ReadJsonInteger(Item, "stepOrder", 1, 128);
        const auto Title = ReadRequiredString(Item, "title", 160U);
        const auto Instruction = ReadRequiredString(Item, "instruction", 4096U);
        const auto Duration = ReadJsonInteger(Item, "durationSeconds", 0, 24 * 60 * 60);
        const auto HasTimer = Support::ReadJsonBool(Item, "hasTimer");
        if (!StepOrder.HasValue() || !Title.HasValue() || !Instruction.HasValue() ||
            !Duration.HasValue() || !HasTimer.HasValue()) {
            return Foundation::Result<Domain::Recipe>::FromError(
                Support::InvalidRequest("菜谱步骤字段无效"));
        }
        RecipeValue.Steps.push_back(Domain::RecipeStep{
            StepOrder.Value(), Title.Value(), Instruction.Value(), Duration.Value(),
            HasTimer.Value()});
    }
    return RecipeValue;
}

Foundation::Result<Domain::Ingredient> ReadIngredientPayload(std::string_view Body) {
    if (Body.empty() || Body.size() > 64U * 1024U) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            Support::InvalidRequest("食材请求体无效"));
    }
    boost::system::error_code ParseError;
    const boost::json::value Parsed = boost::json::parse(Body, ParseError);
    if (ParseError || !Parsed.is_object()) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            Support::InvalidRequest("食材 JSON 无效"));
    }
    static const std::set<std::string> AllowedKeys = {
        "id", "name", "aliases", "category", "defaultUnit", "isPantryStaple",
        "substituteGroup", "storeSkuMapping"};
    for (const auto& Entry : Parsed.as_object()) {
        if (!AllowedKeys.contains(std::string(Entry.key().data(), Entry.key().size()))) {
            return Foundation::Result<Domain::Ingredient>::FromError(
                Support::InvalidRequest("食材包含未知字段"));
        }
    }
    const boost::json::object& Object = Parsed.as_object();
    const auto Id = ReadRequiredString(Object, "id", 128U);
    const auto Name = ReadRequiredString(Object, "name", 160U);
    const auto Category = ReadRequiredString(Object, "category", 64U);
    const auto DefaultUnit = ReadRequiredString(Object, "defaultUnit", 16U);
    const auto Aliases = ReadJsonStringArray(Object, "aliases");
    const auto IsPantryStaple = Support::ReadJsonBool(Object, "isPantryStaple");
    if (!Id.HasValue() || !Name.HasValue() || !Category.HasValue() ||
        !DefaultUnit.HasValue() || !Aliases.HasValue() || !IsPantryStaple.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            Support::InvalidRequest("食材基础字段无效"));
    }
    const auto ReadOptional = [&Object](std::string_view Key, std::size_t MaximumLength) {
        std::string Value;
        const boost::json::value* JsonValue = Object.if_contains(Key);
        if (JsonValue == nullptr) {
            return Foundation::Result<std::string>(Value);
        }
        if (!JsonValue->is_string() || JsonValue->as_string().size() > MaximumLength) {
            return Foundation::Result<std::string>::FromError(
                Support::InvalidRequest("食材可选字段无效"));
        }
        Value = std::string(JsonValue->as_string().c_str());
        return Foundation::Result<std::string>(std::move(Value));
    };
    const auto SubstituteGroup = ReadOptional("substituteGroup", 64U);
    const auto StoreSkuMapping = ReadOptional("storeSkuMapping", 256U);
    if (!SubstituteGroup.HasValue() || !StoreSkuMapping.HasValue()) {
        return Foundation::Result<Domain::Ingredient>::FromError(
            Support::InvalidRequest("食材可选字段无效"));
    }
    Domain::Ingredient IngredientValue;
    IngredientValue.Id = Id.Value();
    IngredientValue.Name = Name.Value();
    IngredientValue.Aliases = Aliases.Value();
    IngredientValue.Category = Category.Value();
    IngredientValue.DefaultUnit = DefaultUnit.Value();
    IngredientValue.IsPantryStaple = IsPantryStaple.Value();
    IngredientValue.SubstituteGroup = SubstituteGroup.Value();
    IngredientValue.StoreSkuMapping = StoreSkuMapping.Value();
    return IngredientValue;
}

}  // namespace Menu::Api::Parsing
