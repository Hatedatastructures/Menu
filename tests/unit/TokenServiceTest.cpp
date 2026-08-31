#include <gtest/gtest.h>

#include <Infrastructure/TokenService.hpp>

#include <string>

TEST(TokenServiceTest, CreatesUniqueOpaqueTokensAndDeterministicHashes) {
    const auto FirstToken = Menu::Infrastructure::TokenService::CreateToken();
    const auto SecondToken = Menu::Infrastructure::TokenService::CreateToken();
    ASSERT_TRUE(FirstToken.HasValue());
    ASSERT_TRUE(SecondToken.HasValue());
    EXPECT_EQ(FirstToken.Value().size(), 64U);
    EXPECT_NE(FirstToken.Value(), SecondToken.Value());

    const auto FirstHash = Menu::Infrastructure::TokenService::HashToken(FirstToken.Value());
    const auto RepeatHash = Menu::Infrastructure::TokenService::HashToken(FirstToken.Value());
    ASSERT_TRUE(FirstHash.HasValue());
    ASSERT_TRUE(RepeatHash.HasValue());
    EXPECT_EQ(FirstHash.Value(), RepeatHash.Value());
    EXPECT_NE(FirstHash.Value(), FirstToken.Value());
}
