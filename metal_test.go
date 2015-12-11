package metal_test

import (
	"encoding/json"
	"gravity/metal"
	"bytes"
	"fmt"
	"testing"
)

func TestSet(t *testing.T) {
	var m = metal.NewMetal()
	m.Set("student.name", "Dev")
	var v, _ = m.Get("student.name")
	if v != "Dev" {
		t.Error("unable to get student.name")
	}
}

func TestSet_array(t *testing.T) {
	var m = metal.NewMetal()
	m.Set("student.marks.@0", 100)
	var v, _ = m.Get("student.marks.@0")
	if v != 100 {
		t.Error("unable to get student.marks.@0")
	}
}

func TestSet_arrayobject(t *testing.T) {
	var m = metal.NewMetal()
	m.Set("student.marks.@0.tamil", 1)
	var v, _ = m.Get("student.marks.@0.tamil")
	if v != 1 {
		t.Error("unable to get student.marks.@0.tamil")
	}
}

func TestSet_array_object_json(t *testing.T) {
	var m = metal.NewMetal()
	m.Set("student.marks.@0", 10)
	m.Set("student.marks.@1", 20)
	m.Set("student.marks.@2.a", 20)
	json, _ := json.Marshal(m.JSON())
	if string(json) != string([]byte(`{"student":{"marks":[10,20,{"a":20}]}}`)) {
		t.Error("got this", string(json), "expecting ...", string([]byte(`{"student":{"marks":[10,20,{"a": 20}]}}`)))
	}
}

func TestSet_Getparent(t *testing.T) {
	var m = metal.NewMetal()
	m.Set("student.name", "Dev")
	m.Set("student.marks.@0", 10)
	m.Set("student.marks.@2.a", 20)
	var v1, _ = m.Get("student.marks.$parent.name")
	var v, _ = m.Get("student.marks.@2.$parent.name")
	if v1 != "Dev" || v != "Dev" {
		t.Error("$parent failed", v1, v)
	}
}

func TestSet_Getlength(t *testing.T) {
	var m = metal.NewMetal()
	m.Set("student.name", "Dev")
	m.Set("student.marks.@0", 10)
	m.Set("student.marks.@1", 10)
	m.Set("student.marks.@2.a", 20)
	var v1, _ = m.Get("student.marks.$length")
	var v, _ = m.Get("student.$length")
	if v1 != 3 || v != 2 {
		t.Error("$length failed", v1, v)
	}
}

func TestParse(t *testing.T) {
	var m = metal.NewMetal();
	m.Set("student.name", 1)
	m.Set("student.marks.@0", 100)
	m.Set("student.marks.@1", 100)
	m.Set("student.marks.@2.a", 100)
	var b = new(bytes.Buffer)
	enc := json.NewEncoder(b)
	var _ = enc.Encode(m.JSON())
	fmt.Println(b);

	var m1 = metal.NewMetal();
	m1.Parse(b.Bytes());
	var b1 = new(bytes.Buffer)
	enc1 := json.NewEncoder(b1)
	var _ = enc1.Encode(m1.JSON())
	fmt.Println(b1);

	if string(b.Bytes()) != string(b1.Bytes()) {
		t.Error("Parsing is not working")
	}
}
