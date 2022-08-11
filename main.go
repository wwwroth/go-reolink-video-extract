package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Cameras []struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginPayload struct {
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

func main() {
	file, _ := os.ReadFile("config-cameras.json")

	cameras := Cameras{}
	json.Unmarshal(file, &cameras)

	for i := 0; i < len(cameras); i++ {
		camera := cameras[i]

		loginPayload := LoginPayload{
			Cmd: "Login",
			Param: LoginPayloadParam{
				User: LoginPayloadUser{
					UserName: camera.Username,
					Password: camera.Password,
				},
			},
		}

		fmt.Println(loginPayload)

		json, err := json.Marshal(loginPayload)
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
		bodyText, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", bodyText)

	}
}