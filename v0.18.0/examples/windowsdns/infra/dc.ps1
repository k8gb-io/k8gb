Install-WindowsFeature AD-Domain-Services -IncludeManagementTools
Import-Module ADDSDeployment
$adminPwd=ConvertTo-SecureString "P@ssw0rd" -AsPlainText -Force
Install-ADDSForest -DomainName k8gb.local -InstallDNS -SafeModeAdministratorPassword $adminPwd -Confirm:$false
Set-DnsServerForwarder -IPAddress "168.63.129.16" -PassThru