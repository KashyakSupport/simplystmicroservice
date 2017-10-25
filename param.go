package simplystmicroservice

import (
	"net/url"
	"strings"
)

// Param is a having kind key value
type Param struct {
	Kind  string
	Key   string
	Value string
}

// NewParam we can trim the url
func NewParam(url *url.URL) *Param {
	path := strings.Trim(url.Path, "/")
	s := strings.Split(path, "/")
	param := new(Param)

	if len(s) >= 2 {
		param.Kind = s[1]
	}
	if len(s) >= 3 {
		param.Key = s[2]
	}
	if len(s) >= 4 {
		param.Value = s[3]
	}
	return param
}

// HasKey is don't know
func (p *Param) HasKey() bool {
	return len(p.Key) > 0

}

// HasValue is dont know
func (p *Param) HasValue() bool {
	return len(p.Value) > 0
}
