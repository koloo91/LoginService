package security

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/koloo91/loginservice/backend/model"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	userKey = "user"
)

type UserClaim struct {
	jwt.StandardClaims
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func GetUserFromContext(ctx *gin.Context) UserClaim {
	value, _ := ctx.Get(userKey)
	userClaim := value.(UserClaim)
	return userClaim
}

func JwtMiddleware(jwtKey []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeaderValue := ctx.GetHeader("Authorization")
		if len(authorizationHeaderValue) == 0 {
			log.Println("missing authorization header")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorVo{Message: "missing authorization header"})
			return
		}

		userClaims := UserClaim{}

		tokenString := strings.ReplaceAll(authorizationHeaderValue, "Bearer ", "")
		token, err := jwt.ParseWithClaims(tokenString, &userClaims, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("unexpected signing method: %v", token.Header["alg"])
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}

			return jwtKey, nil
		})

		if err != nil {
			log.Println("error parsing token")
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, model.ErrorVo{Message: "unexpected error"})
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorVo{Message: "invalid token"})
			return
		}

		ctx.Set(userKey, userClaims)

		ctx.Next()
	}
}
