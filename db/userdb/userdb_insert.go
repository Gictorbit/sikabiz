package userdb

import (
	"context"
)

func (udb *UserDB) InsertUserData(ctx context.Context, user User) error {
	tx, err := udb.GetPgConn().Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Insert user into the users table
	_, err = tx.Exec(ctx, `
		INSERT INTO users (id, name, email, phone_number)
		VALUES ($1, $2, $3, $4)`, user.ID, user.Name, user.Email, user.PhoneNumber)
	if err != nil {
		return err
	}

	// Insert addresses into the addresses and user_address tables
	for _, address := range user.Addresses {
		var addressID int64
		err = tx.QueryRow(ctx, `
			INSERT INTO addresses (street, city, state, zip_code, country)
			VALUES ($1, $2, $3, $4, $5) RETURNING id`,
			address.Street, address.City, address.State, address.ZipCode, address.Country).Scan(&addressID)
		if err != nil {
			return err
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO user_address (user_id, address_id)
			VALUES ($1, $2)`, user.ID, addressID)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}
