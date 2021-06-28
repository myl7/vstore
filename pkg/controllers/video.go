package controllers

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/myl7/vstore/pkg/dao"
	"io"
	"net/http"
	"time"
)

func GetVideoMeta(c *gin.Context) {
	vid, err := ensureParamInt(c.Param("vid"), c)
	if err != nil {
		return
	}

	v := dao.VideoMeta{}
	err = v.Get(vid)
	if err != nil {
		c.String(http.StatusNotFound, "Getting video meta failed")
	} else {
		c.JSON(http.StatusOK, v)
	}
}

func GetVideoStream(c *gin.Context) {
	vid, err := ensureParamInt(c.Param("vid"), c)
	if err != nil {
		return
	}

	v := dao.VideoStream{}
	f, err := v.Get(vid)
	if err != nil {
		c.String(http.StatusNotFound, "Getting video stream failed")
	} else {
		c.Stream(func(w io.Writer) bool {
			defer v.Close()
			defer f.Close()
			_, _ = io.Copy(w, f)
			return false
		})
	}
}

func ListSources(c *gin.Context) {
	res, err := dao.ListSources()
	if err != nil {
		c.String(http.StatusInternalServerError, "")
	} else {
		c.JSON(http.StatusOK, gin.H{"res": res})
	}
}

func AddVideo(c *gin.Context) {
	s := sessions.Default(c)
	user, ok := requireLogin(&s)
	if !ok {
		c.String(http.StatusForbidden, "Login required")
		return
	}

	var body struct {
		Sid         int    `form:"sid"`
		Title       string `form:"title"`
		Description string `form:"description"`
	}
	errMsg := "Invalid data to create a video"
	err := c.ShouldBind(&body)
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		return
	}
	file, err := c.FormFile("file")
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		return
	}
	f, err := file.Open()
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		return
	}
	v := dao.VideoAdd{
		Sid:         body.Sid,
		Uid:         user.Uid,
		Title:       body.Title,
		Description: body.Description,
		File:        f,
	}
	err = v.Add()
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
	} else {
		c.JSON(http.StatusCreated, gin.H{"res": v.Vid})
	}
}

func GetVideoComments(c *gin.Context) {
	vid, err := ensureParamInt(c.Param("vid"), c)
	if err != nil {
		return
	}

	res, err := dao.ListCommentsByVideo(vid)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
	} else {
		c.JSON(http.StatusOK, gin.H{"res": res})
	}
}

func AddVideoComment(c *gin.Context) {
	s := sessions.Default(c)
	user, ok := requireLogin(&s)
	if !ok {
		c.String(http.StatusForbidden, "Login required")
		return
	}

	vid, err := ensureParamInt(c.Param("vid"), c)
	if err != nil {
		return
	}

	errMsg := "Invalid data to create a comment"

	var body struct {
		Text string `json:"text"`
	}
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
		return
	}

	m := dao.CommentAdd{
		Vid:  vid,
		Uid:  user.Uid,
		Text: body.Text,
		Time: time.Now().UTC(),
	}
	err = m.Add()
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
	} else {
		c.JSON(http.StatusCreated, gin.H{"res": m.Mid})
	}
}

func ListUserVideoMeta(c *gin.Context) {
	s := sessions.Default(c)
	user, ok := requireLogin(&s)
	if !ok {
		c.String(http.StatusForbidden, "Login required")
		return
	}

	res, err := dao.ListVideoMetaByUid(user.Uid)
	if err != nil {
		c.String(http.StatusInternalServerError, "")
	} else {
		c.JSON(http.StatusOK, gin.H{"res": res})
	}
}
