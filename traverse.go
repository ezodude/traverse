package traverse

import "strconv"

func Modify(in interface{}, detect func([]string, interface{}) bool, action func(path []string, value interface{}, in map[string]interface{}) error) {
	t := func(path []string, value interface{}) {
		condition := detect(path, value)
		if condition {
			action(path, value, in.(map[string]interface{}))
		}
	}
	Traverse(in, nil, t)
}

func Traverse(in interface{}, path []string, t func([]string, interface{})) {
	if path == nil {
		path = []string{}
	}

	switch in.(type) {
	case map[string]interface{}:
		for k, v := range in.(map[string]interface{}) {
			Traverse(v, append(path, k), t)
		}
	case []interface{}:
		for i, v := range in.([]interface{}) {
			Traverse(v, append(path, strconv.Itoa(i)), t)
		}
	default:
		t(path, in)
	}
}
