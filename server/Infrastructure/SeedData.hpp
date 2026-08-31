#pragma once

#include <Foundation/Result.hpp>
#include <Infrastructure/SqliteDatabase.hpp>

namespace Menu::Infrastructure {

class SeedData final {
public:
    static Foundation::Result<void> InsertIfEmpty(SqliteDatabase& Database);
};

}  // namespace Menu::Infrastructure
