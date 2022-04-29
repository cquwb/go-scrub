package scrub

import (
	"reflect"
	"testing"
)

func TestMyScrub(t *testing.T) {
	type Student struct {
		Name     string
		Password string
		password string
	}
	s := Student{
		Name:     "wb",
		Password: "123456",
		password: "123",
	}
	str := MyScrub(&s, map[string]bool{"password": true})
	t.Log("new str is ", str)

	t.Log("nil is valid", reflect.ValueOf(nil).IsValid())
	var s1 *Student
	s1R := reflect.ValueOf(s1)
	t.Log("nil pointer is valid", s1R.IsValid(), " is zero ", s1R.IsZero(), "is nil", s1R.IsNil())

	type Getter interface {
		GetId() uint32
	}

	var i Getter
	iR := reflect.ValueOf(i)
	t.Log("i pointer is valid", iR.IsValid())
}
