package dao

import "context"

type Source struct {
	Sid  int    `json:"sid"`
	Name string `json:"name"`
}

func ListSources() ([]Source, error) {
	rows, err := GetConn().Query(context.Background(), `
		select sid, name from sources order by name
	`)
	if err != nil {
		return nil, err
	}
	res := make([]Source, 0)
	var s Source
	for rows.Next() {
		err := rows.Scan(&s.Sid, &s.Name)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}
