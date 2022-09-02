package vpndata

import (
	"bytes"
	"fmt"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type VPN struct {
	Name    string
	IP      string
	Country string
	Ping    int
}

func GetVPN() {

	filename := "internal/vpndata/vpn.csv"

	req, err := http.Get("http://www.vpngate.net/api/iphone/")
	if err != nil {
		log.Fatal(err)
	}
	defer req.Body.Close()
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatal(err)
	}
	file, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	ioutil.WriteFile(filename, b, 0644)

}

func CreateVPNConnection(vpn []VPN) {
	sort.SliceStable(vpn, func(i, j int) bool {
		return vpn[i].Ping < vpn[j].Ping
	})
	fmt.Println(vpn)

	vpnName := fmt.Sprintf("Add-VpnConnection -Name \"%v\" -ServerAddress \"%v\" -TunnelType Automatic -L2tpPsk \"vpn\" -Force -AuthenticationMethod MSChapv2 -RememberCredential", vpn[0].Name, vpn[0].IP)
	fmt.Println(vpnName)
	args := strings.Split(vpnName, " ")
	cmd := exec.Command("powershell.exe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("cmd.Run: %s failed: %s\n", err, err)
	}

	startVPNConnection(&vpn[0])
}

func startVPNConnection(vpn *VPN) {
	runVpn := fmt.Sprintf("rasdial \"%v\" \"vpn\" \"vpn\"", vpn.Name)
	args := strings.Split(runVpn, " ")
	cmd := exec.Command("powershell.exe", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

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

	err = cmd.Run()

	if err != nil {
		fmt.Printf("cmd.Run: %s failed: %s\n", err, err)
	}
}

func Ping(vpn *VPN, wg *sync.WaitGroup) {

	ip := "/c" + " ping" + " " + vpn.IP
	args := strings.Split(ip, " ")
	cmd := exec.Command("powershell.exe", args...)

	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = os.Stderr

	err := cmd.Run()

	if err != nil {
		fmt.Printf("cmd.Run: %s failed: %s\n", err, err)
	}

	dec := charmap.CodePage866.NewDecoder()
	reader, err := dec.Bytes(b.Bytes())
	str := string(reader)
	st := strings.Split(str, "\n")
	dnd := strings.Contains(str, "Превышен")
	if len(st) < 4 {
		vpn.Ping = 9999
	} else if dnd {
		vpn.Ping = 9999
	} else {
		if len(st) < 12 {
			wg.Done()
			return
		}
		st1 := strings.Split(st[11], ",")
		st2 := strings.Split(st1[2], " ")
		for i := 0; i < len(st2); i++ {
			num, err := strconv.Atoi(st2[i])
			if err != nil {
				continue
			}
			vpn.Ping = num
			wg.Done()
			return
		}
	}
	wg.Done()
}

func GetVPNlist() []VPN {
	filename := "internal/vpndata/vpnlist.csv"
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	file, err := ioutil.ReadFile("internal/vpndata/vpn.csv")
	if err != nil {
		log.Fatal(err)
	}
	filelines := strings.Split(string(file), "\n")

	var list VPN
	var vpnlist []VPN

	for i := 2; i < len(filelines); i++ {
		vpnstr := ""
		ip := strings.Split(filelines[i], ",")
		if len(ip) > 5 {
			vpnstr = ip[0] + "," + ip[1] + "," + ip[5] + "," + "\n"
			if _, err := f.WriteString(vpnstr); err != nil {
				log.Fatal(err)
			}
		}
	}
	for i := 2; i < len(filelines); i++ {
		ip := strings.Split(filelines[i], ",")
		if len(ip) > 5 {
			list.Name = ip[0]
			list.IP = ip[1]
			list.Country = ip[5]
			vpnlist = append(vpnlist, list)
		}
	}
	return vpnlist
}
