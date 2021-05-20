package dao

import (
	"context"
	"time"
)

type Comment struct {
	Mid      int
	Vid      int
	UserName string
	Text     string
	Time     time.Time
}

func ListCommentsByVideo(vid int) ([]Comment, error) {
	rows, err := GetConn().Query(context.Background(), `
		select mid, vid, u.name, text, time from comments join users u on comments.uid = u.uid where vid = $1
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
	return GetConn().QueryRow(context.Background(), `
		insert into comments (vid, uid, text, time) values ($1, $2, $3, $4) returning mid
	`, m.Vid, m.Uid, m.Text, m.Time).Scan(&m.Mid)
}
