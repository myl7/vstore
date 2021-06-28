package dao

import (
	"context"
	"github.com/jackc/pgx/v4"
	"io"
)

type VideoMeta struct {
	Vid         int    `json:"vid"`
	Source      string `json:"source"`
	Title       string `json:"title"`
	Description string `json:"description"`
}

func (v *VideoMeta) Get(vid int) error {
	conn := GetConn()
	defer conn.Close(context.Background())
	return conn.QueryRow(context.Background(), `
		select vid, s.name, title, description from videos join sources s on s.sid = videos.sid where vid = $1
	`, vid).Scan(&v.Vid, &v.Source, &v.Title, &v.Description)
}

func (v VideoMeta) Put() error {
	conn := GetConn()
	defer conn.Close(context.Background())
	_, err := conn.Exec(context.Background(), `
		update videos set title = $2, description = $3 where vid = $1
	`, v.Vid, v.Title, v.Description)
	return err
}

type VideoStream struct {
	tx   pgx.Tx
	conn *pgx.Conn
	oid  uint32
}

func (v *VideoStream) Get(vid int) (io.ReadSeekCloser, error) {
	if v.tx == nil {
		var err error
		v.conn = GetConn()
		v.tx, err = v.conn.Begin(context.Background())
		if err != nil {
			return nil, err
		}
	}

	err := v.tx.QueryRow(context.Background(), `
		select file from videos where vid = $1
	`, vid).Scan(&v.oid)
	if err != nil {
		return nil, err
	}

	files := v.tx.LargeObjects()
	f, err := files.Open(context.Background(), v.oid, pgx.LargeObjectModeRead)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func (v VideoStream) Close() error {
	defer v.conn.Close(context.Background())
	if v.tx != nil {
		return v.tx.Commit(context.Background())
	}
	return nil
}

type VideoAdd struct {
	Vid         int
	Sid         int
	Uid         int
	Title       string
	Description string
	File        io.ReadCloser
}

func (v *VideoAdd) Add() error {
	defer v.File.Close()
	conn := GetConn()
	defer conn.Close(context.Background())
	tx, err := conn.Begin(context.Background())
	ok := false
	defer func() {
		if ok {
			_ = tx.Commit(context.Background())
		} else {
			_ = tx.Rollback(context.Background())
		}
	}()
	lo := tx.LargeObjects()
	oid, err := lo.Create(context.Background(), 0)
	if err != nil {
		return err
	}
	f, err := lo.Open(context.Background(), oid, pgx.LargeObjectModeWrite)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, v.File)
	if err != nil {
		return err
	}
	err = tx.QueryRow(context.Background(), `
		insert into videos (sid, uid, title, description, file) values ($1, $2, $3, $4, $5) returning vid
	`, v.Sid, v.Uid, v.Title, v.Description, oid).Scan(&v.Vid)
	if err != nil {
		return err
	}
	ok = true
	return nil
}

func ListVideoMetaByUid(uid int) ([]VideoMeta, error) {
	conn := GetConn()
	defer conn.Close(context.Background())
	rows, err := conn.Query(context.Background(), `
		select vid, s.name, title, description from videos join sources s on s.sid = videos.sid where uid = $1
	`, uid)
	if err != nil {
		return nil, err
	}
	res := make([]VideoMeta, 0)
	var v VideoMeta
	for rows.Next() {
		err := rows.Scan(&v.Vid, &v.Source, &v.Title, &v.Description)
		if err != nil {
			return nil, err
		}
		res = append(res, v)
	}
	return res, nil
}
