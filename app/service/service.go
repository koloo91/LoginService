package service

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/koloo91/loginservice/app/model"
	"github.com/koloo91/loginservice/app/repository"
	"github.com/koloo91/loginservice/app/security"
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

func Login(ctx context.Context, db *sql.DB, jwtKey []byte, loginVo *model.LoginVo) (string, error) {
	user, err := repository.GetUserByName(ctx, db, loginVo.Name)
	if err != nil {
		log.Printf("error getting user by name '%s'", err.Error())
		return "", err
	}

	if err := checkPasswordHash(loginVo.Password, user.PasswordHash); err != nil {
		log.Printf("error checking passwords '%s'", err.Error())
		return "", fmt.Errorf("invalid credentials")
	}

	claims := &security.UserClaim{
		Id:      user.Id,
		Name:    user.Name,
		Created: user.Created,
		Updated: user.Updated,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		log.Printf("error signing token '%s'", err.Error())
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

func GetUserById(ctx context.Context, db *sql.DB, id string) (*model.UserVo, error) {
	foundUser, err := repository.GetUserById(ctx, db, id)
	if err != nil {
		log.Printf("error hashing password '%s'", err.Error())
		return nil, err
	}

	return model.UserToVo(foundUser), nil
}
