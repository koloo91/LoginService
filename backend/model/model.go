package model

import (
	"github.com/google/uuid"
	"time"
)

type ErrorVo struct {
	Message string `json:"message"`
}

type RefreshTokenVo struct {
	RefreshToken string `json:"refreshToken"`
}

type LoginVo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	Type         string
}

type LoginResultVo struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Type         string `json:"type"`
}

type RegisterVo struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

type User struct {
	Id           string
	Name         string
	PasswordHash string
	Created      time.Time
	Updated      time.Time
}

type UserVo struct {
	Id      string    `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

func NewUser(name, passwordHash string) *User {
	return &User{
		Id:           uuid.New().String(),
		Name:         name,
		PasswordHash: passwordHash,
		Created:      time.Now(),
		Updated:      time.Now(),
	}
}

func UserToVo(user *User) *UserVo {
	return &UserVo{
		Id:      user.Id,
		Name:    user.Name,
		Created: user.Created,
		Updated: user.Updated,
	}
}
