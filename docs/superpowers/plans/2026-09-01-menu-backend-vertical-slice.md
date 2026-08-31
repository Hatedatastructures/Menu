# Menu 后端第一条垂直切片实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 构建一个可 fresh configure、可运行、可通过真实 HTTP 验证的 C++20/Boost/SQLite 菜谱服务端，并交付确定性“今晚”推荐接口。

**Architecture:** `MenuFoundation` 提供无业务的错误、Result 和 ID；`MenuDomain` 保持纯规则；`MenuApplication` 依赖仓储和推荐端口；`MenuInfrastructure` 在受控 storage executor 中访问 SQLite；`MenuTransportCore` 只提供中立 HTTP 类型和 handler 注入协议；`MenuTransport` 实现 listener；`MenuApi` 只链接 Core 并实现路由和用例适配；`MenuServer` 是唯一 composition root。

**Tech Stack:** C++20, CMake 3.23+, Boost 1.89.0 (Asio/Beast/JSON/System), SQLite 3.53.4 amalgamation, OpenSSL scrypt (认证切片预留), GoogleTest, CTest, PowerShell HTTP smoke。

**Spec:** `docs/superpowers/specs/2026-09-01-menu-backend-vertical-slice-design.md`

## Global Constraints

- 所有 C++ 标识符和文件名采用 PascalCase/大驼峰风格；类、函数、变量、成员和测试名不得使用 snake_case。
- C++ 使用 C++20；服务端网络 I/O 使用 Boost.Asio/Beast 异步模型，不在 event loop 中执行 SQLite、文件或密码哈希等阻塞工作。
- `MenuTransportCore` 和 `MenuTransport` 不得链接 `MenuApi`；`MenuApi` 只能链接 `MenuTransportCore`，不能链接具体 listener；`MenuServer` 才能把 Api handler 注入 listener；configure 必须验证这个 link closure 方向。
- 生产 target 只能依赖生产模块；测试支持不能反向进入 `MenuFoundation`、`MenuDomain`、`MenuApplication`、`MenuInfrastructure`、`MenuTransport` 或 `MenuApi`。
- SQLite 写入使用参数化 SQL 和事务，启动执行版本化迁移，启用 WAL、外键和 busy timeout。
- 依赖、构建、下载缓存和测试输出放在 `I:\code\Menu` 下；下载失败时只对当前进程使用 `http://127.0.0.1:7890` 代理。
- 每个任务完成后运行 `git diff --check`，构建使用 `--parallel 2`，不把旧 build 日志当作当前证据。
- 第一切片只提供已发布 Recipe/Ingredient 读取和确定性推荐；不接入商超、AI、聊天或外部数据库。

---

### Task 1: Bootstrap CMake and Test Harness

**Files:**
- Create: `CMakeLists.txt`
- Create: `CMakePresets.json`
- Create: `cmake/Dependencies.cmake`
- Create: `cmake/TargetBoundaries.cmake`
- Create: `server/CMakeLists.txt`
- Create: `server/Foundation/CMakeLists.txt`
- Create: `server/Domain/CMakeLists.txt`
- Create: `server/Application/CMakeLists.txt`
- Create: `server/Infrastructure/CMakeLists.txt`
- Create: `server/Transport/CMakeLists.txt`
- Create: `server/Api/CMakeLists.txt`
- Create: `server/Main.cpp`
- Create: `tests/CMakeLists.txt`
- Create: `tests/unit/CMakeLists.txt`
- Create: `tests/unit/BootstrapContractTest.cpp`
- Create: `tests/TestMain.cpp`
- Create: `shared/schema/MenuApi.yaml`

**Interfaces:**
- Produces targets `MenuFoundation`, `MenuDomain`, `MenuApplication`, `MenuInfrastructure`, `MenuTransportCore`, `MenuTransport`, `MenuApi`, `MenuServer`, and `MenuUnitTests`.
- Produces presets `WindowsDebug`, `WindowsRelease`, and `WindowsAsan` with `compile_commands.json` and `I:\code\Menu\build\<preset>` binary directories.
- `TargetBoundaries.cmake` exposes `AssertMenuTargetBoundaries()` and fails configure if `MenuTransportCore`/`MenuTransport` closure contains `MenuApi`, if `MenuApi` closure contains concrete `MenuTransport`, or if any target has a direct dependency on itself.

- [x] **Step 1: Write the configure contract test**

Create the first CTest registration and a test executable that includes the future Foundation header. The test is intentionally not implementable yet; this makes the first behavior RED after configure.

```cpp
#include <gtest/gtest.h>
#include <Foundation/Result.hpp>

TEST(BootstrapContractTest, BuildsWithCxx20AndFoundationTarget) {
    const Menu::Foundation::Result<int> Result = 7;
    ASSERT_TRUE(Result.HasValue());
    EXPECT_EQ(Result.Value(), 7);
}
```

- [x] **Step 2: Run configure and record the expected RED**

Run:

```powershell
cmake --preset WindowsDebug
cmake --build --preset WindowsDebug --parallel 2
```

Expected: configure can resolve the declared dependencies, then compilation fails because `Menu/Foundation/Result.hpp` does not exist. Do not add a production header in this task to make the test green.

- [x] **Step 3: Add minimal build wiring and the dependency-direction assertion**

`CMakePresets.json` must use the existing compiler without installing into C drive:

```json
{
  "version": 6,
  "configurePresets": [
    {
      "name": "WindowsDebug",
      "generator": "MinGW Makefiles",
      "binaryDir": "${sourceDir}/build/WindowsDebug",
      "cacheVariables": {
        "CMAKE_BUILD_TYPE": "Debug",
        "CMAKE_CXX_STANDARD": "20",
        "CMAKE_CXX_STANDARD_REQUIRED": "ON",
        "CMAKE_EXPORT_COMPILE_COMMANDS": "ON",
        "CMAKE_C_COMPILER": "C:/msys64/ucrt64/bin/gcc.exe",
        "CMAKE_CXX_COMPILER": "C:/msys64/ucrt64/bin/g++.exe",
        "MENU_ENABLE_TESTS": "ON",
        "FETCHCONTENT_BASE_DIR": "${sourceDir}/.cache/cmake-fetch"
      }
    }
  ],
  "buildPresets": [
    {"name": "WindowsDebug", "configurePreset": "WindowsDebug", "jobs": 2}
  ],
  "testPresets": [
    {"name": "WindowsDebug", "configurePreset": "WindowsDebug", "output": {"outputOnFailure": true}}
  ]
}
```

`cmake/Dependencies.cmake` declares Boost 1.89.0 and GoogleTest with official URLs and SHA256, creates a header-only `MenuBoost` interface exposing Boost.Asio/Beast/System, and builds `MenuBoostJson` from Boost.JSON's supported separate-compilation source so JSON symbols are linked explicitly. `cmake/TargetBoundaries.cmake` traverses `LINK_LIBRARIES` and `INTERFACE_LINK_LIBRARIES` for in-project targets and emits the closure in the configure error/status message.

The transport boundary is explicit: `MenuTransportCore` contains `HttpRequest`, `HttpResponse`, and `HttpHandler`; `MenuTransport` contains only the Beast listener/session implementation; `MenuApi` links Core and implements a handler; the listener accepts the injected handler without including any Api header. Neither Transport target may link `MenuApi`.

- [x] **Step 4: Verify the RED build and configuration properties**

Run:

```powershell
cmake --build --preset WindowsDebug --parallel 2
```

Expected: non-zero exit with the missing `Menu/Foundation/Result.hpp` include. Then inspect `build\WindowsDebug\CMakeCache.txt` and confirm `CMAKE_CXX_STANDARD=20`, `FETCHCONTENT_BASE_DIR` starts with `I:\code\Menu`, and no dependency source is under `C:\Users`.

- [x] **Step 5: Commit the bootstrap**

```powershell
git add CMakeLists.txt CMakePresets.json cmake server tests shared
git commit -m "build: bootstrap Menu CMake targets and tests"
```

### Task 2: Implement Foundation and Domain Rules

**Files:**
- Create: `server/Foundation/Result.hpp`
- Create: `server/Foundation/Error.hpp`
- Create: `server/Foundation/Identifier.hpp`
- Create: `server/Domain/Ingredient.hpp`
- Create: `server/Domain/Recipe.hpp`
- Create: `server/Domain/Recommendation.hpp`
- Create: `server/Domain/RecipeRules.hpp`
- Create: `server/Domain/RecipeRules.cpp`
- Create: `tests/TestFixtures/RecipeFixtures.hpp`
- Create: `tests/unit/RecipeRulesTest.cpp`
- Create: `tests/unit/RecommendationRulesTest.cpp`

**Interfaces:**
- `Menu::Foundation::Result<T>` exposes `HasValue()`, `Value()`, `ErrorValue()` and `FromError(Error)`.
- `Menu::Domain::RecipeRules::ScaleRecipeIngredients(const Recipe&, int)` returns `Result<std::vector<RecipeIngredient>>`.
- `Menu::Domain::RecipeRules::BuildIngredientAvailability(const Recipe&, const std::vector<std::string>&)` returns `IngredientAvailability` with `AvailableCount`, `MissingCount`, and item IDs.
- `Menu::Domain::RecipeRules::ValidateRecipe(const Recipe&)` returns `Result<void>`.
- `Menu::Domain::RuleBasedRecommendation::Rank(const RecommendationRequest&, const std::vector<Recipe>&)` returns at most three `Recommendation` values.
- `Menu::Tests::RecipeFixtures::SingleRecipeWithIngredients(int, std::initializer_list<RecipeIngredient>)` returns a deterministic `Recipe` for unit tests.

- [x] **Step 1: Write the serving-scale RED test**

```cpp
TEST(RecipeRulesTest, ScalesRequiredAndOptionalIngredientsByServingRatio) {
    const Menu::Domain::Recipe Recipe = Menu::Tests::RecipeFixtures::SingleRecipeWithIngredients(
        2, {{"ingredient.rice", 200.0, "g", true}, {"ingredient.salt", 2.0, "g", false}});

    const auto Result = Menu::Domain::RecipeRules::ScaleRecipeIngredients(Recipe, 4);

    ASSERT_TRUE(Result.HasValue());
    ASSERT_EQ(Result.Value().size(), 2U);
    EXPECT_DOUBLE_EQ(Result.Value()[0].Quantity, 400.0);
    EXPECT_TRUE(Result.Value()[0].Required);
    EXPECT_DOUBLE_EQ(Result.Value()[1].Quantity, 4.0);
    EXPECT_FALSE(Result.Value()[1].Required);
}
```

- [x] **Step 2: Run the focused test and verify it fails for the missing domain API**

Run:

```powershell
ctest --preset WindowsDebug -R RecipeRulesTest --output-on-failure
```

Expected: build/configuration failure naming the missing `RecipeRules` symbols, not a passing test or an unrelated compiler error.

- [x] **Step 3: Write RED tests for validation, aliases, availability, and recommendation filtering**

Add one test per behavior: zero servings is rejected; duplicate step order is rejected; a registered alias normalizes to one Ingredient ID; pantry IDs affect missing count; an allergy hard-filters a recipe even when cuisine score is high; three recommendations are the maximum.

- [x] **Step 4: Implement the smallest pure domain model**

Use aggregate structs with explicit fields and no database/network includes:

```cpp
struct RecipeIngredient {
    std::string IngredientId;
    double Quantity = 0.0;
    std::string Unit;
    double ServingFactor = 1.0;
    std::string Preparation;
    bool Required = true;
};

struct RecipeStep {
    int StepOrder = 0;
    std::string Title;
    std::string Instruction;
    int DurationSeconds = 0;
    bool HasTimer = false;
};

struct Recipe {
    std::string Id;
    std::string Slug;
    std::string Name;
    std::string Cuisine;
    std::string Description;
    int PrepMinutes = 0;
    int CookMinutes = 0;
    int Servings = 1;
    int Difficulty = 1;
    std::string ImagePath;
    std::string Status;
    std::vector<RecipeIngredient> Ingredients;
    std::vector<RecipeStep> Steps;
    std::vector<std::string> Allergens;
    std::vector<std::string> Cookware;
};
```

Keep normalization deterministic: trim ASCII whitespace, apply Unicode Chinese aliases from the supplied alias map, and never silently change a quantity or unit.

- [x] **Step 5: Run all domain tests and refactor only after GREEN**

```powershell
cmake --build --preset WindowsDebug --parallel 2
ctest --preset WindowsDebug -R "RecipeRulesTest|RecommendationRulesTest" --output-on-failure
```

Expected: every focused test passes with no warning introduced by Menu code.

- [x] **Step 6: Commit the domain slice**

```powershell
git add server/Foundation server/Domain tests/unit
git commit -m "feat: add recipe domain rules and recommendations"
```

### Task 3: Add SQLite Migrations and Recipe Repository

**Files:**
- Create: `server/Infrastructure/SqliteDatabase.hpp`
- Create: `server/Infrastructure/SqliteDatabase.cpp`
- Create: `server/Infrastructure/MigrationRunner.hpp`
- Create: `server/Infrastructure/MigrationRunner.cpp`
- Create: `server/Infrastructure/SqliteRecipeRepository.hpp`
- Create: `server/Infrastructure/SqliteRecipeRepository.cpp`
- Create: `server/Infrastructure/SeedData.hpp`
- Create: `server/Infrastructure/SeedData.cpp`
- Create: `server/Infrastructure/Migrations/001Initial.sql`
- Create: `tests/TestFixtures/DatabaseFixtures.hpp`
- Create: `tests/unit/SqliteMigrationTest.cpp`
- Create: `tests/unit/SqliteRecipeRepositoryTest.cpp`
- Modify: `cmake/Dependencies.cmake`
- Modify: `server/Infrastructure/CMakeLists.txt`

**Interfaces:**
- `SqliteDatabase::Open(const std::filesystem::path&)`, `Execute(std::string_view)`, `BeginTransaction()`, `Commit()`, `Rollback()`, and `IsReady()`.
- `MigrationRunner::Apply(SqliteDatabase&)` is idempotent and records versions in `SchemaMigration`.
- `SqliteRecipeRepository` implements `RecipeRepository::ListPublished()` and `FindPublishedById(std::string_view)`.
- `SeedData::InsertIfEmpty(SqliteDatabase&)` inserts deterministic recipes/ingredients only when the relevant tables are empty.
- `Menu::Tests::DatabaseFixtures::OpenTemporary()` returns an opened temporary `SqliteDatabase` whose file is deleted by the fixture destructor.

- [x] **Step 1: Write the migration RED test**

```cpp
TEST(SqliteMigrationTest, AppliesSchemaAndEnablesWalAndForeignKeys) {
    auto Database = Menu::Tests::DatabaseFixtures::OpenTemporary();
    ASSERT_TRUE(Menu::Infrastructure::MigrationRunner::Apply(Database).HasValue());

    EXPECT_EQ(Database.ScalarText("PRAGMA journal_mode;"), "wal");
    EXPECT_EQ(Database.ScalarInt("PRAGMA foreign_keys;"), 1);
    EXPECT_EQ(Database.ScalarInt("SELECT COUNT(*) FROM SchemaMigration;"), 1);
}
```

- [x] **Step 2: Run it and verify the failure is caused by missing storage code**

```powershell
ctest --preset WindowsDebug -R SqliteMigrationTest --output-on-failure
```

Expected: compile failure for missing `SqliteDatabase`/`MigrationRunner`, before any schema assertion runs.

- [x] **Step 3: Add the schema and SQLite dependency target**

`001Initial.sql` must create `SchemaMigration`, `Users`, `Preferences`, `Ingredients`, `IngredientAliases`, `Recipes`, `RecipeIngredients`, `RecipeSteps`, `MealPlans`, `MealPlanItems`, `Reminders`, `CookingSessions`, `Feedback`, `MediaAssets`, and `RefreshTokens`. Use `RecipeId` plus `(RecipeId, StepOrder)` uniqueness, foreign keys with explicit delete behavior, and `CHECK` constraints for positive quantities, non-negative durations, valid status, and supported units.

Fetch the SQLite 3.51.1 amalgamation into `I:\code\Menu\.cache\cmake-fetch`, compute SHA256, and append URL/version/path/size/hash to `docs/Downloads.md` before committing.

- [x] **Step 4: Implement connection setup and migration transaction**

Open SQLite with `SQLITE_OPEN_READWRITE | SQLITE_OPEN_CREATE | SQLITE_OPEN_FULLMUTEX`, execute `PRAGMA journal_mode=WAL`, `PRAGMA foreign_keys=ON`, and `PRAGMA busy_timeout=5000`, then run each migration in one transaction. Bind every value through `sqlite3_bind_*`; only fixed migration SQL may be passed as text.

- [x] **Step 5: Write repository RED tests**

Cover: seed data is returned only when `Status='published'`; unknown ID returns an empty optional; a recipe with two ingredients and two ordered steps round-trips; a failed transaction leaves no partial recipe; a second migration run does not add a version row.

- [x] **Step 6: Implement repository and seed data**

Use prepared statements with explicit column lists. Seed at least twelve representative recipes in this slice across Chinese, Western, and Japanese cuisines, with local image paths that will later be backed by `assets/media`; scale content to 60-100 recipes before product release.

- [x] **Step 7: Verify storage tests and database diagnostics**

```powershell
cmake --build --preset WindowsDebug --parallel 2
ctest --preset WindowsDebug -R "SqliteMigrationTest|SqliteRecipeRepositoryTest" --output-on-failure
```

Expected: all storage tests pass; test output reports the temporary database uses WAL and no temporary file is left in the repository.

- [x] **Step 8: Commit the storage slice**

```powershell
git add cmake server/Infrastructure tests/unit docs/Downloads.md
git commit -m "feat: add sqlite migrations and recipe repository"
```

### Task 4: Add Application Service and Storage Executor

**Files:**
- Create: `server/Application/RecipeRepository.hpp`
- Create: `server/Application/RecommendationProvider.hpp`
- Create: `server/Application/RecipeApplicationService.hpp`
- Create: `server/Application/RecipeApplicationService.cpp`
- Create: `server/Application/RuleBasedRecommendationProvider.hpp`
- Create: `server/Application/RuleBasedRecommendationProvider.cpp`
- Create: `server/Application/StorageExecutor.hpp`
- Create: `server/Application/StorageExecutor.cpp`
- Create: `tests/TestFixtures/ApplicationFixtures.hpp`
- Create: `tests/unit/RecipeApplicationServiceTest.cpp`
- Create: `tests/unit/StorageExecutorTest.cpp`

**Interfaces:**
- `RecipeApplicationService::ListPublishedRecipes()` and `FindPublishedRecipe(std::string_view)` are synchronous repository calls intended only for the storage executor.
- `RecipeApplicationService::RecommendTonight(const RecommendationRequest&)` calls the provider with repository data and returns at most three recommendations.
- `StorageExecutor::Submit(Work, Completion)` runs work on an Asio executor and posts completion to the caller executor; it owns shutdown and does not detach tasks that borrow request state.
- `Menu::Tests::ApplicationFixtures::WithSeedRecipes()` returns an application service backed by deterministic in-memory test repositories.

- [x] **Step 1: Write the application RED test**

```cpp
TEST(RecipeApplicationServiceTest, ReturnsOnlyThreeRecipesThatFitUserConstraints) {
    const auto Service = Menu::Tests::ApplicationFixtures::WithSeedRecipes();
    const Menu::Domain::RecommendationRequest Request{
        .Servings = 2,
        .AvailableMinutes = 35,
        .Cuisines = {"中餐"},
        .Allergies = {"花生"},
        .PantryIngredientIds = {"ingredient.rice"},
        .Cookware = {"炒锅"}
    };

    const auto Result = Service.RecommendTonight(Request);

    ASSERT_TRUE(Result.HasValue());
    ASSERT_LE(Result.Value().size(), 3U);
    for (const auto& Recommendation : Result.Value()) {
        EXPECT_EQ(Recommendation.RecipeValue.Cuisine, "中餐");
        EXPECT_EQ(Recommendation.RecipeValue.Allergens.end(),
                  std::find(Recommendation.RecipeValue.Allergens.begin(),
                            Recommendation.RecipeValue.Allergens.end(), "花生"));
    }
}
```

- [x] **Step 2: Run the focused test and verify the expected missing-service RED**

```powershell
ctest --preset WindowsDebug -R RecipeApplicationServiceTest --output-on-failure
```

Expected: failure for missing Application service symbols, not a test that passes against hard-coded fixtures.

- [x] **Step 3: Implement rule provider and service**

Use a fixed scoring tuple `(CuisineMatch, PantryCoverage, RecentPenalty, Difficulty, RecipeId)` so ordering is stable across platforms. Apply allergies, max minutes, cookware, status and serving validity before scoring. The provider must not access SQLite or Asio.

- [x] **Step 4: Write executor ownership test**

Submit 20 tasks from two threads, record completion IDs, assert all 20 complete and no completion runs after `Shutdown()`. Use a value-owned input, not a reference to a stack request, to make the ownership rule observable.

- [x] **Step 5: Implement bounded storage executor**

Use `boost::asio::thread_pool` with one writer strand and a bounded read pool. Each work item stores its input by value and captures a `std::shared_ptr` state for completion; no detached coroutine captures `this` or a borrowed request. Expose `Join()` for deterministic tests and server shutdown.

- [x] **Step 6: Verify application tests and commit**

```powershell
cmake --build --preset WindowsDebug --parallel 2
ctest --preset WindowsDebug -R "RecipeApplicationServiceTest|StorageExecutorTest" --output-on-failure
git diff --check
git add server/Application tests/unit
git commit -m "feat: add recipe application service and storage executor"
```

### Task 5: Implement Asio/Beast HTTP Transport

**Files:**
- Create: `server/Transport/Core/HttpRequest.hpp`
- Create: `server/Transport/Core/HttpResponse.hpp`
- Create: `server/Transport/Core/HttpHandler.hpp`
- Create: `server/Transport/HttpServer.hpp`
- Create: `server/Transport/HttpServer.cpp`
- Create: `server/Transport/HttpSession.hpp`
- Create: `server/Transport/HttpSession.cpp`
- Create: `tests/integration/HttpTransportTest.cpp`
- Create: `tests/integration/CMakeLists.txt`

**Interfaces:**
- `HttpServer::Start(HttpServerOptions, HttpHandler)` asynchronously accepts TCP connections and exposes `LocalPort()`.
- `HttpRequest` and `HttpResponse` are transport-neutral value types in `MenuTransportCore`; `HttpHandler` is an injected callable and contains no `MenuApi` include.
- `HttpSession` owns a Beast request buffer, response queue, per-connection strand, read timeout and write serialization through `std::shared_ptr` lifetime.

- [ ] **Step 1: Write the transport RED test**

Start a real `HttpServer` on port 0 with a handler returning `200 text/plain`, then issue an Asio TCP request and assert the response body. The test should fail to compile until the transport API exists.

- [ ] **Step 2: Run the test and verify it fails for missing transport symbols**

```powershell
ctest --preset WindowsDebug -R HttpTransportTest --output-on-failure
```

Expected: compile failure naming `HttpServer` or `HttpSession`.

- [ ] **Step 3: Implement accept/read/dispatch/write/close**

Use `net::ip::tcp::acceptor`, `beast::http::async_read`, `async_write`, and `net::make_strand`. Enforce a 1 MiB body limit, 8 KiB target limit, 10 second read/write timeout, HTTP/1.1 keep-alive, and explicit `Connection: close` on malformed input. Every asynchronous callback owns the session with `shared_from_this()`.

`MenuTransport` must only invoke the injected `HttpHandler`; it must never include, link, or call `MenuApi`. Add `AssertMenuTargetBoundaries()` to configure and verify that `MenuTransportCore`/`MenuTransport` closures exclude `MenuApi`, that `MenuApi` links Core but not concrete `MenuTransport`, and that `MenuServer` links both.

- [ ] **Step 4: Add malformed and oversized request tests**

Verify 400 for invalid parser input, 413 for a body over 1 MiB, connection close after parser error, and two pipelined requests are serialized rather than concurrently writing the socket.

- [ ] **Step 5: Verify transport integration and commit**

```powershell
cmake --build --preset WindowsDebug --parallel 2
ctest --preset WindowsDebug -R HttpTransportTest --output-on-failure
git diff --check
git add server/Transport tests/integration cmake/TargetBoundaries.cmake
git commit -m "feat: add asynchronous neutral http transport"
```

### Task 6: Add API Routing and Real Server Entry Point

**Files:**
- Create: `server/Api/ApiRouter.hpp`
- Create: `server/Api/ApiRouter.cpp`
- Create: `server/Api/ApiErrors.hpp`
- Create: `server/Api/RecipeDtos.hpp`
- Create: `server/Api/RecipeDtos.cpp`
- Create: `server/Main.cpp`
- Create: `tests/integration/MenuApiTest.cpp`
- Create: `config/Menu.example.json`
- Modify: `server/CMakeLists.txt`
- Modify: `server/Api/CMakeLists.txt`

**Interfaces:**
- `ApiRouter::Handle(HttpRequest, HttpResponseCallback)` maps `/healthz`, `/readyz`, `/api/v1/recipes`, `/api/v1/recipes/{id}`, and `/api/v1/recommendations/tonight`.
- `RecipeDtos::ToJson(const Recipe&)` and `ToJson(const Recommendation&)` use Boost.JSON and fixed external keys.
- `MenuServer` parses config, opens SQLite, applies migrations, seeds data, starts transport and waits for SIGINT/SIGTERM or Windows console close.

- [ ] **Step 1: Write API integration RED tests against a real port**

The test process must start `MenuServer` or the same composition root in-process, then use an external-style Beast HTTP client. Cover:

```text
GET /healthz                         -> 200 {"status":"ok"}
GET /readyz                          -> 200 after migration
GET /api/v1/recipes                  -> 200, only published recipes
GET /api/v1/recipes/{knownId}        -> 200, ingredients and steps
GET /api/v1/recipes/does-not-exist   -> 404, error.code=recipe_not_found
GET /api/v1/recipes?limit=101        -> 400, error.code=invalid_query
POST /api/v1/recommendations/tonight -> 200, <=3 recommendations
```

- [ ] **Step 2: Run integration tests and verify RED**

```powershell
ctest --preset WindowsDebug -R MenuApiTest --output-on-failure
```

Expected: compile or connection failure because the API router/composition root is not implemented.

- [ ] **Step 3: Implement DTOs and route validation**

Parse query values with bounded integer conversion; reject empty/overlong IDs, invalid cuisine values, `limit < 1`, `limit > 100`, missing recommendation body, non-positive servings and negative available minutes. Responses always include `requestId` and set `application/json; charset=utf-8`.

- [ ] **Step 4: Implement asynchronous application dispatch**

`ApiRouter` must copy request inputs into a worker task, call the synchronous Application service only on `StorageExecutor`, and post the serialized response back through the session executor. Do not call SQLite from the read handler. Handle worker failure as `503 storage_unavailable` without exposing SQL text or file paths.

- [ ] **Step 5: Implement server startup and graceful shutdown**

Read `config/Menu.example.json` defaults: `127.0.0.1:8080`, `server/data/menu.db`, `assets/media`, body limit 1 MiB, CORS allowlist empty. Create directories, migrate, seed, start io_context threads sized to `max(2, hardware_concurrency/2)` and join all executors on exit.

`MenuServer` is the only target allowed to compose `MenuApi` and `MenuTransport`: `MenuApi` implements the injected handler and links only `MenuTransportCore`, `MenuTransport` owns the listener, and neither target may include the other's concrete implementation headers.

- [ ] **Step 6: Verify the real process with curl and commit**

```powershell
cmake --build --preset WindowsDebug --parallel 2
ctest --preset WindowsDebug --output-on-failure
Start-Process -FilePath .\build\WindowsDebug\MenuServer.exe -ArgumentList '--config', '.\config\Menu.example.json' -PassThru
curl.exe --fail http://127.0.0.1:8080/healthz
curl.exe --fail http://127.0.0.1:8080/readyz
curl.exe --fail http://127.0.0.1:8080/api/v1/recipes
git diff --check
```

Record the actual PID, status codes, JSON bodies, database path, and stop only that PID. Then commit:

```powershell
git add server config tests/integration
git commit -m "feat: expose recipe and tonight recommendation api"
```

### Task 7: Contract, Security Baseline, and First Release Evidence

**Files:**
- Modify: `shared/schema/MenuApi.yaml`
- Create: `shared/schema/Recipe.schema.json`
- Create: `tests/api/ApiContractTest.cpp`
- Create: `tests/api/CMakeLists.txt`
- Create: `scripts/RunApiSmoke.ps1`
- Modify: `docs/Architecture.md`
- Modify: `docs/Research.md`
- Modify: `docs/Runbook.md`
- Modify: `docs/Downloads.md`

**Interfaces:**
- OpenAPI documents every first-slice endpoint, response/error shape, size limit, and CORS/auth reservation.
- `RunApiSmoke.ps1` starts a named process, sends real HTTP calls, emits machine-readable status, and stops only the recorded process ID.

- [ ] **Step 1: Write contract RED tests**

Validate required JSON keys, `Content-Type`, error code/message/requestId, recipe ingredient quantity/unit, and a maximum of three recommendation entries against `Recipe.schema.json`. Include 401/403 cases as explicit pending security cases until the authentication slice is added; never silently omit the future security contract.

- [ ] **Step 2: Implement schema and security baseline checks**

Add route-level CORS allowlist, request body/target limits, log redaction for authorization headers, and a documented TLS configuration path. Keep access/refresh token implementation out of this slice but reserve `/api/v1/auth/*` under the same error model.

- [ ] **Step 3: Run full first-slice verification**

```powershell
git diff --check
cmake --fresh --preset WindowsDebug
cmake --build --preset WindowsDebug --parallel 2
ctest --preset WindowsDebug --output-on-failure
.\scripts\RunApiSmoke.ps1 -Executable .\build\WindowsDebug\MenuServer.exe -Config .\config\Menu.example.json
```

Expected: fresh configure, build, CTest, and external HTTP smoke all return zero; record exact test count and HTTP status output in the next progress report. If any step fails, add a regression test before changing production code.

- [ ] **Step 4: Commit evidence and update plan status**

```powershell
git add shared scripts docs
git commit -m "test: verify Menu recipe api contract and smoke flow"
```

After this task, the first-slice plan is complete only if the fresh configure/build/CTest/API evidence exists. Linux, Android, Playwright, QML profiling, 60-100 recipe expansion, authentication, admin CRUD and client UI remain separate next plans and must be marked unverified until run.
