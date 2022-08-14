package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Cameras []Camera

type Camera struct {
	IP       string `json:"ip"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Cmd   string            `json:"cmd"`
	Param LoginPayloadParam `json:"param"`
}

type LoginPayloadParam struct {
	User LoginPayloadUser `json:"User"`
}

type LoginPayloadUser struct {
	UserName string `json:"userName"`
	Password string `json:"password"`
}

type LoginResponse []struct {
	Cmd   string `json:"cmd"`
	Code  int    `json:"code"`
	Value struct {
		Token struct {
			LeaseTime int    `json:"leaseTime"`
			Name      string `json:"name"`
		} `json:"Token"`
	} `json:"value"`
}

func (payload LoginRequest) getToken(camera Camera) (string, error) {

	jsonPayload, err := json.Marshal(payload)
	requestJsonString := fmt.Sprint("[", string(jsonPayload), "]")

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

	loginResponse := LoginResponse{}
	json.Unmarshal(bodyText, &loginResponse)

	if loginResponse[0].Code > 0 {
		console(fmt.Sprintf("Camera on IP %s: login failure", camera.IP))
		return "", errors.New("failed to retrieve token")
	}

	token := loginResponse[0].Value.Token.Name
	console(fmt.Sprintf("Camera on IP %s: login success, token fetched, %s", camera.IP, token))

	return token, nil
}

func console(text string) bool {
	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%s | %s \n", t, text)
	return true
}

func main() {
	file, _ := os.ReadFile("config-cameras.json")

	cameras := Cameras{}
	json.Unmarshal(file, &cameras)

	console(fmt.Sprintf("Found %d cameras in configuration.", len(cameras)))

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

		loginRequest.getToken(camera)
	}
}
