package oauth2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/kernle32dll/loginsrv/model"
)

var githubAPI = "https://api.github.com"

func init() {
	RegisterProvider(providerGithub)
}

// GithubUser is used for parsing the github response
type GithubUser struct {
	Login     string `json:"login,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
}

var providerGithub = Provider{
	Name:     "github",
	AuthURL:  "https://github.com/login/oauth/authorize",
	TokenURL: "https://github.com/login/oauth/access_token",
	GetUserInfo: func(token TokenInfo) (model.UserInfo, string, error) {
		gu := GithubUser{}
		url := githubAPI + "/user"
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "token "+token.AccessToken)
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return model.UserInfo{}, "", err
		}
		defer resp.Body.Close()

		if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
			return model.UserInfo{}, "", fmt.Errorf("wrong content-type on github get user info: %v", resp.Header.Get("Content-Type"))
		}

		if resp.StatusCode != 200 {
			return model.UserInfo{}, "", fmt.Errorf("got http status %v on github get user info", resp.StatusCode)
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return model.UserInfo{}, "", fmt.Errorf("error reading github get user info: %v", err)
		}

		err = json.Unmarshal(b, &gu)
		if err != nil {
			return model.UserInfo{}, "", fmt.Errorf("error parsing github get user info: %v", err)
		}

		return model.UserInfo{
			Sub:     gu.Login,
			Picture: gu.AvatarURL,
			Name:    gu.Name,
			Email:   gu.Email,
			Origin:  "github",
		}, string(b), nil
	},
}
