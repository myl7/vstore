package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/myl7/vstore/pkg/dao"
	"github.com/myl7/vstore/pkg/services"
	"net/http"
)

func AuthStart(c *gin.Context) {
	u, _ := services.GetGitHubAuthUrl()
	c.Redirect(302, u)
}

func AuthCallback(c *gin.Context) {
	state := c.Query("state")
	code := c.Query("code")
	if !services.ValidateLoginState(state) {
		c.String(http.StatusUnauthorized, "Auth failed. Please try again.")
		return
	}

	info, err := services.GetGitHubAuthInfo(code, state)
	if err != nil {
		c.String(http.StatusUnauthorized, "Auth failed. Please try again.")
		return
	}

	name, err := services.GetGitHubName(info.AccessToken)
	if err != nil {
		c.String(http.StatusUnauthorized, "Auth failed. Please try again.")
		return
	}

	var user dao.User
	user.Token = info.AccessToken
	user.Name = name
	err = user.Add()
	if err != nil {
		c.String(http.StatusUnauthorized, "Auth failed. Please try again.")
		return
	}

	c.Redirect(302, "/")
}
