package request

import (
	"fmt"
	"net/http"
	"os"
	"sf/src/config"
)

var client = &http.Client{}

func SetHeaders(req *http.Request, TestCookie bool) {
	config.Load()
	Header := make(map[string]string)
	if config.Var.Cookie == "" && TestCookie {
		fmt.Println("Cookie is empty, please login first!")
		os.Exit(1)
	} else {
		Header["Cookie"] = config.Var.Cookie
	}
	Header["sf-minip-info"] = "minip_novel/1.0.70(android;11)/wxmp"
	Header["Content-Type"] = "application/json"
	if config.Var.UserName != "" || config.Var.Password != "" {
		Header["test-sfacg"] = "cookie:" + config.Var.Cookie
		Header["Authorization"] = config.Var.Authorization + config.Var.UserName + "&" + config.Var.Password
		Header["account-sfacg"] = config.Var.UserName + "&" + config.Var.Password
	}
	for k, v := range Header {
		req.Header.Set(k, v)
	}
}
