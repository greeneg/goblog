package model

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

type Blog struct {
	Id      int    `json:"id"`
	Author  string `json:"author"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Image   string `json:"image"`
	Ctime   string `json:"ctime"`
}

func ConnectDatabase() error {
	db, err := sql.Open("sqlite3", "./blogsDemo.db")
	if err != nil {
		return err
	}

	DB = db
	return nil
}

func GetBlogs() ([]Blog, error) {
	rows, err := DB.Query("SELECT * FROM blog")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	blogs := make([]Blog, 0)
	for rows.Next() {
		sBlog := Blog{}
		err = rows.Scan(
			&sBlog.Id,
			&sBlog.Author,
			&sBlog.Title,
			&sBlog.Content,
			&sBlog.Image,
			&sBlog.Ctime,
		)

		if err != nil {
			return nil, err
		}

		blogs = append(blogs, sBlog)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return blogs, nil
}

func GetBlogById(id string) (Blog, error) {
	rec, err := DB.Prepare("SELECT * FROM blog WHERE id = ?")
	if err != nil {
		return Blog{}, err
	}

	b := Blog{}
	qErr := rec.QueryRow(id).Scan(
		&b.Id,
		&b.Author,
		&b.Title,
		&b.Image,
		&b.Content,
		&b.Ctime,
	)
	if qErr != nil {
		if qErr == sql.ErrNoRows {
			return Blog{}, nil
		}
		return Blog{}, qErr
	}

	return b, nil
}

func AddBlog(n Blog) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	ctime := time.Now()
	q, err := t.Prepare("INSERT INTO blog (author, title, content, image, ctime) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return false, err
	}

	_, err = q.Exec(n.Author, n.Title, n.Content, n.Image, ctime)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}

func UpdateBlog(n Blog, id int) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := t.Prepare("UPDATE blog SET author = ?, title = ?, content = ?, image = ? WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer q.Close()

	_, err = q.Exec(n.Author, n.Title, n.Content, n.Image, id)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}

func DeleteBlog(id int) (bool, error) {
	t, err := DB.Begin()
	if err != nil {
		return false, err
	}

	q, err := DB.Prepare("DELETE FROM blog WHERE id = ?")
	if err != nil {
		return false, err
	}

	defer q.Close()

	_, err = q.Exec(id)
	if err != nil {
		return false, err
	}

	t.Commit()

	return true, nil
}
