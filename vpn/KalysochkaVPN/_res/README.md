
## golang check active ethernet connection
https://www.socketloop.com/tutorials/golang-check-whether-a-network-interface-is-up-on-your-machine

HomeSSTP VPN `включено`  
- Run:
```go
func main() {
	ints, err := net.Interfaces()

	for _, interf := range ints {
		fmt.Println(interf)
    }
}
```
Out:
```sh
GOROOT=C:\Go #gosetup
GOPATH=C:\Users\Nice\go #gosetup
C:\Go\bin\go.exe build -o C:\Users\Nice\AppData\Local\Temp\___go_build_main_go.exe D:/Projects/Github/XelaGoDoc/websocket-tutor/xela-ws-tutor-part_2/viewEthernet/main.go #gosetup
C:\Users\Nice\AppData\Local\Temp\___go_build_main_go.exe #gosetup
!!!!!Usage : C:\Users\Nice\AppData\Local\Temp\___go_build_main_go.exe <interface
 name>
{19 1500 VirtualBox Host-Only Network 0a:00:27:00:00:13 up|broadcast|multicast}
{6 1500 Ethernet 2 0a:00:27:00:00:06 up|broadcast|multicast}
{17 1500 Подключение по локальной сети 00:ff:b8:40:32:04 0}
{9 1500 Беспроводная сеть 64:5d:86:7a:93:83 broadcast|multicast}
{4 1500 Подключение по локальной сети* 10 66:5d:86:7a:93:83 broadcast|multicast}

{7 1500 Ethernet 00:d8:61:0e:11:4a up|broadcast|multicast}
{54 1350 HomeSSTP  up|pointtopoint|multicast}
{10 1500 Сетевое подключение Bluetooth 64:5d:86:7a:93:87 broadcast|multicast}
{1 -1 Loopback Pseudo-Interface 1  up|loopback|multicast}
```

HomeSSTP VPN `выключено`   
Out:
```sh
{19 1500 VirtualBox Host-Only Network 0a:00:27:00:00:13 up|broadcast|multicast}
{6 1500 Ethernet 2 0a:00:27:00:00:06 up|broadcast|multicast}
{17 1500 Подключение по локальной сети 00:ff:b8:40:32:04 0}
{9 1500 Беспроводная сеть 64:5d:86:7a:93:83 broadcast|multicast}
{4 1500 Подключение по локальной сети* 10 66:5d:86:7a:93:83 broadcast|multicast}
{7 1500 Ethernet 00:d8:61:0e:11:4a up|broadcast|multicast}
{10 1500 Сетевое подключение Bluetooth 64:5d:86:7a:93:87 broadcast|multicast}
{1 -1 Loopback Pseudo-Interface 1  up|loopback|multicast}

```


HomeSSTP VPN `выключено`  
- Run:
```go
func main() {
		availableInterfaces()
}
```
```
Available network interfaces on this machine :
Name : VirtualBox Host-Only Network
Name : Ethernet 2
Name : Подключение по локальной сети
Name : Беспроводная сеть
Name : Подключение по локальной сети* 10
Name : Ethernet
Name : Сетевое подключение Bluetooth
Name : Loopback Pseudo-Interface 1
```
HomeSSTP VPN `включено`  

Available network interfaces on this machine :  
Name : VirtualBox Host-Only Network  
Name : Ethernet 2  
Name : Подключение по локальной сети  
Name : Беспроводная сеть  
Name : Подключение по локальной сети* 10  
Name : Ethernet  
Name : ___HomeSSTP___  
Name : Сетевое подключение Bluetooth  
Name : Loopback Pseudo-Interface 1  


- Run with params ___HomeSSTP___  :
```go
func main() {
	if len(os.Args) != 2 {
		fmt.Printf("Usage : %s <interface name>\n", os.Args[0])
		os.Exit(0)
	}

	ifName := os.Args[1]

	byNameInterface, err := net.InterfaceByName(ifName)

	if err != nil {
		fmt.Println(err, "["+ifName+"]")
		fmt.Println("-----------------------------")
		availableInterfaces()
		os.Exit(0)
	}

	if strings.Contains(byNameInterface.Flags.String(), "up") {
		fmt.Println("Status : UP")
	} else {
		fmt.Println("Status : DOWN")
	}
}
```
Out:
```

HomeSSTP VPN `выключено`  
Run:
```go
func main() {
		availableInterfaces()
}
```
HomeSSTP VPN `включено`  
```
Status : UP
```
HomeSSTP VPN `выключено` 

```
route ip+net: no such network interface [HomeSSTP]
-----------------------------
Available network interfaces on this machine :
Name : VirtualBox Host-Only Network
Name : Ethernet 2
Name : Подключение по локальной сети
Name : Беспроводная сеть
Name : Подключение по локальной сети* 10
Name : Ethernet
Name : Сетевое подключение Bluetooth
Name : Loopback Pseudo-Interface 1
```

## golang run powershell command
https://github.com/40a/go-powershell  
<img src="./_res/Screenshot 2020-12-04 174542.jpg" width="250" />

Run a file ``ago40-Staart.go``:
```go
func main() {
	// choose a backend
	back := &backend.Local{}

	// start a local powershell process
	shell, err := ps.New(back)
	if err != nil {
		panic(err)
	}
	defer shell.Exit()

	// ... and interact with it
	stdout, stderr, err := shell.Execute("Get-WmiObject -Class Win32_Processor")
	if err != nil {
		panic(err)
	}

	fmt.Println(stdout)
	fmt.Println(stderr)
}
```
Out:

```

Caption           : Intel64 Family 6 Model 158 Stepping 12
DeviceID          : CPU0
Manufacturer      : GenuineIntel
MaxClockSpeed     : 3600
Name              : Intel(R) Core(TM) i9-9900K CPU @ 3.60GHz
SocketDesignation : U3E1
```

Run HomeSSTP VPN `выключено`   
`ago40-VpnnStart.go`
```go
...
	// ... and interact with it
	stdout, stderr, err := shell.Execute(printConnStatus)
	if err != nil {
		panic(err)
	}

	fmt.Println(stdout)
	fmt.Println(stderr)
}

var printConnStatus = `[Console]::OutputEncoding = [System.Text.Encoding]::GetEncoding("cp866")
$vpnName = "HomeSSTP";
$vpn = Get-VpnConnection -Name $vpnName;
$vpn.ConnectionStatus
`
```
out:
```
Disconnect
```
or `Включено` out:
```
Connected
```
- Подключаемся к ВПН  
`ago40-VpnnStart.go`
```go
...
	// ... and interact with it
	stdout, stderr, err := shell.Execute(connVpn)
	if err != nil {
		panic(err)
	}

	fmt.Println(stdout)
	fmt.Println(stderr)
}

var connVpn = `[Console]::OutputEncoding = [System.Text.Encoding]::GetEncoding("cp866")
$vpnName = "HomeSSTP";

$vpnusername = "xela"
$vpnpassword = "jBUTUEg5zwdXag3"
$vpn = Get-VpnConnection | where {$_.Name -eq $vpnname}
$cmd = $env:WINDIR + "\System32\rasdial.exe"
$expression = "$cmd ""$vpnname"" $vpnusername $vpnpassword"
Invoke-Expression -Command $expression
`
```
out:
```
Disconnect
```
## Config in `complect_monitoringStart.json`
```json
{
  "ListActiveInterfaces": true
}
```
Out  
```
Name : VirtualBox Host-Only Network
Name : Ethernet 2
Name : Подключение по локальной сети
Name : Беспроводная сеть
Name : Подключение по локальной сети* 10
Name : Ethernet
Name : HomeSSTP
Name : Сетевое подключение Bluetooth
Name : Loopback Pseudo-Interface 1
```
