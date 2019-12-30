package service

import (
	"context"
	"database/sql"
	"github.com/koloo91/loginservice/app/log"
	"github.com/koloo91/loginservice/app/model"
	"github.com/koloo91/loginservice/app/repository"
)

func GetUserById(ctx context.Context, db *sql.DB, id string) (*model.UserVo, error) {
	foundUser, err := repository.GetUserById(ctx, db, id)
	if err != nil {
		log.Errorf("error hashing password '%s'", err.Error())
		return nil, err
	}

	return model.UserToVo(foundUser), nil
}
