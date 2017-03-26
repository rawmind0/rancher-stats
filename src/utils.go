package main

import (
	"encoding/json"
	"net/http"
	"time"

	influx "github.com/influxdata/influxdb/client/v2"
	log "github.com/Sirupsen/logrus"
)

type RacherMetric interface {
    getPoints(t time.Time) []influx.Point
    getData(url string, accessKey string, secretKey string, admin bool, limit string)
    printJson()
}

func getData(p Params, obj RacherMetric) {
	obj.getData(p.url, p.accessKey, p.secretKey, p.admin, p.limit)

	if p.format == "influx" {
		i := newInflux(p.influxurl, p.influxdb, p.influxuser, p.influxpass)
		t := time.Now()
		i.sendToInflux(obj.getPoints(t))
	} else if p.format == "json" {
		obj.printJson()
	}
}

type Pagination struct {
	 Next 		string `json:"next"`
	 Partial 	bool `json:"partial"`
}
	
func check(e error, m string) {
    if e != nil {
		log.Error("[Error]: ", m , e)
	}
}

func getJSON(url string, accessKey string, secretKey string, target interface{}) error {

	start := time.Now()

	log.Info("Connecting to ", url)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.SetBasicAuth(accessKey, secretKey)
	resp, err := client.Do(req)

	if err != nil {
		log.Error("Error Collecting JSON from API: ", err)
		panic(err)
	}

	respFormatted := json.NewDecoder(resp.Body).Decode(target)

	// Timings recorded as part of internal metrics
 	log.Info("Time to get json: ", float64((time.Since(start))/ time.Millisecond), " ms")

	// Close the response body, the underlying Transport should then close the connection.
	resp.Body.Close()

	// return formatted JSON
	return respFormatted
}


