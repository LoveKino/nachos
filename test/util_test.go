package nachos

import (
	"github.com/LoveKino/nachos/util"
	"testing"
)

func TestGetValueByJsonPath_List(t *testing.T) {
	var list []interface{}
	list = append(list, "str1")
	list = append(list, "str2")
	v1, e1 := util.GetValueByJsonPath(list, "0")
	v2, e2 := util.GetValueByJsonPath(list, "1")

	assertError(t, e1)
	assertError(t, e2)

	assertEqual(t, v1, "str1")
	assertEqual(t, v2, "str2")
}

func TestGetValueByJsonPath_Map(t *testing.T) {
	m := make(map[string]interface{})
	m["a"] = 10
	m["att"] = "12"

	v1, e1 := util.GetValueByJsonPath(m, "a")
	assertError(t, e1)
	assertEqual(t, v1, 10)

	v2, e2 := util.GetValueByJsonPath(m, "att")
	assertError(t, e2)
	assertEqual(t, v2, "12")
}

func TestGetValueByJsonPath_Map2(t *testing.T) {
	m := make(map[string]interface{})
	m["0"] = 10
	v1, e1 := util.GetValueByJsonPath(m, "0")
	assertError(t, e1)
	assertEqual(t, v1, 10)
}

func TestGetValueByJsonPath_MapDep(t *testing.T) {
	m := make(map[string]interface{})
	m2 := make(map[string]interface{})
	m2["n1"] = "abc"
	m["att"] = m2

	v1, e1 := util.GetValueByJsonPath(m, "att.n1")
	assertError(t, e1)
	assertEqual(t, v1, "abc")
}

func TestGetValueByJsonPath_MapDep2(t *testing.T) {
	m := make(map[string]interface{})
	var m2 []interface{}
	m2 = append(m2, "today")
	m["att"] = m2

	v1, e1 := util.GetValueByJsonPath(m, "att.0")
	assertError(t, e1)
	assertEqual(t, v1, "today")
}

func assertEqual(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("%s != %s", a, b)
	}
}

func assertError(t *testing.T, e error) {
	if e != nil {
		t.Error("Got Exception: " + e.Error())
	}
}
