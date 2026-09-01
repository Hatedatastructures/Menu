add_library(MenuApplication STATIC
    AdminIngredientApplicationService.cpp
    AdminRecipeApplicationService.cpp
    CookingSessionApplicationService.cpp
    FeedbackApplicationService.cpp
    MealPlanApplicationService.cpp
    RecipeApplicationService.cpp
    RuleBasedRecommendationProvider.cpp
    StorageExecutor.cpp
)
target_include_directories(MenuApplication PUBLIC "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuApplication PUBLIC MenuDomain MenuFoundation MenuBoost)
target_compile_features(MenuApplication PUBLIC cxx_std_20)
