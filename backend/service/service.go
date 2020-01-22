package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/koloo91/jwt-security"
	"github.com/koloo91/loginservice/model"
	"github.com/koloo91/loginservice/repository"
	"golang.org/x/crypto/bcrypt"
	"log"
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

	refreshTokenString, accessTokenString, err := jwtsecurity.GenerateTokenPair(user.Id, user.Name, user.Created, user.Updated, jwtKey)
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

func Refresh(ctx context.Context, db *sql.DB, jwtKey []byte, refreshTokenClaim jwtsecurity.RefreshTokenClaim) (*model.LoginResult, error) {
	user, err := repository.GetUserById(ctx, db, refreshTokenClaim.Id)
	if err != nil {
		log.Printf("error getting user by name '%s'", err.Error())
		return nil, err
	}

	refreshTokenString, accessTokenString, err := jwtsecurity.GenerateTokenPair(user.Id, user.Name, user.Created, user.Updated, jwtKey)
	if err != nil {
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
