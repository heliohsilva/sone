package repository

import (
	"api/src/models"
	"database/sql"
)

type Posts struct {
	db *sql.DB
}

func NewPostRepository(db *sql.DB) *Posts { return &Posts{db} }

func (repository Posts) Create(post models.Post) (int64, error) {
	stmt, err := repository.db.Prepare("INSERT INTO posts (title, content, user_id) values (?, ?, ?)")

	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(post.Title, post.Content, post.UserID)

	if err != nil {
		return 0, err
	}

	lastInsertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastInsertedID, nil

}

func (repository Posts) FetchByID(postID int64) (models.Post, error) {
	row := repository.db.QueryRow(`
		SELECT p.*, u.nickname
		FROM posts p 
		JOIN users u ON u.id = p.user_id 
		WHERE p.id = ?
	`, postID)

	var post models.Post

	if err := row.Scan(
		&post.ID,
		&post.Title,
		&post.Content,
		&post.UserID,
		&post.Likes,
		&post.CreatedAt,
		&post.UserNick,
	); err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (repository Posts) GetPosts(userID int64) ([]models.Post, error) {
	rows, err := repository.db.Query(`
		SELECT DISTINCT p.*, u.nickname
		FROM posts p
		INNER JOIN users u ON u.id = p.user_id
		LEFT JOIN followers s ON p.user_id = s.user_id
		WHERE p.user_id = ? OR s.follower_id = ?
	`, userID, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserID,
			&post.Likes,
			&post.CreatedAt,
			&post.UserNick,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (repository Posts) UpdatePost(postID int64, post models.Post) error {
	stmt, err := repository.db.Prepare("UPDATE posts SET title=?, content=? where id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()
	if _, err = stmt.Exec(post.Title, post.Content, postID); err != nil {
		return err
	}

	return nil
}

func (repository Posts) DeletePost(postID int64) error {
	stmt, err := repository.db.Prepare("DELETE FROM posts WHERE id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.Exec(postID); err != nil {
		return err
	}

	return nil
}

func (repository Posts) FetchByUser(userID int64) ([]models.Post, error) {
	rows, err := repository.db.Query(`
		SELECT p.*, u.nickname FROM posts p
		JOIN users u ON u.id = p.user_id WHERE p.user_id = ?
	`, userID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {
		var post models.Post

		if err := rows.Scan(
			&post.ID,
			&post.Title,
			&post.Content,
			&post.UserID,
			&post.Likes,
			&post.CreatedAt,
			&post.UserNick,
		); err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (repository Posts) Like(postID int64) error {
	stmt, err := repository.db.Prepare("UPDATE posts SET likes = likes + 1 WHERE id = ?")
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.Exec(postID); err != nil {
		return err
	}

	return nil
}

func (repository Posts) Unlike(postID int64) error {
	stmt, err := repository.db.Prepare(`
	UPDATE posts SET likes = 
	CASE WHEN likes > 0 THEN likes - 1 
	ELSE 0 END
	WHERE id = ?`)
	if err != nil {
		return err
	}

	defer stmt.Close()

	if _, err = stmt.Exec(postID); err != nil {
		return err
	}

	return nil
}
