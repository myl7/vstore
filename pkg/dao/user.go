package dao

import "context"

type User struct {
	Uid   int
	Token string
	Name  string
}

func (u *User) Get(uid int) error {
	return GetConn().QueryRow(context.Background(), `
		select uid, token, name from users where uid = $1
	`, uid).Scan(&u.Uid, &u.Token, &u.Name)
}

func (u *User) Add() error {
	return GetConn().QueryRow(context.Background(), `
		insert into users (token, name) values ($1, $2) returning uid
	`, u.Token, u.Name).Scan(&u.Uid)
}
