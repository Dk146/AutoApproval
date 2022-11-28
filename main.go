package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	res := isApprove("https://api.github.com/repos/Dk146/AutoApproval/pulls/1/files")
	fmt.Println(res)
}

func isApprove(url string) bool {
	raw_pull_url, raw_origin_url := getPullAndOrigin(url)
	fmt.Println(raw_origin_url)
	fmt.Println(raw_pull_url)

	diff := getDiffContents(raw_pull_url, raw_origin_url)
	whiteList, _ := readLines("whitelist.txt")

	fmt.Println(whiteList)
	fmt.Println(diff)

	result := isContain(diff, whiteList)
	fmt.Println(result)
	return result
}

func getDiffContents(raw_pull, raw_origin string) []string {
	map_pull := getFileContent(raw_pull)
	map_origin := getFileContent(raw_origin)

	diff := make([]string, 0)

	for key1 := range map_pull {
		if _, ok := map_origin[key1]; ok {
			if !reflect.DeepEqual(map_pull[key1], map_origin[key1]) {
				diff = append(diff, key1)
			}
		} else {
			diff = append(diff, key1)
		}
	}

	for key2 := range map_origin {
		if _, ok := map_pull[key2]; !ok {
			diff = append(diff, key2)
		}
	}

	return diff
}

func getFileContent(url string) map[string]interface{} {
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)
	var content interface{}
	if err := yaml.Unmarshal([]byte(body), &content); err != nil {
		panic(err)
	}
	map_content := convert(content).(map[string]interface{})
	return map_content
}

func convert(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			m2[k.(string)] = convert(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = convert(v)
		}
	}
	return i
}

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func isContain(diff, whiteList []string) bool {
	check := false
	for _, e1 := range diff {
		check = false
		for _, e2 := range whiteList {
			if e1 == e2 {
				check = true
			}
		}
		if !check {
			return false
		}
	}
	return true
}

func getValueFromArrayJSON(body []byte, key string) string {
	var raw []map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		panic(err)
	}
	res := raw[0][key].(string)
	return res
}

func getValueFromJSON(body []byte, key string, key1 string) string {
	var raw map[string]interface{}
	if err := json.Unmarshal(body, &raw); err != nil {
		panic(err)
	}
	map1 := raw[key].(map[string]interface{})
	res := map1[key1].(string)
	return res
}

func getPullAndOrigin(url string) (string, string) {
	resp, _ := http.Get(url)
	body, _ := ioutil.ReadAll(resp.Body)

	raw_pull_url := getValueFromArrayJSON(body, "raw_url")
	contents_url := getValueFromArrayJSON(body, "contents_url")

	s := strings.Split(contents_url, "?")

	content_origin, _ := http.Get(s[0])
	body_origin, _ := ioutil.ReadAll(content_origin.Body)

	origin_url := getValueFromJSON(body_origin, "_links", "html")
	raw_origin_url := strings.Replace(origin_url, "blob", "raw", 1)

	return raw_pull_url, raw_origin_url
}
