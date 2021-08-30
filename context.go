package common

import (
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// _userID constant for userID key
const (
	_userID = "userId"
)

// SaveUserIDToContext saves user id to context
func SaveUserIDToContext(jwtSecretPassword string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")

		headerParts := strings.Split(authHeader, "Bearer ")

		if len(headerParts) != 2 {
			log.Error("malformed authentication header")
			panic("malformed authentication header")
		}

		accessToken := headerParts[1]

		token, _ := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {

			SECRETKEY := jwtSecretPassword
			return []byte(SECRETKEY), nil
		})

		// 3. get user with id
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

			userID := claims["user_id"]

			ctx.Set(_userID, userID)
		}
	}

}

// GetUserIDFromContext returns user id value from context
func GetUserIDFromContext(ctx *gin.Context) string {
	userID, _ := ctx.Get(_userID)

	return userID.(string)

}
