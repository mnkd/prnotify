package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type UsersMap map[string]string

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func NewUsers(path string) (UsersMap, error) {
	var users UsersMap
	abs, err := filepath.Abs(path)

	if exists(abs) == false {
		return users, nil
	}

	str, err := ioutil.ReadFile(abs)
	if err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] Could not read users.json:", err)
		return users, err
	}

	if err := json.Unmarshal(str, &users); err != nil {
		fmt.Fprintln(os.Stderr, "[ERROR] JSON unmarshal:", err)
		return users, err
	}

	return users, nil
}
