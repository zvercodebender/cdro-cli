package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type ReqType struct {
	username   string
	password   string
	host       string
	script     string
	scriptType string
}

type RunIde struct {
	request string
	dsl     string
	format  string
}

/*********************************************************************************************
 *                 Run the Script on the CD/RO Server
 */
func runScript(reqType ReqType) (string, error) {
	// HTTP call
	var (
		resp   *http.Response
		client = new(http.Client)
	)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client.Transport = tr

	values := map[string]string{"request": "evalDsl", "format": reqType.scriptType, "dsl": reqType.script}
	json_data, err := json.Marshal(values)

	u := fmt.Sprintf("%s/rest/v1.0/server/dsl", reqType.host)
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(json_data))
	req.Header.Set("Content-Type", "application/json; charset=UTF-8")
	req.SetBasicAuth(reqType.username, reqType.password)

	resp, err = client.Do(req)
	if err != nil {
		log.Fatalln(err)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(resp.Status)
	}
	// Decode JSON
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	respString := string(bodyBytes)
	return respString, nil
}

/*********************************************************************************************
 *                 Load the Config YAML file
 */
func loadConfiguration(fileName string) (ReqType, error) {
	configFile, err := os.ReadFile("config.json")
	if err != nil {
		fmt.Println(err)
	}
	jsonStr2 := string(configFile)
	var y map[string]interface{}
	json.Unmarshal([]byte(jsonStr2), &y)
	config := ReqType{}
	config.username = fmt.Sprint(y["username"])
	config.password = fmt.Sprint(y["password"])
	config.host = fmt.Sprint(y["url"])
	return config, nil
}

/*********************************************************************************************
 *                 Load the Script or YAML file
 */
func loadFile(fileName string) (string, error) {
	scriptFile, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	script := string(scriptFile)
	return script, nil
}

/*********************************************************************************************
 *                 Include the Script or YAML file
 */
func includeFile(scriptStr string) (string, error) {
	r, _ := regexp.Compile(`(?m)@@IncludeFile: ([A-Za-z\.\\\/]+)@@`) // Pay attention, no ", instead `
	name, _ := regexp.Compile(` .+`)

	matched_strings := r.FindAllString(scriptStr, -1)

	for i := range matched_strings {
		name := name.FindString(matched_strings[i])
		name = strings.Trim(name, "@")
		name = strings.Trim(name, " ")
		subScriptString, err := loadFile(name)
		if err != nil {
			fmt.Println(err)
			return "", err
		} else {
			scriptStr = strings.Replace(scriptStr, matched_strings[i], subScriptString, 1)
		}
	}
	return scriptStr, nil
}

/*********************************************************************************************
 *********************************************************************************************
 **
 **                       M A I N   L O O P
 **
 *********************************************************************************************
 *********************************************************************************************
 */
func main() {
	usernamePtr := flag.String("username", "", "Username")
	passwdPtr := flag.String("password", "", "Password")
	urlPtr := flag.String("url", "", "CDRO URL")
	typePtr := flag.String("type", "groovy", "DSL or YAML files")
	scriptPtr := flag.String("file", "ERROR", "Groovy or YAML file to run")
	cfgPrt := flag.String("conf", "NONE", "Config file")
	flag.Parse()

	request := ReqType{}
	if *cfgPrt != "NONE" {
		request, _ = loadConfiguration(*cfgPrt)
	}
	scriptStr, err := loadFile(*scriptPtr)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	scriptStr, err = includeFile(scriptStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	request.script = scriptStr
	if *cfgPrt != "NONE" {
		fmt.Println("load config file")
	}
	if *usernamePtr != "" || request.username == "" {
		request.username = *usernamePtr
	}
	if *passwdPtr != "" || request.password == "" {
		request.password = *passwdPtr
	}
	if *urlPtr != "" || request.host == "" {
		request.host = *urlPtr
	}
	request.scriptType = *typePtr

	resType, err := runScript(request)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
	} else {
		fmt.Println(resType)
	}

}
