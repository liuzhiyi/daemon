package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	Username string
	Password string
}

func main() {
	user := User{
		Username: "admin",
		Password: "123456",
	}
	postToken(user)
	user.Password = "12354"
	postToken(user)
	getToken()
}

func postToken(user User) {
	if ss, err := json.Marshal(user); err != nil {
		panic(err.Error())
	} else {
		r := bytes.NewReader(ss)
		if rsp, err := http.Post("http://127.0.0.1:3000/token", "json", r); err == nil {
			processResult(rsp)
			rsp.Body.Close()
		} else {
			fmt.Println(err.Error())
		}
	}
}

func processResult(rsp *http.Response) {
	str, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println(string(str))
	}
}

func getToken() {
	rsp, err := http.Get("http://127.0.0.1:3000/token")
	if err == nil {
		processResult(rsp)
		rsp.Body.Close()
	}
}
