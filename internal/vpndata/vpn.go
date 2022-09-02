package vpndata

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func GetVPN() {

}

func CreateVPNConnection() {
	//site := "https://www.vpngate.net/en/"
	vpnName := "Add-VpnConnection -Name \"123\" -ServerAddress \"219.100.37.60\" -TunnelType L2TP -L2tpPsk \"vpn\" -Force -AuthenticationMethod MSChapv2 -RememberCredential"

	args := strings.Split(vpnName, " ")
	cmd := exec.Command("powershell.exe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	err := cmd.Run()

	if err != nil {
		fmt.Printf("cmd.Run: %s failed: %s\n", err, err)
	}

	startVPNConnection()
}

func startVPNConnection() {
	runVpn := "rasdial \"123\" \"vpn\" \"vpn\""
	args := strings.Split(runVpn, " ")
	cmd := exec.Command("powershell.exe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	err := cmd.Run()

	if err != nil {
		fmt.Printf("cmd.Run: %s failed: %s\n", err, err)
	}
}

func DisconnectVPNConnection() {
	removeVPN := "rasdial \"123\" /DISCONNECT"
	args := strings.Split(removeVPN, " ")
	cmd := exec.Command("powershell.exe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	err := cmd.Run()

	if err != nil {
		fmt.Printf("cmd.Run: %s failed: %s\n", err, err)
	}
}

func RemoveVPNConnection() {

	removeVPN := "Remove-VpnConnection -Name \"123\""
	args := strings.Split(removeVPN, " ")

	cmd := exec.Command("powershell.exe", args...)
	pipe, err := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	pipe.Write([]byte("Y"))
	pipe.Close()
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	cmd.Run()

	err = cmd.Run()

	if err != nil {
		fmt.Printf("cmd.Run: %s failed: %s\n", err, err)
	}
}

func Ping() {

	args := strings.Split("/c ping google.ru", " ")
	cmd := exec.Command("powershell.exe", args...)

	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = os.Stderr
	cmd.Run()

	err := cmd.Run()

	if err != nil {
		fmt.Printf("cmd.Run: %s failed: %s\n", err, err)
	}

	dec := charmap.CodePage866.NewDecoder()
	reader, err := dec.Bytes(b.Bytes())
	str := string(reader)
	st := strings.Split(str, "\n")
	if len(st) < 5 {
		fmt.Println("недоступен")
	} else {
		st1 := strings.Split(st[11], ",")
		st2 := strings.Split(st1[2], " ")
		var timePing int
		for i := 0; i < len(st2); i++ {
			num, err := strconv.Atoi(st2[i])
			if err != nil {
				continue
			}
			timePing = num
		}
		fmt.Println("Средний пинг:", timePing)
	}
}
