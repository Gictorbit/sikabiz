package db

import (
	"github.com/gictorbit/sikabiz/db/userdb"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserDBConn interface {
	GetPgConn() *pgxpool.Pool
	userdb.UserDBPgConn
}

var _ UserDBConn = &UserDataBase{}

type UserDataBase struct {
	pgConn *pgxpool.Pool
	*userdb.UserDB
}

func (tdb *UserDataBase) GetPgConn() *pgxpool.Pool {
	return tdb.pgConn
}

func NewUserDB(pgURL string) (*UserDataBase, error) {
	userDbConn, err := userdb.ConnectToUserDB(pgURL)
	if err != nil {
		return nil, err
	}
	return &UserDataBase{
		UserDB: userdb.NewUserDB(userDbConn),
	}, nil
}
