package main

import (
	"os"
	"time"
	"sync"
	"os/signal"

	log "github.com/Sirupsen/logrus"
)

func get(p Params, wg *sync.WaitGroup) {
	wg.Add(3)
	go func() {
		defer wg.Done()
		go func() {
			defer wg.Done()
			var acc = newAccounts()
			getData(p, acc)
		}()
		go func() {
			defer wg.Done()
			var pro Projects
			getData(p, &pro)
		}()
	}()
}

func main() {
	var params Params 
	var wg sync.WaitGroup

	params.init()

	ticker := time.NewTicker(time.Second * time.Duration(params.refresh))
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, os.Kill)

	get(params, &wg)

	for {
        select {
        case <-ticker.C:
            get(params, &wg)
        case <- exit:
        	log.Info("Exit signal detected. Waiting for running jobs...")
        	wg.Wait()
        	log.Info("Done")
            return
        }
    }
}