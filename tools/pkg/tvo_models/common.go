package tvomodels

const DefaultLimit = 10
const MaxLimit = 50
const DefaultOffset = 0
const DefaultFrom = 0
const DefaultPerPage = 10
const MaxPerPage = 50
const DefaultPage = 0

const (
	USER      RoleId = 1
	CREATOR   RoleId = 2
	MODERATOR RoleId = 99
	ADMIN     RoleId = 100
)

// TokenData структура с данными из токена
type TokenData struct {
	UserID     int64  `json:"id"`
	UserPhone  string `json:"phone"`
	UserEmail  string `json:"email"`
	UserRoleID RoleId `json:"role_id"`
	RawToken   string `json:"raw_token"`
}
