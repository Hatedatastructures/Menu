#include "Infrastructure/MigrationRunner.hpp"

#include <MigrationSql.hpp>

#include <chrono>
#include <cstdint>
#include <string>

namespace Menu::Infrastructure {
namespace {

Foundation::Error StorageError(sqlite3* Handle) {
    const char* Message = Handle == nullptr ? nullptr : sqlite3_errmsg(Handle);
    return Foundation::Error(
        Foundation::ErrorCode::StorageUnavailable,
        Message == nullptr ? "SQLite 操作失败" : std::string(Message));
}

Foundation::Result<void> RecordMigration(SqliteDatabase& Database, std::int64_t Version) {
    sqlite3_stmt* Statement = nullptr;
    constexpr char Sql[] =
        "INSERT INTO SchemaMigration (Version, AppliedAt) VALUES (?, ?);";
    const int PrepareResult = sqlite3_prepare_v2(
        Database.NativeHandle(), Sql, -1, &Statement, nullptr);
    if (PrepareResult != SQLITE_OK) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }

    const int VersionResult = sqlite3_bind_int64(Statement, 1, Version);
    const auto Timestamp = std::chrono::duration_cast<std::chrono::seconds>(
        std::chrono::system_clock::now().time_since_epoch()).count();
    const std::string TimestampText = std::to_string(Timestamp);
    const int TimestampResult = sqlite3_bind_text(
        Statement, 2, TimestampText.c_str(), -1, SQLITE_TRANSIENT);
    if (VersionResult != SQLITE_OK || TimestampResult != SQLITE_OK) {
        sqlite3_finalize(Statement);
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }

    const int StepResult = sqlite3_step(Statement);
    sqlite3_finalize(Statement);
    if (StepResult != SQLITE_DONE) {
        return Foundation::Result<void>::FromError(StorageError(Database.NativeHandle()));
    }
    return Foundation::Result<void>();
}

}  // namespace

Foundation::Result<void> MigrationRunner::Apply(SqliteDatabase& Database) {
    const auto SchemaResult = Database.Execute(
        "CREATE TABLE IF NOT EXISTS SchemaMigration ("
        "Version INTEGER PRIMARY KEY, AppliedAt TEXT NOT NULL);");
    if (!SchemaResult.HasValue()) {
        return SchemaResult;
    }

    const auto AppliedCount = Database.ScalarInt(
        "SELECT COUNT(*) FROM SchemaMigration WHERE Version = 1;");
    if (!AppliedCount.HasValue()) {
        return Foundation::Result<void>::FromError(AppliedCount.ErrorValue());
    }
    if (AppliedCount.Value() != 0) {
        return Foundation::Result<void>();
    }

    const auto BeginResult = Database.BeginTransaction();
    if (!BeginResult.HasValue()) {
        return BeginResult;
    }

    const auto MigrationResult = Database.Execute(Generated::InitialMigrationSql);
    if (!MigrationResult.HasValue()) {
        Database.Rollback();
        return MigrationResult;
    }

    const auto RecordResult = RecordMigration(Database, 1);
    if (!RecordResult.HasValue()) {
        Database.Rollback();
        return RecordResult;
    }

    const auto CommitResult = Database.Commit();
    if (!CommitResult.HasValue()) {
        Database.Rollback();
        return CommitResult;
    }
    return Foundation::Result<void>();
}

}  // namespace Menu::Infrastructure
