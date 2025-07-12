package handlers

import (
	"context"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"main/internal/auth/tools"
	jwtManager "main/internal/lib/jwt"
	"main/internal/models"
	"main/internal/repository"
	"main/tools/pkg/cache"
	"main/tools/pkg/helpers"
	tvomodels "main/tools/pkg/tvo_models"
	"main/tools/validator"

	"main/internal/dto"
	httputils "main/tools/pkg/http_utils"
	"main/tools/pkg/logger"
	tvoerrors "main/tools/pkg/tvo_errors"
)

// AuthHandlers
type AuthHandlers struct {
	logger          *logger.Logger
	jwt             *jwtManager.JWTManager
	userRepository  repository.UserRepository
	tokenRepository repository.UserTokenRepository
	roleRepository  repository.RoleRepository
	cache           cache.CacheClient
	secret          string
}

var ErrNotAdmin = errors.New("available only to admin")
var ErrInvalidPassword = errors.New("invalid password")
var ErrPhoneTaken = errors.New("phone already taken")
var AuthHandler *AuthHandlers

// NewAuthHandlers конструктор для обработчиков IDM методов
func NewAuthHandlers(logger *logger.Logger,
	jwt *jwtManager.JWTManager,
	userRepository repository.UserRepository,
	tokenRepository repository.UserTokenRepository,
	roleRepository repository.RoleRepository,
	client cache.CacheClient, secret string) *AuthHandlers {
	AuthHandler = &AuthHandlers{
		logger:          logger,
		jwt:             jwt,
		userRepository:  userRepository,
		tokenRepository: tokenRepository,
		roleRepository:  roleRepository,
		cache:           client,
		secret:          secret,
	}
	return AuthHandler
}

// Registration прокси метод для отправки его в сервис IDM
// @Summary user Registration
// @Description Register a new user
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration Request Body"
// @Success 200 {object} dto.RegisterResponse "Registration successful"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Router /idm/registration [post]
func (h *AuthHandlers) Registration(c *fiber.Ctx) (interface{}, error) {
	var request dto.RegisterRequest

	ctx := c.Context()
	if err := httputils.ParseRequestBody(c, &request, "Registration", h.logger); err != nil {
		return nil, tvoerrors.ErrInvalidRequestData
	}
	if request.Code != "98765" {
		return nil, status.Error(codes.InvalidArgument, "invalid code") //nolint
	}
	if err := validator.ValidPhone(request.Phone); err != nil {
		log.Error("Invalid phone number", "error", err)
		return nil, status.Error(codes.InvalidArgument, "invalid phone number") //nolint
	}

	//phoneExists, err := h.userRepository.PhoneExists(ctx, request.Phone)
	//if err != nil {
	//	if !errors.Is(err, tvoerrors.ErrNotFound) {
	//		log.Error("Find phone error", "error", err)
	//		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	//	}
	//}
	//
	//if phoneExists {
	//	log.Error("Phone exists", "phone", request.Phone, "error", ErrPhoneTaken)
	//	return nil, status.Error(codes.InvalidArgument, "phone already taken") //nolint
	//}

	_, err := h.userRepository.CreateUser(ctx, request.Phone, request.Password)
	if err != nil {
		log.Error("Error creating user ", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	return &dto.RegistrationResponse{
		Message: "User registration successful",
	}, nil
}

// Login прокси метод для отправки его в сервис IDM
// @Summary Logs in a user
// @Description Logs in a user with the provided phone number and password
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Request Body"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /idm/login [post]
func (h *AuthHandlers) Login(c *fiber.Ctx) (interface{}, error) {
	var request dto.LoginRequest

	ctx := c.Context()

	if err := httputils.ParseRequestBody(c, &request, "Login", h.logger); err != nil {
		return nil, tvoerrors.ErrInvalidRequestData
	}

	if err := validator.ValidPhone(request.Phone); err != nil {
		log.Error("Validation failed", "error", tvoerrors.ErrInvalidPhone)
		return nil, status.Error(codes.InvalidArgument, "invalid phone number") //nolint
	}

	user, err := h.userRepository.UserByPhone(ctx, request.Phone)
	if err != nil {
		log.Error("Error finding user", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	if !tools.PasswordsMatch(user.Password, request.Password, h.secret, user.Salt) {
		log.Error("Invalid password", "error", ErrInvalidPassword)
		return nil, status.Error(codes.Unauthenticated, "Invalid Username or Password") //nolint
	}

	refreshToken := tools.GenerateRefreshToken()
	accessToken, err := h.jwt.Generate(user)
	if err != nil {
		log.Error("Error generate token", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	userToken, err := h.tokenRepository.Create(ctx, user.ID, accessToken, refreshToken, h.jwt.GetTokenTTL(), h.jwt.GetRefreshTTL())
	if err != nil {
		log.Error("Error creating token", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	data := helpers.JsonEncodeString(user)
	if err = h.cache.Set(ctx, accessToken, data, h.jwt.GetTokenTTL()); err != nil {
		log.Error("Error store token", "error", err)
	}

	return &dto.LoginResponse{
		AccessToken:  userToken.Token,
		RefreshToken: userToken.RefreshToken,
	}, nil
}

// Refresh прокси метод для отправки его в сервис IDM
// @Summary Refreshes an access token
// @Description Refreshes an access token using the provided refresh token
// @Security ApiKeyAuth
// @Tags Authentication
// @Accept json
// @Produce json
// @Param RefreshToken body dto.RefreshRequest true "Refresh token"
// @Success 200 {object} dto.RefreshResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /idm/refresh [post]
func (h *AuthHandlers) Refresh(c *fiber.Ctx) (interface{}, error) {
	var request dto.RefreshRequest

	ctx := c.Context()

	if err := httputils.ParseRequestBody(c, &request, "Refresh", h.logger); err != nil {
		return nil, tvoerrors.ErrInvalidRequestData
	}

	userToken, err := h.tokenRepository.GetRefreshToken(ctx, request.RefreshToken)
	if err != nil {
		if errors.Is(err, tvoerrors.ErrNotFound) {
			log.Error("GetRefreshToken refresh token not found", "token",
				request.RefreshToken, "error", err)
			return nil, status.Error(codes.InvalidArgument, "invalid token") //nolint
		}
		log.Error("Error get refresh token", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	user, err := h.userRepository.UserById(ctx, userToken.UserID)
	if err != nil {
		log.Error("Error getting user", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	refreshToken := tools.GenerateRefreshToken()
	accessToken, err := h.jwt.Generate(user)
	if err != nil {
		log.Error("Error generate token", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	userToken, err = h.tokenRepository.Create(ctx, user.ID, accessToken, refreshToken, h.jwt.GetTokenTTL(),
		h.jwt.GetRefreshTTL())
	if err != nil {
		log.Error("Error creating token", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	data := helpers.JsonEncodeString(user)
	if err = h.cache.Set(ctx, accessToken, data, h.jwt.GetTokenTTL()); err != nil {
		log.Error("Error store token", "error", err)
	}

	return &dto.RefreshTokenResponse{
		AccessToken:  userToken.Token,
		RefreshToken: userToken.RefreshToken,
	}, nil
}

// Recovery прокси метод для отправки его в сервис IDM
// @Summary Handles the password recovery request
// @Description Handles the password recovery request by verifying the OTP and updating the user's password
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.RecoveryRequest true "Password recovery request body"
// @Success 200 {object} dto.RecoveryResponse "Password changed successfully"
// @Failure 400 {object} dto.ErrorResponse "Bad request"
// @Failure 401 {object} dto.ErrorResponse "Unauthorized"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /idm/recovery [post]
func (h *AuthHandlers) Recovery(c *fiber.Ctx) (interface{}, error) {
	var request dto.RecoveryRequest

	ctx := c.Context()
	if err := httputils.ParseRequestBody(c, &request, "Recovery", h.logger); err != nil {
		return nil, tvoerrors.ErrInvalidRequestData
	}

	user, err := h.userRepository.UserByPhone(ctx, request.Phone)
	if err != nil {
		log.Error("Error getting user", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	if err = h.tokenRepository.TokenReset(ctx, user.ID); err != nil {
		log.Error("Error resetting token", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	if err = h.removeUserTokens(ctx, user.ID); err != nil {
		log.Error("Error removing tokens", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	if err = h.userRepository.UpdatePassword(ctx, user.ID, request.Password); err != nil {
		log.Error("Error updating password", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	return &dto.RecoveryPasswordResponse{
		Message: "Password changed successfully",
	}, nil
}

// Logout прокси метод для отправки его в сервис IDM
// @Summary Logs out a user
// @Description Logs out a user by deleting the provided token
// @Tags Authentication
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} dto.LogoutResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /idm/logout [post]
func (h *AuthHandlers) Logout(c *fiber.Ctx) (interface{}, error) {
	ctx := httputils.CtxWithAuthToken(c)
	tokenData, err := httputils.GetTokenDataFromCtx(ctx)
	if err != nil {
		log.Error("Failed to get token data", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get claims from token") //nolint
	}

	if err = h.tokenRepository.DeleteRefreshToken(ctx, tokenData.RawToken); err != nil {
		log.Error("Error delete refresh token", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	if _, err = h.cache.Del(ctx, tokenData.RawToken); err != nil {
		log.Error("Error caching token", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	return &dto.LogoutResponse{
		Message: "Successfully logged out",
	}, nil
}

// DeleteUser прокси метод для отправки его в сервис IDM
// @Summary Delete a user
// @Description Delete a user from the system
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.DeleteUserRequest true "Request body"
// @Success 200 {object} dto.DeleteUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /idm/delete_user [post]
func (h *AuthHandlers) DeleteUser(c *fiber.Ctx) (interface{}, error) {
	var request dto.DeleteUserRequest

	if err := httputils.ParseRequestBody(c, &request, "DeleteUser", h.logger); err != nil {
		return nil, tvoerrors.ErrInvalidRequestData
	}

	ctx := httputils.CtxWithAuthToken(c)
	tokenData, err := httputils.GetTokenDataFromCtx(ctx)
	if err != nil {
		log.Error("Failed to get token data", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get claims from token") //nolint
	}

	user, err := h.userRepository.UserById(ctx, request.UserID)
	if err != nil {
		log.Error("Error fetching user", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	isAdmin := tokenData.UserRoleID == tvomodels.ADMIN

	if (isAdmin && tokenData.UserID == request.UserID) || (!isAdmin && tokenData.UserID != request.UserID) {
		log.Error("Delete user not allowed")
		return nil, status.Error(codes.PermissionDenied, "not allowed") //nolint
	}

	if err = h.userRepository.DeleteUser(ctx, user.ID); err != nil {
		log.Error("Error deleting user", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	if err = h.removeUserTokens(ctx, user.ID); err != nil {
		log.Error("Error removing tokens", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	return &dto.DeleteUserResponse{
		Message: "User deleted",
	}, nil
}

// DigupUser прокси метод для отправки его в сервис IDM
// Only administrators are allowed to perform this action.
// @Summary Dig up a user
// @Description Dig up a user from the database
// @Tags User
// @Accept json
// @Produce json
// @Param request body dto.DigupUserRequest true "Request body"
// @Success 200 {object} dto.DigupUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Security ApiKeyAuth
// @Router /idm/digup_user [post]
func (h *AuthHandlers) DigupUser(c *fiber.Ctx) (interface{}, error) {
	var request dto.DigupUserRequest

	if err := httputils.ParseRequestBody(c, &request, "DigupUser", h.logger); err != nil {
		return nil, tvoerrors.ErrInvalidRequestData
	}

	ctx := httputils.CtxWithAuthToken(c)
	tokenData, err := httputils.GetTokenDataFromCtx(ctx)
	if err != nil {
		log.Error("Failed to get token data", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get claims from token") //nolint
	}

	if tokenData.UserRoleID != tvomodels.ADMIN {
		log.Error("Only admin can dig up", "error", ErrNotAdmin)
		return nil, status.Error(codes.PermissionDenied, "something went wrong") //nolint
	}
	_, err = h.userRepository.DigUpUser(ctx, request.UserID)
	if err != nil {
		log.Error("Error digup user", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	return &dto.DigupUserResponse{
		Message: "User digup successful",
	}, nil
}

// UpdateUser прокси метод для отправки его в сервис IDM
// @Summary Change user's password or phone number
// @Description This endpoint allows users to change their password or phone number by providing either a new password or a new phone number along with a verification code.
// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.UpdateUserRequest true "Request body"
// @Success 200 {object} dto.UpdateUserResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /idm/update [post]
func (h *AuthHandlers) UpdateUser(c *fiber.Ctx) (interface{}, error) {
	var request dto.UpdateUserRequest

	if err := httputils.ParseRequestBody(c, &request, "UpdateUser", h.logger); err != nil {
		return nil, tvoerrors.ErrInvalidRequestData
	}

	ctx := httputils.CtxWithAuthToken(c)

	if (request.Phone == "" || request.Code == "") && request.Password == "" {
		log.Error("Request are empty", "error", tvoerrors.ErrInvalidRequestData)
		return nil, status.Error(codes.InvalidArgument, "Request are empty") //nolint
	}

	tokenData, err := httputils.GetTokenDataFromCtx(ctx)
	if err != nil {
		log.Error("Failed to get token data", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get claims from token") //nolint
	}

	if request.Phone != "" && request.Code != "" {
		if err = validator.ValidPhone(request.Phone); err != nil {
			log.Error("Invalid phone number", "error", tvoerrors.ErrInvalidPhone)
			return nil, status.Error(codes.InvalidArgument, "Invalid phone number") //nolint
		}

		exists, err := h.userRepository.PhoneExists(ctx, request.Phone)
		if err != nil {
			if !errors.Is(err, tvoerrors.ErrNotFound) {
				log.Error("Find phone error", "error", err)
				return nil, status.Error(codes.Internal, "something went wrong") //nolint
			}
		}

		if exists {
			log.Error("Phone exists", "phone", request.Phone, "error", ErrPhoneTaken)
			return nil, status.Error(codes.InvalidArgument, "phone already taken") //nolint
		}

		if err = h.userRepository.UpdatePhone(ctx, request.Phone, tokenData.UserID); err != nil {
			log.Error("Error updating phone", "error", err)
			return nil, status.Error(codes.Internal, "something went wrong") //nolint
		}
	}

	if request.Password != "" {
		if err = h.userRepository.UpdatePassword(ctx, tokenData.UserID, request.Password); err != nil {
			log.Error("Error updating password", "error", err)
			return nil, status.Error(codes.Internal, "something went wrong") //nolint
		}

		if err = h.removeUserTokens(ctx, tokenData.UserID); err != nil {
			log.Error("Error removing tokens", "error", err)
			return nil, status.Error(codes.Internal, "something went wrong") //nolint
		}
	}

	return &dto.UpdateUserResponse{
		Message: "Update user successfully",
	}, nil
}

// ChangeRole прокси метод для отправки его в сервис IDM
// @Summary Change user's role
// @Description This endpoint allows admins to change the role of a user by providing the user's ID and the new role.
// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.ChangeRoleRequest true "Request body"
// @Success 200 {object} dto.ChangeRoleResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /idm/change_role [post]
func (h *AuthHandlers) ChangeRole(c *fiber.Ctx) (interface{}, error) {
	var req dto.ChangeRoleRequest

	if err := httputils.ParseRequestBody(c, &req, "ChangeRole", h.logger); err != nil {
		return nil, tvoerrors.ErrInvalidRequestData
	}
	ctx := httputils.CtxWithAuthToken(c)

	tokenData, err := httputils.GetTokenDataFromCtx(ctx)
	if err != nil {
		log.Error("Failed to get token data", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get claims from token") //nolint
	}

	if tokenData.UserID == req.UserID {
		log.Error("Failed to change role", "error", "can't change your role")
		return nil, status.Error(codes.InvalidArgument, "can't change your role") //nolint
	}

	if tokenData.UserRoleID != tvomodels.ADMIN {
		log.Error("Failed to change role", "error", ErrNotAdmin)
		return nil, status.Error(codes.PermissionDenied, "can't change role") //nolint
	}

	role, err := h.roleRepository.RoleByName(ctx, req.Role)
	if err != nil {
		log.Error("Error get role", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	if err = h.userRepository.ChangeRole(ctx, req.UserID, int64(role.ID)); err != nil {
		log.Error("Error filed to change role", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	if err = h.removeUserTokens(ctx, req.UserID); err != nil {
		log.Error("Error removing tokens", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	return &dto.ChangeRoleResponse{
		Message: "Role changed",
	}, nil
}

// ResetToken прокси метод для отправки его в сервис IDM
// @Summary Reset user tokens
// @Description This endpoint allows users (or admins) to reset their authentication tokens,
// effectively logging them out from all active sessions.
// @Tags User
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param request body dto.ResetTokenRequest true "Request body"
// @Success 200 {object} dto.ResetTokenResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /idm/reset_token [post]
func (h *AuthHandlers) ResetToken(c *fiber.Ctx) (interface{}, error) {
	var req dto.ResetTokenRequest

	if err := httputils.ParseRequestBody(c, &req, "ResetToken", h.logger); err != nil {
		return nil, tvoerrors.ErrInvalidRequestData
	}

	ctx := httputils.CtxWithAuthToken(c)
	var userId int64

	tokenData, err := httputils.GetTokenDataFromCtx(ctx)
	if err != nil {
		log.Error("Failed to get token data", "error", err)
		return nil, status.Error(codes.Internal, "Failed to get claims from token") //nolint
	}

	if tokenData.UserRoleID == tvomodels.ADMIN {
		userId = req.UserID
	} else {
		userId = tokenData.UserID
	}

	if err = h.tokenRepository.TokenReset(ctx, userId); err != nil {
		log.Error("Error resetting token", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	if err = h.removeUserTokens(ctx, userId); err != nil {
		log.Error("Error removing tokens", "error", err)
		return nil, status.Error(codes.Internal, "something went wrong") //nolint
	}

	return &dto.ResetTokenResponse{
		Message: "All tokens reset",
	}, nil

	//res, err := ResetToken(ctx, &idm_pb.ResetTokenRequest{
	//	UserId: req.UserID,
	//})

	//if err != nil {
	//	h.logger.Error("Idm.ResetToken error", "method", "ResetToken", "error", err)
	//	return nil, err //nolint
	//}
}

func (h *AuthHandlers) Ping(c *fiber.Ctx) (interface{}, error) {
	return &dto.PingResponse{
		Status:  true,
		Message: "pong",
	}, nil
}

// CheckToken checks if the provided token is valid
func (s *AuthHandlers) CheckToken(ctx context.Context, request *dto.CheckTokenRequest) (*dto.CheckTokenResponse, error) {
	token := request.Token
	if token == "" {
		return nil, status.Error(codes.InvalidArgument, "invalid request") //nolint
	}

	user := &models.User{}
	var err error
	var errorMessage string
	emptyCache := false

	// проверяем токен в кеше
	cacheData, err := s.cache.Get(ctx, token)
	if err == nil {
		// проверяем валидность данных в кеше
		if err := helpers.JsonDecode(cacheData, user); err != nil || user.ID == 0 {
			emptyCache = true
			log.Debug("invalid cache data for token", "token", token, "data", string(cacheData))
		}
	}
	if emptyCache || err != nil {
		log.Debug("check user data in database", "token", token)
		tokenTTL := s.jwt.GetTokenTTL()
		data := make([]byte, 0)

		// проверяем токен в БД
		userToken, err := s.tokenRepository.GetUserIdByToken(ctx, token)
		if err != nil {
			errorMessage = "invalid token"
			log.Error("get user by token error", "error", err, "token", token)
		} else {
			user, err = s.userRepository.UserById(ctx, userToken.UserID)
			if err != nil {
				errorMessage = "invalid token data"
				log.Error("get user by id error", "error", err, "user_id", userToken.UserID)
			} else {
				data = helpers.JsonEncode(user)
			}
		}

		// добавляем токен в кеш (если валидный - кладем информацию о пользователе, если невалидный - кладем null)
		if err = s.cache.Set(ctx, token, data, tokenTTL); err != nil {
			log.Error("Error store token", "error", err)
		}
	}

	if user.ID == 0 {
		// для случая, если пользователь не найден в БД
		if err := helpers.JsonDecode(cacheData, user); err != nil || user.ID == 0 {
			errorMessage = "invalid cache data"
			log.Error("can't decode cached data", "data", string(cacheData), "error", err)
		}
	}

	// обновляем дату последнего визита пользователя
	if errorMessage == "" && user.ID > 0 {
		if err := s.userRepository.UpdateLastVisit(ctx, user.ID); err != nil {
			log.Error("can't update user last visit date", "user_id", user.ID, "error", err)
		}
	}

	log.Debug("ChekToken success", "user", user, "error", errorMessage)

	// возвращаем ответ
	return &dto.CheckTokenResponse{
		IsValid: user.ID != 0,
		Error:   errorMessage,
		UserId:  user.ID,
		RoleId:  int64(user.RoleID),
		Phone:   user.Phone,
	}, nil
}

// removeUserTokens removes active tokens associated with the given user ID.
func (s *AuthHandlers) removeUserTokens(ctx context.Context, userId int64) error {
	tokens, err := s.tokenRepository.ActiveTokens(ctx, userId)
	if err != nil {
		return tvoerrors.Wrap("Error fetching active tokens", err)
	}

	if len(tokens) == 0 {
		return nil
	}

	tokensToDelete := make([]string, 0, len(tokens))
	for _, t := range tokens {
		tokensToDelete = append(tokensToDelete, t.Token)
	}

	_, err = s.cache.Del(ctx, tokensToDelete...)
	if err != nil {
		return tvoerrors.Wrap("Error remove token", err)
	}
	return nil
}
