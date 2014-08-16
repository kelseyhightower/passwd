// Copyright (c) 2014 Kelsey Hightower. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package user

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"syscall"
	"testing"
)

var records = []string{
	"nobody:*:-2:-2:Unprivileged User:/var/empty:/usr/bin/false",
	"root:*:0:0:System Administrator:/var/root:/bin/sh",
	"daemon:*:1:1:System Services:/var/root:/usr/bin/false",
	"nginx:*:500:500:Nginx User:/var/nginx:/usr/bin/false",
	"kelsey:*:1000:1000:Kelsey Hightower,,5555555555,@kelseyhightower:/home/kelsey:/bin/bash",
}

func newUserDatabase(users []string) (string, error) {
	f, err := ioutil.TempFile("", "")
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = f.WriteString("# Comment\n")
	if err != nil {
		return "", err
	}
	_, err = f.WriteString(strings.Join(users, "\n"))
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

func TestCurrent(t *testing.T) {
	uid := strconv.Itoa(syscall.Getuid())
	users := []string{
		fmt.Sprintf("current:*:%s:%s:Current User:/home/current:/bin/bash", uid, uid),
	}
	db, err := newUserDatabase(users)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(db)
	userDatabase = db
	want := &User{
		Username: "current",
		Uid:      uid,
		Gid:      uid,
		Name:     "Current User",
		HomeDir:  "/home/current",
	}
	got, err := Current()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

var passwdLookupTests = []struct {
	username string
	user     *User
}{
	{
		"kelsey",
		&User{
			Username: "kelsey",
			Uid:      "1000",
			Gid:      "1000",
			Name:     "Kelsey Hightower",
			HomeDir:  "/home/kelsey",
		},
	},
	{"unknown", nil},
}

func TestLookup(t *testing.T) {
	db, err := newUserDatabase(records)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(db)
	userDatabase = db
	for _, tt := range passwdLookupTests {
		got, err := Lookup(tt.username)
		switch tt.username {
		case "unknown":
			want := "user: unknown user unknown"
			if err.Error() != want {
				t.Errorf("want %s, got %s", want, err.Error())
			}
		default:
			if err != nil {
				t.Error(err)
			}
		}
		if !reflect.DeepEqual(tt.user, got) {
			t.Errorf("want %v, got %v", tt.user, got)
		}
	}
}

var passwdLookupIdTests = []struct {
	uid  string
	user *User
}{
	{
		"500",
		&User{
			Username: "nginx",
			Uid:      "500",
			Gid:      "500",
			Name:     "Nginx User",
			HomeDir:  "/var/nginx",
		},
	},
	{"30000", nil},
	{"unknown", nil},
}

func TestLookupId(t *testing.T) {
	db, err := newUserDatabase(records)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(db)
	userDatabase = db
	for _, tt := range passwdLookupIdTests {
		got, err := LookupId(tt.uid)
		switch tt.uid {
		case "unknown":
			if err == nil {
				t.Error("Expected non-nil err")
			}
		case "30000":
			want := "user: unknown userid 30000"
			if err.Error() != want {
				t.Errorf("want %s, got %s", want, err.Error())
			}
		default:
			if err != nil {
				t.Error(err)
			}
		}
		if !reflect.DeepEqual(tt.user, got) {
			t.Errorf("want %v, got %v", tt.user, got)
		}
	}
}

func TestMissingPasswd(t *testing.T) {
	userDatabase = "/doesnotexist"
	_, err := lookupPasswd(500, "nginx", true)
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestInvalidUserDatabase(t *testing.T) {
	users := []string{
		fmt.Sprintf("kelsey:*:1000:1000:Current User:/home/current"),
	}
	db, err := newUserDatabase(users)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(db)
	userDatabase = db
	_, err = Lookup("kelsey")
	if err != InvalidUserDatabaseError {
		t.Errorf("Exected InvalidUserDatabaseError")
	}
}
