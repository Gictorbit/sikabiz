package userdb

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"
)

type UserDBPgConn interface {
	GetPgConn() *pgxpool.Pool

	InsertUserData(ctx context.Context, user User) error
	GetUserById(ctx context.Context, userID string) (*User, error)
}

var _ UserDBPgConn = &UserDB{}

type UserDB struct {
	postgresConn *pgxpool.Pool
}

func (udb *UserDB) GetPgConn() *pgxpool.Pool {
	return udb.postgresConn
}

func NewUserDB(rawConn *pgxpool.Pool) *UserDB {
	return &UserDB{
		postgresConn: rawConn,
	}
}

func ConnectToUserDB(databaseURL string) (*pgxpool.Pool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	rawConn, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		return nil, err
	}
	return rawConn, nil
}
