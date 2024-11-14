package middleware

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"test-task/shared/cache"
	"test-task/shared/config"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/satori/go.uuid"
)

type UserTokenData struct {
	Id        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}
type IMiddleware interface {
	AuthHandler() gin.HandlerFunc
}

// Middleware is
type Middleware struct {
	Config config.IConfig
}

var AccessTokenKey string
var RefreshTokenKey string

func NewMiddlewareService(cf config.IConfig) IMiddleware {
	AccessTokenKey = cf.App().AccessTokenKey
	RefreshTokenKey = cf.App().RefreshTokenKey
	return &Middleware{
		Config: cf,
	}
}

// AuthHandler is used for User authentication
func (m *Middleware) AuthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		bearerToken := c.Request.Header.Get("Authorization")

		if !strings.HasPrefix(bearerToken, "Bearer ") {
			c.JSON(401, gin.H{"message": "Your request is not authorized", "status": http.StatusUnauthorized})
			c.Abort()
			return
		}

		token := strings.Split(bearerToken, "Bearer ")
		if len(token) < 2 {
			c.JSON(401, gin.H{"message": "An authorization token was not supplied", "status": http.StatusUnauthorized})
			c.Abort()
			return
		}

		// Validate token
		valid, err := ValidateToken(token[1], m.Config.App().AccessTokenKey)
		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				c.JSON(401, gin.H{"message": "The authorization token is expired", "status": http.StatusUnauthorized})
				c.Abort()
				return
			} else {
				c.JSON(401, gin.H{"message": "invalid authorization token", "status": http.StatusUnauthorized})
				c.Abort()
				return
			}
		}
		c.Set("userData", valid.Claims.(jwt.MapClaims)["userData"])

		var expTime int
		expInfo := valid.Claims.(jwt.MapClaims)["exp"]
		if expInfo != nil {
			val := expInfo.(float64)
			expTime = int(val)
		}
		c.Set("exp", expTime)
		userObject, err := GetUserDataFromToken(c)
		if err != nil {
			c.JSON(401, gin.H{"message": "something went wrong", "status": http.StatusUnauthorized})
			c.Abort()
			return
		}
		cacheKey := fmt.Sprintf("%s_%d", userObject.Id.String(), expTime)

		value, err := cache.GetToken(context.TODO(), cacheKey)
		if err != nil {
			c.JSON(500, gin.H{"message": "something went wrong", "status": http.StatusInternalServerError})
			c.Abort()
			return
		}

		if token[1] == value {

			c.JSON(401, gin.H{"message": "invalid authorization token", "status": http.StatusUnauthorized})
			c.Abort()
			return
		}

		c.Next()
	}
}

func GenerateToken(userData interface{}) (string, error) {
	// Create the token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set some claims
	claims := make(jwt.MapClaims)
	claims["userData"] = userData
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	token.Claims = claims
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(AccessTokenKey))
	return tokenString, err
}

func GenerateRefreshToken(id uuid.UUID) (string, error) {
	// Create the refresh token
	token := jwt.New(jwt.SigningMethodHS256)
	// Set the claims for the refresh token with just the id
	claims := make(jwt.MapClaims)
	claims["id"] = id.String() // Only include id in the claims
	// Set an expiration for the refresh token (e.g., 7 day)
	// claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	token.Claims = claims

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(RefreshTokenKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func ValidateRefreshTokenAndExtractID(t string) (uuid.UUID, error) {
	token, err := ValidateToken(t, RefreshTokenKey)
	if err != nil {
		return uuid.Nil, err
	}

	// Extract the id claim from the token
	claims := token.Claims.(jwt.MapClaims)
	idStr, ok := claims["id"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("id not found in the token")
	}

	// Convert id from string to uuid
	id, err := uuid.FromString(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid id format: %v", err)
	}

	// Return the id
	return id, nil
}

func ValidateToken(t string, k string) (*jwt.Token, error) {
	token, err := jwt.Parse(t, func(token *jwt.Token) (interface{}, error) {
		return []byte(k), nil
	})

	return token, err
}

// GetUserDataFromToken extracts the user's UUID from the token's payload
func GetUserDataFromToken(c *gin.Context) (UserTokenData, error) {
	var userData UserTokenData
	userInfo, userExists := c.Get("userData")
	if !userExists {
		return UserTokenData{}, nil
	}
	data := userInfo.(map[string]interface{})

	userId, err := uuid.FromString(data["id"].(string))
	if err != nil {
		return UserTokenData{}, err
	}
	Email := data["email"].(string)
	userData.Id = userId
	userData.Email = Email
	return userData, nil
}
