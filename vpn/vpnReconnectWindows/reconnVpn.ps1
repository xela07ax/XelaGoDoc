[Console]::OutputEncoding = [System.Text.Encoding]::GetEncoding("cp866")
$vpnName = "HomeSSTP";
$vpn = Get-VpnConnection -Name $vpnName;
if($vpn.ConnectionStatus -eq "Disconnected"){


$vpnusername = "username"
$vpnpassword = "p@$$w0Rd"
$vpn = Get-VpnConnection | where {$_.Name -eq $vpnname}
if ($vpn.ConnectionStatus -eq "Disconnected")
{
$cmd = $env:WINDIR + "\System32\rasdial.exe"
$expression = "$cmd ""$vpnname"" $vpnusername $vpnpassword"
Invoke-Expression -Command $expression
}

}