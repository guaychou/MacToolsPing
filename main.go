package main

import (
	"github.com/getlantern/systray"
	"github.com/sparrc/go-ping"
	"time"
	"errors"
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
				time.Sleep(3 * time.Second)
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
	pinger.Count = 1
	pinger.Run()                 // blocks until finished
	stats := pinger.Statistics() // get send/receive/rtt stats
	if stats!=nil{
		return 	"Ping: "+stats.Rtts[0].String(),nil
	}
	return "", errors.New("Network Error")
}