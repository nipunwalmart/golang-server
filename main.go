// package main

// import (
// 	"errors"
// 	"fmt"
// 	"golang-yaml/validation"
// 	"io/ioutil"

// 	m "github.com/hashicorp/go-multierror"
// 	"gopkg.in/yaml.v2"
// )

// func main() {
// 	// even if one file fails the overall validation should fail.
// 	files := make([]string, 0)
// 	files = append(files, "../fold2/exmp.yaml")
// 	files = append(files, "../fold2/exmp2.yaml")

// 	uniqueName := map[string]string{}
// 	var errorList error
// 	for _, file := range files {
// 		yamlData, err := ioutil.ReadFile(file)
// 		if err != nil {
// 			errorList = m.Append(errorList, err)
// 			continue
// 		}
// 		var playbook validation.PlaybookStruct
// 		err = yaml.Unmarshal(yamlData, &playbook)
// 		if err != nil {
// 			errorList = m.Append(errorList, err)
// 			continue
// 		}

// 		nameValidation := validation.ValidateName(file, playbook.Name)
// 		if nameValidation == nil && uniqueName[playbook.Name] != "" {
// 			errorList = m.Append(errorList, errors.New("name is already in use at : "+uniqueName[playbook.Name]))
// 		} else {
// 			uniqueName[playbook.Name] = file
// 		}

// 		errorList = m.Append(errorList, validation.ValidateYamlFile(file, yamlData))
// 	}

// 	listOfErrors := errorList.(*m.Error)
// 	if len(listOfErrors.Errors) > 0 {
// 		fmt.Println(errorList)
// 	} else {
// 		fmt.Println("Validation passed!")
// 	}

// 	fmt.Println()
// }
