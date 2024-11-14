package v1Service

import (
	"context"
	"errors"
	"fmt"
	"test-task/model"
	v1repo "test-task/repository/v1"
	v1req "test-task/resources/request/v1"
	v1resp "test-task/resources/response/v1"
	"test-task/shared/cache"
	u "test-task/shared/common"
	"test-task/shared/database"
	"test-task/shared/log"
	"test-task/shared/utils"
	"test-task/shared/utils/crypto"
	msg "test-task/shared/utils/message"
	"test-task/shared/utils/middleware"
	"time"

	"net/http"

	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type IAuthService interface {
	SignUpUser(req v1req.SignUpRequest) map[string]interface{}
	SignInUser(req v1req.SignInRequest) map[string]interface{}
	GetUserDetails(userId uuid.UUID) map[string]interface{}
	SignOutUser(ctx context.Context, userID uuid.UUID, expiry int, token string) map[string]interface{}
	RefreshToken(req v1req.RefreshTokenRequest) map[string]interface{}
}

type AuthService struct {
	UserRepo  v1repo.IUserRepository
	TokenRepo v1repo.ITokenRepository
}

func NewAuthService() IAuthService {
	userRepo := v1repo.NewUserWriter()
	tokenRepo := v1repo.NewTokenWriter()
	return &AuthService{
		UserRepo:  userRepo,
		TokenRepo: tokenRepo,
	}
}

func (as *AuthService) SignUpUser(req v1req.SignUpRequest) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Auth Service Called(SignUpUser).")
	conn := database.NewConnection()
	var user model.User

	//adding the request data to user model
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Email = req.Email
	user.TimeStamp()
	user.Password = utils.HashedPassword(req.Password)

	user.ID = uuid.NewV1()

	existingUser, err := as.UserRepo.GetUserByEmail(conn, user.Email)
	if err != nil {
		return u.ResponseErrorWithCode(http.StatusBadRequest, msg.InvalidRequest)
	}

	if existingUser != nil && existingUser.ID != uuid.Nil {
		return u.ResponseErrorWithCode(http.StatusBadRequest, msg.EmailInUse)
	}

	//Call article repository
	err = as.UserRepo.CreateUser(conn, &user)
	if err != nil {
		log.GetLog().Info("ERROR(from repo) : ", err.Error())
		return u.ResponseErrorWithCode(http.StatusBadRequest, msg.InvalidRequest)
	}

	response := u.ResponseSuccessWithObj(msg.SignUpSuccess, nil)
	return response
}

func (as *AuthService) SignInUser(req v1req.SignInRequest) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Auth Service Called(SignInUser).")
	conn := database.NewConnection()

	// Step 1: Retrieve the user by email
	existingUser, err := as.UserRepo.GetUserByEmail(conn, req.Email)
	if err != nil {
		// Return an error if the user retrieval fails
		log.GetLog().Info("ERROR : ", err.Error())
		return u.ResponseErrorWithCode(http.StatusBadRequest, msg.InvalidRequest)
	}

	if existingUser == nil || existingUser.ID == uuid.Nil {
		// If user does not exist, return an error response
		log.GetLog().Info("WARN : ", "Email not found...")
		return u.ResponseErrorWithCode(http.StatusBadRequest, msg.EmailNotRegistered)
	}
	//  Compare the provided password with the hashed password in the database
	err = bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(req.Password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		// If passwords do not match, return an error
		log.GetLog().Error("ERROR : ", "Password mismatch")
		return u.ResponseErrorWithCode(http.StatusBadRequest, msg.InvalidPassword)
	}

	// Step 3: Generate an authentication token
	accessToken, err := crypto.GenerateAuthToken(existingUser.ID, existingUser.Email, existingUser.CreatedAt)
	if err != nil {
		// If user does not exist, return an error response
		log.GetLog().Info("ERROR Generating token : ", err.Error())
		return u.ResponseErrorWithCode(http.StatusBadRequest, msg.InternalServer)
	}

	refreshToken, err := middleware.GenerateRefreshToken(existingUser.ID)
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error generating refresh token")
		return u.ResponseErrorWithCode(http.StatusInternalServerError, msg.InternalServer)
	}

	refreshTokenData := &model.UserRefreshToken{
		UserID:       existingUser.ID,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	if err := as.TokenRepo.SaveRefreshToken(conn, refreshTokenData); err != nil {
		log.GetLog().Info("ERROR : ", "Error saving refresh token")
		return u.ResponseErrorWithCode(http.StatusInternalServerError, msg.InternalServer)
	}
	// Step 4: Prepare the response with user details and the generated token
	SignInResp := v1resp.SigninResponse{RefreshToken: refreshToken, AccessToken: accessToken}

	// Prepare the final response with success message and token
	response := u.ResponseSuccessWithObj(msg.SignInSuccess, SignInResp)

	return response
}

func (as *AuthService) GetUserDetails(userId uuid.UUID) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Auth Service Called(GetUserDetails).")
	conn := database.NewConnection()

	existingUser, err := as.UserRepo.GetUserById(conn, userId)
	if err != nil {
		// Return an error if the user retrieval fails
		log.GetLog().Info("ERROR : ", err.Error())
		return u.ResponseErrorWithCode(http.StatusBadRequest, msg.InvalidRequest)
	}

	if existingUser == nil || existingUser.ID == uuid.Nil {
		log.GetLog().Info("WARN : ", "User not found.")
		return u.ResponseErrorWithCode(http.StatusNotFound, msg.UserNotFound)
	}

	userData := v1resp.UserResponse{
		Id:        existingUser.ID,
		FirstName: existingUser.FirstName,
		LastName:  existingUser.LastName,
		Email:     existingUser.Email,
		CreatedAt: existingUser.CreatedAt,
	}

	return u.ResponseSuccessWithObj(msg.UserProfileFetched, userData)
}

func (as *AuthService) SignOutUser(ctx context.Context, userID uuid.UUID, expiry int, token string) map[string]interface{} {
	log.GetLog().Info("INFO : ", "Auth Service Called(SignOut).")

	userIdString := userID.String()
	redisKey := fmt.Sprintf("%s_%d", userIdString, expiry)
	UnixTime := time.Unix(int64(expiry), 0)
	expiryTime := time.Until(UnixTime)
	//setting the cache key-value
	err := cache.SetToken(ctx, redisKey, token, expiryTime)
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		return u.ResponseErrorWithCode(http.StatusInternalServerError, msg.InternalServer)
	}
	return u.ResponseSuccessWithCode("", nil)
}

func (as *AuthService) RefreshToken(req v1req.RefreshTokenRequest) map[string]interface{} {
	conn := database.NewConnection()

	// Validate the refresh token and extract the user ID
	userID, err := middleware.ValidateRefreshTokenAndExtractID(req.RefreshToken)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			log.GetLog().Info("ERROR : ", "Token expired")
			return u.ResponseErrorWithCode(http.StatusUnauthorized, msg.RefreshTokenExpired)
		}

		log.GetLog().Info("ERROR : ", "Invalid refresh token")
		return u.ResponseErrorWithCode(http.StatusUnauthorized, msg.InvalidRefreshToken)
	}

	refreshTokenDetails, err := as.TokenRepo.FindTokenData(conn, userID, req.RefreshToken)
	if err != nil || refreshTokenDetails == nil {
		log.GetLog().Info("ERROR : ", "Refresh token not found in database")
		return u.ResponseErrorWithCode(http.StatusBadRequest, msg.InvalidRefreshToken)
	}

	user, err := as.UserRepo.GetUserById(conn, userID)
	if err != nil || user == nil {
		log.GetLog().Info("ERROR : ", "User not found")
		return u.ResponseErrorWithCode(http.StatusUnauthorized, msg.UserNotFound)
	}

	accessToken, err := middleware.GenerateToken(middleware.UserTokenData{Id: userID, Email: user.Email, CreatedAt: user.CreatedAt})
	if err != nil {
		log.GetLog().Info("ERROR : ", "Error generating access token")
		return u.ResponseErrorWithCode(http.StatusInternalServerError, msg.InternalServer)
	}

	tokenResp := v1resp.RefreshTokenResponse{AccessToken: accessToken}
	return u.ResponseSuccessWithObj(msg.TokenRefreshSuccess, tokenResp)
}
