package credential

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func InitCred() *oauth2.Config {

	config := loadEnvConfig(`envConfig.json`)

	// Credentials which stores google ids.
	type Credentials struct {
		Cid     string `json:"client_id"`
		Csecret string `json:"client_secret"`
	}
	var cred Credentials

	file, err := ioutil.ReadFile("creds.json")
	checkError(err)
	json.Unmarshal(file, &cred)

	conf := &oauth2.Config{
		ClientID:     cred.Cid,
		ClientSecret: cred.Csecret,
		RedirectURL:  `http://` + config.APIHandlerIP + `:` + config.APIHandlerPort + `/auth`,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
		Endpoint: google.Endpoint,
	}

	return conf
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

type config struct {
	APIHandlerIP   string `json:"APIHandlerIP"`
	WebsocketIP    string `json:"WebsocketIP"`
	APIHandlerPort string `json:"APIHandlerPort"`
	WebsocketPort  string `json:"WebsocketPort"`
}

func loadEnvConfig(file string) config {
	var config config
	configFile, err := os.Open(file)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
}
