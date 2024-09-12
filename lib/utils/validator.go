package utils

import (
	"fmt"
	"net/mail"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

type ValidationError struct {
	Message  string
	Criteria string
	Field    string
}

func NewValidationError(m string, c string, f string) *ValidationError {
	return &ValidationError{Message: m, Criteria: c, Field: f}
}

func (e *ValidationError) Error() string {
	return e.Message
}

const (
	VALIDATE   = "validate"
	REQUIRED   = "required"
	MIN_LENGTH = "min_length"
	MAX_LENGTH = "max_length"
	MIN        = "min"
	MAX        = "max"
	REGEX      = "regex"
	IN         = "in"
	EMAIL      = "email"
	URL        = "url"
)

type Validator struct {
	entity interface{}
	fields map[string]*fieldProperty
}

type fieldProperty struct {
	Kind        string
	Type        string
	Tags        []string
	ValueString string
	ValueInt    int64
	ValueUint   uint64
	SliceLength int
	IsNil       bool
	IsZero      bool
}

func NewValidator(entity interface{}) *Validator {
	return &Validator{
		entity: entity,
		fields: map[string]*fieldProperty{},
	}
}

func (e *Validator) Validate() error {
	e.parseFields()
	for f := range e.fields {
		tags := e.fields[f].Tags
		kind := e.fields[f].Kind
		_type := e.fields[f].Type
		for _, c := range tags {
			if c == "" {
				continue
			}
			if c == REQUIRED {
				if ((kind == "int" || kind == "int64" || kind == "uint" || kind == "uint64") && (e.fields[f].ValueInt == 0 && e.fields[f].ValueUint == 0)) || (kind == "string" && e.fields[f].ValueString == "") || (kind == "slice" && e.fields[f].SliceLength == 0) || (kind == "ptr" && e.fields[f].IsNil) || (_type == "primitive.ObjectID" && e.fields[f].IsZero) {
					return NewValidationError(fmt.Sprintf("%s is required", f), c, f)
				}

			} else if c == EMAIL {
				if kind != "string" {
					return NewValidationError(fmt.Sprintf("The '%s' criteria doesn't work for field of type '%s'", c, kind), c, f)
				}
				if e.fields[f].ValueString == "" {
					continue
				}
				if _, err := mail.ParseAddress(e.fields[f].ValueString); err != nil {
					return NewValidationError(fmt.Sprintf("%s is invalid", f), c, f)
				}
			} else if c == URL {
				if kind != "string" {
					return NewValidationError(fmt.Sprintf("The '%s' criteria doesn't work for field of type '%s'", c, kind), c, f)
				}
				if e.fields[f].ValueString == "" {
					continue
				}
				r, err := url.ParseRequestURI(e.fields[f].ValueString)
				if err != nil || r.Host == "" {
					return NewValidationError(fmt.Sprintf("%s is invalid", f), c, f)
				}
			} else if strings.Contains(c, "=") {
				cr := strings.Split(c, "=")
				crs, crv := cr[0], cr[1]
				if crs == MIN {
					if kind != "int" && kind != "int64" && kind != "uint" && kind != "uint64" {
						return NewValidationError(fmt.Sprintf("The '%s' criteria doesn't work for field of type '%s'", crs, kind), c, f)
					}
					v_int, err := strconv.ParseInt(crv, 10, 64)
					if err != nil {
						return NewValidationError(fmt.Sprintf("Criteria value %s invalid", crv), c, f)
					}
					if e.fields[f].ValueInt < int64(v_int) && e.fields[f].ValueUint < uint64(v_int) {
						return NewValidationError(fmt.Sprintf("%s must be at least %d", f, v_int), c, f)
					}
				} else if crs == MAX {
					if kind != "int" && kind != "int64" && kind != "uint" && kind != "uint64" {
						return NewValidationError(fmt.Sprintf("The '%s' criteria doesn't work for field of type '%s'", crs, kind), c, f)
					}
					v_int, err := strconv.ParseInt(crv, 10, 64)
					if err != nil {
						return NewValidationError(fmt.Sprintf("Criteria value %s invalid", crv), c, f)
					}
					if e.fields[f].ValueInt > int64(v_int) || e.fields[f].ValueUint > uint64(v_int) {
						return NewValidationError(fmt.Sprintf("%s must be no more than %d", f, v_int), c, f)
					}
				} else if crs == MIN_LENGTH {
					if kind != "string" && kind != "slice" {
						return NewValidationError(fmt.Sprintf("The '%s' criteria doesn't work for field of type '%s'", crs, kind), c, f)
					}
					v_int, err := strconv.ParseInt(crv, 10, 32)
					if err != nil {
						return NewValidationError(fmt.Sprintf("Criteria value %s invalid", crv), c, f)
					}
					if kind == "string" && utf8.RuneCountInString(e.fields[f].ValueString) < int(v_int) {
						return NewValidationError(fmt.Sprintf("The length of %s must be at least %d characters", f, v_int), c, f)
					}
					if kind == "slice" && e.fields[f].SliceLength < int(v_int) {
						return NewValidationError(fmt.Sprintf("The length of %s must be at least %d items", f, v_int), c, f)
					}
				} else if crs == MAX_LENGTH {
					if kind != "string" && kind != "slice" {
						return NewValidationError(fmt.Sprintf("The '%s' criteria doesn't work for field of type '%s'", crs, kind), c, f)
					}
					v_int, err := strconv.ParseInt(crv, 10, 32)
					if err != nil {
						return NewValidationError(fmt.Sprintf("Criteria value %s invalid", crv), c, f)
					}
					if kind == "string" && utf8.RuneCountInString(e.fields[f].ValueString) > int(v_int) {
						return NewValidationError(fmt.Sprintf("The length of %s must be no more than %d characters", f, v_int), c, f)
					}
					if kind == "slice" && e.fields[f].SliceLength > int(v_int) {
						return NewValidationError(fmt.Sprintf("The length of %s must be no more than %d items", f, v_int), c, f)
					}
				} else if crs == REGEX {
					if kind != "string" {
						return NewValidationError(fmt.Sprintf("The '%s' criteria doesn't work for field of type '%s'", crs, kind), c, f)
					}
					if e.fields[f].ValueString == "" {
						continue
					}
					ex, err := regexp.Compile(crv)
					if err != nil {
						return NewValidationError("The regex was provided is invalid", c, f)
					}
					str_fields := strings.Fields(e.fields[f].ValueString)
					for _, v := range str_fields {
						match := ex.MatchString(v)
						if !match {
							return NewValidationError(fmt.Sprintf("%s is invalid", f), c, f)
						}
					}
				} else if crs == IN {
					if kind != "string" {
						return NewValidationError(fmt.Sprintf("The '%s' criteria doesn't work for field of type '%s'", crs, kind), c, f)
					}
					if e.fields[f].ValueString == "" {
						continue
					}
					crv = strings.ReplaceAll(crv, "[", "")
					crv = strings.ReplaceAll(crv, "]", "")
					items := strings.Split(crv, "|")
					isIn := false
					for _, v := range items {
						if strings.Trim(v, " ") == e.fields[f].ValueString {
							isIn = true
							break
						}
					}
					if !isIn {
						return NewValidationError(fmt.Sprintf("%s is invalid", f), c, f)
					}
				}
			}
		}
	}
	return nil
}

func (e *Validator) parseFields() {
	v := reflect.ValueOf(e.entity).Elem()
	t := reflect.TypeOf(e.entity).Elem()
	rv := reflect.ValueOf(e.entity)
	iv := reflect.Indirect(rv)
	for j := 0; j < v.NumField(); j++ {
		field_name := v.Type().Field(j).Name
		field, _ := t.FieldByName(field_name)
		tags := field.Tag.Get(VALIDATE)
		if len(tags) == 0 {
			continue
		}
		value := iv.FieldByName(field_name)
		fp := &fieldProperty{Kind: v.Field(j).Kind().String(), Type: v.Field(j).Type().String(), Tags: strings.Split(tags, ",")}
		if fp.Kind == "int" || fp.Kind == "int64" {
			fp.ValueInt = value.Int()
		}
		if fp.Kind == "uint" || fp.Kind == "uint64" {
			fp.ValueUint = value.Uint()
		}
		if fp.Kind == "string" {
			fp.ValueString = value.String()
		}
		if fp.Type == "primitive.ObjectID" {
			fp.IsZero = value.IsZero()
		}
		if fp.Kind == "slice" {
			fp.SliceLength = value.Len()
		}
		if fp.Kind == "ptr" {
			fp.IsNil = value.IsNil()
		}
		e.fields[field_name] = fp
	}
}
