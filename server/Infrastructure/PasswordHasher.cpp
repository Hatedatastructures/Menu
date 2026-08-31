#include "Infrastructure/PasswordHasher.hpp"

#include <openssl/crypto.h>
#include <openssl/evp.h>
#include <openssl/rand.h>

#include <array>
#include <cstddef>
#include <cstdint>
#include <optional>
#include <string>
#include <string_view>
#include <utility>
#include <vector>

namespace Menu::Infrastructure {
namespace {

constexpr std::uint64_t ScryptCost = 32768;
constexpr std::uint64_t ScryptBlockSize = 8;
constexpr std::uint64_t ScryptParallelism = 1;
constexpr std::uint64_t ScryptMaxMemory = 64U * 1024U * 1024U;
constexpr std::size_t SaltSize = 16;
constexpr std::size_t HashSize = 32;

constexpr char HexDigits[] = "0123456789abcdef";

Foundation::Error CryptoError(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::StorageUnavailable, std::move(Message));
}

std::string EncodeHex(const unsigned char* Bytes, std::size_t Size) {
    std::string Encoded;
    Encoded.reserve(Size * 2U);
    for (std::size_t Index = 0; Index < Size; ++Index) {
        const unsigned char Byte = Bytes[Index];
        Encoded.push_back(HexDigits[(Byte >> 4U) & 0x0FU]);
        Encoded.push_back(HexDigits[Byte & 0x0FU]);
    }
    return Encoded;
}

int HexValue(char Character) {
    if (Character >= '0' && Character <= '9') {
        return Character - '0';
    }
    if (Character >= 'a' && Character <= 'f') {
        return Character - 'a' + 10;
    }
    if (Character >= 'A' && Character <= 'F') {
        return Character - 'A' + 10;
    }
    return -1;
}

std::optional<std::vector<unsigned char>> DecodeHex(
    std::string_view Encoded,
    std::size_t ExpectedSize) {
    if (Encoded.size() != ExpectedSize * 2U) {
        return std::nullopt;
    }
    std::vector<unsigned char> Bytes(ExpectedSize);
    for (std::size_t Index = 0; Index < ExpectedSize; ++Index) {
        const int High = HexValue(Encoded[Index * 2U]);
        const int Low = HexValue(Encoded[Index * 2U + 1U]);
        if (High < 0 || Low < 0) {
            return std::nullopt;
        }
        Bytes[Index] = static_cast<unsigned char>((High << 4) | Low);
    }
    return Bytes;
}

Foundation::Result<std::array<unsigned char, HashSize>> DeriveHash(
    std::string_view Password,
    const std::vector<unsigned char>& Salt) {
    std::array<unsigned char, HashSize> HashValue{};
    const int ResultCode = EVP_PBE_scrypt(
        Password.data(),
        Password.size(),
        Salt.data(),
        Salt.size(),
        ScryptCost,
        ScryptBlockSize,
        ScryptParallelism,
        ScryptMaxMemory,
        HashValue.data(),
        HashValue.size());
    if (ResultCode != 1) {
        return Foundation::Result<std::array<unsigned char, HashSize>>::FromError(
            CryptoError("密码派生失败"));
    }
    return HashValue;
}

}  // namespace

Foundation::Result<std::string> PasswordHasher::Hash(std::string_view Password) {
    std::array<unsigned char, SaltSize> Salt{};
    if (RAND_bytes(Salt.data(), static_cast<int>(Salt.size())) != 1) {
        return Foundation::Result<std::string>::FromError(CryptoError("随机盐生成失败"));
    }
    const std::vector<unsigned char> SaltVector(Salt.begin(), Salt.end());
    const auto HashResult = DeriveHash(Password, SaltVector);
    if (!HashResult.HasValue()) {
        return Foundation::Result<std::string>::FromError(HashResult.ErrorValue());
    }
    return "scrypt$" + std::to_string(ScryptCost) + "$" +
           std::to_string(ScryptBlockSize) + "$" +
           std::to_string(ScryptParallelism) + "$" +
           EncodeHex(Salt.data(), Salt.size()) + "$" +
           EncodeHex(HashResult.Value().data(), HashResult.Value().size());
}

Foundation::Result<bool> PasswordHasher::Verify(
    std::string_view Password,
    std::string_view EncodedHash) {
    if (Password.empty() || Password.size() > 128U) {
        return Foundation::Result<bool>::FromError(
            Foundation::Error(Foundation::ErrorCode::InvalidArgument, "密码长度无效"));
    }
    std::array<std::string_view, 6> Parts{};
    std::size_t PartIndex = 0;
    std::size_t Start = 0;
    while (PartIndex < Parts.size()) {
        const std::size_t End = EncodedHash.find('$', Start);
        Parts[PartIndex++] = EncodedHash.substr(
            Start, End == std::string_view::npos ? EncodedHash.size() - Start : End - Start);
        if (End == std::string_view::npos) {
            break;
        }
        Start = End + 1U;
    }
    if (PartIndex != Parts.size() || Parts[0] != "scrypt" || Parts[1] != "32768" ||
        Parts[2] != "8" || Parts[3] != "1") {
        return false;
    }
    const auto Salt = DecodeHex(Parts[4], SaltSize);
    const auto ExpectedHash = DecodeHex(Parts[5], HashSize);
    if (!Salt.has_value() || !ExpectedHash.has_value()) {
        return false;
    }
    const auto DerivedResult = DeriveHash(Password, *Salt);
    if (!DerivedResult.HasValue()) {
        return Foundation::Result<bool>::FromError(DerivedResult.ErrorValue());
    }
    return CRYPTO_memcmp(
               DerivedResult.Value().data(), ExpectedHash->data(), HashSize) == 0;
}

}  // namespace Menu::Infrastructure
