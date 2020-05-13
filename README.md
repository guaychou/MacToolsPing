# Mac Tools Ping

The purpose of this tools is monitoring your ping latency to google.com <br /> 
This tools will appears in your system tray.

- Tested in Mac OS Mojave 10.14.6

### Screenshot

![Screenshot](images/img1.png)

### How to build from source:

Requirements:
- Golang 1.13 or above
- Go module on

```cassandraql
$ git clone https://github.com/guaychou/MacToolsPing.git
$ cd MacToolsPing
$ chmod +x build.sh
$ ./build.sh
```
Download release [here](https://github.com/guaychou/MacToolsPing/releases/download/v1.0/PingService.zip)

Credits to :

- [getlantern](https://github.com/getlantern/systray)
- [sparrc](github.com/sparrc/go-ping)
