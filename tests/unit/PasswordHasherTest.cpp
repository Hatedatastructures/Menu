#include <gtest/gtest.h>

#include <Infrastructure/PasswordHasher.hpp>

#include <string>

TEST(PasswordHasherTest, HashesAndVerifiesPasswordWithoutStoringPlaintext) {
    const std::string Password = "Menu-Cook-Password-2026";
    const auto HashResult = Menu::Infrastructure::PasswordHasher::Hash(Password);

    ASSERT_TRUE(HashResult.HasValue());
    EXPECT_NE(HashResult.Value(), Password);
    EXPECT_TRUE(HashResult.Value().starts_with("scrypt$"));

    const auto VerifyResult = Menu::Infrastructure::PasswordHasher::Verify(
        Password, HashResult.Value());
    ASSERT_TRUE(VerifyResult.HasValue());
    EXPECT_TRUE(VerifyResult.Value());
}

TEST(PasswordHasherTest, RejectsWrongOrMalformedPasswordHash) {
    const auto HashResult = Menu::Infrastructure::PasswordHasher::Hash("Correct password");
    ASSERT_TRUE(HashResult.HasValue());

    const auto WrongPassword = Menu::Infrastructure::PasswordHasher::Verify(
        "Wrong password", HashResult.Value());
    const auto MalformedHash = Menu::Infrastructure::PasswordHasher::Verify(
        "Correct password", "scrypt$invalid");

    ASSERT_TRUE(WrongPassword.HasValue());
    EXPECT_FALSE(WrongPassword.Value());
    ASSERT_TRUE(MalformedHash.HasValue());
    EXPECT_FALSE(MalformedHash.Value());
}

TEST(PasswordHasherTest, RejectsPasswordInputAboveConfiguredLimitBeforeDerivation) {
    const auto HashResult = Menu::Infrastructure::PasswordHasher::Hash(
        "Correct password for the test");
    ASSERT_TRUE(HashResult.HasValue());

    const std::string LongPassword(129, 'x');
    const auto VerifyResult = Menu::Infrastructure::PasswordHasher::Verify(
        LongPassword, HashResult.Value());

    ASSERT_FALSE(VerifyResult.HasValue());
    EXPECT_EQ(VerifyResult.ErrorValue().CodeValue(),
              Menu::Foundation::ErrorCode::InvalidArgument);
}
