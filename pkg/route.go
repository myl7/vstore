package pkg

import (
	"github.com/gin-gonic/gin"
	"github.com/myl7/vstore/pkg/controllers"
)

func Route(r *gin.Engine) {
	r.GET("/api/videos/:vid/meta", controllers.GetVideoMeta)
	r.GET("/api/videos/:vid/stream", controllers.GetVideoStream)
	r.GET("/api/videos/:vid/comments", controllers.GetVideoComments)
	r.POST("/api/videos/:vid/comments", controllers.AddVideoComment)
	r.GET("/api/sources", controllers.ListSources)
	r.GET("/api/videos", controllers.ListUserVideoMeta)
	r.POST("/api/videos", controllers.AddVideo)
	r.GET("/auth/start", controllers.AuthStart)
	r.GET("/auth/cb", controllers.AuthCallback)
}
