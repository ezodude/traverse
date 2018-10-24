package traverse_test

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"testing"

	"github.com/ezodude/traverse"
)

func ReadJson(path string) (map[string]interface{}, error) {
	origin, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var m map[string]interface{}
	err = json.Unmarshal(origin, &m)
	if err != nil {
		return nil, err
	}

	return m, nil
}

func assertEqualMaps(tb testing.TB, actual, expected map[string]interface{}, msg string, v ...interface{}) {
	ac, err := json.Marshal(actual)
	if err != nil {
		panic(err)
	}
	ex, err := json.Marshal(expected)
	if err != nil {
		panic(err)
	}

	a, e := strings.Split(strings.TrimSpace(string(ac)), ""),
		strings.Split(strings.TrimSpace(string(ex)), "")

	sort.Strings(a)
	sort.Strings(e)
	condition := strings.EqualFold(strings.Join(a, ""), strings.Join(e, ""))

	if !condition {
		_, file, line, _ := runtime.Caller(1)
		tb.Fatalf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
	}
}

func TestTraversal(t *testing.T) {
	origin, err := ReadJson("testdata/origin.json")
	if err != nil {
		panic(err)
	}

	expected, err := ReadJson("testdata/transformed.json")
	if err != nil {
		panic(err)
	}

	detect := func(path []string, value interface{}) bool {
		parent := path[:len(path)-1]
		return parent[len(parent)-1:][0] == "Context"
	}

	action1 := func(path []string, value interface{}, in map[string]interface{}) error {
		if len(path) < 3 {
			return nil
		}
		node := in[path[0]]
		key := path[len(path)-1]

		switch node.(type) {
		case map[string]interface{}:
			node.(map[string]interface{})[key] = value
		case []interface{}:
			index, err := strconv.Atoi(path[1])
			if err != nil {
				return err
			}
			entry := node.([]interface{})[index]
			entry.(map[string]interface{})[key] = value
		}
		return nil
	}

	action2 := func(path []string, value interface{}, in map[string]interface{}) error {
		if len(path) < 3 {
			return nil
		}
		node := in[path[0]]

		switch node.(type) {
		case map[string]interface{}:
			delete(node.(map[string]interface{}), "Context")
		case []interface{}:
			index, err := strconv.Atoi(path[1])
			if err != nil {
				return err
			}
			entry := node.([]interface{})[index]
			delete(entry.(map[string]interface{}), "Context")
		}
		return nil
	}

	traverse.Modify(origin, detect, action1)
	traverse.Modify(origin, detect, action2)
	assertEqualMaps(t, origin, expected, "Modification unsuccessful \nactual:%+v\nexpected:%+v", origin, expected)
}
