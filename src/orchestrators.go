package main

import (
)

type Orchestration struct {
	Orchestration		string `json:"orchestration"`
	Total 		int `json:"total"`
	Active 		int `json:"active"`
	Hashost 	int	`json:"hashost"`
	Usage 		int	`json:"usage"`
	Hosts		int `json:"hosts"`
	Containers	int `json:"containers"`
}

func (k *Orchestration) setTotal(t int){
	k.Total = t
}

func (k *Orchestration) addTotal(){
	k.Total ++
}

func (k *Orchestration) addHosts(i int){
	k.Hosts += i
}

func (k *Orchestration) addContainers(i int){
	k.Containers += i
}

func (k *Orchestration) getOrchestration() string{
	return k.Orchestration
}

func (k *Orchestration) getTotal() int{
	return k.Total
}

func (k *Orchestration) getActive() int{
	return k.Active
}

func (k *Orchestration) getHashost() int{
	return k.Hashost
}

func (k *Orchestration) getHosts() int{
	return k.Hosts
}

func (k *Orchestration) getContainers() int{
	return k.Containers
}

func (k *Orchestration) getUsage() int{
	return k.Usage
}

func (k *Orchestration) addActive(){
	k.Active ++
}

func (k *Orchestration) addHashost(){
	k.Hashost ++
}

func (k *Orchestration) setUsage(t int){
	k.Usage = k.Total*100/t
}
