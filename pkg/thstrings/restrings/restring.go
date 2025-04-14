package restrings

import (
	"fmt"
	"regexp"

	"github.com/k1nky/tookhook/pkg/thstrings"
)

type String struct {
	regexp *regexp.Regexp
	RegExp string
}

func New(s string) *String {
	return &String{
		RegExp: s,
	}
}

func (rs *String) UnmarshalYAML(unmarshal func(interface{}) error) error {
	return unmarshal(&rs.RegExp)
}

func (rs String) MarshalYAML() (interface{}, error) {
	return rs.RegExp, nil
}

func (rs *String) Compile() (err error) {
	if thstrings.IsEmpty(rs.RegExp) {
		return nil
	}
	if rs.regexp, err = regexp.Compile(rs.RegExp); err != nil {
		err = fmt.Errorf("invalid expr value: %w %w", err, thstrings.ErrCompile)
	}
	return
}

func (rs *String) Match(data []byte) bool {
	if rs.regexp == nil {
		return true
	}
	return rs.regexp.Match(data)
}

func (rs *String) IsEmpty() bool {
	return rs.RegExp == ""
}
