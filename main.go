package main

import (
	"errors"
	"github.com/getlantern/systray"
	"github.com/sparrc/go-ping"
	"strconv"
	"time"
)
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
	mQuit := systray.AddMenuItem("Quit", "Quits this app")
	go func() {
		for {
			select {
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
	pinger,err := ping.NewPinger("www.google.com")
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