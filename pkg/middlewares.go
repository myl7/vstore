package pkg

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

func MountMiddlewares(r *gin.Engine) {
	r.Use(cors.Default())

	s, err := redis.NewStore(10, "tcp", "localhost:6379", "", []byte(os.Getenv("SECRET")))
	if err != nil {
		log.Fatalln(err)
	}
	r.Use(sessions.Sessions("login", s))
}
