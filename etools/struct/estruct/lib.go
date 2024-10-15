package estruct

import (
	"github.com/pkg/errors"
	"reflect"
	"strings"
)

func NotNil(structs ...any) error {
	for _, s := range structs {
		err := notNil(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func NotZero(structs ...any) error {
	for _, s := range structs {
		err := notZero(s)
		if err != nil {
			return err
		}
	}

	return nil
}

func notNil(Struct any) error {
	value := reflect.ValueOf(Struct)

	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag.Get("estruct")
		switch tag {
		case "-":
		default:
			field := value.Field(i)
			if strings.HasPrefix(field.Type().String(), "*") && field.IsNil() {
				return errors.Errorf("[struct => %s]: nil value field => %s", value.Type(), value.Type().Field(i).Name)
			}
		}
	}

	return nil
}

func notZero(Struct any) error {
	value := reflect.ValueOf(Struct)

	for i := 0; i < value.NumField(); i++ {
		tag := value.Type().Field(i).Tag.Get("estruct")
		switch tag {
		case "-":
		default:
			field := value.Field(i)
			if field.IsZero() {
				return errors.Errorf("[struct => %s]: zero value field => %s", value.Type(), value.Type().Field(i).Name)
			}
		}
	}

	return nil
}
