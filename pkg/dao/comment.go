package dao

import (
	"context"
	"time"
)

type Comment struct {
	Mid      int       `json:"mid"`
	Vid      int       `json:"vid"`
	UserName string    `json:"user_name"`
	Text     string    `json:"text"`
	Time     time.Time `json:"time"`
}

func ListCommentsByVideo(vid int) ([]Comment, error) {
	conn := GetConn()
	defer conn.Close(context.Background())
	rows, err := conn.Query(context.Background(), `
		select mid, vid, u.name as user_name, text, time from comments join users u on comments.uid = u.uid where vid = $1
	`, vid)
	if err != nil {
		return nil, err
	}
	res := make([]Comment, 0)
	var m Comment
	for rows.Next() {
		err := rows.Scan(&m.Mid, &m.Vid, &m.UserName, &m.Text, &m.Time)
		if err != nil {
			return nil, err
		}
		res = append(res, m)
	}
	return res, nil
}

type CommentAdd struct {
	Mid  int
	Vid  int
	Uid  int
	Text string
	Time time.Time
}

func (m *CommentAdd) Add() error {
	conn := GetConn()
	defer conn.Close(context.Background())
	return conn.QueryRow(context.Background(), `
		insert into comments (vid, uid, text, time) values ($1, $2, $3, $4) returning mid
	`, m.Vid, m.Uid, m.Text, m.Time).Scan(&m.Mid)
}
