/*
Copyright 2020 IBM All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package main

import (
	"asset-transfer-basic/controllers"
	_ "asset-transfer-basic/controllers"
	"asset-transfer-basic/middlewares"
	"asset-transfer-basic/models"

	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	log.Println("============ application-golang starts ============")

	err := models.InitDB() //
	if err != nil {
		fmt.Printf("init db failed,err:%v\n", err)
		return
	}
	log.Printf("============ database connected  ============")

	log.Println("============ application-golang ends ============")

	r := gin.Default()
	r.Use(middlewares.CORSMiddleware())
	r.Use(middlewares.TlsHandler())

	r.GET("/GetAllAssets", controllers.GetAllAssets)

	r.POST("/VerifyPath", controllers.VerifyPath)

	r.RunTLS(":8081", "./tlsCert/cert.pem", "./tlsCert/key.pem") //

}
