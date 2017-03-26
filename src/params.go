package main

import (
	"os"
	"flag"

	log "github.com/Sirupsen/logrus"
)

type Params struct {
		url string
		accessKey string
		secretKey string
		format string
		influxurl string
		influxdb string
		influxuser string
		influxpass string
		admin bool
		limit string
		refresh int
}

func (p *Params) init() {
	flag.StringVar(&p.url, "url", os.Getenv("RANCHER_URL"), "Rancher url. Or env RANCHER_URL")
	flag.StringVar(&p.accessKey, "accessKey", os.Getenv("RANCHER_ACCESS_KEY"), "Rancher access key. Or env RANCHER_ACCESS_KEY")
	flag.StringVar(&p.secretKey, "secretKey", os.Getenv("RANCHER_SECRET_KEY"), "Rancher secret key. Or env RANCHER_SECRET_KEY")
	flag.StringVar(&p.format, "format", "influx", "Output format. influx | json")
	flag.StringVar(&p.influxurl, "influxurl", "http://localhost:8086", "Influx url connection")
	flag.StringVar(&p.influxdb, "influxdb", "", "Influx db name")
	flag.StringVar(&p.influxuser, "influxuser", "", "Influx username")
	flag.StringVar(&p.influxpass, "influxpass", "", "Influx password")
	flag.BoolVar(&p.admin, "admin", false, "Admin flag to get stats")
	flag.StringVar(&p.limit, "limit", "1000", "Limit query results")
	flag.IntVar(&p.refresh, "refresh", 120, "Get metrics every refresh seconds")

	flag.Parse()

	p.checkParams()
}

func (p *Params) checkParams() {
	if ( len(p.url) == 0 || len(p.accessKey) == 0 || len(p.secretKey) == 0 ) { 
		flag.Usage()
		log.Info("Check your url, accessKey and/or secretKey params.")
		os.Exit(1) 
	}
	if p.format != "influx" && p.format != "json"{
		flag.Usage()
		log.Info("Check your format params. influx | yml | json ")
		os.Exit(1) 
	}
	if p.format == "influx" {
		if ( len(p.influxdb) == 0 || len(p.influxurl) == 0 ) { 
			flag.Usage()
			log.Info("Check your influxdb and/or influxurl params.")
			os.Exit(1) 
		}
	}
}