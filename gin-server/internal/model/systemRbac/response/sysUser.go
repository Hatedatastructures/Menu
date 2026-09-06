package response

type RegisterRes struct {
	Id          int64  `json:"id"`          // 用户ID
	UUID        string `json:"uUID"`        // 用户UUID
	UserName    string `json:"userName"`    // 用户登录名
	NickName    string `json:"nickName"`    // 用户昵称
	HeaderImg   string `json:"headerImg"`   // 用户头像
	AuthorityId uint   `json:"authorityId"` // 用户角色ID
	Phone       string `json:"phone"`       // 用户手机号
	Email       string `json:"email"`       // 用户邮箱
	Enable      int    `json:"enable"`      // 是否启用
}

type LoginRes struct {
	Uuid        string              `json:"uuid"`        // 用户UUID
	UserName    string              `json:"userName"`    // 用户登录名
	NickName    string              `json:"nickName"`    // 用户昵称
	HeaderImg   string              `json:"headerImg"`   // 用户头像
	AuthorityId uint                `json:"authorityId"` // 用户角色ID
	Authorities []LoginResAuthority `json:"authorities"` // []
	Phone       string              `json:"phone"`       // 用户手机号
	Email       string              `json:"email"`       // 用户邮箱
	Enable      int                 `json:"enable"`      // 是否启用
	Token       string              `json:"token"`       // JWT token
	ExpiresAt   int64               `json:"expiresAt"`   // 过期时间
}

type LoginResAuthority struct {
	AuthorityId   uint   `json:"authorityId"`   // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
	ParentId      uint   `json:"parentId"`      // 父角色ID
	DefaultRouter string `json:"defaultRouter"` // 默认菜单
}

type GetUserInfoListRes struct {
	Page     int                      `json:"page"`
	Size     int                      `json:"size"`
	Username string                   `json:"username"`
	NickName string                   `json:"nickName"`
	Phone    string                   `json:"phone"`
	Email    string                   `json:"email"`
	List     []GetUserInfoListResList `json:"list"`  // []
	Total    int64                    `json:"total"` // 总数
}

type GetUserInfoListResList struct {
	Id          int64                           `json:"id"`          // 用户ID
	Uuid        string                          `json:"uuid"`        // 用户UUID
	UserName    string                          `json:"userName"`    // 用户登录名
	NickName    string                          `json:"nickName"`    // 用户昵称
	HeaderImg   string                          `json:"headerImg"`   // 用户头像
	AuthorityId uint                            `json:"authorityId"` // 用户角色ID
	Authority   GetUserInfoListResListAuthority `json:"authority"`
	Authorities []GetUserInfoListResListRole    `json:"authorities"` // []
	Roles       []GetUserInfoListResListRole    `json:"roles"`       // []
	Phone       string                          `json:"phone"`       // 用户手机号
	Email       string                          `json:"email"`       // 用户邮箱
	Enable      int                             `json:"enable"`      // 是否启用
}

type GetUserInfoListResListAuthority struct {
	AuthorityId   uint   `json:"authorityId"`   // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
	DefaultRouter string `json:"defaultRouter"` // 默认菜单
}

type GetUserInfoListResListRole struct {
	AuthorityId   uint   `json:"authorityId"`   // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
}

type GetUserInfRes struct {
	Id          int64                  `json:"id"`          // 用户ID
	Uuid        string                 `json:"uuid"`        // 用户UUID
	UserName    string                 `json:"userName"`    // 用户登录名
	NickName    string                 `json:"nickName"`    // 用户昵称
	HeaderImg   string                 `json:"headerImg"`   // 用户头像
	AuthorityId uint                   `json:"authorityId"` // 用户角色ID
	Authority   GetUserInfResAuthority `json:"authority"`
	Authorities []GetUserInfResAuthori `json:"authorities"` // []
	Authori     []GetUserInfResAuthori `json:"authori"`     // []
	Phone       string                 `json:"phone"`       // 用户手机号
	Email       string                 `json:"email"`       // 用户邮箱
	Enable      int                    `json:"enable"`      // 是否启用
}

type GetUserInfResAuthority struct {
	AuthorityId   uint   `json:"authorityId"`   // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
	DefaultRouter string `json:"defaultRouter"` // 默认菜单
}

type GetUserInfResAuthori struct {
	AuthorityId   uint   `json:"authorityId"`   // 角色ID
	AuthorityName string `json:"authorityName"` // 角色名
}

type FindUserByIdRes struct {
	Id          int64  `json:"id"`          // 用户ID
	Uuid        string `json:"uuid"`        // 用户UUID
	UserName    string `json:"userName"`    // 用户登录名
	NickName    string `json:"nickName"`    // 用户昵称
	HeaderImg   string `json:"headerImg"`   // 用户头像
	AuthorityId uint   `json:"authorityId"` // 用户角色ID
	Phone       string `json:"phone"`       // 用户手机号
	Email       string `json:"email"`       // 用户邮箱
	Enable      int    `json:"enable"`      // 是否启用
}

type FindUserByUuidRes struct {
	Id          int64  `json:"id"`          // 用户ID
	Uuid        string `json:"uuid"`        // 用户UUID
	UserName    string `json:"userName"`    // 用户登录名
	NickName    string `json:"nickName"`    // 用户昵称
	HeaderImg   string `json:"headerImg"`   // 用户头像
	AuthorityId uint   `json:"authorityId"` // 用户角色ID
	Phone       string `json:"phone"`       // 用户手机号
	Email       string `json:"email"`       // 用户邮箱
	Enable      int    `json:"enable"`      // 是否启用
}
