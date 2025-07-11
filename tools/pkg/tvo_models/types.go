package tvomodels

// RoleId тип для роли пользователя
type RoleId int64

// Boolean тип для булевских параметров dto
type Boolean int

const (
	Null Boolean = iota // EnumIndex = 0
	Yes                 // 1
	No                  // 2
)

// IsNull проверка значения на null
func (b Boolean) IsNull() bool {
	return b == Null
}

// ToBool преобразование значения к bool
func (b Boolean) ToBool() bool {
	return b == Yes
}

// IntToBoolean преобразует int к Boolean
func IntToBoolean(value int) Boolean {
	switch value {
	case 0:
		return Null
	case 1:
		return Yes
	case 2:
		return No
	default:
		return Null
	}
}

// MaskValue нормализует роль пользователя, заменяя MODERATOR и ADMIN на USER
func (r RoleId) MaskValue() RoleId {
	switch r {
	case MODERATOR, ADMIN:
		return USER
	default:
		return r
	}
}
