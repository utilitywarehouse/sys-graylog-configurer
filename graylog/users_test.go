package graylog

import (
	"os"
	"testing"
)

func TestUsers(t *testing.T) {
	apiURL := "http://127.0.0.1:9000/api"
	password := os.Getenv("ADMIN_PASSWORD")
	if os.Getenv("GRAYLOG_CONFIGURER_REAL_API") == "" {
		t.Skip("Skipping tests against real graylog api")
	}
	cl := NewClient(apiURL, password)
	err := cl.ApiReachable()
	if err != nil {
		t.Fatalf("API not reachable: %s", err)
	}

	t.Run("get user", func(t *testing.T) {
		cl := NewClient(apiURL, password)
		name := "testerGet"
		u, _ := cl.getUser(name)
		if u.Username != name {
			t.Error("can't get user")
		}
	})
	t.Run("create user", func(t *testing.T) {
		cl := NewClient(apiURL, password)
		name := "testerCreate"
		_, err := cl.getUser(name)
		if err == nil {
			t.Error("user already exists")
		}
		cu, _ := cl.createUser(name)
		gu, _ := cl.getUser(name)

		if cu.Username != gu.Username {
			t.Error("can't create user")
		}
		cl.deleteUser(name)
	})
	t.Run("delete user", func(t *testing.T) {
		cl := NewClient(apiURL, password)
		name := "testDelete"
		cl.createUser(name)
		err := cl.deleteUser(name)

		if err != nil {
			t.Error("can't delete user")
		}
		_, err = cl.getUser(name)
		if err == nil {
			t.Error("can't delete user")
		}
	})
	t.Run("delete user", func(t *testing.T) {
		cl := NewClient(apiURL, password)
		name := "testEdit"
		u, _ := cl.createUser(name)
		u.Email = "changed"
		err := cl.editUser(u)
		if err != nil {
			t.Error("can't edit user")
		}

		u, err = cl.getUser(name)
		if u.Email != "changed" {
			t.Error("can't edit user")
		}
		cl.deleteUser(name)
	})
}

func TestAdmins(t *testing.T) {
	apiURL := "http://127.0.0.1:9000/api"
	password := os.Getenv("ADMIN_PASSWORD")
	if os.Getenv("GRAYLOG_CONFIGURER_REAL_API") == "" {
		t.Skip("Skipping tests against real graylog api")
	}
	cl := NewClient(apiURL, password)
	err := cl.ApiReachable()
	if err != nil {
		t.Fatalf("API not reachable: %s", err)
	}

	t.Run("isAdmin works as expected", func(t *testing.T) {
		u := User{Roles: []string{"Admin"}}
		if !u.isAdmin() {
			t.Error("admin user not detected")
		}
		u = User{Roles: []string{"Not Admin"}}
		if u.isAdmin() {
			t.Error("false admin user not detected")
		}
	})
	t.Run("creates admin if not exist", func(t *testing.T) {
		cl := NewClient(apiURL, password)
		name := "testerNewAdmin"
		cl.SetAdmin(name)
		u, _ := cl.getUser(name)

		if !u.isAdmin() {
			t.Error("admin user not created")
		}
		cl.deleteUser(name)
	})
	t.Run("promotes user to admin", func(t *testing.T) {
		cl := NewClient(apiURL, password)
		name := "testerAdmin"
		u, _ := cl.createUser(name)
		u.Roles = []string{}
		cl.editUser(u)
		u, _ = cl.getUser(name)
		if u.isAdmin() {
			t.Error("test failed to create non-admin user")
		}

		cl.SetAdmin(name)
		u, _ = cl.getUser(name)
		if !u.isAdmin() {
			t.Error("user not promoted to admin")
		}
		cl.deleteUser(name)
	})
}
