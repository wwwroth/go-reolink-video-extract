package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Cameras []Camera

type Camera struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Cmd   string `json:"cmd"`
	Param LoginPayloadParam `json:"param"`
}

type LoginPayloadParam struct {
	User LoginPayloadUser `json:"User"`
}

type LoginPayloadUser struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

func (payload LoginRequest) getToken(camera Camera) string {
	json, err := json.Marshal(payload)
	requestJsonString := fmt.Sprint("[", string(json), "]")

	urlString := fmt.Sprint("http://", camera.IP, "/cgi-bin/api.cgi?cmd=Login")

	client := &http.Client{}
	var data = strings.NewReader(requestJsonString)
	req, err := http.NewRequest("POST", urlString, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyText)
}

func main() {
	file, _ := os.ReadFile("config-cameras.json")

	cameras := Cameras{}
	json.Unmarshal(file, &cameras)

	for i := 0; i < len(cameras); i++ {
		camera := cameras[i]

		loginRequest := LoginRequest{
			Cmd: "Login",
			Param: LoginPayloadParam{
				User: LoginPayloadUser{
					UserName: camera.Username,
					Password: camera.Password,
				},
			},
		}

		token := loginRequest.getToken(camera)

		fmt.Println(token)
	}
}