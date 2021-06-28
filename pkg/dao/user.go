package dao

import "context"

type User struct {
	Uid   int
	Token string
	Name  string
}

func (u *User) Get(uid int) error {
	conn := GetConn()
	defer conn.Close(context.Background())
	return conn.QueryRow(context.Background(), `
		select uid, token, name from users where uid = $1
	`, uid).Scan(&u.Uid, &u.Token, &u.Name)
}

func (u *User) AddOrUpdateToken() error {
	conn := GetConn()
	defer conn.Close(context.Background())
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}

	ok := false
	defer func() {
		if ok {
			_ = tx.Commit(context.Background())
		} else {
			_ = tx.Rollback(context.Background())
		}
	}()

	err = tx.QueryRow(context.Background(), `
		select uid from users where name = $1
	`, u.Name).Scan(&u.Uid)
	if err == nil {
		_, err := tx.Exec(context.Background(), `
			update users set token = $2 where uid = $1
		`, u.Uid, u.Token)
		return err
	} else {
		return tx.QueryRow(context.Background(), `
			insert into users (token, name) values ($1, $2) returning uid
		`, u.Token, u.Name).Scan(&u.Uid)
	}
}
