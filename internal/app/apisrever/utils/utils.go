package utils

import (
	"fmt"
	"net/http"
)

//Извекает данные из формы, передаем требуемы поля

func ParseFormFields(r *http.Request, requiredFields []string) (map[string]string, error) {

	if err := r.ParseForm(); err != nil {
		return nil, err
	}

	formData := make(map[string]string)

	for _, field := range requiredFields {
		value := r.FormValue(field)

		if value == "" {
			return nil, fmt.Errorf("field %s is missing", field) //определить свою ошибку
		}
		formData[field] = value
	}

	return formData, nil

}
