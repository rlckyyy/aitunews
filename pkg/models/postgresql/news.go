package postgresql

import (
	"database/sql"
	"errors"
	_ "github.com/lib/pq"
	"relucky.net/aitunews/pkg/models"
)

type NewsModel struct {
	DB *sql.DB
}

func (m *NewsModel) Insert(title, content, category string) (int, error) {
	stmt := `
        INSERT INTO news (title, content, category, created)
        VALUES($1, $2, $3, CURRENT_TIMESTAMP)
        RETURNING id`

	var id int
	err := m.DB.QueryRow(stmt, title, content, category).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (m *NewsModel) Get(id int) (*models.News, error) {
	stmt := `SELECT id, title, content, category, created FROM news
	WHERE id = $1`

	row := m.DB.QueryRow(stmt, id)

	n := &models.News{}

	err := row.Scan(&n.ID, &n.Title, &n.Content, &n.Category, &n.Created)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return n, nil
}

func (m *NewsModel) Latest() ([]*models.News, error) {
	stmt := `SELECT id, title, content, category, created FROM news
		ORDER BY created DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newsList := []*models.News{}

	for rows.Next() {
		n := &models.News{}

		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Category, &n.Created)
		if err != nil {
			return nil, err
		}

		newsList = append(newsList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return newsList, nil
}

func (m *NewsModel) GetByCategory(category string) ([]*models.News, error) {
	stmt := `
        SELECT id, title, content, category, created FROM news
        WHERE category = $1
        ORDER BY created DESC`

	rows, err := m.DB.Query(stmt, category)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	newsList := []*models.News{}

	for rows.Next() {
		n := &models.News{}

		err := rows.Scan(&n.ID, &n.Title, &n.Content, &n.Category, &n.Created)
		if err != nil {
			return nil, err
		}

		newsList = append(newsList, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return newsList, nil
}

func (m *NewsModel) Delete(id int) error {
	stmt := `DELETE FROM news WHERE id = $1`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return err
	}

	return nil
}
