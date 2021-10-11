package sweb

import (
	"errors"
	"net/http"
	"strconv"
)

var (
	ErrEmptyKeys   = errors.New("keys is empty")
	ErrKeyNotFound = errors.New("key not found")
)

func ParseForm(r *http.Request) error {
	return r.ParseForm()
}

func GetKeysFunc(r *http.Request) func(key string) (string, error) {
	if err := r.ParseForm(); err != nil {
		return func(key string) (string, error) {
			return "", err
		}
	}
	return func(key string) (string, error) {
		return r.Form.Get(key), nil
	}
}

func UrlParams(r *http.Request, keys ...string) ([]string, error) {
	if len(keys) == 0 {
		return nil, ErrEmptyKeys
	}
	getKeyFunc := GetKeysFunc(r)
	values := make([]string, len(keys))
	var err error
	for i, v := range keys {
		values[i], err = getKeyFunc(v)
		if err != nil {
			return nil, err
		}
	}
	return values, nil
}

// URLParamString must called ParseForm
func URLParamString(r *http.Request, key string) (string, error) {
	value := r.Form.Get(key)
	if value == "" {
		return "", ErrKeyNotFound
	}
	return value, nil
}

// URLParamInt must called ParseForm
func URLParamInt(r *http.Request, key string) (int, error) {
	value := r.Form.Get(key)
	if value == "" {
		return 0, ErrKeyNotFound
	}
	return strconv.Atoi(value)
}
