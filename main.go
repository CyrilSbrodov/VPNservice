package main

import "VPNservice/internal/vpndata"

func main() {
	//vpndata.Ping()
	vpndata.CreateVPNConnection()
	vpndata.DisconnectVPNConnection()
	vpndata.RemoveVPNConnection()
}
