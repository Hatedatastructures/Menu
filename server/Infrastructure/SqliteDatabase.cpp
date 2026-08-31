#include "Infrastructure/SqliteDatabase.hpp"

#include <utility>

namespace Menu::Infrastructure {
namespace {

Foundation::Error StorageError(sqlite3* Handle) {
    const char* Message = Handle == nullptr ? nullptr : sqlite3_errmsg(Handle);
    const std::string SafeMessage = Message == nullptr ? "SQLite 操作失败" : Message;
    return Foundation::Error(Foundation::ErrorCode::StorageUnavailable, SafeMessage);
}

}  // namespace

SqliteDatabase::~SqliteDatabase() {
    Close();
}

SqliteDatabase::SqliteDatabase(SqliteDatabase&& Other) noexcept
    : Handle(std::exchange(Other.Handle, nullptr)),
      DatabasePath(std::move(Other.DatabasePath)) {}

SqliteDatabase& SqliteDatabase::operator=(SqliteDatabase&& Other) noexcept {
    if (this != &Other) {
        Close();
        Handle = std::exchange(Other.Handle, nullptr);
        DatabasePath = std::move(Other.DatabasePath);
    }
    return *this;
}

Foundation::Result<SqliteDatabase> SqliteDatabase::Open(
    const std::filesystem::path& PathValue) {
    sqlite3* OpenedHandle = nullptr;
    const std::string PathText = PathValue.string();
    const int ResultCode = sqlite3_open_v2(
        PathText.c_str(),
        &OpenedHandle,
        SQLITE_OPEN_READWRITE | SQLITE_OPEN_CREATE | SQLITE_OPEN_FULLMUTEX,
        nullptr);
    if (ResultCode != SQLITE_OK) {
        const Foundation::Error ErrorValue = StorageError(OpenedHandle);
        if (OpenedHandle != nullptr) {
            sqlite3_close_v2(OpenedHandle);
        }
        return Foundation::Result<SqliteDatabase>::FromError(ErrorValue);
    }

    SqliteDatabase Database;
    Database.Handle = OpenedHandle;
    Database.DatabasePath = PathValue;
    const auto ConfigureResult = Database.Configure();
    if (!ConfigureResult.HasValue()) {
        const Foundation::Error ErrorValue = ConfigureResult.ErrorValue();
        Database.Close();
        return Foundation::Result<SqliteDatabase>::FromError(ErrorValue);
    }
    return Database;
}

Foundation::Result<void> SqliteDatabase::Configure() const {
    if (!IsReady()) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 未打开"));
    }
    const auto JournalResult = Execute("PRAGMA journal_mode = WAL;");
    if (!JournalResult.HasValue()) {
        return JournalResult;
    }
    const auto ForeignKeyResult = Execute("PRAGMA foreign_keys = ON;");
    if (!ForeignKeyResult.HasValue()) {
        return ForeignKeyResult;
    }
    const auto SynchronousResult = Execute("PRAGMA synchronous = NORMAL;");
    if (!SynchronousResult.HasValue()) {
        return SynchronousResult;
    }
    if (sqlite3_busy_timeout(Handle, 5000) != SQLITE_OK) {
        return Foundation::Result<void>::FromError(StorageError(Handle));
    }
    return Foundation::Result<void>();
}

Foundation::Result<void> SqliteDatabase::Execute(std::string_view Sql) const {
    if (!IsReady()) {
        return Foundation::Result<void>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 未打开"));
    }
    const std::string SqlText(Sql);
    char* ErrorMessage = nullptr;
    const int ResultCode = sqlite3_exec(Handle, SqlText.c_str(), nullptr, nullptr, &ErrorMessage);
    if (ResultCode != SQLITE_OK) {
        if (ErrorMessage != nullptr) {
            sqlite3_free(ErrorMessage);
        }
        return Foundation::Result<void>::FromError(StorageError(Handle));
    }
    return Foundation::Result<void>();
}

Foundation::Result<std::string> SqliteDatabase::ScalarText(std::string_view Sql) const {
    if (!IsReady()) {
        return Foundation::Result<std::string>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 未打开"));
    }
    sqlite3_stmt* Statement = nullptr;
    const std::string SqlText(Sql);
    const int PrepareResult = sqlite3_prepare_v2(
        Handle, SqlText.c_str(), static_cast<int>(SqlText.size()), &Statement, nullptr);
    if (PrepareResult != SQLITE_OK) {
        return Foundation::Result<std::string>::FromError(StorageError(Handle));
    }

    const int StepResult = sqlite3_step(Statement);
    if (StepResult != SQLITE_ROW) {
        sqlite3_finalize(Statement);
        return Foundation::Result<std::string>::FromError(StorageError(Handle));
    }
    const unsigned char* TextValue = sqlite3_column_text(Statement, 0);
    const std::string Value = TextValue == nullptr
                                  ? std::string()
                                  : reinterpret_cast<const char*>(TextValue);
    sqlite3_finalize(Statement);
    return Value;
}

Foundation::Result<std::int64_t> SqliteDatabase::ScalarInt(std::string_view Sql) const {
    if (!IsReady()) {
        return Foundation::Result<std::int64_t>::FromError(
            Foundation::Error(Foundation::ErrorCode::StorageUnavailable, "SQLite 未打开"));
    }
    sqlite3_stmt* Statement = nullptr;
    const std::string SqlText(Sql);
    const int PrepareResult = sqlite3_prepare_v2(
        Handle, SqlText.c_str(), static_cast<int>(SqlText.size()), &Statement, nullptr);
    if (PrepareResult != SQLITE_OK) {
        return Foundation::Result<std::int64_t>::FromError(StorageError(Handle));
    }

    const int StepResult = sqlite3_step(Statement);
    if (StepResult != SQLITE_ROW) {
        sqlite3_finalize(Statement);
        return Foundation::Result<std::int64_t>::FromError(StorageError(Handle));
    }
    const std::int64_t Value = sqlite3_column_int64(Statement, 0);
    sqlite3_finalize(Statement);
    return Value;
}

Foundation::Result<void> SqliteDatabase::BeginTransaction() const {
    return Execute("BEGIN IMMEDIATE TRANSACTION;");
}

Foundation::Result<void> SqliteDatabase::Commit() const {
    return Execute("COMMIT;");
}

Foundation::Result<void> SqliteDatabase::Rollback() const {
    return Execute("ROLLBACK;");
}

bool SqliteDatabase::IsReady() const noexcept {
    return Handle != nullptr;
}

const std::filesystem::path& SqliteDatabase::Path() const noexcept {
    return DatabasePath;
}

sqlite3* SqliteDatabase::NativeHandle() const noexcept {
    return Handle;
}

void SqliteDatabase::Close() noexcept {
    if (Handle != nullptr) {
        sqlite3_close_v2(Handle);
        Handle = nullptr;
    }
    DatabasePath.clear();
}

}  // namespace Menu::Infrastructure
