package controller

import (
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/koloo91/loginservice/backend/docs"
	"github.com/koloo91/loginservice/backend/model"
	"github.com/koloo91/loginservice/backend/security"
	"github.com/koloo91/loginservice/backend/service"
	"github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
	"time"
)

func SetupRoutes(db *sql.DB, jwtKey []byte, validateExpirationDate bool) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(unhandledErrorHandler())

	log.Println("Setting up routes")

	{
		apiGroup := router.Group("/api")
		apiGroup.POST("/token/refresh", refreshToken(db, jwtKey))
		apiGroup.POST("/register", register(db))
		apiGroup.POST("/login", login(db, jwtKey))
		apiGroup.GET("/profile", security.JwtMiddleware(jwtKey, validateExpirationDate), profile())

		apiGroup.GET("/alive", alive())
	}

	router.GET("/swagger/*any", ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "GIN_MODE"))

	return router
}

func unhandledErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Println(err)
				ctx.AbortWithStatusJSON(http.StatusBadRequest, model.ErrorVo{Message: "unexpected error"})
			}
		}()
		ctx.Next()
	}
}

// Refresh token godoc
// @Summary Refresh token
// @ID refresh_token
// @Accept json
// @Produce json
// @Param loginVo body model.RefreshTokenVo true "refresh token json"
// @Success 200 {object} model.LoginResultVo
// @Failure 400 {object} model.ErrorVo
// @Router /api/token/refresh [post]
func refreshToken(db *sql.DB, jwtKey []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var refreshTokenVo model.RefreshTokenVo
		if err := ctx.ShouldBindJSON(&refreshTokenVo); err != nil {
			ctx.JSON(http.StatusBadRequest, "Invalid json")
			return
		}

		refreshTokenClaim := security.RefreshTokenClaim{}

		token, err := jwt.ParseWithClaims(refreshTokenVo.RefreshToken, &refreshTokenClaim, func(token *jwt.Token) (interface{}, error) {
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

		if time.Unix(refreshTokenClaim.ExpiresAt, 0).Sub(time.Now()).Seconds() <= 0 {
			log.Println("token expired")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, model.ErrorVo{Message: "token expired"})
			return
		}

		loginResult, err := service.Refresh(ctx.Request.Context(), db, jwtKey, refreshTokenClaim)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.ErrorVo{Message: "invalid refresh token"})
			return
		}

		ctx.JSON(http.StatusOK, model.LoginResultVo{
			AccessToken:  loginResult.AccessToken,
			RefreshToken: loginResult.RefreshToken,
			Type:         "Bearer",
		})
	}
}

// Register godoc
// @Summary Registers a new user
// @ID register
// @Accept json
// @Produce json
// @Param registerVo body model.RegisterVo true "register json"
// @Success 200 {object} model.UserVo
// @Failure 400 {object} model.ErrorVo
// @Router /api/register [post]
func register(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var registerVo model.RegisterVo
		if err := ctx.ShouldBindJSON(&registerVo); err != nil {
			ctx.JSON(http.StatusBadRequest, "Invalid json")
			return
		}

		userVo, err := service.Register(ctx.Request.Context(), db, &registerVo)
		if err != nil {
			if err, ok := err.(*pq.Error); ok {
				ctx.JSON(http.StatusBadRequest, model.ErrorVo{Message: err.Message})
				return
			}
			ctx.JSON(http.StatusBadRequest, model.ErrorVo{Message: err.Error()})
			return
		}

		ctx.JSON(http.StatusCreated, userVo)
	}
}

// Login godoc
// @Summary Login a user
// @ID login
// @Accept json
// @Produce json
// @Param loginVo body model.LoginVo true "login json"
// @Success 200 {object} model.LoginResultVo
// @Failure 400 {object} model.ErrorVo
// @Router /api/login [post]
func login(db *sql.DB, jwtKey []byte) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var loginVo model.LoginVo
		if err := ctx.ShouldBindJSON(&loginVo); err != nil {
			ctx.JSON(http.StatusBadRequest, model.ErrorVo{Message: "invalid json"})
			return
		}

		loginResult, err := service.Login(ctx.Request.Context(), db, jwtKey, &loginVo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.ErrorVo{Message: "invalid credentials"})
			return
		}

		ctx.JSON(http.StatusOK, model.LoginResultVo{
			AccessToken:  loginResult.AccessToken,
			RefreshToken: loginResult.RefreshToken,
			Type:         "Bearer",
		})
	}
}

// Profile godoc
// @Summary Returns the profile of the logged in user
// @ID profile
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} security.AccessTokenClaim
// @Failure 401 {object} model.ErrorVo
// @Router /api/profile [get]
func profile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := security.GetAccessTokenFromContext(ctx)
		ctx.JSON(http.StatusOK, model.UserVo{
			Id:      user.Id,
			Name:    user.Name,
			Created: user.Created,
			Updated: user.Updated,
		})
	}
}

// Alive godoc
// @Summary Checks if the service is running
// @ID alive
// @Produce text/plain
// @Success 204 {string} string	""
// @Router /api/alive [get]
func alive() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusNoContent, "")
	}
}
