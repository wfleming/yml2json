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
	if arrYaml, ok := yml.([]interface{}); ok {
		sanitized, err := sanitizeYamlArr(arrYaml)
		if nil != err {
			return nil, err
		}
		return json.Marshal(sanitized)
	} else if mapYaml, ok := yml.(map[interface{}]interface{}); ok {
		sanitized, err := sanitizeYamlMap(mapYaml)
		if nil != err {
			return nil, err
		}
		return json.Marshal(sanitized)
	} else {
		return nil, fmt.Errorf("Unexpected type of YAML: %T", yml)
	}
}

func sanitizeYamlArr(yml []interface{}) ([]interface{}, error) {
	sanitized := make([]interface{}, len(yml))

	for idx, val := range yml {
		if mapVal, ok := val.(map[interface{}]interface{}); ok {
			sanitizedMapVal, err := sanitizeYamlMap(mapVal)
			if nil != err {
				return nil, err
			}
			sanitized[idx] = sanitizedMapVal
		} else {
			sanitized[idx] = val
		}
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
		if mapVal, ok := val.(map[interface{}]interface{}); ok {
			sanitizedMapVal, err := sanitizeYamlMap(mapVal)
			if nil != err {
				return nil, err
			}
			sanitized[strKey] = sanitizedMapVal
		} else {
			sanitized[strKey] = val
		}
	}

	return sanitized, nil
}
