package main

import (
	log "github.com/Sirupsen/logrus"
)

type HostData struct {
	Id          	string `json:"id"`
	State       	string `json:"state"`
	Accountid       string `json:"accountId"`
}

type HostResult struct {
	Total 		int `json:"total"`
	Active 		int `json:"active"`
}

type Hosts struct {
	Page 			[]HostData `json:"data"`
	Pagination 		Pagination `json:"pagination"`
	Data 			[]HostData
	Result 			HostResult
}

func (c *Hosts) setTotal(t int){
	c.Result.Total = t
}

func (c *Hosts) getTotal() int {
	return c.Result.Total
}

func (c *Hosts) addTotal(){
	c.Result.Total ++
}

func (c *Hosts) addActive(){
	c.Result.Active ++
}

func (c *Hosts) getActive() int {
	return c.Result.Active
}

func (c *Hosts) getResult() HostResult {
	return c.Result
}

func (a *Hosts) getData(url string, accessKey string, secretKey string, admin bool, limit string) {
	var uri string

	if admin {
		uri = "/hosts?all=true&limit="+limit
	} else {
		uri = "/hosts?limit="+limit
	}

	a.getAllPages(url+uri, accessKey, secretKey)

	a.setTotal(len(a.Data))
	for i := range a.Data {
		if a.Data[i].State == "active" {a.addActive()}
    }

}

func (a *Hosts) getAllPages(url string, accessKey string, secretKey string) {
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

func (a *Hosts) countById(id string) int {
	count := 0
	for i := range a.Data {
		if a.Data[i].Accountid == id {count++}
    }
    return count
}

