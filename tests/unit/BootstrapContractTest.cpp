#include <gtest/gtest.h>

#include <Foundation/Result.hpp>

TEST(BootstrapContractTest, BuildsWithCxx20AndFoundationTarget) {
    const Menu::Foundation::Result<int> Result = 7;
    ASSERT_TRUE(Result.HasValue());
    EXPECT_EQ(Result.Value(), 7);
}
