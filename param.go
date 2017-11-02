package hellodatastore

// p := NewParam(r.URL)
type Param struct {
	Kind  string
	Key   string
	Value string
	LName string
	UName string
	Pass  string
	Email string
}

/*
func NewParam(url *url.URL) *Param {
	path := strings.Trim(url.Path, "/")
	//log.Printf("%#v", url.Path)
	//log.Printf("%#v", path)
	s := strings.Split(path, "/")
	//log.Printf("%#v", s)
	param := new(Param)
	log.Printf("%#v", param)
	if len(s) >= 2 {
		log.Printf("%#v", len(s))
		log.Printf("%#v", s[1])
		log.Printf("%#v", param.Kind)
		param.Kind = s[1]
		log.Printf("%#v", param.Kind)
	}
	if len(s) >= 3 {
		param.Key = s[2]
	}
	/*if len(s) >= 4 {
		param.Value = s[3]
	}*/
/*
	if len(s) >= 4 {
		param.LName = s[3]
	}
	if len(s) >= 5 {
		param.UName = s[4]
	}
	if len(s) >= 6 {
		param.Pass = s[5]
	}
	if len(s) >= 7 {
		param.Pass = s[6]

	}
	fmt.Println(param)
	log.Printf("%#v", param)
	return param
}
func (p *Param) HasKey() bool {
	return len(p.Key) > 0
}

/*func (p *Param) HasValue() bool {
	return len(p.Value) > 0
}*/
