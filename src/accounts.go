package main

import (
	"fmt"
	"time"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	influx "github.com/influxdata/influxdb/client/v2"
)

type AccountData struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	State       string `json:"state"`
	Kind        string `json:"kind"`
}

type AccountKind struct {
	Kind		string `json:"kind"`
	Total 		int `json:"total"`
	Active 		int `json:"active"`
}

type AccountResult struct {
	Kinds		[]AccountKind `json:"kinds"`
	Total 		int `json:"total"`
	Active 		int `json:"active"`
}

func (k *AccountKind) setTotal(t int){
	k.Total = t
}

func (k *AccountKind) addTotal(){
	k.Total ++
}

func (k *AccountKind) addActive(){
	k.Active ++
}

func (r *AccountResult) setTotal(t int){
	r.Total = t
}

func (r *AccountResult) addTotal(){
	r.Total ++
}

func (r *AccountResult) addTotalByKind(s string){
	for i := range r.Kinds {
		if r.Kinds[i].Kind == s {
			r.Kinds[i].addTotal()
			break
		}
	}
}

func (r *AccountResult) addActive(){
	r.Active ++
}

func (r *AccountResult) addActiveByKind(s string){
	for i := range r.Kinds {
		if r.Kinds[i].Kind == s {
			r.Kinds[i].addActive()
			break
		}
	}
}

func (r *AccountResult) addActiveKind(id int){
	r.Kinds[id].addActive()
}

func (r *AccountResult) addTotalKind(id int){
	r.Kinds[id].addTotal()
}

func (r *AccountResult) existKind(s string) bool {
	for i := range r.Kinds {
		if r.Kinds[i].Kind == s {
			return true
		}
	}
	return false
}

func (r *AccountResult) getKindId(s string) int {
	for i := range r.Kinds {
		if r.Kinds[i].Kind == s {
			return i
		}
	}
	return -1
}

func (r *AccountResult) addKind(s string) int {
	var k = &AccountKind{ Kind: s, Total: 0, Active: 0,}
	r.Kinds = append(r.Kinds, *k)
	return len(r.Kinds)-1
}

// Data is used to store data from all the relevant endpoints in the API
type Accounts struct {
	Page 			[]AccountData `json:"data"`
	Pagination 		Pagination `json:"pagination"`
	Data 			[]AccountData
	Result			*AccountResult `json:"accounts"`
	Points 			[]influx.Point
}

func newAccounts() *Accounts {
	var a = &Accounts{
		Result: &AccountResult{
			Kinds: []AccountKind{},
		},
	}
	return a
}

func (a *Accounts) getData(url string, accessKey string, secretKey string, admin bool, limit string) {
	var uri string

	if admin {
        uri = "/accounts?all=true&limit="+limit
    } else {
        uri = "/accounts?limit="+limit
    }

	a.getAllPages(url+uri, accessKey, secretKey)

	a.Result.setTotal(len(a.Data))
	for i := range a.Data {
		//switch a.Data[i].Kind {
		//	case "admin", "user", "superadmin", "project", "agent":
		//		a.Result.addTotalByType(a.Data[i].Kind)
				//a.Result.Kind[a.Data[i].Kind]["total"]++
		//		if a.Data[i].State == "active" {a.Result.addActiveByType(a.Data[i].Kind)}
		//}
		if len(a.Data[i].Kind) > 0 {

			id := a.Result.getKindId(a.Data[i].Kind)
			if id < 0 { id = a.Result.addKind(a.Data[i].Kind) }

			if a.Data[i].State == "active" {
				a.Result.addActiveKind(id)
				a.Result.addActive()
			}
			a.Result.addTotalKind(id)
		}
    }

}

func (a *Accounts) getAllPages(url string, accessKey string, secretKey string) {
	for i := 0; i==0 || a.Pagination.Partial ; i++ {
		err := getJSON(url, accessKey, secretKey, a)
		if err != nil {
			log.Error("Error getting JSON from URL ", url)
		}
		log.Debugf("JSON Fetched for: "+url+": ", a)
		if a.Pagination.Partial {
			url = a.Pagination.Next
			log.Info("Next page..."+a.Pagination.Next)
		}
		a.Data = append(a.Data, a.Page...)
	}
}

func (a *Accounts) addPoint(m *influx.Point){
	a.Points = append(a.Points, *m)
}

func (a *Accounts) getPointsByKind(ti time.Time) {
	n := "accounts"

	for index := range a.Result.Kinds {
		v := map[string]interface{}{
            "active":  a.Result.Kinds[index].Active,
            "total": a.Result.Kinds[index].Total,
    	}
    	t := map[string]string{
    		"kind": a.Result.Kinds[index].Kind,
    	}

		m, err := influx.NewPoint(n,t,v,ti)
		check(err, "Getting accounts points by king")
		a.addPoint(m)	
	}
}

func (a *Accounts) getPoints(ti time.Time) []influx.Point {

	var t = map[string]string{}
	var n = "accounts"

	v := map[string]interface{}{
        "active":  a.Result.Active,
        "total": a.Result.Total,
    }

	m, err := influx.NewPoint(n,t,v,ti)
	check(err, "Getting accounts points")
	a.addPoint(m)

	a.getPointsByKind(ti)

	return a.Points
	
}

func (a *Accounts) printJson() {
	j, err := json.Marshal(a.Result)
	if err != nil {
    	log.Error("json")
	}
	fmt.Println(string(j))

}


