package service

import (
	"bitbucket.org/Koloo/lgn/app/log"
	"bitbucket.org/Koloo/lgn/app/model"
	"bitbucket.org/Koloo/lgn/app/repository"
	"context"
	"database/sql"
)

func GetUserById(ctx context.Context, db *sql.DB, id string) (*model.UserVo, error) {
	foundUser, err := repository.GetUserById(ctx, db, id)
	if err != nil {
		log.Errorf("error hashing password '%s'", err.Error())
		return nil, err
	}

	return model.UserToVo(foundUser), nil
}
