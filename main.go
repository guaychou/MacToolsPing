package main

import (
	"errors"
	"github.com/getlantern/systray"
	"github.com/sparrc/go-ping"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"
	yaml "gopkg.in/yaml.v2"
)

const url_ip = "https://www.myexternalip.com/raw"
const configFileName="config.yml"
var pingAddress string

type Config struct {
	Server struct {
		Address string `yaml:"address"`
	}
}

func readConfig()(*Config,error){
	config := &Config{}
	//Dir, err := os.Getwd()
	Dir,err:=os.Executable()
	log.Print(Dir)
	if err != nil {
		log.Fatal(err)
	}
	configPath:=getConfigDirectory()+"/"+configFileName
	file,err:=os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}
	return config, nil
}

func getConfigDirectory()string{
	Dir,err:=os.Executable()
	log.Print(Dir)
	if err != nil {
		log.Fatal(err)
	}
	configPath := filepath.Dir(Dir)
	return configPath
}

func init(){
	cfg,err:=readConfig()
	if err!=nil{
		ioutil.WriteFile("/tmp/error.log", []byte(err.Error()),0644)
		log.Fatal(err)
	}
	pingAddress=cfg.Server.Address
}

func main(){
	systray.Run(onReady, onExit)
}

func onReady(){
	go func() {
		var result string
		var err error
		for {
			result,err=pingGoogle()
			if err!=nil{
				systray.SetTitle(err.Error())
			}else{
				systray.SetTitle(result)
				time.Sleep(2 * time.Second)
			}
		}
	}()
	systray.AddSeparator()
	pingTray:=systray.AddMenuItem("Ping Address","Your Ping Address")
	googleAddress:=pingTray.AddSubMenuItem("Google","Google Address")
	cloudflareAddress:=pingTray.AddSubMenuItem("Cloudflare","Cloudflare Address")
	publicIp:=systray.AddMenuItem("","Your Public IP")
	systray.AddSeparator()
	mQuit := systray.AddMenuItem("Quit", "Quits this app")

	go func() {
		for  {
			result,err:=getPublicIp()
			if err!=nil{
				publicIp.SetTitle(err.Error())

			}else{
				publicIp.SetTitle(result)
				time.Sleep(30 * time.Minute)
			}
		}
	}()
	go func() {
		for {
			select {
			case <-googleAddress.ClickedCh:
				pingAddress="www.google.com"
				writeConfigState("www.google.com")

			case <-cloudflareAddress.ClickedCh:
				pingAddress="1.1.1.1"
				writeConfigState("1.1.1.1")

			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	// Cleaning stuff here.
}

func pingGoogle() (string,error) {
	pinger,err := ping.NewPinger(pingAddress)
	log.Print(pingAddress)
	if err != nil {
		return "", errors.New("Network Error")
	}
	runPinger(pinger)
	result,err:=getPingLatency(pinger)
	if err!=nil{
		return "", errors.New("Network Error")
	}
	return 	strconv.Itoa(result) +" ms",nil
}

func runPinger(pinger *ping.Pinger) {
	pinger.Count = 1
	pinger.Timeout=2 * time.Second
	pinger.Debug=true
	pinger.Run()
}

func getPingLatency(pinger *ping.Pinger)(int,error){
	if pinger.Statistics().Rtts==nil{
		return -1,errors.New("Result zero")
	}
	return int(pinger.Statistics().Rtts[0].Milliseconds()),nil
}

func getPublicIp()(string,error){
	resp,err:=http.Get(url_ip)
	if err!=nil{
		return "", errors.New("Network Error")
	}
	defer resp.Body.Close()
	if resp.StatusCode==http.StatusOK{
		body, err := ioutil.ReadAll(resp.Body)
		if err!=nil{
			log.Fatal(err)
		}
		return "IP: "+string(body),nil
	}
	return "",errors.New("Network Error")
}

func writeConfigState(address string){
	var config Config
	config.Server.Address=address
	d, err := yaml.Marshal(&config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	err = ioutil.WriteFile(configFileName, d, 0644)
	if err != nil {
		log.Fatal(err)
	}
}