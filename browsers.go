package main

import (
	"encoding/xml"
	"fmt"
)

type Host struct {
	Text         string `xml:",chardata"`
	Name         string `xml:"name,attr"`
	Port         string `xml:"port,attr"`
	Count        string `xml:"count,attr"`
	Udid         string `xml:"udid,attr"`
	Label        string `xml:"label,attr"`
	Ws           string `xml:"ws,attr"`
	WsAppiumLogs string `xml:"ws_appium_logs,attr"`
	WsDeviceLogs string `xml:"ws_device_logs,attr"`
	WsHostCmd    string `xml:"ws_host_cmd,attr"`
}

type Region struct {
	Text string `xml:",chardata"`
	Name string `xml:"name,attr"`
	Host []Host `xml:"host"`
}

type Version struct {
	Text     string `xml:",chardata"`
	Number   string `xml:"number,attr"`
	Platform string `xml:"platform,attr"`
	Region   Region `xml:"region"`
}

type Browser struct {
	Text           string  `xml:",chardata"`
	Name           string  `xml:"name,attr"`
	DefaultVersion string  `xml:"defaultVersion,attr"`
	Version        Version `xml:"version"`
}

type Browsers struct {
	XMLName xml.Name  `xml:"browsers"`
	Text    string    `xml:",chardata"`
	Browser []Browser `xml:"browser"`
}

func getIphoneBrowser(devices []string) Browser {
	hosts := []Host{}
	for _, device := range devices {
		host := Host{
			Name:         localIP,
			Port:         localPort,
			Count:        "1",
			Udid:         device,
			Label:        device,
			Ws:           fmt.Sprintf("ws://%s:%s/appium/device/syslog/%s", localIP, localPort, device),
			WsAppiumLogs: fmt.Sprintf("ws://%s:%s/ws/appium/logs/%s", localIP, localPort, device),
			WsDeviceLogs: fmt.Sprintf("ws://%s:%s/ws/device/logs/%s", localIP, localPort, device),
			WsHostCmd:    fmt.Sprintf("ws://%s:%s/ws/host/cmd/%s", localIP, localPort, device),
		}
		hosts = append(hosts, host)
	}

	region := Region{
		Name: hostName,
		Host: hosts,
	}

	iphoneVersion := Version{
		Number:   iosVer,
		Platform: "iOS",
		Region:   region,
	}
	iphoneBrowser := Browser{
		Name:           "iPhone",
		DefaultVersion: iosVer,
		Version:        iphoneVersion,
	}

	return iphoneBrowser
}

func getAndroidBrowser(devices []string) Browser {
	hosts := []Host{}
	for _, device := range devices {
		host := Host{
			Name:         localIP,
			Port:         localPort,
			Count:        "1",
			Udid:         device,
			Label:        device,
			Ws:           fmt.Sprintf("ws://%s:%s/appium/device/logcat/%s", localIP, localPort, device),
			WsAppiumLogs: fmt.Sprintf("ws://%s:%s/ws/appium/logs/%s", localIP, localPort, device),
			WsDeviceLogs: fmt.Sprintf("ws://%s:%s/ws/device/logs/%s", localIP, localPort, device),
			WsHostCmd:    fmt.Sprintf("ws://%s:%s/ws/host/cmd/%s", localIP, localPort, device),
		}
		hosts = append(hosts, host)
	}

	region := Region{
		Name: hostName,
		Host: hosts,
	}

	iphoneVersion := Version{
		Number:   androidVer,
		Platform: "Android",
		Region:   region,
	}
	iphoneBrowser := Browser{
		Name:           "Android",
		DefaultVersion: androidVer,
		Version:        iphoneVersion,
	}

	return iphoneBrowser
}
