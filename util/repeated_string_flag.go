package util

import (
	"strings"
)

type RepeatedStringFlag struct {
	vs           []string
	defaultValue string
}

func NewRepeatedStringFlag(defaultValue string) RepeatedStringFlag {
	f := RepeatedStringFlag{
		vs:           []string{},
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
	return nil
}

func (f *RepeatedStringFlag) IsSet() bool {
	return true
}
