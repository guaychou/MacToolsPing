package main

import (
	"errors"
	"github.com/getlantern/systray"
	"github.com/sparrc/go-ping"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const url_ip = "https://www.myexternalip.com/raw"
const configFileName="config.yml"
var pingAddress string

func init(){
	cfg,err:=readConfig()
	if err!=nil{
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

			case <-cloudflareAddress.ClickedCh:
				pingAddress="1.1.1.1"

			case <-mQuit.ClickedCh:
				systray.Quit()
				return
			}
		}
	}()
}

func onExit() {
	writeConfigState(pingAddress)
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