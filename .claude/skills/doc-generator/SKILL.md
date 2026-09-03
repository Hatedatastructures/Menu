---
name: doc-generator
description: Generate YAML API definitions from requirements, then run autoApi.exe to generate Go backend + TS frontend code. Use when user says "generate API docs", "write api.yaml", or asks to create API interfaces from a requirement document. Parses markdown requirement documents and outputs YAML API specs to docs_yaml/ (or a directory passed to autoApi.exe).
---

# Doc Generator

根据需求文档生成 **YAML API 定义**,并用 `autoApi.exe` 生成前后端代码。

当前工具的输入格式是 **YAML**(已替换旧的行式 DSL)。旧的 `API_NAME METHOD PATH *pkg ^file { req{} }` 写法已不再支持。

## 工作流程

1. 运行 [`scripts/find-model.js`](scripts/find-model.js) 拿到 model 文件列表,读取业务相关 model 作参照,避免"幻觉字段"。
2. 按下方 YAML 格式编写 API 定义,写入 `docs_yaml/<模块>.yaml`(或用户指定、随后传给 `autoApi.exe` 的目录)。
3. 运行 `autoApi.exe <目录>`(增量)或 `autoApi.exe <目录> -w`(强制重生成),生成到 `example/`(见 `config.yaml` 的 `gen.root`)。

## 支持的 HTTP 方法

只支持四种: GET、POST、PUT、DELETE。

## YAML API 定义格式

每个 API 是一个顶层映射,键为 API 名(建议 `动词+名词`,如 `GetCharacterList`)。

```yaml
# 接口描述注释(会进入生成的 Go 注释,多行会被拍平为单行)
API_Name:
  method: GET|POST|PUT|DELETE   # 必须
  path: /xxx/:id(int64)?page=int&size=int   # 必须; 支持路径/查询参数
  package: pkg名                # 生成到 server/internal/{service,api/v1,router}/<pkg名>
  filename: file名              # 生成文件名(file 同名的 API 归入同一个文件/路由组)
  permission: true              # 可选: 标记需要鉴权(默认 false; 只影响 Public/Private 分组)
  desc: 接口描述                # 可选: 覆盖行上方的注释
  req:                          # 可选: 请求体字段
    field: type %规则 # 描述
    ...
  res:                          # 可选: 响应字段
    field: type # 描述
    ...
```

### 字段类型

`字段名: 类型 %规则 #描述` —— 类型与规则之间用 `%` 分隔,`#` 后是字段描述。

| 写法 | DSL/Go 类型 | 说明 |
|------|-------------|------|
| `name: string` | `string` | 基础类型,见下方别名表 |
| `tags: [string]` | `[]string` | 流式数组(标量可省略引号) |
| `ids: "[int64]"` | `[]int64` | 数组也可写引号标量 |
| `avatar: "[]byte"` | `[]byte` | **`[]byte` 必须双引号**,否则 YAML 报错 |
| `pointer: "*int64"` | `*int64` | **指针必须双引号**(裸 `*int64` 是 YAML 锚点别名,会卡死解析器/内存暴涨) |
| `extra: json` | `map[string]interface{}` | map 用 `json`/`jsonb`/`object` 别名 |
| `mp: map[string]interface{}` | `map[string]interface{}` | 也可写完整 `map[K]V` |
| `anything: any` | `interface{}` | 任意类型 |
| 嵌套映射(见下) | 内联结构体 | `req/res/字段` 下再叠一层映射 |
| 序列(见下) | 结构体数组 | 字段值是 `- 键: 类型` 列表 |

**内联结构体:**
```yaml
res:
  basicInfo:
    fullName: string
    age: int64
```

**结构体数组:**
```yaml
res:
  list:
    - id: int64
      name: string
```

### 路径参数与查询参数

直接在 `path` 里声明,不再需要单列字段:

- 路径参数:`/characters/:id` 或带类型 `/characters/:id(int64)` → req 里生成 `Id int64` 并从 URL 取值。
- 查询参数:`/types?page=int&size=int` → req 里生成 `Page int` / `Size int` 并用 `c.Query(...)` 取值。

```yaml
GetCollectionCharacters:
  method: GET
  path: /collections/:id(int)/characters?page=int&size=int
  package: collection
  res:
    list:
      - id: int64
        name: string
    total: int64
```

### 校验规则(白名单)

规则写在类型后 `%规则1,规则2`,多条用逗号分隔(括号内逗号不拆分)。

```yaml
    action: string %enum(start,stop,restart) # → validate:"required,oneof=start stop restart"
    code: string %regex(^[a-z]+,[0-9]+$) # → 正则原样保留(含逗号)
    phone: string %regex(^1[3-9]\d{9}$) # 手机号
    email: string %email
    score: int %gte(0),lte(100)
    name: string %min(1),max(32)
```

- **白名单规则**: `enum`/`oneof`, `regex`/`regexp`, `min`/`max`/`gt`/`gte`/`lt`/`lte`, `len`/`eq`/`ne`/`eqfield`/`nefield`, `email`/`phone`/`url`/`uuid`, `alphanum`/`numeric`/`contains`。
- `enum(1,2,3)` → `oneof=1 2 3`(空格分隔)。
- `regex(...)` 原样进入 `validate:"regex=..."`,包含逗号也不破坏。
- **白名单外的规则直接忽略**,不会注入生成的 validate tag。
- **只要写了任意规则,该字段就是必填**(生成 `required,…`)。

## 类型别名表

| DSL 写法 | Go 类型 | 说明 |
|----------|---------|------|
| string / str / varchar / text / char / uuid | `string` | |
| bool / boolean | `bool` | |
| int / integer | `int` | |
| int8 / tinyint | `int8` | |
| int16 / smallint | `int16` | |
| int32 | `int32` | |
| int64 / long / bigint | `int64` | |
| uint | `uint` | |
| uint64 / bigint unsigned | `uint64` | |
| float / float32 | `float32` | |
| float64 / double / decimal / numeric / number | `float64` | |
| time / time.time / datetime / timestamp / date | `time.Time` | |
| bytes / blob / binary / []byte | `[]byte` | 须加引号 |
| json / jsonb / object | `map[string]interface{}` | |
| any / interface / null | `interface{}` | |

> 未知类型不会被拒绝,但会在解析时打印 ⚠ 告警并按 `string` 处理——写错类型请留意 CLI 输出。

## model 定义(可选)

顶层 `model:` 块定义 gorm 模型(生成 `server/internal/model/<pkg>/XxxModel.go`):

```yaml
model:
  Goods:
    package: mall
    fields:
      Id: int64 # pk
      Name: string # 名称
      Price: float64
```

- 字段名用 **PascalCase**。
- `# pk` 注释标记主键。
- `CreatedAt/UpdatedAt/DeletedAt` 由模板自动生成,**不要再写**。

## 输出格式要求

直接输出 YAML 定义,不要输出:`#`/`##` 等 markdown 标题、规则说明、分类标题。单文件注释(`# 模块名`)和 `# =====` 分隔线可以保留。

## 完整示例

一个模块 = 一个文件,包含 CRUD 一套:

```yaml
# 电商商品模块
model:
  Goods:
    package: mall
    fields:
      Id: int64 # pk
      Name: string
      Price: float64

# 创建商品
CreateGoods:
  method: POST
  path: /goods/create
  package: mall
  filename: goods
  permission: true
  req:
    name: string %min(1),max(64) # 商品名称
    price: float64 %gt(0) # 价格
    status: int %enum(0,1,2) # 状态
  res:
    id: int64

# 商品列表
ListGoods:
  method: GET
  path: /goods/list?page=int&size=int
  package: mall
  filename: goods
  req:
    keyword: string # 关键词
    status: int # 状态
  res:
    list:
      - id: int64
        name: string
        price: float64
        status: int
    total: int64

# 删除商品
DeleteGoods:
  method: DELETE
  path: /goods/delete
  package: mall
  filename: goods
  permission: true
  req:
    id: int64 # 商品ID
```

## 详细说明

见 [references/apis-spec.md](references/apis-spec.md)(标准 CRUD 模板)与 [references/example.md](references/example.md)(完整业务示例)。

## Model 文件参考

生成 API 时应参考项目中的 model 文件,确保 req/res 字段与真实数据结构一致。

```bash
node .claude/skills/doc-generator/scripts/find-model.js
```

脚本会:向上查找 `server` 目录 → 读取 `internal/model/<pkg>/` 下第一层 `.go` 文件(排除 `common/enum/system/request/response` 子目录与 `_test.go`)→ 返回绝对路径数组。按 `/api/v1/<pkg>` 定位业务模块后,读取对应 model 文件再写 API,避免"幻觉字段"。