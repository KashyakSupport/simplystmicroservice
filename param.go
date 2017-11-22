package hellodatastore

import (
	"fmt"
	"log"
	"net/url"
	"strings"
)

// Param is cotaning user details
type Param struct {
	Kind  string
	UName string
	Fname string
	LName string
	Email string
	Pass  string
}

// NewParam is cotaning new user details
func NewParam(url *url.URL) *Param {
	path := strings.Trim(url.Path, "/")
	s := strings.Split(path, "/")
	param := new(Param)

	if len(s) >= 2 {
		param.Kind = s[1]
		log.Printf("%#v", param.Kind)
	}
	if len(s) >= 3 {
		param.UName = s[2]
	}
	if len(s) >= 4 {
		param.Fname = s[3]
	}
	if len(s) >= 5 {
		param.LName = s[4]
	}
	if len(s) >= 6 {
		param.Email = s[5]
	}
	if len(s) >= 7 {
		param.Pass = s[6]

	}
	fmt.Println(param)
	log.Printf("%#v", param)
	return param
}

// HasKey is used to check user details
func (p *Param) HasKey() bool {
	return len(p.UName) > 0
}

/*func (p *Param) HasValue() bool {
	return len(p.Value) > 0
}*/
