// Copyright 2021 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/casdoor/casdoor/util"
	"xorm.io/core"
)

type User struct {
	Owner       string `xorm:"varchar(100) notnull pk" json:"owner"`
	Name        string `xorm:"varchar(100) notnull pk" json:"name"`
	CreatedTime string `xorm:"varchar(100)" json:"createdTime"`

	Id            string `xorm:"varchar(100)" json:"id"`
	Type          string `xorm:"varchar(100)" json:"type"`
	Password      string `xorm:"varchar(100)" json:"password"`
	DisplayName   string `xorm:"varchar(100)" json:"displayName"`
	Avatar        string `xorm:"varchar(255)" json:"avatar"`
	Email         string `xorm:"varchar(100)" json:"email"`
	Phone         string `xorm:"varchar(100)" json:"phone"`
	Affiliation   string `xorm:"varchar(100)" json:"affiliation"`
	Tag           string `xorm:"varchar(100)" json:"tag"`
	IsAdmin       bool   `json:"isAdmin"`
	IsGlobalAdmin bool   `json:"isGlobalAdmin"`
	IsForbidden   bool   `json:"isForbidden"`
	Hash          string `xorm:"varchar(100)" json:"hash"`
	PreHash       string `xorm:"varchar(100)" json:"preHash"`

	Github string `xorm:"varchar(100)" json:"github"`
	Google string `xorm:"varchar(100)" json:"google"`
	QQ     string `xorm:"qq varchar(100)" json:"qq"`
	WeChat string `xorm:"wechat varchar(100)" json:"wechat"`
}

func GetGlobalUsers() []*User {
	users := []*User{}
	err := adapter.Engine.Desc("created_time").Find(&users)
	if err != nil {
		panic(err)
	}

	return users
}

func GetUsers(owner string) []*User {
	users := []*User{}
	err := adapter.Engine.Desc("created_time").Find(&users, &User{Owner: owner})
	if err != nil {
		panic(err)
	}

	return users
}

func getUser(owner string, name string) *User {
	user := User{Owner: owner, Name: name}
	existed, err := adapter.Engine.Get(&user)
	if err != nil {
		panic(err)
	}

	if existed {
		return &user
	} else {
		return nil
	}
}

func GetUser(id string) *User {
	owner, name := util.GetOwnerAndNameFromId(id)
	return getUser(owner, name)
}

func UpdateUser(id string, user *User) bool {
	owner, name := util.GetOwnerAndNameFromId(id)
	if getUser(owner, name) == nil {
		return false
	}

	user.UpdateUserHash()

	affected, err := adapter.Engine.ID(core.PK{owner, name}).Cols("display_name", "avatar", "affiliation", "tag", "is_admin", "is_global_admin", "is_forbidden", "hash").Update(user)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func UpdateUserForOriginal(user *User) bool {
	affected, err := adapter.Engine.ID(core.PK{user.Owner, user.Name}).Cols("display_name", "password", "phone", "avatar", "is_forbidden", "hash", "pre_hash").Update(user)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddUser(user *User) bool {
	user.Id = util.GenerateId()

	organization := GetOrganizationByUser(user)
	user.UpdateUserPassword(organization)

	user.UpdateUserHash()
	user.PreHash = user.Hash

	affected, err := adapter.Engine.Insert(user)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddUsers(users []*User) bool {
	if len(users) == 0 {
		return false
	}

	organization := GetOrganizationByUser(users[0])
	for _, user := range users {
		user.UpdateUserPassword(organization)

		user.UpdateUserHash()
		user.PreHash = user.Hash
	}

	affected, err := adapter.Engine.Insert(users)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func AddUsersSafe(users []*User) bool {
	batchSize := 1000

	if len(users) == 0 {
		return false
	}

	affected := false
	for i := 0; i < (len(users)-1)/batchSize+1; i++ {
		start := i * batchSize
		end := (i + 1) * batchSize
		if end > len(users) {
			end = len(users)
		}

		tmp := users[start:end]
		fmt.Printf("Add users: [%d - %d].\n", start, end)
		if AddUsers(tmp) {
			affected = true
		}
	}

	return affected
}

func DeleteUser(user *User) bool {
	affected, err := adapter.Engine.ID(core.PK{user.Owner, user.Name}).Delete(&User{})
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func GetUserByField(organizationName string, field string, value string) *User {
	user := User{Owner: organizationName}
	existed, err := adapter.Engine.Where(fmt.Sprintf("%s=?", field), value).Get(&user)
	if err != nil {
		panic(err)
	}

	if existed {
		return &user
	} else {
		return nil
	}
}

func HasUserByField(organizationName string, field string, value string) bool {
	return GetUserByField(organizationName, field, value) != nil
}

func GetUserByFields(organization string, field string) *User {
	// check username
	user := GetUserByField(organization, "name", field)
	if user != nil {
		return user
	}

	// check email
	user = GetUserByField(organization, "email", field)
	if user != nil {
		return user
	}

	// check phone
	user = GetUserByField(organization, "phone", field)
	if user != nil {
		return user
	}

	return nil
}

func SetUserField(user *User, field string, value string) bool {
	if field == "password" {
		organization := GetOrganizationByUser(user)
		user.UpdateUserPassword(organization)
		value = user.Password
	}

	affected, err := adapter.Engine.Table(user).ID(core.PK{user.Owner, user.Name}).Update(map[string]interface{}{field: value})
	if err != nil {
		panic(err)
	}

	user = getUser(user.Owner, user.Name)
	user.UpdateUserHash()
	_, err = adapter.Engine.ID(core.PK{user.Owner, user.Name}).Cols("hash").Update(user)
	if err != nil {
		panic(err)
	}

	return affected != 0
}

func LinkUserAccount(user *User, field string, value string) bool {
	return SetUserField(user, field, value)
}

func GetUserField(user *User, field string) string {
	// https://socketloop.com/tutorials/golang-how-to-get-struct-field-and-value-by-name
	u := reflect.ValueOf(user)
	f := reflect.Indirect(u).FieldByName(field)
	return f.String()
}

func GetMaskedUser(user *User) *User {
	if user.Password != "" {
		user.Password = "***"
	}
	return user
}

func GetMaskedUsers(users []*User) []*User {
	for _, user := range users {
		user = GetMaskedUser(user)
	}
	return users
}

func calculateHash(user *User) string {
	s := strings.Join([]string{user.Id, user.Password, user.DisplayName, user.Avatar, user.Phone}, "|")
	return util.GetMd5Hash(s)
}

func (user *User) UpdateUserHash() {
	hash := calculateHash(user)
	user.Hash = hash
}

func (user *User) UpdateUserPassword(organization *Organization) {
	if organization.PasswordType == "salt" {
		user.Password = getSaltedPassword(user.Password, organization.PasswordSalt)
	}
}

func (user *User) GetId() string {
	return fmt.Sprintf("%s/%s", user.Owner, user.Name)
}
