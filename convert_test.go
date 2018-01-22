package struct_converter

import (
	"testing"
	"fmt"
	"encoding/json"
)

type User struct {
	FirstName  string
	LastName   string
	Tags       []int
	Age			int
}
type UserInfo struct {
	FirstName  *string
	LastName   *string
	Tags       []int
	Age         int64
}


func TestConvert_Convert(t *testing.T) {
	u := User{
		"a",
		"b",
		[]int{1,2,3},
		18,
	}

	conv := New(new(User), new(UserInfo))
	uinfo := UserInfo{}
	fmt.Println(uinfo)
	conv.Convert(&u, &uinfo)
	bytes, _ := json.Marshal(uinfo)
	fmt.Println(string(bytes))

	firstName := u.FirstName
	uinfo2 := UserInfo{
		&(firstName),
		&u.LastName,
		[]int{1,2,3},
		18,
	}
	conv = New(new(UserInfo), new(User))
	u2 := User{}
	fmt.Println(u2)
	conv.Convert(&uinfo2, &u2)
	bytes, _ = json.Marshal(u2)
	fmt.Println(string(bytes))


}