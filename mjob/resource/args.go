package resource

import "github.com/ghodss/yaml"

func ParseCLIArgs(args []string) ([]interface{}, error) {
	ifaces := make([]interface{}, len(args))
	for i, arg := range args {
		var iarg interface{}
		if err := yaml.Unmarshal([]byte(arg), &iarg); err != nil {
			return nil, err
		}
		ifaces[i] = iarg
	}
	return ifaces, nil
}
