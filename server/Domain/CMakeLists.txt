add_library(MenuDomain STATIC
    PlanRules.cpp
    RecipeRules.cpp
    Recommendation.cpp
)
target_include_directories(MenuDomain PUBLIC "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuDomain INTERFACE MenuFoundation)
target_compile_features(MenuDomain PUBLIC cxx_std_20)
