package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var stubIos = []string{"00008101-000D15113C05001E", "00008030-001129DA3CB8802E", "00008020-00095C200ABB002E"}
var stubAndroid = []string{"36271FDHS003R7", "RF8T20L58HF", "4950395453313498"}

func setHostName() {
	h, err := os.Hostname()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	hostName = h
	fmt.Println(hostName)
}

func setLocalIP() {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	localIP = localAddr.IP.String()
	fmt.Println(localIP)
}

func initByOs() {
	sys := runtime.GOOS
	switch sys {
	case "darwin":
		// stub
		getDevicesCmd = `ioreg -p IOUSB -w0 | sed 's/[^o]*o //; s/@.*$//' | grep -v '^Root.*'`
	case "linux":
		getDevicesCmd = `usb-devices | awk -v RS='' -v ORS='\n\n' '/Vendor=05ac/ && /Product=iPhone/ || /Vendor=(04e8|18d1|1949|0fce)/' | grep -o -P '(?<=SerialNumber=).*' | awk -F, '{ if (length($1)>20) {print substr($1,1,8) "-" substr($1,9)} else {print}}'`
	default:
		fmt.Printf("%s.\n", sys)
		os.Exit(1)
	}
}

func getDevices() ([]string, []string) {
	if isStub {
		return stubIos, stubAndroid
	}

	cmd := exec.Command("bash", "-c", getDevicesCmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	devices := strings.Split(out.String(), "\n")
	ioss := []string{}
	androids := []string{}
	for _, device := range devices {
		if strings.TrimSpace(device) != "" {
			if strings.Contains(device, "-") {
				ioss = append(ioss, device)
			} else {
				androids = append(androids, device)
			}
		}
	}
	return ioss, androids
}

func makeXml(name string, v any) {
	xmlFile, err := os.Create(name)
	if err != nil {
		fmt.Println("Error creating XML file: ", err)
		os.Exit(1)
		return
	}
	encoder := xml.NewEncoder(xmlFile)
	encoder.Indent("", "\t")
	err = encoder.Encode(v)
	if err != nil {
		fmt.Println("Error encoding XML to file: ", err)
		os.Exit(1)
		return
	}
	fmt.Println("Quota XML", name, "created")
}
