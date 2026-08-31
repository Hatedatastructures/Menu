#include "Infrastructure/TokenService.hpp"

#include <openssl/evp.h>
#include <openssl/rand.h>

#include <array>
#include <chrono>
#include <string>
#include <utility>
#include <vector>

namespace Menu::Infrastructure {
namespace {

constexpr char HexDigits[] = "0123456789abcdef";

Foundation::Error CryptoError(std::string Message) {
    return Foundation::Error(Foundation::ErrorCode::StorageUnavailable, std::move(Message));
}

std::string EncodeHex(const unsigned char* Bytes, std::size_t Size) {
    std::string Encoded;
    Encoded.reserve(Size * 2U);
    for (std::size_t Index = 0; Index < Size; ++Index) {
        Encoded.push_back(HexDigits[(Bytes[Index] >> 4U) & 0x0FU]);
        Encoded.push_back(HexDigits[Bytes[Index] & 0x0FU]);
    }
    return Encoded;
}

}  // namespace

Foundation::Result<std::string> TokenService::CreateToken(std::size_t ByteCount) {
    if (ByteCount == 0U || ByteCount > 128U) {
        return Foundation::Result<std::string>::FromError(
            CryptoError("Token 长度无效"));
    }
    std::vector<unsigned char> Bytes(ByteCount);
    if (RAND_bytes(Bytes.data(), static_cast<int>(Bytes.size())) != 1) {
        return Foundation::Result<std::string>::FromError(
            CryptoError("Token 随机数生成失败"));
    }
    return EncodeHex(Bytes.data(), Bytes.size());
}

Foundation::Result<std::string> TokenService::HashToken(std::string_view Token) {
    std::array<unsigned char, EVP_MAX_MD_SIZE> Digest{};
    unsigned int DigestLength = 0;
    if (EVP_Digest(
            Token.data(), Token.size(), Digest.data(), &DigestLength,
            EVP_sha256(), nullptr) != 1) {
        return Foundation::Result<std::string>::FromError(
            CryptoError("Token 摘要计算失败"));
    }
    return EncodeHex(Digest.data(), DigestLength);
}

std::int64_t TokenService::CurrentUnixSeconds() noexcept {
    return std::chrono::duration_cast<std::chrono::seconds>(
               std::chrono::system_clock::now().time_since_epoch())
        .count();
}

}  // namespace Menu::Infrastructure
