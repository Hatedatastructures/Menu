# 完整业务示例(YAML 版)

> 角色收藏夹模块 CRUD —— 演示路径参数、查询参数、结构体数组、`[]int64` 数组等写法。
> 这是旧行式 DSL 的 YAML 等价写法(旧 `*collection ^collection` → 新 `package: collection` / `filename: collection`)。

```yaml
# 获取角色收藏夹列表
GetCollectionList:
  method: GET
  path: /users/:userId(int)/collections?page=int&size=int
  package: collection
  filename: collection
  res:
    list:
      - id: int64
        name: string # 收藏夹名称
        description: string # 描述
        characterCount: int64 # 角色数量
        cover: string # 封面URL
        isPublic: bool # 是否公开
        createdAt: string # 创建时间
    total: int64

# 获取收藏夹详情
GetCollection:
  method: GET
  path: /collections/:id(int)
  package: collection
  filename: collection
  res:
    id: int64
    name: string # 收藏夹名称
    description: string # 描述
    cover: string # 封面URL
    isPublic: bool # 是否公开
    characterCount: int64 # 角色数量
    createdAt: string # 创建时间
    updatedAt: string # 更新时间

# 创建角色收藏夹
CreateCollection:
  method: POST
  path: /collections
  package: collection
  filename: collection
  permission: true
  req:
    name: string %min(1),max(32) # 收藏夹名称
    description: string # 描述
    isPublic: bool # 是否公开
  res:
    id: int64

# 更新角色收藏夹
UpdateCollection:
  method: PUT
  path: /collections/:id(int)
  package: collection
  filename: collection
  permission: true
  req:
    name: string %max(32) # 收藏夹名称
    description: string # 描述
    isPublic: bool # 是否公开

# 删除角色收藏夹
DeleteCollection:
  method: DELETE
  path: /collections/:id(int)
  package: collection
  filename: collection
  permission: true

# 批量删除角色收藏夹
BatchDeleteCollections:
  method: DELETE
  path: /collections
  package: collection
  filename: collection
  permission: true
  req:
    ids: [int64] # 收藏夹ID数组

# 批量更新收藏夹
BatchUpdateCollections:
  method: PUT
  path: /collections/batch
  package: collection
  filename: collection
  permission: true
  req:
    list:
      - id: int64
        name: string # 收藏夹名称
        isPublic: bool # 是否公开

# 修改收藏夹状态(启用/禁用)
UpdateCollectionStatus:
  method: PUT
  path: /collections/status
  package: collection
  filename: collection
  permission: true
  req:
    ids: [int64] # 收藏夹ID数组
    isPublic: bool # 是否公开

# 向收藏夹添加角色
AddCharacterToCollection:
  method: POST
  path: /collections/characters
  package: collection
  filename: collection
  permission: true
  req:
    characterId: int64 # 角色ID
  res:
    success: bool

# 从收藏夹移除角色
RemoveCharacterFromCollection:
  method: DELETE
  path: /collections/:id(int)/characters/:characterId(int)
  package: collection
  filename: collection
  permission: true

# 移动角色到其他收藏夹
MoveCharacter:
  method: POST
  path: /characters/move
  package: collection
  filename: collection
  permission: true
  req:
    characterId: int64 # 角色ID
    fromCollectionId: int64 # 源收藏夹ID
    toCollectionId: int64 # 目标收藏夹ID

# 批量移除角色
BatchRemoveCharacters:
  method: DELETE
  path: /collections/:id(int)/characters
  package: collection
  filename: collection
  permission: true
  req:
    characterIds: [int64] # 角色ID数组

# 调整角色顺序
SortCharacters:
  method: PUT
  path: /collections/:id(int)/characters/sort
  package: collection
  filename: collection
  permission: true
  req:
    characterIds: [int64] # 角色ID数组(顺序即为排序)

# 检查角色是否已收藏
CheckCharacterCollected:
  method: GET
  path: /characters/:characterId(int)/collected
  package: collection
  filename: collection
  res:
    collected: bool # 是否已收藏
    collectionId: int64 # 收藏夹ID

# 获取收藏夹中的角色列表
GetCollectionCharacters:
  method: GET
  path: /collections/:id(int)/characters?page=int&size=int
  package: collection
  filename: collection
  res:
    list:
      - id: int64
        name: string
        avatar: string
        age: int64
        race: string
        description: string
        tags: [string]
        addedAt: string # 添加时间
    total: int64
```