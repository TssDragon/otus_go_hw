package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	ErrInterfaceIsNotStruct = errors.New("interface is not struct")
	ErrValidateMinimumValue = errors.New("value must be greater")
	ErrValidateMaximumValue = errors.New("value must be less")
	ErrValidateStringLength = errors.New("wrong length of string")
	ErrValidateContain      = errors.New("value does not contain in rule")
	ErrValidateRegexp       = errors.New("string does not satisfy regexp")

	ErrValidateTagNotImplemented = errors.New("validate tag has no implement")
	ErrParseValidateValue        = errors.New("parse validate value error")
)

type ValidationOperator int

const (
	Operator ValidationOperator = iota
	Min
	Max
	Len
	In
	Regexp
)

type ValidationRule struct {
	Operator ValidationOperator
	Value    interface{}
}

type ValidationRules []ValidationRule

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err)
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	str := make([]string, 0, len(v))
	for _, e := range v {
		str = append(str, e.Error())
	}
	return strings.Join(str, "\n")
}

func Validate(v interface{}) error {
	if !isInterfaceStruct(v) {
		return ErrInterfaceIsNotStruct
	}

	p := reflect.ValueOf(v)
	t := p.Type()

	numField := t.NumField()

	var validateErrors ValidationErrors
	for i := 0; i < numField; i++ {
		fieldInfo := t.Field(i)
		fieldValue := p.Field(i)

		validationRules, err := prepareValidationRules(fieldInfo)
		if err != nil {
			return err
		}
		if len(validationRules) == 0 {
			continue
		}

		err = validateFieldValueByRules(fieldValue, validationRules)
		if err != nil {
			validateError := ValidationError{Err: err, Field: fieldInfo.Name}
			validateErrors = append(validateErrors, validateError)
		}
	}

	if len(validateErrors) == 0 {
		return nil
	}
	return validateErrors
}

func isInterfaceStruct(v interface{}) bool {
	p := reflect.ValueOf(v)
	return p.Kind() == reflect.Struct
}

func prepareValidationRules(fieldInfo reflect.StructField) (ValidationRules, error) {
	structTag := fieldInfo.Tag
	rawValidationRules := structTag.Get("validate")
	if len(rawValidationRules) == 0 {
		return ValidationRules{}, nil
	}

	return makeValidationRulesFromRawData(rawValidationRules)
}

func makeValidationRulesFromRawData(rawData string) (ValidationRules, error) {
	result := ValidationRules{}
	rules := strings.Split(rawData, "|")
	for _, rule := range rules {
		splitOperatorAndValue := strings.Split(rule, ":")
		if len(splitOperatorAndValue) != 2 {
			continue
		}

		operator, err := detectValidationOperator(splitOperatorAndValue[0])
		if err != nil {
			return ValidationRules{}, err
		}

		value, err := prepareValidateValue(operator, splitOperatorAndValue[1])
		if err != nil {
			return ValidationRules{}, err
		}

		result = append(result, ValidationRule{
			Operator: operator,
			Value:    value,
		})
	}

	return result, nil
}

func detectValidationOperator(s string) (ValidationOperator, error) {
	switch s {
	case "min":
		return Min, nil
	case "max":
		return Max, nil
	case "len":
		return Len, nil
	case "in":
		return In, nil
	case "regexp":
		return Regexp, nil
	}
	return 0, ErrValidateTagNotImplemented
}

func prepareValidateValue(operator ValidationOperator, s string) (interface{}, error) {
	switch operator { //nolint:exhaustive
	case Min:
		fallthrough
	case Max:
		fallthrough
	case Len:
		return strconv.Atoi(s)
	case In:
		return strings.Split(s, ","), nil
	case Regexp:
		return regexp.Compile(s)
	default:
		return nil, ErrValidateTagNotImplemented
	}
}

func validateFieldValueByRules(value reflect.Value, validationRules ValidationRules) error {
	switch value.Type().Kind() { //nolint:exhaustive
	case reflect.Int:
		return validateIntFieldByRules(value.Int(), validationRules)
	case reflect.String:
		return validateStringFieldByRules(value.String(), validationRules)
	case reflect.Slice:
		return validateSliceFiledByRules(value, validationRules)
	default:
		return nil
	}
}

func validateIntFieldByRules(value int64, rules ValidationRules) error {
	for _, rule := range rules {
		switch rule.Operator { //nolint:exhaustive
		case Min:
			if value < int64(rule.Value.(int)) {
				return ErrValidateMinimumValue
			}
		case Max:
			if value > int64(rule.Value.(int)) {
				return ErrValidateMaximumValue
			}
		case In:
			if !validateValueContainsInRule(value, rule.Value) {
				return ErrValidateContain
			}
		default:
			continue
		}
	}
	return nil
}

func validateStringFieldByRules(value string, rules ValidationRules) error {
	for _, rule := range rules {
		switch rule.Operator { //nolint:exhaustive
		case Len:
			if len(value) != rule.Value.(int) {
				return ErrValidateStringLength
			}
		case In:
			if !validateValueContainsInRule(value, rule.Value) {
				return ErrValidateContain
			}
		case Regexp:
			if !validateValueByRegexp(value, rule.Value.(*regexp.Regexp)) {
				return ErrValidateRegexp
			}
		default:
			continue
		}
	}
	return nil
}

func validateSliceFiledByRules(value reflect.Value, rules ValidationRules) error {
	for i := 0; i < value.Len(); i++ {
		err := validateFieldValueByRules(value.Index(i), rules)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateValueContainsInRule(value interface{}, slice interface{}) bool {
	s := reflect.ValueOf(slice)
	strValue := fmt.Sprintf("%v", value)
	for i := 0; i < s.Len(); i++ {
		if s.Index(i).String() == strValue {
			return true
		}
	}
	return false
}

func validateValueByRegexp(value string, regexp *regexp.Regexp) bool {
	return regexp.FindStringIndex(value) != nil
}
