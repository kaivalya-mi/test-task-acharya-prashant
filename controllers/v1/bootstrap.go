package v1Ctl

import (
	v1Service "test-task/services/v1"
	"test-task/shared/utils/middleware"
	validator "test-task/validator"
)

func AuthController(validatorService validator.IAPIValidatorService, authService v1Service.IAuthService, middlwareService middleware.IMiddleware) *AuthCtl {
	authCtl := AuthCtl{
		AuthService:       authService,
		APIValidator:      validatorService,
		MiddlewareService: middlwareService,
	}

	return &authCtl
}
