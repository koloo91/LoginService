package repository

import (
	"context"
	"database/sql"
	"github.com/koloo91/loginservice/model"
	"strings"
	"time"
)

func CreateUser(ctx context.Context, db *sql.DB, user *model.User) error {
	statement, err := db.Prepare("INSERT INTO app_user(id, name, password_hash) VALUES($1, $2, $3)")
	if err != nil {
		return err
	}

	if _, err := statement.ExecContext(ctx, user.Id, user.Name, user.PasswordHash); err != nil {
		return err
	}

	return nil
}

func GetUserById(ctx context.Context, db *sql.DB, userId string) (*model.User, error) {
	row := db.QueryRowContext(ctx, "SELECT id, name, password_hash, created, updated FROM app_user WHERE id = $1", userId)

	var id, name, passwordHash string
	var created, updated time.Time

	if err := row.Scan(&id, &name, &passwordHash, &created, &updated); err != nil {
		return nil, err
	}

	return &model.User{
		Id:           id,
		Name:         name,
		PasswordHash: passwordHash,
		Created:      created,
		Updated:      updated,
	}, nil
}

func GetUserByName(ctx context.Context, db *sql.DB, userName string) (*model.User, error) {
	row := db.QueryRowContext(ctx, "SELECT id, name, password_hash, created, updated FROM app_user WHERE name = $1", strings.ToLower(userName))

	var id, name, passwordHash string
	var created, updated time.Time

	if err := row.Scan(&id, &name, &passwordHash, &created, &updated); err != nil {
		return nil, err
	}

	return &model.User{
		Id:           id,
		Name:         name,
		PasswordHash: passwordHash,
		Created:      created,
		Updated:      updated,
	}, nil

}
