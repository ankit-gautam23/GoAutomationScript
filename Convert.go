package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
)

func main() {
	// check if file path argument is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide a file path.")
		return
	}

	// read file contents
	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	// convert method signatures and return statements
	output := convert(string(content))

	// write output to file
	err = ioutil.WriteFile("output.go", []byte(output), 0644)
	if err != nil {
		fmt.Println("Error writing output file:", err)
		return
	}

	fmt.Println("Conversion complete. Output file: output.go")
}

func convert(content string) string {
	// regular expression to match C# method signatures
	re := regexp.MustCompile(`(?ms)^(public|private|protected|internal)\s+(static\s+)?([a-zA-Z0-9_<>]+)\s+([a-zA-Z0-9_]+)\((.?)\)\s\{([\s\S]*?)^\}`)

	// find all matches
	matches := re.FindAllStringSubmatch(content, -1)

	// iterate over matches and convert to Go
	var output string
	for _, match := range matches {
		access := match[1]
		static := match[2]
		returnType := match[3]
		methodName := match[4]
		parameters := match[5]
		body := match[6]

		// convert method signature
		output += fmt.Sprintf("%s func %s(", access, methodName)

		// convert parameters
		parameterList := strings.Split(parameters, ",")
		for i, param := range parameterList {
			param = strings.TrimSpace(param)
			paramParts := strings.Split(param, " ")
			if len(paramParts) == 2 {
				paramName := paramParts[1]
				output += fmt.Sprintf("%s %s", paramName, paramParts[0])
			} else {
				output += param
			}

			if i < len(parameterList)-1 {
				output += ", "
			}
		}

		output += fmt.Sprintf(") %s {\n", returnType)

		// convert body
		body = strings.TrimSpace(body)
		if body != "" {
			body = strings.ReplaceAll(body, "return", "return ")
			output += body + "\n"
		}

		output += "}\n\n"
	}

	return output
}
