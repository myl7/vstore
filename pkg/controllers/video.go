package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/myl7/vstore/pkg/dao"
	"io"
	"net/http"
)

func GetVideoMeta(vid int, c *gin.Context) {
	v := dao.VideoMeta{}
	err := v.Get(vid)
	if err != nil {
		c.String(http.StatusNotFound, "Getting video meta failed")
	} else {
		c.JSON(http.StatusOK, v)
	}
}

func GetVideoStream(vid int, c *gin.Context) {
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
		Title:       body.Title,
		Description: body.Description,
		File:        f,
	}
	err = v.Add()
	if err != nil {
		c.String(http.StatusBadRequest, errMsg)
	} else {
		c.JSON(http.StatusOK, gin.H{"res": v.Vid})
	}
}
