package validation

import (
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"
)

func traverseData1(data interface{}, isVali *bool) {
	isValid := true
	isVali = &isValid
	switch v := data.(type) {
	case map[interface{}]interface{}:
		for key, value := range v {
			// fmt.Printf("Key: %v, Value: %v\n", key, value)
			if key == "name" {
				switch v1 := value.(type) {
				case int:
					fmt.Printf("%v Cannot Be %T", key, v1)
					isValid = false
					return
				case string:
					for _, ch := range v1 {
						if ch == ' ' {
							fmt.Printf("%v Cannot have spaces in it", key)
							isValid = false
							return
						}
					}

				}
			} else if key == "owner" {

				switch v1 := value.(type) {

				case int:
					fmt.Printf("%v Cannot Be %T", key, v1)
					isValid = false
					return
				case string:
					res := make([]rune, 0)
					haveSeen := false
					for _, ch := range v1 {
						if ch == '@' {
							haveSeen = true
						}
						if haveSeen {
							res = append(res, ch)
						}
					}

					if string(res) != "@walmart.com" {
						fmt.Println("owner must have @walmart.com as suffix")
						isValid = false
						return
					}

				}
			}
			traverseData1(value, isVali)
		}
	case []interface{}:
		for _, element := range v {
			// fmt.Println(element)
			traverseData1(element, isVali)
		}
	}
}

func validateEmail1(s string) bool {
	res := make([]rune, 0)
	hasSeen := false
	for _, ch := range s {
		if ch == '@' {
			hasSeen = true
		}
		if hasSeen {
			res = append(res, ch)
		}
	}
	if string(res) != "@walmart.com" {
		fmt.Println("Email/owner Must End with @walmart.com as Suffix")
		return false
	}
	return true
}

func validateName1(s string) bool {
	for _, ch := range s {
		if ch == ' ' {
			fmt.Println("Name Cannot Have Spaces")
			return false
		}
		if ch >= '0' && ch <= '9' {
			fmt.Println("Name Cannot be Numeric")
			return false
		}
	}
	return true
}

func ValidateYamlFile1(file string) bool {
	yamlData, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return false
	}

	var data interface{}
	err = yaml.Unmarshal(yamlData, &data)
	if err != nil {
		fmt.Println(err)
		return false
	}

	_, ok := data.(map[interface{}]interface{})
	if !ok {
		fmt.Println("It is not covertable to a map Structure!")
		return false
	}
	isValid := true
	traverseData1(data, &isValid)
	return isValid
}

func processMultipleFiles(fileNames []string) error {
	var resultErr error
	for _, fileName := range fileNames {
		if err := processFile(fileName); err != nil {
			resultErr = multierror.Append(resultErr, err)
		}
	}
	return resultErr
}

func processFile(fileName string) error {
	// Perform some processing on the file
	// If an error occurs, return it
	str := "absolutely fine"
	return errors.New(str)
}

func Maini() {
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	if err := processMultipleFiles(files); err != nil {
		fmt.Println("Encountered errors:")
		fmt.Println(err.Error())
	} else {
		fmt.Println("All files processed successfully.")
	}
}
