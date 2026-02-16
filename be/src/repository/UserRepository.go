package repository

import (
	"api/src/models"
	"database/sql"
	"fmt"
)

type users struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *users {
	return &users{db}
}

func (repository users) CreateUser(user models.User) (int64, error) {
	stmt, err := repository.db.Prepare("INSERT INTO users (name, nickname, email, password) VALUES (?, ?, ?, ?)")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(user.Name, user.Nickname, user.Email, user.Password)

	if err != nil {
		return 0, err
	}

	insertedId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	return insertedId, nil
}

func (repository users) UpdateUser(user models.User) (int64, error) {
	stmt, err := repository.db.Prepare("UPDATE users SET name=?, nickname=?, email=? WHERE id=?")

	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(user.Name, user.Nickname, user.Email, user.ID)

	if err != nil {
		return 0, err
	}

	updatedRow, err := result.RowsAffected()

	return updatedRow, nil

}

func (repository users) GetUser(userID int) (models.User, error) {
	stmt, err := repository.db.Prepare("SELECT * FROM users WHERE id=?")

	if err != nil {
		return models.User{}, err
	}

	var user models.User
	err = stmt.QueryRow(userID).Scan(&user.ID, &user.Name, &user.Nickname, &user.Email, &user.Password, &user.CreatedAt)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repository users) GetUsers(name string) ([]models.User, error) {
	var rows *sql.Rows
	fmt.Print(name)

	if name != "" {
		stmt, err := repository.db.Prepare("SELECT * FROM users WHERE name=? OR nickname=?")

		if err != nil {
			return []models.User{}, err
		}

		rows, err = stmt.Query(name, name)

		if err != nil {
			return []models.User{}, err
		}
	} else {
		stmt, err := repository.db.Prepare("SELECT * FROM users")

		if err != nil {
			return []models.User{}, err
		}

		rows, err = stmt.Query()

		if err != nil {
			return []models.User{}, err
		}
	}

	var users1 []models.User

	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Name, &user.Nickname, &user.Email, &user.Password, &user.CreatedAt)

		if err != nil {
			return []models.User{}, nil
		}

		users1 = append(users1, user)
	}

	return users1, nil
}

func (repository users) DeleteUser(userID int) (int64, error) {
	stmt, err := repository.db.Prepare("DELETE FROM users WHERE id=?")

	if err != nil {
		return 0, err
	}

	result, err := stmt.Exec(userID)

	if err != nil {
		return 0, err
	}

	deletedUser, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return deletedUser, nil
}

func (repository users) GetByEmail(email string) (models.User, error) {
	query, err := repository.db.Prepare("SELECT id, password FROM users WHERE email=?")

	if err != nil {
		return models.User{}, err
	}

	var user models.User

	err = query.QueryRow(email).Scan(&user.ID, &user.Password)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (repository users) Follow(userID int64, followerID int64) error {
	stmt, err := repository.db.Prepare(
		"INSERT IGNORE INTO followers (user_id, follower_id) values (?, ?)",
	)

	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.Exec(userID, followerID); err != nil {
		return nil
	}

	return nil
}

func (repository users) UnfollowUser(userID int64, followerID int64) error {
	stmt, err := repository.db.Prepare(
		"DELETE FROM followers where user_id = ? and follower_id = ?",
	)

	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.Exec(userID, followerID); err != nil {
		return err
	}

	return nil
}

func (repository users) FetchFollowers(userID int64) ([]models.User, error) {
	rows, err := repository.db.Query(
		`
		SELECT u.id, u.name, u.nickname, u.email, u.createdAt
		FROM users u INNER JOIN followers s ON u.id = s.follower_id
		WHERE s.user_id=?
		`, userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User

		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repository users) FetchFollowing(userID int64) ([]models.User, error) {
	rows, err := repository.db.Query(
		`
		SELECT u.id, u.name, u.nickname, u.email, u.createdAt
		FROM users u INNER JOIN followers s ON u.id = s.user_id
		WHERE s.follower_id=?
		`, userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.User

	for rows.Next() {
		var user models.User
		if err = rows.Scan(
			&user.ID,
			&user.Name,
			&user.Nickname,
			&user.Email,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (repository users) FetchPassword(userID int64) (string, error) {
	row, err := repository.db.Query("select passworc from users where id = ?", userID)

	if err != nil {
		return "", err
	}

	defer row.Close()

	var user models.User

	if row.Next() {
		if err = row.Scan(&user.Password); err != nil {
			return "", err
		}
	}

	return user.Password, nil
}

func (repository users) UpdatePassword(userID int64, hashedPassword string) error {
	stmt, err := repository.db.Prepare("UPDATE users SET password = ? WHERE id = ?")

	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err = stmt.Exec(hashedPassword, userID); err != nil {
		return err
	}

	return nil
}
