package main

import (
	log2 "log"
	"net/http"

	"github.com/gin-gonic/gin"

	"GoCampusLogin/entity"
	"GoCampusLogin/utils/httputil"
	"GoCampusLogin/utils/loginutil"
)

func main() {
	router := gin.Default()
	router.GET("/login", func(context *gin.Context) {
		httputil.RemoveAllCookie() // 清空所有 cookie
		username := context.Query("username")
		password := context.Query("password")
		if username == "" || password == "" || len(username) < 10 {
			log2.Printf("a bad request: %v\n", username)
			context.JSON(http.StatusBadRequest, &entity.Result{Code: -1, Msg: entity.FAIL, Data: nil})
			return
		}
		cookies, err := loginutil.Login(username, password)
		if err != nil {
			log2.Fatalf("login fail: %v", err)
			context.JSON(http.StatusBadRequest, &entity.Result{Code: -1, Msg: entity.FAIL, Data: nil})
			return
		}
		context.JSON(http.StatusOK, &entity.Result{Code: 1, Msg: entity.SUCCESS, Data: cookies})
	})
	router.Run(":9999")
}
