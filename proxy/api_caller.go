package proxy

import (
	"log"
	"regexp"
	"sort"
	"strings"
)

//var API_PREF string = "api_pref"

const (
	API_PREF_PARMA_NAME  string = "api_pref"
	API_PREF_TYPE_REQ           = "req"
	API_PREF_TYPE_COOKIE        = "cookie"
	API_PREF_TYPE_HEADER        = "header"
)

var prefTypes = []string{API_PREF_TYPE_REQ, API_PREF_TYPE_COOKIE, API_PREF_TYPE_HEADER}

type Caller []*CallerItem

type CallerItem struct {
	Note   string         `json:"note"`
	Ip     string         `json:"ip"`
	IpReg  *regexp.Regexp `json:"-"`
	Enable bool           `json:"enable"`
	Pref   []string       `json:"pref"`
	Ignore []string       `json:"ignore"`
}

func NewCaller() Caller {
	return make([]*CallerItem, 0)
}

func NewCallerItem(ip string) (*CallerItem, error) {
	item := &CallerItem{
		Ip:     ip,
		Pref:   make([]string, 0),
		Ignore: make([]string, 0),
	}
	var err error
	err = item.Init()

	return item, err
}

func NewCallerItemMust(ip string) *CallerItem {
	item, _ := NewCallerItem(ip)
	item.Enable = true
	return item
}

func (citem *CallerItem) Init() (err error) {
	citem.IpReg, err = regexp.Compile(strings.Replace(strings.Replace(citem.Ip, ".", `\.`, -1), "*", `\d+`, -1))
	if err != nil {
		log.Println("ip wrong:", citem.Ip)
	}
	if citem.Ignore == nil {
		citem.Ignore = make([]string, 0)
	}
	return err
}

func (citem *CallerItem) IsHostIgnore(host_name string) bool {
	return In_StringSlice(host_name, citem.Ignore)
}

const IP_ALL string = "*.*.*.*"

func (caller *Caller) Init() (err error) {
	has_all := false
	for _, citem := range *caller {
		err := citem.Init()
		if err != nil {
			return err
		}
		if citem.Ip == IP_ALL {
			has_all = true
		}
	}
	if !has_all {
		citem := NewCallerItemMust(IP_ALL)
		citem.Note = "default all"
		citem.Enable = true
		citem.Init()
		caller.AddNewCallerItem(citem)
	}
	return nil
}

func (caller *Caller) GetPrefHostName(allowNames []string, cpf *CallerPrefConf) string {

	if len(allowNames) == 0 || len(*caller) == 0 {
		return StrSliceRandItem(allowNames)
	}

	for _, prefType := range prefTypes {
		if len(cpf.prefHostName[prefType]) > 0 {
			pref := StrSliceIntersectGetOne(cpf.prefHostName[prefType], allowNames)
			if pref != "" {
				return pref
			}
		}
	}
	item := caller.getCallerItemByIp(cpf.ip)
	if item != nil && len(item.Pref) > 0 {
		pref := StrSliceIntersectGetOne(item.Pref, allowNames)
		if pref != "" {
			return pref
		}
	}
	return StrSliceRandItem(allowNames)
}

func (caller Caller) Sort() {
	sort.Sort(caller)
}
func (caller Caller) Len() int {
	return len(caller)
}

/**
*让 127.0.0.1 排在127.0.0.* 前面
 */
func (caller Caller) Less(i, j int) bool {
	aPos := strings.Index(caller[i].Ip, "*")
	if aPos == -1 {
		return true
	}
	bPos := strings.Index(caller[j].Ip, "*")
	if bPos == -1 {
		return false
	}

	return aPos > bPos
}

func (caller Caller) Swap(i, j int) {
	caller[i], caller[j] = caller[j], caller[i]
}

var Default_Caller= &CallerItem{Ip: IP_ALL, Enable: true, Note: "default"}

func init() {
	Default_Caller.Init()
}

func (caller Caller) getCallerItemByIp(ip string) *CallerItem {
	for _, item := range caller {
		if !item.Enable {
			continue
		}
		if item.Ip == ip || item.IpReg.MatchString(ip) {
			return item
		}
	}
	return Default_Caller
}

func (caller *Caller) AddNewCallerItem(item *CallerItem) {
	*caller = append(*caller, item)
	caller.Sort()
}
