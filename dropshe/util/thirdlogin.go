package util

import (
	"fmt"
	"github.com/astaxie/beego"
	"net/http"
	"golang.org/x/oauth2"
	"io/ioutil"
	"encoding/json"

)

const htmlIndex = `<html><body>
<a href="/GoogleLogin">Log in with Google</a>
</body></html>
`

var endpotin = oauth2.Endpoint{
	AuthURL:  "https://accounts.google.com/o/oauth2/auth",
	TokenURL: "https://accounts.google.com/o/oauth2/token",
}

var googleOauthConfig = &oauth2.Config{
	ClientID:     "715700887174-af7eoh5ceg7qm9f4d57lnlg9asm8qar5.apps.googleusercontent.com",
	ClientSecret: "G2ahpc261Ru6AAahXfVRlr_v",
	RedirectURL:  "https://app.dropshe.com/order/drop",
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
		"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint: endpotin,
}

const oauthStateString = "random"

func ThirdpartyLogin() {

	http.HandleFunc("/login/oauth", HandleGoogleLogin)
	//http.HandleFunc("/GoogleCallback", HandleGoogleCallback)
	//fmt.Println(http.ListenAndServe(":8000", nil))
}



func HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	//拼接google登录url
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	fmt.Println(url)
	//307
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

//google账号登录完成后的处理
func HandleGoogleCallback(w http.ResponseWriter, r *http.Request)(map[string][]map[string]string) {
	state := r.FormValue("state")
	if state != oauthStateString {
		fmt.Printf("invalid oauth state, expected '%s', got '%s'\n", oauthStateString, state)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}
	fmt.Println(state)
	//申请授权码
	code := r.FormValue("code")
	//fmt.Println(code)
	//利用授权码向认证服务器申请令牌
	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)
	//fmt.Println(token)
	if err != nil {
		fmt.Println("Code exchange failed with '%s'\n", err)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err!=nil{
		beego.Error(err)
	}
	//得到用户的json
	/*类似
	Content: {
	 "id": "114512230444013345330",
	 "email": "wangshubo1989@126.com",
	 "verified_email": true,
	 "name": "王书博",
	 "given_name": "书博",
	 "family_name": "王",
	 "picture": "https://lh3.googleusercontent.com/-XdUIqdMkCWA/AAAAAAAAAAI/AAAAAAAAAAA/4252rscbv5M/photo.jpg",
	 "locale": "zh-CN"
	}
		*/
	var J map[string][]map[string]string
	if err := json.Unmarshal(body, &J);
		err != nil {
		beego.Error(err)
	}
	return J
}
