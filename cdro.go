/* Main Package */
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

/*
Configuration file structure
*/
type ReqType struct {
	username   string
	password   string
	host       string
	script     string
	scriptType string
}

/*
RunIde structure holds the parameters that need to be sent to CD/RO for uploading and running Groovy DSL an yaml files
*/
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
	jsonData, err := json.Marshal(values)

	u := fmt.Sprintf("%s/rest/v1.0/server/dsl", reqType.host)
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewBuffer(jsonData))
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
	configFile, err := os.ReadFile(fileName)
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
	r, _ := regexp.Compile(`(?m)@@IncludeFile: (.+)@@`)
	name, _ := regexp.Compile(` .+`)

	matchedStrings := r.FindAllString(scriptStr, -1)

	for i := range matchedStrings {
		name := name.FindString(matchedStrings[i])
		name = strings.Trim(name, "@")
		name = strings.Trim(name, " ")
		subScriptString, err := loadFile(name)
		if err != nil {
			fmt.Println(err)
			return "", err
		}
		subScriptString = strings.Replace(subScriptString, "\\", "\\\\", -1)
		subScriptString = strings.Replace(subScriptString, "'", "\\'", -1)
		scriptStr = strings.Replace(scriptStr, matchedStrings[i], subScriptString, 1)
	}
	return scriptStr, nil
}

/*********************************************************************************************
 *                 Include the Script or YAML file
 */
func includeValues(scriptStr string, valuesList arrayFlags) (string, error) {
	r, _ := regexp.Compile(`(?m)@@IncludeValue: (.+)@@`)
	name, _ := regexp.Compile(` .+`)

	matchedStrings := r.FindAllString(scriptStr, -1)
	for i := range matchedStrings {
		fmt.Println("valuesList ", valuesList)
		for j := range valuesList {
			fmt.Println("valuesList[", j, "] = ", valuesList[j])
			myValue := strings.Split(valuesList[j], "=")
			key := myValue[0]
			value := myValue[1]

			fmt.Println(key, " = ", value)

			name := name.FindString(matchedStrings[i])
			name = strings.Trim(name, "@")
			name = strings.Trim(name, " ")
			fmt.Println("Looking for ", name)
			if name == key {
				fmt.Printf("Replace '%s' with content from '%s'\n", matchedStrings[i], value)
				scriptStr = strings.Replace(scriptStr, matchedStrings[i], value, 1)
			}
		}
	}
	return scriptStr, nil
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

var valuesList arrayFlags

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
	cfgTest := flag.Bool("test", false, "don't apply")
	cfgVerbose := flag.Bool("verbose", false, "Show extra output")
	flag.Var(&valuesList, "value", "Some description for this param.")
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
	scriptStr, err = includeValues(scriptStr, valuesList)
	if err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	if *cfgVerbose {
		fmt.Println("===================")
		fmt.Println(scriptStr)
		fmt.Println("===================")
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

	if !*cfgTest {
		resType, err := runScript(request)
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		} else {
			fmt.Println(resType)
		}
	}

}
