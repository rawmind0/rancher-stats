[![](https://images.microbadger.com/badges/image/rawmind/rancher-stats.svg)](https://microbadger.com/images/rawmind/rancher-stats "Get your own image badge on microbadger.com")

rancher-stats
==============

This image run rancher-stats app. It comes from [rawmind/alpine-base][alpine-base].

## Build

```
docker build -t rawmind/rancher-stats:<version> .
```

## Versions

- `0.0.1` [(Dockerfile)](https://github.com/rawmind0/rancher-stats/blob/0.0.1/Dockerfile)


## Usage

This image run rancher-stats service. Rancher-stats get metrics from your rancher server and send them to a influx in order to be viewed by a grafana. 

```
Usage of rancher-stats:
  -accessKey string
    	Rancher access key. Or env RANCHER_ACCESS_KEY
  -admin
    	Admin flag to get stats
  -format string
    	Output format. influx | json (default "influx")
  -influxdb string
    	Influx db name
  -influxpass string
    	Influx password
  -influxurl string
    	Influx url connection (default "http://localhost:8086")
  -influxuser string
    	Influx username
  -limit string
    	Limit query results (default "1000")
  -refresh int
    	Get metrics every refresh seconds (default 120)
  -secretKey string
    	Rancher secret key. Or env RANCHER_SECRET_KEY
  -url string
    	Rancher url. Or env RANCHER_URL (default "https://try.rancher.com/v2-beta")
```

[alpine-base]: https://github.com/rawmind0/alpine-base
