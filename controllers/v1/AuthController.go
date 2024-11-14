package v1Ctl

import (
	"context"
	"strings"
	v1req "test-task/resources/request/v1"
	v1Service "test-task/services/v1"
	u "test-task/shared/common"
	"test-task/shared/log"
	msg "test-task/shared/utils/message"
	"test-task/shared/utils/middleware"

	"net/http"
	valid "test-task/validator"

	"github.com/gin-gonic/gin"
)

type AuthCtl struct {
	AuthService       v1Service.IAuthService
	APIValidator      valid.IAPIValidatorService
	MiddlewareService middleware.IMiddleware
}

// SignUp is made for signning up a user
// @router /api/v1/sign-up [post]
func (ac *AuthCtl) SignUp(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Auth Controller Called(SignUp).")
	var req v1req.SignUpRequest

	//decode the request body into struct and failed if any error occurs
	if err := c.BindJSON(&req); err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		u.Respond(c.Writer, http.StatusBadRequest, u.ResponseErrorWithCode(u.CodeBadRequest, msg.InvalidRequest))
		return
	}

	// Struct field validation
	if resp, ok := ac.APIValidator.ValidateStruct(req, "SignUpRequest"); !ok {
		log.GetLog().Info("ERROR : ", "Struct validation error")
		u.Respond(c.Writer, http.StatusBadRequest, u.ResponseErrorWithCode(u.CodeBadRequest, resp))
		return
	}

	//call service
	resp := ac.AuthService.SignUpUser(req)
	statusCode := u.GetHTTPStatusCode(resp["res_code"])

	//return response using api helper
	u.Respond(c.Writer, statusCode, resp)
}

// SignIn is made for signning in a user
// @router /api/v1/sign-in [post]
func (ac *AuthCtl) SignIn(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Auth Controller Called(SignIn).")
	var req v1req.SignInRequest

	//decode the request body into struct and failed if any error occurs
	if err := c.BindJSON(&req); err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		u.Respond(c.Writer, http.StatusBadRequest, u.ResponseErrorWithCode(u.CodeBadRequest, msg.InvalidRequest))
		return
	}

	// Struct field validation
	if resp, ok := ac.APIValidator.ValidateStruct(req, "SignInRequest"); !ok {
		log.GetLog().Info("ERROR : ", "Struct validation error")
		u.Respond(c.Writer, http.StatusBadRequest, u.ResponseErrorWithCode(u.CodeBadRequest, resp))
		return
	}

	//call service
	resp := ac.AuthService.SignInUser(req)
	statusCode := u.GetHTTPStatusCode(resp["res_code"])

	//return response using api helper
	u.Respond(c.Writer, statusCode, resp)

}

// GetProfile is made for getting user details
// @router /api/v1/user-profile [get]
func (ac *AuthCtl) GetProfile(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Auth Controller Called(GetProfile).")

	userData, err := middleware.GetUserDataFromToken(c)
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		u.Respond(c.Writer, http.StatusBadRequest, u.ResponseErrorWithCode(u.CodeBadRequest, msg.SomethingWrong))
		return
	}

	//call service
	resp := ac.AuthService.GetUserDetails(userData.Id)
	statusCode := u.GetHTTPStatusCode(resp["res_code"])

	//return response using api helper
	u.Respond(c.Writer, statusCode, resp)
}

// SignOut is made for singout a user - blacklist the token passed by user till its expiry
// @router /api/v1/sign-out [post]
func (ac *AuthCtl) SignOut(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Auth Controller Called(SignOut).")

	userData, err := middleware.GetUserDataFromToken(c)
	if err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		u.Respond(c.Writer, http.StatusBadRequest, u.ResponseErrorWithCode(u.CodeBadRequest, msg.SomethingWrong))
		return
	}

	//getting the expiry time of jwt
	var expiryTime int
	exp := c.Keys["exp"]
	if exp != nil {
		expiryTime = exp.(int)
	}

	//getting the token from header
	bearerToken := c.Request.Header.Get("Authorization")
	token := strings.Split(bearerToken, "Bearer ")
	if len(token) < 2 {
		u.Respond(c.Writer, http.StatusBadRequest, u.ResponseErrorWithCode(u.CodeBadRequest, msg.SomethingWrong))
		return
	}

	//call service
	resp := ac.AuthService.SignOutUser(context.Background(), userData.Id, expiryTime, token[1])

	//return response using api helper
	u.Respond(c.Writer, http.StatusNoContent, resp)
}

func (ac *AuthCtl) RefreshToken(c *gin.Context) {
	log.GetLog().Info("INFO : ", "Auth Controller Called(RefreshToken).")

	// Declare the request struct to bind JSON body
	var req v1req.RefreshTokenRequest

	// Decode the request body into the struct and fail if any error occurs
	if err := c.BindJSON(&req); err != nil {
		log.GetLog().Info("ERROR : ", err.Error())
		u.Respond(c.Writer, http.StatusBadRequest, u.ResponseErrorWithCode(u.CodeBadRequest, msg.InvalidRequest))
		return
	}

	// Validate the struct fields (e.g., check for empty fields)
	if resp, ok := ac.APIValidator.ValidateStruct(req, "RefreshTokenRequest"); !ok {
		log.GetLog().Info("ERROR : ", "Struct validation error")
		u.Respond(c.Writer, http.StatusBadRequest, u.ResponseErrorWithCode(u.CodeBadRequest, resp))
		return
	}

	// Call the service to handle the token refresh logic
	resp := ac.AuthService.RefreshToken(req)
	statusCode := u.GetHTTPStatusCode(resp["res_code"])

	// Return the response using api helper
	u.Respond(c.Writer, statusCode, resp)
}
