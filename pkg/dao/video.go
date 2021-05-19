package dao

import (
	"context"
	"github.com/jackc/pgx/v4"
	"io"
)

type VideoMeta struct {
	Vid         int
	Source      string
	Title       string
	Description string
}

func (v *VideoMeta) Get(vid int) error {
	return GetConn().QueryRow(context.Background(), `
		select vid, s.name, title, description from videos join sources s on s.sid = videos.sid where vid = $1
	`, vid).Scan(&v.Vid, &v.Source, &v.Title, &v.Description)
}

func (v VideoMeta) Put() error {
	_, err := GetConn().Exec(context.Background(), `
		update videos set title = $2, description = $3 where vid = $1
	`, v.Vid, v.Title, v.Description)
	return err
}

type VideoStream struct {
	tx  pgx.Tx
	oid uint32
}

func (v *VideoStream) Get(vid int) (io.ReadSeekCloser, error) {
	if v.tx == nil {
		var err error
		v.tx, err = GetConn().Begin(context.Background())
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
	if v.tx != nil {
		return v.tx.Commit(context.Background())
	}
	return nil
}
