package main

import (
	log "github.com/Sirupsen/logrus"
)

type ContainerData struct {
	Id          string `json:"id"`
	State       string `json:"state"`
	Accountid   string `json:"accountId"`
}

type ContainerResult struct {
	Total 		int `json:"total"`
	Active 		int `json:"active"`
}

type Containers struct {
	Page 			[]ContainerData `json:"data"`
	Pagination 		Pagination `json:"pagination"`
	Data 			[]ContainerData
	Result 			ContainerResult
}

func (c *Containers) setTotal(t int){
	c.Result.Total = t
}

func (c *Containers) getTotal() int {
	return c.Result.Total
}

func (c *Containers) addTotal(){
	c.Result.Total ++
}

func (c *Containers) addActive(){
	c.Result.Active ++
}

func (c *Containers) getActive() int {
	return c.Result.Active
}

func (c *Containers) getResult() ContainerResult {
	return c.Result
}

func (a *Containers) getData(url string, accessKey string, secretKey string, admin bool, limit string) {
	var uri string

	if admin {
        uri = "/containers?all=true&limit="+limit
    } else {
        uri = "/containers?limit="+limit
    }

    a.getAllPages(url+uri, accessKey, secretKey)

	a.setTotal(len(a.Data))
	for i := range a.Data {
		if a.Data[i].State == "active" {a.addActive()}
    }

}

func (a *Containers) getAllPages(url string, accessKey string, secretKey string) {
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

func (a *Containers) countById(id string) int {
	count := 0
	for i := range a.Data {
		if a.Data[i].Accountid == id {count++}
    }
    return count
}


