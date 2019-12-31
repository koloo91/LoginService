package controller

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/koloo91/loginservice/app/model"
	"github.com/koloo91/loginservice/app/security"
	"github.com/koloo91/loginservice/app/service"
	_ "github.com/koloo91/loginservice/docs"
	"github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"net/http"
)

func SetupRoutes(db *sql.DB, jwtKey []byte) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(unhandledErrorHandler())

	log.Println("Setting up routes")

	{
		apiGroup := router.Group("/api")
		apiGroup.POST("/register", register(db))
		apiGroup.POST("/login", login(db, jwtKey))
		apiGroup.GET("/profile", security.JwtMiddleware(jwtKey), profile())

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

		token, err := service.Login(ctx.Request.Context(), db, jwtKey, &loginVo)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, model.ErrorVo{Message: "invalid credentials"})
			return
		}

		ctx.JSON(http.StatusOK, model.LoginResultVo{Token: token, Type: "Bearer"})
	}
}

// Profile godoc
// @Summary Returns the profile of the logged in user
// @ID profile
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} model.UserClaim
// @Failure 401 {object} model.ErrorVo
// @Router /api/profile [get]
func profile() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		user := security.GetUserFromContext(ctx)
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
