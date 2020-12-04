package main

import (
	ps "./go-powershell"
	"./go-powershell/backend"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/xela07ax/toolsXela/tp"
	"html/template"
	"log"
	"net"
	"os"
	"time"
)

type Config struct {
	ListActiveInterfaces bool
	VpnName string
	User string
	Password string
	TimeSleepSecond time.Duration
}

const configDir = "complect_monitoringStart.json"

func main() {
	log.Println("****#####   Kolyaska v.1 - VPN Checker  #####****")
	// Открываем конфиг
	fi,err := tp.OpenReadFile(configDir)
	if err != nil {
		fmt.Printf("Ошибка при открытии конфигурации %s: %s\n",configDir,err)
		tp.ExitWithSecTimeout(1)
	}
	var cfg Config
	err = json.Unmarshal(fi,&cfg)
	if err != nil {
		fmt.Printf("Ошибка чтения JSON %s: %s\n",configDir,err)
		tp.ExitWithSecTimeout(1)
	}
	// Networks
	if cfg.ListActiveInterfaces {
		availableInterfaces()
		os.Exit(0)
	}
	//...
	for {
		// Бем в цикле кадый раз проверять
		_, err = net.InterfaceByName(cfg.VpnName)
		if err != nil {
			// Сюда попали потому, что подключение не существует
			// Подключаемся же
			connect(Templater{
				User:     cfg.User,
				Password: cfg.Password,
			})
			//fmt.Println(err, "["+cfg.VpnName+"]")
			//fmt.Println("-----------------------------")
			//availableInterfaces()
			//os.Exit(0)
		}
		time.Sleep(cfg.TimeSleepSecond*time.Second)
	}
	//fmt.Println(byNameInterface)
	log.Println("****#####   Good by!  #####****")
}


func connect(auth Templater)  {
	log.Println("****#####   Connect starting   #####****")
	// choose a backend
	back := &backend.Local{}

	// start a local powershell process
	shell, err := ps.New(back)
	if err != nil {
		panic(err)
	}
	defer shell.Exit()


	// my templater Go
	buf := new(bytes.Buffer)
	err = connVpnTemplate.Execute(buf, auth)
	if err != nil {
		panic(err)
	}
	//fmt.Println(buf.String())
	// ... and interact with it
	stdout, stderr, err := shell.Execute(buf.String())
	if err != nil {
		panic(err)
	}
	fmt.Println("stdout:")
	fmt.Println(stdout)
	fmt.Println("stderr:")
	fmt.Println(stderr)
	log.Println("****#####   Connect succesfully   #####****")

}

func availableInterfaces() {
	interfaces, err := net.Interfaces()

	if err != nil {
		fmt.Print(err)
		os.Exit(0)
	}

	fmt.Println("Available network interfaces on this machine : ")
	for _, i := range interfaces {
		fmt.Printf("Name : %v \n", i.Name)
	}
}
type Templater struct {
	User string
	Password string
}
var connVpnTemplate = template.Must(template.New("").Parse(`[Console]::OutputEncoding = [System.Text.Encoding]::GetEncoding("cp866")
$vpnName = "HomeSSTP";

$vpnusername = "{{.User}}"
$vpnpassword = "{{.Password}}"
$vpn = Get-VpnConnection | where {$_.Name -eq $vpnname}
$cmd = $env:WINDIR + "\System32\rasdial.exe"
$expression = "$cmd ""$vpnname"" $vpnusername $vpnpassword"
Invoke-Expression -Command $expression
`))
