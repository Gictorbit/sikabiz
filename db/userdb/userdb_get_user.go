package userdb

import "context"

const getUserByIDQuery = `
	SELECT 
	    u.id,
	    u.name, 
	    u.email, 
	    u.phone_number,
	    a.street, 
	    a.city,
	    a.state,
	    a.zip_code,
	    a.country
	FROM users u
	LEFT JOIN user_address ua ON u.id = ua.user_id
	LEFT JOIN addresses a ON ua.address_id = a.id
	WHERE u.id = $1
`

func (udb *UserDB) GetUserById(ctx context.Context, userID string) (*User, error) {
	var user User

	rows, err := udb.GetPgConn().Query(ctx, getUserByIDQuery, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Iterate over the result set
	for rows.Next() {
		var address Address
		if err := rows.Scan(
			&user.ID, &user.Name, &user.Email, &user.PhoneNumber,
			&address.Street, &address.City, &address.State, &address.ZipCode, &address.Country,
		); err != nil {
			return nil, err
		}

		// Append the address to the user's addresses slice
		user.Addresses = append(user.Addresses, address)
	}

	// Check for any error that occurred during row iteration
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &user, nil
}
