#pragma once

#include <Foundation/Result.hpp>

#include <sqlite3.h>

#include <cstdint>
#include <filesystem>
#include <string>
#include <string_view>

namespace Menu::Infrastructure {

class SqliteDatabase final {
public:
    SqliteDatabase() = default;
    ~SqliteDatabase();

    SqliteDatabase(const SqliteDatabase&) = delete;
    SqliteDatabase& operator=(const SqliteDatabase&) = delete;

    SqliteDatabase(SqliteDatabase&& Other) noexcept;
    SqliteDatabase& operator=(SqliteDatabase&& Other) noexcept;

    static Foundation::Result<SqliteDatabase> Open(const std::filesystem::path& Path);

    Foundation::Result<void> Execute(std::string_view Sql) const;
    Foundation::Result<std::string> ScalarText(std::string_view Sql) const;
    Foundation::Result<std::int64_t> ScalarInt(std::string_view Sql) const;
    Foundation::Result<void> BeginTransaction() const;
    Foundation::Result<void> Commit() const;
    Foundation::Result<void> Rollback() const;

    [[nodiscard]] bool IsReady() const noexcept;
    [[nodiscard]] const std::filesystem::path& Path() const noexcept;
    [[nodiscard]] sqlite3* NativeHandle() const noexcept;

    void Close() noexcept;

private:
    Foundation::Result<void> Configure() const;

    sqlite3* Handle = nullptr;
    std::filesystem::path DatabasePath;
};

}  // namespace Menu::Infrastructure
