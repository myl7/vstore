package services

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/myl7/vstore/pkg/dao"
	"github.com/myl7/vstore/pkg/utils"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func GetGitHubAuthUrl() (string, error) {
	u, err := url.Parse("https://github.com/login/oauth/authorize")
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Add("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	state := utils.GetRandStr(16)
	err = dao.GetKV().Set(context.Background(), "login-state:"+state, "1", 30*time.Minute).Err()
	if err != nil {
		return "", err
	}
	q.Add("state", state)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func ValidateLoginState(state string) bool {
	err := dao.GetKV().Get(context.Background(), "login-state:"+state).Err()
	if err != nil {
		return false
	}
	return true
}

type GitHubAuthInfo struct {
	AccessToken string
	Scope       []string
	TokenType   string
}

func GetGitHubAuthInfo(code string, state string) (GitHubAuthInfo, error) {
	d := url.Values{}
	d.Add("client_id", os.Getenv("GITHUB_CLIENT_ID"))
	d.Add("client_secret", os.Getenv("GITHUB_CLIENT_SECRET"))
	d.Add("state", state)
	d.Add("code", code)
	u := "https://github.com/login/oauth/access_token"
	req, err := http.NewRequest("POST", u, strings.NewReader(d.Encode()))
	if err != nil {
		return GitHubAuthInfo{}, nil
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len(d.Encode())))

	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return GitHubAuthInfo{}, nil
	}
	defer res.Body.Close()
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return GitHubAuthInfo{}, nil
	}
	var body struct {
		AccessToken string `json:"access_token"`
		Scope       string
		TokenType   string `json:"token_type"`
	}
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return GitHubAuthInfo{}, nil
	}
	return GitHubAuthInfo{
		AccessToken: body.AccessToken,
		Scope:       strings.Split(body.Scope, ","),
		TokenType:   body.TokenType,
	}, nil
}

func GetGitHubName(token string) (string, error) {
	u := "https://api.github.com/user"
	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Authorization", "token "+token)

	c := http.Client{}
	res, err := c.Do(req)
	if err != nil {
		return "", err
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	var body map[string]json.RawMessage
	err = json.Unmarshal(bodyBytes, &body)
	if err != nil {
		return "", err
	}
	name, ok := body["login"]
	if !ok {
		return "", errors.New("no `login` field in the response")
	}
	return string(name), nil
}

func SetSessionUser(s *sessions.Session, user dao.User) error {
	(*s).Set("vstore-uid", user.Uid)
	(*s).Set("vstore-token", user.Token)
	err := (*s).Save()
	if err != nil {
		return err
	}
	return nil
}

func GetSessionUser(s *sessions.Session) (dao.User, bool) {
	uidVal := (*s).Get("vstore-uid")
	if uidVal == nil {
		return dao.User{}, false
	}
	uid := uidVal.(int)

	tokenVal := (*s).Get("vstore-token")
	if tokenVal == nil {
		return dao.User{}, false
	}
	token := tokenVal.(string)

	var user dao.User
	err := user.Get(uid)
	if err != nil {
		return dao.User{}, false
	}
	if token != user.Token {
		return dao.User{}, false
	}
	return user, true
}
