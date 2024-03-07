package postgresql

import (
	"database/sql"
	"relucky.net/aitunews/pkg/models"
)

type CommentsModel struct {
	DB *sql.DB
}

func (m *CommentsModel) Insert(text string, userId, newsId int) error {
	stmt := `
        INSERT INTO comments (user_id, news_id, text)
        VALUES($1, $2, $3)
        RETURNING id`

	var id int
	err := m.DB.QueryRow(stmt, userId, newsId, text).Scan(&id)
	if err != nil {
		return err
	}
	return nil
}

func (m *NewsModel) GetComments(newsId int) ([]*models.Comments, error) {
	stmt := `SELECT id, user_id, text FROM comments WHERE news_id = $1`
	rows, err := m.DB.Query(stmt, newsId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	commentList := []*models.Comments{}
	for rows.Next() {
		n := &models.Comments{}
		err := rows.Scan(&n.ID, &n.UserId, &n.Text)
		if err != nil {
			return nil, err
		}
		commentList = append(commentList, n)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return commentList, nil
}

func (m *CommentsModel) Delete(commentId int) error {
	stmt := `DELETE FROM comments WHERE id = $1`
	_, err := m.DB.Exec(stmt, commentId)
	if err != nil {
		return err
	}
	return nil
}
func (m *CommentsModel) GetNewsId(commentId int) (int, error) {
	stmt := `SELECT news_id FROM comments WHERE id = $1`
	var newsId int
	err := m.DB.QueryRow(stmt, commentId).Scan(&newsId)
	if err != nil {
		return 0, err
	}
	return newsId, nil
}

func (m *CommentsModel) GetAuthorId(commentId int) (int, error) {
	stmt := `SELECT user_id FROM comments WHERE id = $1`
	var userId int
	err := m.DB.QueryRow(stmt, commentId).Scan(&userId)
	if err != nil {
		return 0, err
	}
	return userId, nil
}
