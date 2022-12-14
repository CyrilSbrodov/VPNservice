package main

import (
	"VPNservice/internal/vpndata"
	"fmt"
	"sync"
)

func main() {
	wg := sync.WaitGroup{}

	vpndata.GetVPN()

	vpn := vpndata.GetVPNlist()

	for i := 0; i < len(vpn); i++ {
		wg.Add(1)
		go vpndata.Ping(&vpn[i], &wg)
	}
	fmt.Println("Выполнение программы. Ожидайте.")
	wg.Wait()
	name := vpndata.CreateVPNConnection(vpn)
	vpndata.DisconnectVPNConnection(name)
	//serv := server.NewApp()
	//serv.Run()
}
