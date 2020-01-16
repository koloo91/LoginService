package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/koloo91/loginservice/model"
	"github.com/koloo91/loginservice/repository"
	"github.com/koloo91/loginservice/security"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

func Register(ctx context.Context, db *sql.DB, registerVo *model.RegisterVo) (*model.UserVo, error) {
	hashedPassword, err := hashPassword(registerVo.Password)
	if err != nil {

		log.Printf("error hashing password '%s'", err.Error())
		return nil, err
	}

	user := model.NewUser(registerVo.Name, hashedPassword)
	if err := repository.CreateUser(ctx, db, user); err != nil {
		log.Printf("error creating user '%s'", err.Error())
		return nil, err
	}

	return model.UserToVo(user), nil
}

func Login(ctx context.Context, db *sql.DB, jwtKey []byte, loginVo *model.LoginVo) (*model.LoginResult, error) {
	user, err := repository.GetUserByName(ctx, db, loginVo.Name)
	if err != nil {
		log.Printf("error getting user by name '%s'", err.Error())
		return nil, err
	}

	if err := checkPasswordHash(loginVo.Password, user.PasswordHash); err != nil {
		log.Printf("error checking passwords '%s'", err.Error())
		return nil, fmt.Errorf("invalid credentials")
	}

	refreshTokenClaims := &security.RefreshTokenClaim{
		Id: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		log.Printf("error signing refresh token '%s'", err.Error())
		return nil, err
	}

	accessTokenClaims := &security.AccessTokenClaim{
		Id:      user.Id,
		Name:    user.Name,
		Created: user.Created,
		Updated: user.Updated,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		log.Printf("error signing access token '%s'", err.Error())
		return nil, err
	}

	return &model.LoginResult{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		Type:         "Bearer",
	}, nil
}

func Refresh(ctx context.Context, db *sql.DB, jwtKey []byte, refreshTokenClaim security.RefreshTokenClaim) (*model.LoginResult, error) {
	user, err := repository.GetUserById(ctx, db, refreshTokenClaim.Id)
	if err != nil {
		log.Printf("error getting user by name '%s'", err.Error())
		return nil, err
	}

	refreshTokenClaims := &security.RefreshTokenClaim{
		Id: user.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshTokenClaims)
	refreshTokenString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		log.Printf("error signing refresh token '%s'", err.Error())
		return nil, err
	}

	accessTokenClaims := &security.AccessTokenClaim{
		Id:      user.Id,
		Name:    user.Name,
		Created: user.Created,
		Updated: user.Updated,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims)
	accessTokenString, err := accessToken.SignedString(jwtKey)
	if err != nil {
		log.Printf("error signing access token '%s'", err.Error())
		return nil, err
	}

	return &model.LoginResult{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
		Type:         "Bearer",
	}, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
