package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/myl7/vstore/pkg/dao"
	"github.com/myl7/vstore/pkg/services"
	"net/http"
	"strconv"
)

func ensureParamInt(param string, c *gin.Context) (int, error) {
	res, err := strconv.Atoi(param)
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid video ID")
		return res, err
	}
	return res, nil
}

func requireLogin(s *sessions.Session) (dao.User, bool) {
	user, ok := services.GetSessionUser(s)
	if !ok {
		return dao.User{}, false
	}
	return user, true
}
