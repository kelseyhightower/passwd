# passwd

[![Build Status](https://travis-ci.org/kelseyhightower/passwd.png?branch=master)](https://travis-ci.org/kelseyhightower/passwd) [![GoDoc](https://godoc.org/github.com/kelseyhightower/passwd/user?status.svg)](https://godoc.org/github.com/kelseyhightower/passwd/user)

Drop in replacement for the os/user package. Useful when building without cgo and want to fallback to /etc/passwd.

## Usage

```
package main

import (
	"fmt"
	"log"

	"github.com/kelseyhightower/passwd/user"
)

func main() {
	u, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	u, err = user.Lookup("root")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s:%s:%s:%s:%s\n",
		u.Username, u.Uid, u.Gid, u.Name, u.HomeDir)
	u, err = user.LookupId("0")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s:%s:%s:%s:%s\n",
		u.Username, u.Uid, u.Gid, u.Name, u.HomeDir)
}
```
