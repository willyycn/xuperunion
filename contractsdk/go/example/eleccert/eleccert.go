/*
 * Copyright (c) 2019. Baidu Inc. All Rights Reserved.
 */

package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/xuperchain/xuperunion/contractsdk/go/code"
	"github.com/xuperchain/xuperunion/contractsdk/go/driver"
)

type elecCert struct {
	User *User
}

// User TODO
type User struct {
	Owner     string
	UserFiles map[string]*UserFile
}

// UserFile TODO
type UserFile struct {
	Timestamp int64
	Hashval   []byte
}

func newElecCert() *elecCert {
	return &elecCert{}
}

func (e *elecCert) putFile(user string, filehash string, ts int64) *User {
	userFile := &UserFile{
		Timestamp: ts,
		Hashval:   []byte(filehash),
	}

	if e.User != nil {
		e.User.Owner = user
		e.User.UserFiles[filehash] = userFile
		return e.User
	}

	u := &User{
		Owner:     user,
		UserFiles: map[string]*UserFile{},
	}
	u.UserFiles[filehash] = userFile

	e.User = u

	return u
}

func (e *elecCert) getFile(user string, filehash string) (*UserFile, error) {
	if e.User != nil {
		if userFile, ok := e.User.UserFiles[filehash]; ok {
			return userFile, nil
		}
		return nil, fmt.Errorf("User's file:%v no exist", filehash)
	}

	return nil, fmt.Errorf("User:%v no exist", user)
}

func (e *elecCert) setContext(ctx code.Context, user string) {
	value, err := ctx.GetObject([]byte(user))
	if err != nil {
	} else {
		userStruc := &User{}
		err = json.Unmarshal(value, userStruc)
		if err != nil {
		}
		e.User = userStruc
	}
}

func (e *elecCert) Initialize(ctx code.Context) code.Response {
	user, ok := ctx.Args()["owner"].(string)
	if !ok {
		return code.Errors("Missing key: owner")
	}

	e.setContext(ctx, user)

	return code.OK(nil)
}

func (e *elecCert) Save(ctx code.Context) code.Response {
	user, ok := ctx.Args()["owner"].(string)
	if !ok {
		return code.Errors("Missing key: owner")
	}
	filehash, ok := ctx.Args()["filehash"].(string)
	if !ok {
		return code.Errors("Missing key: filehash")
	}
	ts, ok := ctx.Args()["timestamp"].(string)
	if !ok {
		return code.Errors("Missing key: filehash")
	}

	e.setContext(ctx, user)
	tsInt, _ := strconv.ParseInt(ts, 10, 64)
	userStruc := e.putFile(user, filehash, tsInt)
	userJSON, _ := json.Marshal(userStruc)

	err := ctx.PutObject([]byte(user), userJSON)
	if err != nil {
		return code.Errors("Invoke method PutObject error")
	}

	return code.Response{
		Status:  200,
		Message: string(userJSON),
	}
}

func (e *elecCert) Query(ctx code.Context) code.Response {
	user, ok := ctx.Args()["owner"].(string)
	if !ok {
		return code.Errors("Missing key: owner")
	}
	filehash, ok := ctx.Args()["filehash"].(string)
	if !ok {
		return code.Errors("Missing key: filehash")
	}

	e.setContext(ctx, user)

	userFile, err := e.getFile(user, filehash)
	if err != nil {
		return code.Errors("Query not exist")
	}

	userFileJSON, _ := json.Marshal(userFile)
	return code.Response{
		Status:  200,
		Message: string(userFileJSON),
	}
}

func main() {
	driver.Serve(newElecCert())
}
