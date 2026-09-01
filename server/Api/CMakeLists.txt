add_library(MenuApi STATIC
    AdminRoutes.cpp
    ApiErrors.cpp
    ApiMiddleware.cpp
    ApiParsing.cpp
    ApiParsingJson.cpp
    ApiParsingRecipe.cpp
    ApiParsingAuth.cpp
    ApiParsingRecommendation.cpp
    ApiParsingWorkflow.cpp
    AuthDtos.cpp
    AuthRoutes.cpp
    ApiRouter.cpp
    RecipeDtos.cpp
    RecipeRoutes.cpp
    WorkflowDtos.cpp
    WorkflowRoutes.cpp
)
target_include_directories(MenuApi PUBLIC "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuApi PUBLIC
    MenuApplication MenuTransportCore MenuBoost MenuBoostJson)
target_compile_features(MenuApi PUBLIC cxx_std_20)
