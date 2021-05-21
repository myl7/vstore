package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/myl7/vstore/pkg/dao"
	"github.com/myl7/vstore/pkg/services"
	"net/http"
)

func AuthStart(c *gin.Context) {
	s := sessions.Default(c)
	_, ok := requireLogin(&s)
	if ok {
		c.Redirect(302, "/")
		return
	}

	u, _ := services.GetGitHubAuthUrl()
	c.Redirect(302, u)
}

func AuthCallback(c *gin.Context) {
	errMsg := "Auth failed. Please try again."

	state := c.Query("state")
	code := c.Query("code")
	if !services.ValidateLoginState(state) {
		c.String(http.StatusUnauthorized, errMsg)
		return
	}

	info, err := services.GetGitHubAuthInfo(code, state)
	if err != nil {
		c.String(http.StatusUnauthorized, errMsg)
		return
	}

	name, err := services.GetGitHubName(info.AccessToken)
	if err != nil {
		c.String(http.StatusUnauthorized, errMsg)
		return
	}

	var user dao.User
	user.Token = info.AccessToken
	user.Name = name
	err = user.AddOrUpdateToken()
	if err != nil {
		c.String(http.StatusUnauthorized, errMsg)
		return
	}

	s := sessions.Default(c)
	err = services.SetSessionUser(&s, user)
	if err != nil {
		c.String(http.StatusUnauthorized, errMsg)
		return
	}

	c.Redirect(302, "/")
}
