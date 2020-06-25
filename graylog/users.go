package graylog

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

type User struct {
	Email       string   `json:"email"`
	Fullname    string   `json:"full_name"`
	Password    string   `json:"password"`
	Permissions []string `json:"permissions"`
	Roles       []string `json:"roles"`
	Username    string   `json:"username"`
}

func (cl *Client) SetAdmin(name string) error {
	u, err := cl.getUser(name)
	if err != nil {
		_, err = cl.createUser(name)
		return err
	}
	u.Roles = append(u.Roles, "Admin")
	return cl.editUser(u)
}

func (u *User) isAdmin() bool {
	for _, r := range u.Roles {
		if r == "Admin" {
			return true
		}
	}
	return false
}

func (cl *Client) getUser(name string) (User, error) {
	var u User

	err := cl.callAPI("GET", "/users/"+name, nil, &u)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (cl *Client) createUser(name string) (User, error) {
	u := User{
		Email:       name + "@utilitywarehouse.co.uk",
		Fullname:    name,
		Permissions: []string{},
		Roles:       []string{"Admin"},
		Username:    name,
	}

	p, err := generatePassword()
	if err != nil {
		return u, fmt.Errorf("can't generate secure password: %s", err)
	}
	u.Password = p

	err = cl.callAPI("POST", "/users", u, nil)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (cl *Client) deleteUser(name string) error {
	return cl.callAPI("DELETE", "/users/"+name, nil, nil)
}

func (cl *Client) editUser(u User) error {
	return cl.callAPI("PUT", "/users/"+u.Username, u, nil)
}

func generatePassword() (string, error) {
	b := make([]byte, 30)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
