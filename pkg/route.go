package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/myl7/vstore/pkg/controllers"
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

func Route(r *gin.Engine) {
	r.GET("/api/videos/:vid/meta", func(c *gin.Context) {
		vid, err := ensureParamInt(c.Param("vid"), c)
		if err != nil {
			return
		}

		controllers.GetVideoMeta(vid, c)
	})
	r.GET("/api/videos/:vid/stream", func(c *gin.Context) {
		vid, err := ensureParamInt(c.Param("vid"), c)
		if err != nil {
			return
		}

		controllers.GetVideoStream(vid, c)
	})
	r.GET("/auth/start", controllers.AuthStart)
	r.GET("/auth/cb", controllers.AuthCallback)
}
