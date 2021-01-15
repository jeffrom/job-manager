// Package config handles job-manager configuration helpers.
package config

import (
	"reflect"

	"github.com/imdario/mergo"
	"github.com/kelseyhightower/envconfig"
)

// MergeEnvFlags takes a pointer to a struct populated by flags, and returns a
// copy with values populated by first merging flags, then environment, then
// defaults.
func MergeEnvFlags(cfgFlags interface{}, defs interface{}) (interface{}, error) {
	cfgEnv := reflect.New(reflect.ValueOf(cfgFlags).Elem().Type()).Interface()
	if err := envconfig.Process("", cfgEnv); err != nil {
		return nil, err
	}
	// fmt.Printf("flags: %+v\n", cfgFlags)
	// fmt.Printf("env: %+v\n", cfgEnv)

	cfg := reflect.New(reflect.ValueOf(cfgFlags).Elem().Type()).Interface()
	if err := mergo.Merge(cfg, cfgFlags); err != nil {
		return nil, err
	}
	if err := mergo.Merge(cfg, cfgEnv); err != nil {
		return nil, err
	}
	if err := mergo.Merge(cfg, defs); err != nil {
		return nil, err
	}
	return cfg, nil
}
