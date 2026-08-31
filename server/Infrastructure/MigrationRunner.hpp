#pragma once

#include <Foundation/Result.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

namespace Menu::Infrastructure {

class MigrationRunner final {
public:
    static Foundation::Result<void> Apply(SqliteDatabase& Database);
};

}  // namespace Menu::Infrastructure
