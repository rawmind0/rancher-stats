package main

import (
	"fmt"
	"time"
	"sort"
	"strings"
	"encoding/json"

	log "github.com/Sirupsen/logrus"
	influx "github.com/influxdata/influxdb/client/v2"
)

type projectData struct {
	Id          	string `json:"id"`
	Name        	string `json:"name"`
	State       	string `json:"state"`
	Orchestration   string `json:"orchestration"`
	Containers 		int
	Hosts 			int
}

type ProjectResult struct {
	Total 		int `json:"total"`
	Active 		int `json:"active"`
	Hashost		int `json:"hashost"`
	Containers 	ContainerResult `json:"containers"`
	Hosts 		HostResult `json:"hosts"`
	Orchestrators 	[]Orchestration `json:"orchestrators"`
}

func (r *ProjectResult) getOrchestratorId(s string) int {
	for i := range r.Orchestrators {
		if r.Orchestrators[i].Orchestration == s {
			return i
		}
	}
	return -1
}

func (r *ProjectResult) addActiveOrchestration(id int){
	r.Orchestrators[id].addActive()
}

func (r *ProjectResult) addTotalOrchestration(id int){
	r.Orchestrators[id].addTotal()
}

func (r *ProjectResult) addHashostOrchestration(id int){
	r.Orchestrators[id].addHashost()
}

func (r *ProjectResult) addHostsOrchestration(id int, h int){
	r.Orchestrators[id].addHosts(h)
}

func (r *ProjectResult) addContainersOrchestration(id int, c int){
	r.Orchestrators[id].addContainers(c)
}

func (r *ProjectResult) addOrchestration(s string) int {
	var k = &Orchestration{ Orchestration: s, Total: 0, Active: 0, Hashost: 0, Usage:0 }
	r.Orchestrators = append(r.Orchestrators, *k)
	return len(r.Orchestrators)-1
}

func (a *ProjectResult) addDataOrchestration(o string, state string, host int, cont int) {
	if len(o) > 0 {
		id := a.getOrchestratorId(o)
		if id < 0 { id = a.addOrchestration(o) }

		a.addTotalOrchestration(id)
		if state == "active" {a.addActiveOrchestration(id)}
		if host > 0 {a.addHashostOrchestration(id)}
		a.addHostsOrchestration(id, host)
		a.addContainersOrchestration(id, cont)
	}
}

func (a *ProjectResult) addUsageOrchestration(t int) {
    for id := range a.Orchestrators {
    	a.Orchestrators[id].setUsage(t)
	}

}

type Projects struct {
	Page 			[]projectData `json:"data"`
	Pagination 		Pagination `json:"pagination"`
	Data 			[]projectData
	Containers		Containers
	Hosts			Hosts 
	Results			ProjectResult
	Points 			[]influx.Point
}

type By func(p1, p2 *projectData) bool

type projectSorter struct {
	Data 	[]projectData
	by      func(p1, p2 *projectData) bool // Closure used in the Less method.
}

func (by By) Sort(d []projectData) {
	ps := &projectSorter{
		Data: 	d,
		by:     by, 
	}
	sort.Sort(ps)
}

func (s *projectSorter) Len() int {
	return len(s.Data)
}

func (s *projectSorter) Swap(i, j int) {
	s.Data[i], s.Data[j] = s.Data[j], s.Data[i]
}

func (s *projectSorter) Less(i, j int) bool {
	return s.by(&s.Data[i], &s.Data[j])
}

func (a *Projects) sortByHosts() {

	host := func(p1, p2 *projectData) bool {
		return p1.Hosts < p2.Hosts
	}

	By(host).Sort(a.Data)
}

func (a *Projects) sortByContainers() {

	container := func(p1, p2 *projectData) bool {
		return p1.Containers < p2.Containers
	}

	By(container).Sort(a.Data)
}

func (r *Projects) setTotal(i int){
	r.Results.Total = i
}

func (r *Projects) getTotal() int {
	return r.Results.Total
}

func (r *Projects) addTotal(){
	r.Results.Total ++
}

func (r *Projects) getActive() int {
	return r.Results.Active
}

func (r *Projects) addActive(){
	r.Results.Active ++
}

func (r *Projects) getHashost() int {
	return r.Results.Hashost
}

func (r *Projects) addHashost(){
	r.Results.Hashost ++
}

func (r *Projects) getDataContainers(url string, accessKey string, secretKey string, admin bool, limit string){

	r.Containers.getData(url, accessKey, secretKey, admin, limit)
	r.Results.Containers = r.Containers.getResult()
}

func (r *Projects) getActiveContainers() int {
	return r.Containers.getActive()
}

func (r *Projects) getTotalContainers() int {
	return r.Containers.getTotal()
}

func (r *Projects) countContainerById(id string) int {
	return r.Containers.countById(id)
}

func (r *Projects) getDataHosts(url string, accessKey string, secretKey string, admin bool, limit string){

	r.Hosts.getData(url, accessKey, secretKey, admin, limit)
	r.Results.Hosts = r.Hosts.getResult()
}

func (r *Projects) getActiveHosts() int {
	return r.Hosts.getActive()
}

func (r *Projects) getTotalHosts() int {
	return r.Hosts.getTotal()
}

func (r *Projects) countHostById(id string) int {
	return r.Hosts.countById(id)
}

func (a *Projects) getData(url string, accessKey string, secretKey string, admin bool, limit string) {
	var uri string

	if admin {
		uri = "/projects?all=true&limit="+limit
	} else {
		uri = "/projects?limit="+limit
	}

	a.getAllPages(url+uri, accessKey, secretKey)

	a.getDataContainers(url, accessKey, secretKey, admin, limit)

	a.getDataHosts(url, accessKey, secretKey, admin, limit)

	for i := range a.Data {
		if len(a.Data[i].Name) == 0 {a.Data[i].Name = "NONAME"}
		a.Data[i].Name = strings.Replace(a.Data[i].Name, " ", "+", -1)
		if a.Data[i].State == "active" {a.addActive()}
		a.Data[i].Containers = a.countContainerById(a.Data[i].Id)
		a.Data[i].Hosts = a.countHostById(a.Data[i].Id)
		if a.Data[i].Hosts > 0 { a.addHashost() }
		a.Results.addDataOrchestration(a.Data[i].Orchestration,a.Data[i].State,a.Data[i].Hosts,a.Data[i].Containers)
    }

    a.setTotal(len(a.Data))
    a.Results.addUsageOrchestration(a.getTotal())

}

func (a *Projects) getAllPages(url string, accessKey string, secretKey string) {
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

func (a *Projects) addPoint(m *influx.Point){
	a.Points = append(a.Points, *m)
}

func (a *Projects) getPointsByHosts(ti time.Time) {
	n := "projects"

	a.sortByHosts()
	l := len(a.Data) - 1
	for i := l ; i > (l - 20) && i >= 0 ; i-- {
		v := map[string]interface{}{
            "hosts":  a.Data[i].Hosts,
            "containers": a.Data[i].Containers,
    	}
    	t := map[string]string{
    		"top": "hosts",
    		"id": a.Data[i].Id,
    		"name": a.Data[i].Name,
    		"orchestration": a.Data[i].Orchestration,
    		"state": a.Data[i].State,
    	}
		m, err := influx.NewPoint(n,t,v,ti)
		check(err, "Getting projects points by hosts")
		a.addPoint(m)
	}
}

func (a *Projects) getPointsByContainers(ti time.Time) {
	n := "projects"

	a.sortByContainers()
	l := len(a.Data) - 1
	for i := l ; i > (l - 20) && i >= 0 ; i-- {
		v := map[string]interface{}{
            "hosts":  a.Data[i].Hosts,
            "containers": a.Data[i].Containers,
    	}
    	t := map[string]string{
    		"top": "containers",
    		"id": a.Data[i].Id,
    		"name": a.Data[i].Name,
    		"orchestration": a.Data[i].Orchestration,
    		"state": a.Data[i].State,
    	}
		m, err := influx.NewPoint(n,t,v,ti)
		check(err, "Getting projects points by containers")
		a.addPoint(m)
	}
}

func (a *Projects) getPointsByOrchestrators(ti time.Time) {
	n := "projects"

	for id := range a.Results.Orchestrators {
		v := map[string]interface{}{
            "usage":  a.Results.Orchestrators[id].getUsage(),
            "total": a.Results.Orchestrators[id].getTotal(),
            "hashost": a.Results.Orchestrators[id].getHashost(),
            "hosts": a.Results.Orchestrators[id].getHosts(),
            "containers": a.Results.Orchestrators[id].getContainers(),
    	}
    	t := map[string]string{
    		"orchestration": a.Results.Orchestrators[id].getOrchestration(),
    	}

		m, err := influx.NewPoint(n,t,v,ti)
		check(err, "Getting projects points by orchestrators")
		a.addPoint(m)	
	}
}

func (a *Projects) getPointsByData(ti time.Time) {
	n := "projects"

	for i := range a.Data {
		v := map[string]interface{}{
            "hosts":  a.Data[i].Hosts,
            "containers": a.Data[i].Containers,
    	}
    	t := map[string]string{
    		"id": a.Data[i].Id,
    		"orchestration": a.Data[i].Orchestration,
    		"state": a.Data[i].State,
    	}

		m, err := influx.NewPoint(n,t,v,ti)
		check(err, "Getting projects points by data")
		a.addPoint(m)	
	}
}

func (a *Projects) getPoints(ti time.Time) []influx.Point {
	var t = map[string]string{}
	var n = "projects"
	v := map[string]interface{}{
        "active":  a.getActive(),
        "total": a.getTotal(),
        "hashost": a.getHashost(),
        "containers": a.getTotalContainers(),
        "hosts": a.getTotalHosts(),
    }

	m, err := influx.NewPoint(n,t,v,ti)
	check(err, "Getting projects points")
	a.addPoint(m)

	a.getPointsByData(ti)
	a.getPointsByOrchestrators(ti)
	a.getPointsByContainers(ti)
	a.getPointsByHosts(ti)

	return a.Points
	
}

func (a *Projects) printJson() {
	j, err := json.Marshal(a.Results)
	if err != nil {
    	log.Error("json")
	}
	fmt.Println(string(j))

}

