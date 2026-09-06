# API 规范参考(YAML 版)

> 本规范适用于 `auto_api` 工具的 **YAML 输入格式**。旧行式 DSL(`&描述 %规则 $ *pkg ^file`)已废弃。

## 基础结构

```yaml
# 接口说明(建议标明前台/后台/中台)
API_NAME:
  method: GET|POST|PUT|DELETE
  path: /path/:id(int)?query=int
  package: pkg名
  filename: file名
  permission: true   # 可选
  desc: 接口描述     # 可选
  req:               # 可选
    field: type %规则 # 描述
  res:               # 可选
    field: type # 描述
```

- 一句话场景、前台/后台/中台写在上方注释里(会进入生成的 Go 注释)。
- 没有 req/res 就不写 `req:` / `res:` 键。
- `permission: true` → 鉴权分组;缺省为公开接口。

## 类型系统

- 基础类型:`string` `int` `int8` `int16` `int32` `int64` `uint` `uint64` `float32` `float64` `bool` `datetime` `time.Time`
- 数组:
  - 流式:`tags: [string]`、`ids: [int64]`
  - 引号标量:`scores: "[int64]"`、`avatar: "[]byte"`
- 指针(必须引号):`pointer: "*int64"`
- map:`extra: json`(别名)或 `mp: map[string]interface{}`
- 嵌套结构体:字段值是嵌套映射
- 结构体数组:字段值是 `- 键: 类型` 序列

---

## 标准 CRUD 模板

### 1. 创建 API (POST)

```yaml
# 创建XXX-后台使用
CreateXXX:
  method: POST
  path: /xxx/create
  package: xxx
  filename: xxx
  permission: true
  req:
    name: string %min(1),max(32) # 名称
    description: string # 描述
  res:
    id: int64
```

### 2. 查询/搜索 API (GET)

```yaml
# 查询XXX列表-后台使用
GetXXXList:
  method: GET
  path: /xxx/list?keyword=string&page=int&size=int
  package: xxx
  filename: xxx
  req:
    status: int # 状态过滤(可选)
  res:
    list:
      - id: int64
        name: string
        description: string
        createdAt: time.Time
    total: int64
```

### 3. 详情查询 API (GET)

```yaml
# 获取XXX详情-后台使用
GetXXX:
  method: GET
  path: /xxx/detail
  package: xxx
  filename: xxx
  req:
    id: int64 # ID
  res:
    id: int64
    name: string
    description: string
    createdAt: time.Time
    updatedAt: time.Time
```

### 4. 单个更新 API (PUT)

```yaml
# 更新XXX-后台使用
UpdateXXX:
  method: PUT
  path: /xxx/update
  package: xxx
  filename: xxx
  permission: true
  req:
    id: int64 # ID
    name: string %max(32) # 名称
    description: string # 描述
```

### 5. 批量更新 API (PUT)

```yaml
# 批量更新XXX-后台使用
BatchUpdateXXX:
  method: PUT
  path: /xxx/batch
  package: xxx
  filename: xxx
  permission: true
  req:
    list:
      - id: int64 # ID
        name: string # 名称
        status: int # 状态
```

### 6. 单个删除 API (DELETE)

```yaml
# 删除XXX-后台使用
DeleteXXX:
  method: DELETE
  path: /xxx/delete
  package: xxx
  filename: xxx
  permission: true
  req:
    id: int64 # ID
```

### 7. 批量删除 API (DELETE)

```yaml
# 批量删除XXX-后台使用
BatchDeleteXXX:
  method: DELETE
  path: /xxx/delete
  package: xxx
  filename: xxx
  permission: true
  req:
    ids: [int64] # ID数组
```

### 8. 下拉选项/枚举 API (GET)

```yaml
# 获取XXX下拉选项-前台使用
GetXXXOptions:
  method: GET
  path: /xxx/options
  package: xxx
  filename: xxx
  res:
    list:
      - value: int64
        label: string
```

---

## 补充说明

1. **列表接口建议包含**:`page`(页码)、`size`(每页数量)、`keyword`(关键词)、`total`(总数)。查询参数声明在 `path` 里即可:`/xxx/list?page=int&size=int`。
2. **批量删除**:DELETE + req 里的 `ids: [int64]`。
3. **校验规则白名单**:`enum`/`oneof`、`regex`/`regexp`、`min`/`max`/`gt`/`gte`/`lt`/`lte`、`len`/`eq`/`ne`/`eqfield`/`nefield`、`email`/`phone`/`url`/`uuid`、`alphanum`/`numeric`/`contains`。白名单外忽略;有规则即必填;`enum` → `oneof`;`regex` 原样保留。
4. **指针/`[]byte` 必须双引号**,如 `"*int64"`、`"[]byte"`。
5. **包名与文件名尽量统一**,便于维护与路由分组。