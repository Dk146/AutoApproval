package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"

	"gopkg.in/yaml.v2"
)

func main() {
	diff := compareFile("file1.yaml", "file2.yaml")
	whiteList, _ := readLines("whitelist.txt")
	fmt.Println(whiteList)
	fmt.Println(diff)
	result := isApprove(diff, whiteList)
	fmt.Println(result)
}

func compareFile(path1, path2 string) []string {
	yamlFile1, _ := os.ReadFile(path1)
	yamlFile2, _ := os.ReadFile(path2)

	var body1 interface{}
	var body2 interface{}

	if err := yaml.Unmarshal([]byte(yamlFile1), &body1); err != nil {
		panic(err)
	}
	if err := yaml.Unmarshal([]byte(yamlFile2), &body2); err != nil {
		panic(err)
	}

	map1, _ := convert(body1).(map[string]interface{})
	map2, _ := convert(body2).(map[string]interface{})

	diff := make([]string, 0)

	for key1 := range map1 {
		if _, ok := map2[key1]; ok {
			if !reflect.DeepEqual(map1[key1], map2[key1]) {
				diff = append(diff, key1)
			}
		} else {
			diff = append(diff, key1)
		}
	}

	for key2 := range map2 {
		if _, ok := map1[key2]; !ok {
			diff = append(diff, key2)
		}
	}

	return diff
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

func isApprove(diff, whiteList []string) bool {
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
