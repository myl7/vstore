package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/myl7/vstore/pkg"
	"github.com/myl7/vstore/pkg/dao"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	defer dao.CloseConn()

	r := gin.Default()

	pkg.Route(r)

	err = r.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
