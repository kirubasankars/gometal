package metal

import (
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type Metal struct {
	attributes map[string]interface{}
	parent     *Metal
	array      bool
	length     int
}

func NewMetal() *Metal {
	return &Metal{make(map[string]interface{}), nil, false, 0}
}

func (m *Metal) IsArray() bool {
	return m.array
}

func (m *Metal) Raw() interface{} {
	if m.array == false {
		object := make(map[string]interface{})
		for k, v := range m.attributes {
			if pmetal, ok := v.(*Metal); ok {
				object[k] = pmetal.Raw()
			} else {
				object[k] = v
			}
		}
		return object
	} else {
		array := make([]interface{}, m.length)
		for k, v := range m.attributes {
			if pmetal, ok := v.(*Metal); ok {
				idx, _ := strconv.Atoi(k[1:])
				array[idx] = pmetal.Raw()
			} else {
				idx, _ := strconv.Atoi(k[1:])
				array[idx] = v
			}
		}
		return array
	}
}

func (m *Metal) Get(property string) interface{} {
	dot := strings.Index(property, ".")
	if dot > -1 {
		var path, remaingPath string = property[0:dot], property[dot+1:]
		pathValue := m.Get(path)
		if pathValue != nil {
			if m1, ok := pathValue.(*Metal); ok {
				return m1.Get(remaingPath)
			} else {
				if remaingPath != "" {
					errors.New(path + " is not a object." + remaingPath + " can't be accessed")
				}
			}
		}
		if m.parent != nil && path == "$parent" {
			if m.parent.array == true {
				return m.parent.parent.Get(remaingPath)
			}
			return m.parent.Get(remaingPath)
		}
	}

	if property == "$length" {
		return m.length
	}

	if m.parent != nil && property == "$parent" {
		if m.parent.array == true {
			return m.parent.parent
		}
		return m.parent
	}

	if v, ok := m.attributes[property]; ok {
		return v
	}

	return nil
}

func (m *Metal) Set(property string, value interface{}) error {
	dot := strings.Index(property, ".")
	if dot > -1 {
		path, remaingPath := property[0:dot], property[dot+1:]
		pathValue := m.Get(path)
		if pathValue == nil {
			pathValue = NewMetal()
			m.Set(path, pathValue)
		}
		pathValue.(*Metal).Set(remaingPath, value)
	} else {
		if property[0] == '@' {
			if _, err := strconv.Atoi(property[1:]); err != nil {
				return errors.New("array should accessed by index")
			}
			m.array = true
		}
		var m1, ok = value.(*Metal)
		if ok {
			m1.parent = m
		}

		if _, ok = m.attributes[property]; !ok {
			m.length = m.length + 1
		}
		m.attributes[property] = value
	}
	return nil
}

func (m *Metal) Push(v interface{}) {
	m.array = true
	m.attributes["@"+strconv.Itoa(m.length)] = v
	m.length++
}

func (m *Metal) Properties() map[string]interface{} {
	var attrs = make(map[string]interface{})
	for k, v := range m.attributes {
		attrs[k] = v
	}
	return attrs
}

func (m *Metal) Parse(data []byte) {
	var d interface{}
	if err := json.Unmarshal(data, &d); err != nil {
		panic(err)
	}
	parseData("", d, m)
	return
}

func parseData(key interface{}, value interface{}, m *Metal) {
	switch x := value.(type) {
	case map[string]interface{}:
		for k1, v1 := range value.(map[string]interface{}) {
			if _, ok := m.attributes[k1]; !ok {
				m.length = m.length + 1
			}
			switch x := v1.(type) {
			case string, float32, float64, int, bool:
				m.attributes[k1] = v1
			default:
				sm := NewMetal()
				m.attributes[k1] = sm
				parseData(k1, v1, sm)
				_ = x
			}
		}
	case []interface{}:
		var array = value.([]interface{})
		m.array = true
		for idx := range array {
			var item = array[idx]
			if _, ok := m.attributes["@"+strconv.Itoa(idx)]; !ok {
				m.length = m.length + 1
			}
			switch x := item.(type) {
			case string, float32, float64, int, bool:
				m.attributes["@"+strconv.Itoa(idx)] = item
			default:
				var sm = NewMetal()
				sm.parent = m
				m.attributes["@"+strconv.Itoa(idx)] = sm
				parseData("", item, sm)
				_ = x
			}
		}
	default:
		_ = x
	}
}
