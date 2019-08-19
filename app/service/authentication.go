package service

import (
	"bitbucket.org/Koloo/lgn/app/log"
	"bitbucket.org/Koloo/lgn/app/model"
	"bitbucket.org/Koloo/lgn/app/repository"
	"context"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func Register(ctx context.Context, db *sql.DB, registerVo *model.RegisterVo) (*model.UserVo, error) {
	hashedPassword, err := hashPassword(registerVo.Password)
	if err != nil {
		log.Errorf("error hashing password '%s'", err.Error())
		return nil, err
	}

	user := model.NewUser(registerVo.Name, hashedPassword)
	if err := repository.CreateUser(ctx, db, user); err != nil {
		log.Errorf("error creating user '%s'", err.Error())
		return nil, err
	}

	return model.UserToVo(user), nil
}

func Login(ctx context.Context, db *sql.DB, jwtKey []byte, loginVo *model.LoginVo) (string, error) {
	user, err := repository.GetUserByName(ctx, db, loginVo.Name)
	if err != nil {
		log.Errorf("error getting user by name '%s'", err.Error())
		return "", err
	}

	if err := checkPasswordHash(loginVo.Password, user.PasswordHash); err != nil {
		log.Errorf("error checking passwords '%s'", err.Error())
		return "", fmt.Errorf("invalid credentials")
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &model.UserClaim{
		Id:   user.Id,
		Name: user.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Errorf("error signing token '%s'", err.Error())
		return "", err
	}

	return tokenString, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
