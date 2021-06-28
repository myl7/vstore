package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/myl7/vstore/pkg"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	r := gin.Default()

	pkg.MountMiddlewares(r)
	pkg.Route(r)

	err = r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
