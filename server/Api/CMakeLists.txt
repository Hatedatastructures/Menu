add_library(MenuApi STATIC
    ApiErrors.cpp
    ApiRouter.cpp
    RecipeDtos.cpp
)
target_include_directories(MenuApi PUBLIC "${PROJECT_SOURCE_DIR}/server")
target_link_libraries(MenuApi PUBLIC
    MenuApplication MenuTransportCore MenuBoost MenuBoostJson)
target_compile_features(MenuApi PUBLIC cxx_std_20)
