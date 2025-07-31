package dto

// ErrorResponse represents a JSON error response.
type ErrorResponse struct {
	Message string
}

// DeleteUserRequest represents the request structure for the DeleteUser endpoint.
type DeleteUserRequest struct {
	UserID int64 `json:"user_id"  example:"1"`
}

// DeleteUserResponse represents the response structure for the DeleteUser endpoint.
type DeleteUserResponse struct {
	Message string `json:"message"`
}

// DigupUserRequest represents the request structure for the DigupUser endpoint.
type DigupUserRequest struct {
	UserID int64 `json:"user_id" example:"1"`
}

// DigupUserResponse represents the response structure for the DigupUser endpoint.
type DigupUserResponse struct {
	Message string `json:"message"`
}

// LoginRequest represents the structure of the login request
type LoginRequest struct {
	Phone    string `json:"phone" example:"79999999999"`
	Password string `json:"password" example:"123456"`
}

// LoginResponse represents the structure of the login response
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type LoginTelegramRequest struct {
	Phone      string `json:"phone" example:"79999999999"`
	TelegramId int64  `json:"telegram_id" example:"123"`
	Timestamp  int64  `json:"timestamp" example:"1744985930"`
}

// LoginTelegramResponse represents the structure of the login response
type LoginTelegramResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// LogoutResponse represents the structure of the logout response
type LogoutResponse struct {
	Message string `json:"message"`
}

// RecoveryRequest represents the request structure for password recovery.
type RecoveryRequest struct {
	Phone    string `json:"phone" example:"79999999999"`
	Password string `json:"password" example:"123456"`
	Code     string `json:"code" example:"12345"`
}

// RecoveryResponse represents the response structure for password recovery.
type RecoveryResponse struct {
	Message string `json:"message"`
}

// RefreshRequest represents the structure of the refresh request
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" example:"23f5383a94d671c0ddc82f7e635f6719362e74c67eecd22aec3df18833929a02c7423fe0f8d96283b4389ada5883f77700f64d64918702b7527e57a7e0c444ce"`
}

// RefreshResponse represents the structure of the refresh response
type RefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterRequest struct {
	Phone    string `json:"phone" example:"79999999999"`
	Password string `json:"password" example:"123456"`
	Email    string `json:"email" example:"test@test.com"`
	Code     string `json:"code" example:"12345"`
}

type RegisterResponse struct {
	Message string `json:"message" example:"User registration successful"`
}

// ResetTokenRequest represents the structure of the request payload for resetting tokens associated with a user account.
type ResetTokenRequest struct {
	UserID int64 `json:"user_id,omitempty" example:"1"`
}

// ResetTokenResponse represents the structure of the response payload for resetting tokens associated with a user account.
type ResetTokenResponse struct {
	Message string `json:"message"`
}

// ChangeRoleRequest represents the structure of the request payload for changing a user's role.
type ChangeRoleRequest struct {
	UserID int64  `json:"user_id" example:"1"`
	Role   string `json:"role" example:"user" enums:"user,creator,moderator,admin"`
}

// ChangeRoleResponse represents the structure of the response payload for changing a user's role.
type ChangeRoleResponse struct {
	Message string `json:"message"`
}

// UpdateUserRequest represents the structure of the request payload for changing phone number.
type UpdateUserRequest struct {
	Password string `json:"password,omitempty" example:"12345"`
	Phone    string `json:"phone,omitempty" example:"89999999999"`
	Code     string `json:"code,omitempty" example:"45236"`
}

// UpdateUserResponse represents the structure of the response payload for changing phone number.
type UpdateUserResponse struct {
	Message string `json:"message"`
}

type CheckTokenResponse struct {
	IsValid bool
	Error   string
	UserId  int64
	RoleId  int64
	Phone   string
}

type CheckTokenRequest struct {
	Token string
}

type PingResponse struct {
	Status  bool
	Message string
}

type RecoveryPasswordResponse struct {
	Message string
}

type RegistrationResponse struct {
	Message string
}

type RefreshTokenResponse struct {
	AccessToken  string
	RefreshToken string
}
