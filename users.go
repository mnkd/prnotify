package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
)

type UsersMap map[string]string

type UsersManager struct {
	UsersMap UsersMap
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func (manager UsersManager) ConvertGitHubToSlack(name string) string {
	var username = name
	if v, ok := manager.UsersMap[name]; ok {
		username = v
	}
	return username
}

func NewUsersMap() (UsersMap, error) {
	users := make(UsersMap)

	usr, err := user.Current()
	if err != nil {
		fmt.Fprintln(os.Stderr, "[Error] Could not get current user:", err)
		return users, err
	}

	path := filepath.Join(usr.HomeDir, "/.config/prnotify/users.json")
	if exists(path) == false {
		fmt.Fprintln(os.Stderr, "[Warning] File not exists:", path)
		return users, nil
	}

	str, err := ioutil.ReadFile(path)
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

func NewUsersManager() (UsersManager, error) {
	var manager = UsersManager{}
	var err error
	manager.UsersMap, err = NewUsersMap()
	if err != nil {
		return manager, err
	}
	return manager, nil
}
