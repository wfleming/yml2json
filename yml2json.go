// The yml2json tool converts YAML to JSON.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"yaml"
)

func main() {
	input, err := getInput()
	if nil == input {
		printUsage()
		os.Exit(64)
	} else if nil != err {
		fmt.Fprintf(os.Stderr, "Error reading in YAML: %s\n", err)
		os.Exit(1)
	}

	var yml interface{}
	err = yaml.Unmarshal(input, &yml)

	if nil != err {
		fmt.Fprintf(os.Stderr, "Error parsing YAML: %s\n", err)
		os.Exit(1)
	}

	jsonRaw, err := yamlToJSON(yml)
	if nil != err {
		fmt.Fprintf(os.Stderr, "Error converting to JSON: %s\n", err)
		os.Exit(1)
	}

	fmt.Println(string(jsonRaw))
}

func getInput() ([]byte, error) {
	if len(os.Args) > 1 {
		filename := os.Args[1]
		return ioutil.ReadFile(filename)
	} else if stat, _ := os.Stdin.Stat(); (stat.Mode() & os.ModeCharDevice) == 0 {
		return ioutil.ReadAll(os.Stdin)
	} else {
		return nil, nil
	}
}

func printUsage() {
	fmt.Printf("yml2json version %s\n\n", version)
	fmt.Println("YAML goes in, JSON comes out.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("\tyml2json file.json")
	fmt.Println("\tcat file.json | yml2json")
}

func yamlToJSON(yml interface{}) ([]byte, error) {
	sanitized, err := sanitizeYaml(yml)
	if nil != err {
		return nil, err
	}
	return json.Marshal(sanitized)
}

func sanitizeYaml(yml interface{}) (interface{}, error) {
	if arrYaml, ok := yml.([]interface{}); ok {
		return sanitizeYamlArr(arrYaml)
	} else if mapYaml, ok := yml.(map[interface{}]interface{}); ok {
		return sanitizeYamlMap(mapYaml)
	} else {
		return yml, nil
	}
}

func sanitizeYamlArr(yml []interface{}) ([]interface{}, error) {
	sanitized := make([]interface{}, len(yml))

	for idx, val := range yml {
		sanitizedVal, err := sanitizeYaml(val)
		if nil != err {
			return nil, err
		}
		sanitized[idx] = sanitizedVal
	}

	return sanitized, nil
}

func sanitizeYamlMap(yml map[interface{}]interface{}) (map[string]interface{}, error) {
	sanitized := make(map[string]interface{})

	for key, val := range yml {
		var strKey string
		if castKey, ok := key.(string); ok {
			strKey = castKey
		} else {
			strKey = fmt.Sprintf("%v", key)
		}

		sanitizedVal, err := sanitizeYaml(val)
		if err != nil {
			return nil, err
		}
		sanitized[strKey] = sanitizedVal
	}

	return sanitized, nil
}
