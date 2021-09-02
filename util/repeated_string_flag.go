package util

import (
	"strings"
)

type RepeatedStringFlag struct {
	vs           []string
	isSet        bool
	defaultValue string
}

func NewRepeatedStringFlag(defaultValue string) RepeatedStringFlag {
	f := RepeatedStringFlag{
		vs:           []string{},
		isSet:        true,
		defaultValue: defaultValue,
	}
	return f
}

func (f *RepeatedStringFlag) Type() string {
	return "string"
}

func (f *RepeatedStringFlag) String() string {
	if len(f.vs) == 0 {
		return f.defaultValue
	}
	return strings.Join(f.vs, ",")
}

func (f *RepeatedStringFlag) Strings() []string {
	if len(f.vs) == 0 {
		return []string{f.defaultValue}
	}
	return f.vs
}

func (f *RepeatedStringFlag) Set(s string) error {
	f.vs = append(f.vs, s)
	f.isSet = true
	return nil
}

func (f *RepeatedStringFlag) IsSet() bool {
	return f.isFlagSet()
}

func (f *RepeatedStringFlag) isFlagSet() bool {
	return f.isSet
}
