# Menu 后端第一条垂直切片设计

状态：已根据用户提供的产品约束确认，作为 Phase 1-2 的实现规格。

日期：2026-09-01

## 目标

交付一个可以从零配置、启动并通过真实 HTTP 请求验证的 C++20 服务端，以及可被网站和 Qt 客户端直接消费的菜谱/推荐契约。用户能够获得人工核对的结构化菜谱，服务端可以按偏好、时间、过敏原和家中已有食材返回最多三个今晚方案。

本切片完成后，以下命令链路必须存在：

```text
fresh CMake configure
        |
build MenuServer + unit/integration tests
        |
start MenuServer with SQLite migrations/seed
        |
GET /healthz, /readyz, /api/v1/recipes, /api/v1/recipes/{id}
        |
validate JSON/error status with an external HTTP client
```

## 不在本切片内

- Android Qt kit、Android SDK/NDK/JDK/Gradle 的安装和 APK 打包。
- React 管理面板和 Qt Quick 客户端页面；它们只先消费本切片稳定的 API。
- 计划、做饭会话、提醒和反馈的完整用例；数据模型和 API 兼容边界在后续 spec 中追加。
- 外部商超、AI 推荐、社交、聊天和版权内容抓取。

## 模块与依赖

```text
MenuFoundation
    <- MenuDomain
    <- MenuApplication
    <- MenuInfrastructure
    <- MenuTransport

MenuDomain <- MenuApplication
MenuApplication <- MenuInfrastructure
MenuApplication <- MenuApi
MenuTransport <- MenuApi
MenuApi <- MenuServer
```

### MenuFoundation

提供 `Error`、`Result<T>`、配置值、稳定 ID、UTC 时间、输入大小限制和日志接口。它只依赖 C++ 标准库和 Boost.System，不知道 HTTP、SQLite 或菜谱。

### MenuDomain

提供以下值和实体：

- `Ingredient`：稳定 ID、中文名称、别名、类别、默认单位、是否常备、替代组和可为空的 `StoreSkuMapping`。
- `RecipeIngredient`：Ingredient ID、数量、单位、份量系数、处理方式和必需/可选标记。
- `RecipeStep`：顺序、标题、说明、预计秒数和是否需要计时。
- `Recipe`：稳定 ID、slug、名称、菜系、描述、难度、准备/烹饪分钟、基准份量、图片资源、发布状态、食材和步骤。
- `RecommendationRequest` 与 `Recommendation`：推荐输入和渲染所需的已有/缺少食材数量。

领域规则是纯函数或无 I/O 的对象方法：

1. `NormalizeIngredientAlias` 统一大小写、空白和已注册中文别名。
2. `ScaleRecipeIngredients` 按正份量计算，并保持单位与可选标记。
3. `BuildIngredientAvailability` 将用户库存和配方 Ingredient ID 比较，输出已有/缺少。
4. `ValidateRecipe` 拒绝零份量、负时长、重复步骤序号、空必需食材和不支持的单位。

### MenuApplication

定义同步执行的存储端口和推荐端口；同步只是端口调用方式，不表示它可以运行在 Asio event loop：Infrastructure 调用必须由 worker executor 调度。

```cpp
class RecipeRepository {
public:
    virtual ~RecipeRepository() = default;
    virtual Result<std::vector<Recipe>> ListPublished() = 0;
    virtual Result<std::optional<Recipe>> FindPublishedById(std::string_view id) = 0;
};

class RecommendationProvider {
public:
    virtual ~RecommendationProvider() = default;
    virtual Result<std::vector<Recommendation>> Recommend(
        const RecommendationRequest& request,
        const std::vector<Recipe>& recipes) = 0;
};
```

默认实现 `RuleBasedRecommendationProvider` 按以下顺序过滤和排序：过敏原/忌口硬过滤，时间和厨具硬过滤，偏好菜系加分，已有食材覆盖率加分，难度和近期重复度作为次级排序，最终最多取三项。推荐过程不访问网络，不调用模型。

### MenuInfrastructure

`SqliteDatabase` 负责同一连接的打开、busy timeout、WAL、外键、事务和错误转换。`MigrationRunner` 按数字前缀执行迁移并在 `SchemaMigration` 表记录版本。`SqliteRecipeRepository` 只使用准备语句和绑定参数，禁止把用户输入拼进 SQL。

数据库 writer 由 Asio `strand`/单一 actor 顺序化；读取和种子导入在 storage executor 执行。首版可以使用一个受控 SQLite 连接，但这个约束必须是显式接口，后续才能安全扩展读取连接池。

### MenuTransport 与 MenuApi

`MenuTransport` 用 Boost.Beast + Boost.Asio 异步接收 HTTP/1.1 请求，连接 session 由 `std::enable_shared_from_this` 管理，读写串行化到 per-connection strand，设置请求体大小上限和超时。

`MenuApi` 将请求转换为边界 DTO，做路径/查询/JSON 校验，把阻塞的 Application 调用投递到 worker executor，再将响应投递回连接 executor。统一错误格式：

```json
{
  "error": {
    "code": "recipe_not_found",
    "message": "菜谱不存在",
    "requestId": "..."
  }
}
```

第一切片公开：

| Method | Path | 成功 | 失败 |
| --- | --- | --- | --- |
| GET | `/healthz` | 200，进程存活 | 无 |
| GET | `/readyz` | 200，迁移/种子就绪 | 503 |
| GET | `/api/v1/recipes` | 200，分页/筛选后的已发布菜谱 | 400 参数错误 |
| GET | `/api/v1/recipes/{id}` | 200，含规范化食材和步骤 | 404 不存在，400 非法 ID |
| POST | `/api/v1/recommendations/tonight` | 200，最多三个方案 | 400 输入错误 |

列表查询支持 `cuisine`、`maxMinutes`、`difficulty`、`limit`，默认 `limit=20`，上限 100。推荐请求至少包含 `servings` 和 `availableMinutes`，可选 `cuisines`、`allergies`、`pantryIngredientIds`、`cookware`。

## 数据库首版表

首个迁移建立完整模型的基础表，即使第一切片暂时只读 Recipe/Ingredient，也不在后续用破坏性重构替换 ID：

`Users`、`Preferences`、`Ingredients`、`IngredientAliases`、`Recipes`、`RecipeIngredients`、`RecipeSteps`、`MealPlans`、`MealPlanItems`、`Reminders`、`CookingSessions`、`Feedback`、`MediaAssets`、`RefreshTokens`、`SchemaMigration`。

Recipe/Ingredient 的外键、唯一 slug、步骤 `(RecipeId, StepOrder)`、RecipeIngredient 的数量/单位检查由 SQL 约束和领域校验双重保护。时间统一存 UTC ISO-8601 或整数秒，数量使用 REAL 仅表示可量化食材，展示层负责本地化。

## 错误和安全

- 请求行、header、query 和 body 都有明确上限，JSON 深度/数组长度有上限。
- 非法 JSON、未知字段、非法 ID、无效单位和负数时长返回 400，不泄漏 SQL 错误。
- SQLite 错误只写脱敏日志并映射为稳定 `storage_unavailable`；请求响应不包含文件路径。
- 首切片先保留认证接口的路由位置和统一错误模型；完整登录使用 OpenSSL scrypt、短期 access token、数据库哈希 refresh token 和撤销状态，在认证切片中实现。
- CORS 只允许配置白名单来源；默认只监听本机或配置指定地址，管理 API 后续默认需要管理员授权。

## 测试策略

严格遵循 RED-GREEN-REFACTOR：

1. 先写领域测试并运行到“符号不存在/行为失败”。
2. 最小实现后单测转绿。
3. 再写迁移/仓储测试，验证 WAL、事务回滚、参数化查询和重复迁移幂等。
4. 再写 API 集成测试，启动真实 listener，用 TCP/HTTP 客户端验证状态码、JSON、404、400、超大 body 和并发读。
5. 每个阶段运行 `git diff --check`、fresh configure、构建和 CTest；不使用旧 build 目录作为证据。

最低测试集合：

- IngredientAliasNormalizationTest
- RecipeServingScaleTest
- RecipeValidationTest
- RecommendationRuleTest
- SqliteMigrationTest
- SqliteRecipeRepositoryTest
- MenuApiHealthTest
- MenuApiRecipeTest
- MenuApiInvalidInputTest
- MenuApiConcurrentReadTest

## 运行时验收

必须记录：构建编译器、CMake preset、CTest 用例数、服务端 PID、实际 curl 请求和响应状态、数据库路径、端口、并发请求数量、p95/p99、错误率。Windows 本地结果只代表 Windows；Linux、Android、桌面截图和 120Hz 测量单独报告。
