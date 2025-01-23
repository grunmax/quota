package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	getDevicesCmd string
	devices       []string
	iphones       []string
	androids      []string
	hostName      string
	localIP       string
	localPort     string
	checkIP       string
	checkPort     string
	iosVer        string
	androidVer    string
	isStub        bool
	fMake         string
	checkFile     string
	showDevices   bool
)

func init() {
	fmt.Println("")
	fmt.Println("Quota version 0.5")
	fmt.Println("")
	flag.StringVar(&hostName, "host", "1", "Host name of local GGR")
	flag.StringVar(&localPort, "port", "4444", "Port of local GGR")
	flag.StringVar(&localIP, "ip", "0.0.0.0", "IP of local GGR")
	flag.StringVar(&checkPort, "xport", "0", "Port of check GGR")
	flag.StringVar(&checkIP, "xip", "0.0.0.0", "IP of check GGR")
	flag.StringVar(&iosVer, "iver", "14.4", "IOS version for devices")
	flag.StringVar(&androidVer, "aver", "11", "IOS version for devices")
	flag.BoolVar(&isStub, "stub", false, "Use stub udids data")
	flag.StringVar(&fMake, "make", "", "Make xml quota with name")
	flag.BoolVar(&showDevices, "devices", false, "Show devices")
	flag.StringVar(&checkFile, "check", "", "Check quota xml from ip")
	flag.Parse()

	initByOs()
	if hostName == "1" {
		setHostName()
	}
	if localIP == "0.0.0.0" {
		setLocalIP()
	}
}

func main() {
	iosDevices, androidDevices := getDevices()

	if showDevices {
		fmt.Println("Ios devices:")
		fmt.Println(iosDevices)
		fmt.Println("Android devices:")
		fmt.Println(androidDevices)
		fmt.Println("")
		os.Exit(0)
	}

	check(checkFile)

	if fMake == "" {
		os.Exit(0)
	}

	browsers := []Browser{}
	iphoneBrowser := getIphoneBrowser(iosDevices)
	androidBrowser := getAndroidBrowser(androidDevices)
	browsers = append(browsers, iphoneBrowser)
	browsers = append(browsers, androidBrowser)

	root := &Browsers{
		Browser: browsers,
	}

	makeXml(fMake, root)
}
