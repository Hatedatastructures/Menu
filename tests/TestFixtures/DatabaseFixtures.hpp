#pragma once

#include <Infrastructure/SqliteDatabase.hpp>

#include <chrono>
#include <cstdint>
#include <filesystem>
#include <stdexcept>
#include <string>
#include <string_view>
#include <utility>

namespace Menu::Tests::DatabaseFixtures {

class TemporaryDatabase final {
public:
    TemporaryDatabase(std::filesystem::path PathValue, Infrastructure::SqliteDatabase DatabaseValue)
        : Path(std::move(PathValue)), Database(std::move(DatabaseValue)) {}

    TemporaryDatabase(const TemporaryDatabase&) = delete;
    TemporaryDatabase& operator=(const TemporaryDatabase&) = delete;

    TemporaryDatabase(TemporaryDatabase&& Other) noexcept
        : Path(std::move(Other.Path)), Database(std::move(Other.Database)) {}

    TemporaryDatabase& operator=(TemporaryDatabase&& Other) noexcept {
        if (this != &Other) {
            Cleanup();
            Path = std::move(Other.Path);
            Database = std::move(Other.Database);
        }
        return *this;
    }

    ~TemporaryDatabase() {
        Cleanup();
    }

    operator Infrastructure::SqliteDatabase&() {
        return Database;
    }

    operator const Infrastructure::SqliteDatabase&() const {
        return Database;
    }

    [[nodiscard]] Foundation::Result<std::string> ScalarText(std::string_view Sql) const {
        return Database.ScalarText(Sql);
    }

    [[nodiscard]] Foundation::Result<std::int64_t> ScalarInt(std::string_view Sql) const {
        return Database.ScalarInt(Sql);
    }

private:
    void Cleanup() noexcept {
        Database.Close();
        if (!Path.empty()) {
            std::error_code Error;
            std::filesystem::remove(Path, Error);
            std::filesystem::remove(Path.string() + "-wal", Error);
            std::filesystem::remove(Path.string() + "-shm", Error);
        }
    }

    std::filesystem::path Path;
    Infrastructure::SqliteDatabase Database;
};

inline TemporaryDatabase OpenTemporary() {
    const auto Timestamp = std::chrono::steady_clock::now().time_since_epoch().count();
    const std::filesystem::path Path =
        std::filesystem::current_path() /
        ("MenuTest-" + std::to_string(Timestamp) + ".db");
    std::filesystem::create_directories(Path.parent_path());

    auto Result = Infrastructure::SqliteDatabase::Open(Path);
    if (!Result.HasValue()) {
        throw std::runtime_error(Result.ErrorValue().MessageValue());
    }
    return TemporaryDatabase(Path, std::move(Result).Value());
}

}  // namespace Menu::Tests::DatabaseFixtures
