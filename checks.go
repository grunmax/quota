package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"golang.org/x/net/websocket"
)

type Value struct {
	SessionID string `json:"sessionId"`
}

type SessionBody struct {
	Value Value `json:"value"`
}

func waitkey() {
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')
}

func checkWS(ip_ string, port_ string, url_ string, doRead bool) {
	origin := "http://localhost/"
	parts := strings.Split(url_, "/")
	parts = parts[3:]
	s := strings.Join(parts[:], "/")
	remote := fmt.Sprintf(`ws://%s:%s/%s`, ip_, port_, s)
	fmt.Println("connect to: " + remote)
	ws, err := websocket.Dial(remote, "", origin)
	if err != nil {
		fmt.Println("connect error")
		waitkey()
	} else {
		fmt.Println("connected")
	}
	if !doRead {
		return
	}
	err = ws.SetReadDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		fmt.Println("SetReadDeadline failed:", err)
		return
	}
	var msg = make([]byte, 1024)
	var n int
	n, err = ws.Read(msg)
	if err != nil {
		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Println("read timeout:", err)
			waitkey()
		} else {
			fmt.Println("read error:", err)
			waitkey()
		}
	}
	fmt.Println("Received: ", len(msg[:n]))

}

func checkSession(url_ string, body_ string) {
	var jsonStr = []byte(body_)
	req, err := http.NewRequest("POST", url_, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("connect error:", err)
		waitkey()
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	var sessionBody SessionBody
	if err := json.Unmarshal(data, &sessionBody); err != nil {
		panic(err)
	}
	if resp.Status == "200 OK" {
		fmt.Println(sessionBody.Value.SessionID)
	} else {
		fmt.Println(resp.Status)
		waitkey()
	}
}

func checkAndroidSession(ip_ string, port_ string, udid string, aver_ string) {
	url := fmt.Sprintf("http://%s:%s/wd/hub/session", ip_, port_)
	reqBody := fmt.Sprintf(`{"capabilities":{"alwaysMatch":{"appium:automationName":"uiautomator2","appium:deviceName":"Android","appium:labels":"%s","appium:noReset":false,"appium:noSign":true,"browserVersion":"%s","platformName":"Android"},"firstMatch":[{}]}}`, udid, aver_)
	checkSession(url, reqBody)
}

func checkIosSession(ip_ string, port_ string, udid string, iver_ string) {
	url := fmt.Sprintf("http://%s:%s/wd/hub/session", ip_, port_)
	reqBody := fmt.Sprintf(`{"capabilities":{"alwaysMatch":{"appium:automationName":"XCuiTest","appium:clearSystemFiles":true,"appium:connectHardwareKeyboard":true,"appium:derivedDataPath":"/jenkins/workspace/ios-regression-2424/build/wda/c1272d9e41d1406ae1dfd9a28ddcc932d0997d11/WebDriverAgent","appium:deviceName":"iPhone","appium:fullReset":false,"appium:newCommandTimeout":3600,"appium:noReset":true,"appium:orientation":"PORTRAIT","appium:platformVersion":"15.8.1","appium:processArguments":{"args":[],"env":{}},"appium:showXcodeLog":true,"appium:udid":"%s","appium:useJSONSource":true,"appium:usePrebuiltWDA":true,"browserVersion":"%s","platformName":"ios"},"firstMatch":[{}]}}`, udid, iver_)
	checkSession(url, reqBody)
}

func readXml(name string) {
	xmlFile, err := os.Open(name)
	if err != nil {
		panic(err)
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	var browsers Browsers
	err = xml.Unmarshal(byteValue, &browsers)
	if err != nil {
		panic(err)
	}

	for _, browser := range browsers.Browser {
		hosts := browser.Version.Region.Host
		for _, host := range hosts {
			fmt.Println("-> Check for " + host.Udid)
			if strings.Contains(host.Udid, "-") {
				fmt.Println("Session:")
				checkIosSession(checkIP, checkPort, host.Udid, iosVer)
				fmt.Println("Appium logs:")
				checkWS(checkIP, checkPort, host.WsAppiumLogs, true)
				fmt.Println("Device logs:")
				checkWS(checkIP, checkPort, host.WsDeviceLogs, true)
				fmt.Println("Host cmd: " + host.WsHostCmd)
				checkWS(checkIP, checkPort, host.WsHostCmd, false)

			} else {
				fmt.Println("Session:")
				checkAndroidSession(checkIP, checkPort, host.Udid, androidVer)
				fmt.Println("Appium logs:")
				checkWS(checkIP, checkPort, host.WsAppiumLogs, true)
				fmt.Println("Device logs:")
				checkWS(checkIP, checkPort, host.WsDeviceLogs, true)
				fmt.Println("Host cmd: " + host.WsHostCmd)
				checkWS(checkIP, checkPort, host.WsHostCmd, false)
			}
		}
	}
}

func check(name string) {
	if name != "" {
		if checkIP == "0.0.0.0" {
			fmt.Println("Check IP is not set")
			os.Exit(1)
		}
		if checkPort == "0" {
			fmt.Println("Check port is not set")
			os.Exit(1)
		}
		fmt.Println("Checking " + name)
		readXml(name)
		os.Exit(0)
	}
}
