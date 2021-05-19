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
