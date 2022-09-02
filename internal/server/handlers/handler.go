package handlers

import (
	"VPNservice/internal/vpndata"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"golang.org/x/text/encoding/charmap"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type Handler interface {
	Register(router *mux.Router)
	Ping(vpn *vpndata.VPN, wg *sync.WaitGroup)
	GetVPNlist() []vpndata.VPN
}

type handler struct {
}

func (h *handler) GetVPNlist() []vpndata.VPN {
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

	var list vpndata.VPN
	var vpnlist []vpndata.VPN

	for i := 2; i < len(filelines); i++ {
		vpnstr := ""
		ip := strings.Split(filelines[i], ",")
		if len(ip) > 5 {
			if ip[5] == "Russian Federation" {
				continue
			} else {
				vpnstr = ip[0] + "," + ip[1] + "," + ip[5] + "," + "\n"
				if _, err := f.WriteString(vpnstr); err != nil {
					log.Fatal(err)
				}
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

func (h *handler) Ping(vpn *vpndata.VPN, wg *sync.WaitGroup) {
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

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc("/", h.GetAll)
	router.HandleFunc("/servers", h.GetMSK)
	router.HandleFunc("/sh", h.GetSH)
}

func NewHandler() Handler {
	return &handler{}
}

func (h *handler) GetAll(w http.ResponseWriter, r *http.Request) {
	wg := sync.WaitGroup{}

	vpn := h.GetVPNlist()

	for i := 0; i < len(vpn); i++ {
		wg.Add(1)
		go vpndata.Ping(&vpn[i], &wg)
	}

	fmt.Println("Выполнение программы. Ожидайте.")
	wg.Wait()

	resultJson, err := json.MarshalIndent(vpn, " ", " ")
	if err != nil {
		errors.New(fmt.Sprintf("не удалось перекодировать данные. ошибка: %v", err))
	}
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resultJson)
}

func (h *handler) GetMSK(w http.ResponseWriter, r *http.Request) {

	//resultJson, err := json.MarshalIndent(collect.ProductMSK, " ", " ")
	//if err != nil {
	//	errors.New(fmt.Sprintf("не удалось перекодировать данные. ошибка: %v", err))
	//}
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//w.WriteHeader(http.StatusOK)
	//_, _ = w.Write(resultJson)
}

func (h *handler) GetSH(w http.ResponseWriter, r *http.Request) {

	//resultJson, err := json.MarshalIndent(collect.ProductSH, " ", " ")
	//if err != nil {
	//	errors.New(fmt.Sprintf("не удалось перекодировать данные. ошибка: %v", err))
	//}
	//w.Header().Set("Access-Control-Allow-Origin", "*")
	//w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	//w.WriteHeader(http.StatusOK)
	//_, _ = w.Write(resultJson)
}
