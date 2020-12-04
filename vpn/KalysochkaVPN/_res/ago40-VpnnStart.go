package main

import (
	"fmt"

	ps "./go-powershell"
	"./go-powershell/backend"
)

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

var printConnStatus = `[Console]::OutputEncoding = [System.Text.Encoding]::GetEncoding("cp866")
$vpnName = "HomeSSTP";
$vpn = Get-VpnConnection -Name $vpnName;
$vpn.ConnectionStatus
`