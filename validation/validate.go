package validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	m "github.com/hashicorp/go-multierror"
	"gopkg.in/yaml.v2"
)

// validates owner field and returns error
// if there exists space, or empty value, or invalid suffix
func validateOwner(file string, s string) error {
	var errorList error
	if s == "" {
		errString := "owner cannot be empty, in the file : " + file
		errorList = m.Append(errorList, errors.New(errString))
		return errorList
	}

	spaceExists := strings.ContainsRune(s, ' ')
	hasCorrectSuffix := strings.HasSuffix(s, "@walmart.com")

	if spaceExists {
		errString := "owner cannot have spaces, in the file : " + file
		errorList = m.Append(errorList, errors.New(errString))
	}
	if !hasCorrectSuffix {
		errString := "owner must have @walmart.com as suffix , in the file : " + file
		errorList = m.Append(errorList, errors.New(errString))
	}
	return errorList
}

// validates name and return error
// if there exists space, empty field, numeric character.
func ValidateName(file string, s string) error {
	var errorList error
	if s == "" {
		errString := "name field cannot be empty, in the file : " + file
		errorList = m.Append(errorList, errors.New(errString))
		return errorList
	}

	spaceExists := strings.ContainsRune(s, ' ')
	numericExists := regexp.MustCompile(`\d`).MatchString(s)

	if spaceExists {
		errString := "name cannot have spaces, in the file : " + file
		errorList = m.Append(errorList, errors.New(errString))
	}
	if numericExists {
		errString := "name cannot have numeric characters, in the file : " + file
		errorList = m.Append(errorList, errors.New(errString))
	}

	return errorList
}

func validateEmpty(s string, s1 string, isMandatoryField *bool) string {
	if s == "" {
		*isMandatoryField = false
		return s1
	}
	return ""
}

// returns error mentioning all the mandatory fields which are not mentioned in a yaml file.
func validateMandatoryField(file string, playbook PlaybookStruct) error {
	var errorList error
	isMandatoryField := true
	errString := "you must include "

	errString = errString + validateEmpty(playbook.Title, "title, ", &isMandatoryField)
	errString = errString + validateEmpty(playbook.Description, "description, ", &isMandatoryField)
	errString = strings.TrimSuffix(errString, ", ")
	if errString != "you must include " {
		errorList = m.Append(errorList, errors.New(errString+", in the file : "+file))
	}

	for idx, val := range playbook.Tasks {
		errString := "you must include "
		errString = errString + validateEmpty(val.Title, "tasks.title, ", &isMandatoryField)
		errString = errString + validateEmpty(val.Description, "tasks.description, ", &isMandatoryField)
		errString = errString + validateEmpty(val.Type1, "tasks.type, ", &isMandatoryField)
		errString = strings.TrimSuffix(errString, ", ")
		errString = strings.TrimSuffix(errString, "you must include  ")
		if errString != "you must include " {
			errString = errString + " at task : " + fmt.Sprint(idx+1)
			errorList = m.Append(errorList, errors.New(errString+", in the file : "+file))
		}

		errString = "you must include "
		errString = errString + validateEmpty(val.Process.Org, "tasks.process.org, ", &isMandatoryField)
		errString = errString + validateEmpty(val.Process.Project, "tasks.process.project, ", &isMandatoryField)
		errString = errString + validateEmpty(val.Process.Repo, "tasks.process.repo, ", &isMandatoryField)
		errString = errString + validateEmpty(val.Process.Entrypoint, "tasks.process.entrypoint, ", &isMandatoryField)
		errString = errString + validateEmpty(val.Process.RepoBranchOrTag, "tasks.process.repoBranchOrTag, ", &isMandatoryField)
		errString = strings.TrimSuffix(errString, ", ")
		errString = strings.TrimSuffix(errString, "you must include  ")
		if errString != "you must include " {
			errString = errString + " at task : " + fmt.Sprint(idx+1)
			errorList = m.Append(errorList, errors.New(errString+", in the file : "+file))
		}
	}

	for idx, val := range playbook.Tasks {
		if val.Type1 == "concord" {
			if val.Process.Arguments["backupName"] == "" {
				errString := "In concord, you must have arguments and backupName at task: " + fmt.Sprint(idx+1)
				errorList = m.Append(errorList, errors.New(errString+", in the file : "+file))

				isMandatoryField = false
			}
		}
	}
	return errorList

}

func TrimSuffix(errString, s string) {
	panic("unimplemented")
}

// checks if the tasks field is empty and returns error if it is empty.
func validateForNullTasks(file string, t []TasksStruct) error {
	var errorList error
	if len(t) == 0 {
		errString := "you cannot have empty tasks"
		errorList = m.Append(errorList, errors.New(errString+", in the file : "+file))
	}
	return errorList
}

func ValidateYamlFile(file string, yamlData []byte) error {
	// TODO : Metadata -> we have not considered it as mandatory field for now.
	var errorList error

	var playbook PlaybookStruct
	err := yaml.Unmarshal(yamlData, &playbook)
	if err != nil {
		errorList = m.Append(errorList, err)
		return errorList
	}

	emailValidation := validateOwner(file, playbook.Owner)
	nameValidation := ValidateName(file, playbook.Name)
	nullTasksValidation := validateForNullTasks(file, playbook.Tasks)
	mandatoryFieldValidation := validateMandatoryField(file, playbook)

	errorList = m.Append(errorList, emailValidation)
	errorList = m.Append(errorList, nameValidation)
	errorList = m.Append(errorList, nullTasksValidation)
	errorList = m.Append(errorList, mandatoryFieldValidation)

	return errorList
}
