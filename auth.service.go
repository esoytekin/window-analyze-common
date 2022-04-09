package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

// AuthenticationService service for authentication
type AuthenticationService interface {
	GetUserID(ctx *gin.Context) string
	GetAuthorizationHeader(ctx *gin.Context) string
	HasRole(ctx *gin.Context, role Role) bool
	HasPermission(ctx *gin.Context, roles []Role) bool
	GetUser(ctx *gin.Context) (*User, error)
}

var authService AuthenticationService

func NewAuthService(authApi string) AuthenticationService {
	if authService == nil {
		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			logrus.Debug("returning new instance")
		}

		authService = &AuthenticationServiceImpl{authAPI: authApi}

	}

	return authService
}

func NewAuthServiceWithF(getUser func(userID string) (*User, error)) AuthenticationService {

	if authService == nil {
		if logrus.IsLevelEnabled(logrus.DebugLevel) {
			logrus.Debug("returning new instance")
		}

		authService = &AuthenticationServiceImpl{getUserF: getUser}

	}

	return authService
}

type AuthenticationServiceImpl struct {
	authAPI  string
	getUserF func(userID string) (*User, error)
}

func (s *AuthenticationServiceImpl) GetUserID(ctx *gin.Context) string {
	return GetUserIDFromContext(ctx)
}

func (s *AuthenticationServiceImpl) GetAuthorizationHeader(ctx *gin.Context) string {
	return ctx.GetHeader("Authorization")
}

func (s *AuthenticationServiceImpl) HasRole(ctx *gin.Context, role Role) bool {
	user, err := s.GetUser(ctx)

	if err != nil {
		return false
	}

	for _, r := range user.Role {
		if r == string(role) {
			return true
		}
	}

	return false
}

func (s *AuthenticationServiceImpl) HasPermission(ctx *gin.Context, roles []Role) bool {
	user, err := s.GetUser(ctx)

	if err != nil {
		logrus.Error(err.Error())
		return false
	}

	for _, role := range user.Role {
		for _, r := range roles {
			if string(r) == role {
				return true
			}
		}
	}

	return false
}

func (s *AuthenticationServiceImpl) GetUser(ctx *gin.Context) (*User, error) {
	userID := s.GetUserID(ctx)
	if s.getUserF != nil {
		return s.getUserF(userID)
	}
	authHeader := s.GetAuthorizationHeader(ctx)
	return getUser(authHeader, userID, s.authAPI)
}

func getUser(authHeader, userID, authApi string) (*User, error) {

	var user User

	var client http.Client
	authAPIAddr := fmt.Sprintf("%s/user/%s", authApi, userID)
	r, _ := http.NewRequest("GET", authAPIAddr, nil)
	r.Header.Set("Authorization", authHeader)
	response, err := client.Do(r)

	if err != nil {
		return nil, err
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != 200 {
		return nil, errors.New(string(responseData))
	}
	err = json.Unmarshal(responseData, &user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
