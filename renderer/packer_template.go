package renderer

import (
	"fmt"
	"io/ioutil"
)

// PackerTemplate contains the packer.json and Autounattend.xml text/templates
type PackerTemplate struct {
	PackerTpl       string
	AutounattendTpl string
	VagrantfileTpl  string
}

// NewPackerTemplateWithOverrides creates a fully functioning template using
// any provided template file overrides if they exist.
func NewPackerTemplateWithOverrides(packerTplPath string, autounattendTplPath string, vagrantfileTplPath string) *PackerTemplate {
	tpl := NewPackerTemplate()
	packerTpl, err := ioutil.ReadFile(packerTplPath)
	if err == nil {
		tpl.PackerTpl = string(packerTpl)
	} else {
		fmt.Println(fmt.Sprintf("WARN: Couldn't open '%s', defaulting to internal inductor template", packerTplPath))
	}
	autounattendTpl, err := ioutil.ReadFile(autounattendTplPath)
	if err == nil {
		tpl.AutounattendTpl = string(autounattendTpl)
	} else {
		fmt.Println(fmt.Sprintf("WARN: Couldn't open '%s', defaulting to internal inductor template", autounattendTplPath))
	}
	vagrantfileTpl, err := ioutil.ReadFile(vagrantfileTplPath)
	if err == nil {
		tpl.VagrantfileTpl = string(vagrantfileTpl)
	} else {
		fmt.Println(fmt.Sprintf("WARN: Couldn't open '%s', defaulting to internal inductor template", vagrantfileTplPath))
	}
	return tpl
}

// NewPackerTemplate creates a new fully functioning PackerTemplate with
// hardcoded default templates.
func NewPackerTemplate() *PackerTemplate {
	pt := &PackerTemplate{}
	pt.PackerTpl = `{
  "builders": [
    {
      "type": "vmware-iso",
      "iso_url": "{{.IsoURL}}",
      "iso_checksum_type": "{{.IsoChecksumType}}",
      "iso_checksum": "{{.IsoChecksum}}",
      "headless": {{.Headless}},
      "boot_wait": "2m",
      "ssh_username": "{{.Username}}",
      "ssh_password": "{{.Password}}",
      "ssh_wait_timeout": "2h",
      "shutdown_command": "shutdown /s /t 10 /f /d p:4:1 /c \"Packer Shutdown\"",
      "guest_os_type": "{{.VmwareGuestOsType}}",
      "tools_upload_flavor": "windows",
      "disk_size": {{.DiskSize}},
      "vnc_port_min": 5900,
      "vnc_port_max": 5980,
      "floppy_files": [
        "Autounattend.xml",
        "./scripts/fixnetwork.ps1",
        "./scripts/microsoft-updates.bat",
        "./scripts/win-updates.ps1",
        "./scripts/openssh.ps1"
      ],
      "vmx_data": {
        "RemoteDisplay.vnc.enabled": "false",
        "RemoteDisplay.vnc.port": "5900",
        "memsize": "{{.RAM}}",
        "numvcpus": "{{.CPU}}",
        "scsi0.virtualDev": "lsisas1068"
      }
    },
    {
      "type": "virtualbox-iso",
      "iso_url": "{{.IsoURL}}",
      "iso_checksum_type": "{{.IsoChecksumType}}",
      "iso_checksum": "{{.IsoChecksum}}",
      "headless": {{.Headless}},
      "boot_wait": "2m",
      "ssh_username": "{{.Username}}",
      "ssh_password": "{{.Password}}",
      "ssh_wait_timeout": "2h",
      "shutdown_command": "shutdown /s /t 10 /f /d p:4:1 /c \"Packer Shutdown\"",
      "guest_os_type": "{{.VirtualboxGuestOsType}}",
      "disk_size": {{.DiskSize}},
      "floppy_files": [
        "Autounattend.xml",
        "./scripts/fixnetwork.ps1",
        "./scripts/microsoft-updates.bat",
        "./scripts/win-updates.ps1",
        "./scripts/openssh.ps1",
        "./scripts/oracle-cert.cer"
      ],
      "vboxmanage": [
        [
          "modifyvm",
          "{{"{{"}}.Name{{"}}"}}",
          "--memory",
          "{{.RAM}}"
        ],
        [
          "modifyvm",
          "{{"{{"}}.Name{{"}}"}}",
          "--cpus",
          "{{.CPU}}"
        ]
      ]
    }
  ],
  "provisioners": [
    {
      "type": "shell",
      "remote_path": "/tmp/script.bat",
      "execute_command": "{{"{{"}}.Vars{{"}}"}} cmd /c C:/Windows/Temp/script.bat",
      "scripts": [
        "./scripts/vm-guest-tools.bat",
        "./scripts/vagrant-ssh.bat",
        "./scripts/disable-auto-logon.bat",
        "./scripts/enable-rdp.bat",
        "./scripts/compile-dotnet-assemblies.bat",
        "./scripts/compact.bat"
      ]
    }
  ],
  "post-processors": [
    {
			"type": "vagrant",
      "keep_input_artifact": false,
      "output": "{{.OSName}}_{{"{{"}}.Provider{{"}}"}}.box",
      "vagrantfile_template": "Vagrantfile"
    }
  ]
}
`
	pt.AutounattendTpl = `<?xml version="1.0" encoding="utf-8"?>
<unattend xmlns="urn:schemas-microsoft-com:unattend">
    <servicing/>
    <settings pass="windowsPE">
        <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <DiskConfiguration>
                <Disk wcm:action="add">
                    <CreatePartitions>
                        <CreatePartition wcm:action="add">
                            <Order>1</Order>
                            <Type>Primary</Type>
                            <Extend>true</Extend>
                        </CreatePartition>
                    </CreatePartitions>
                    <ModifyPartitions>
                        <ModifyPartition wcm:action="add">
                            <Extend>false</Extend>
                            <Format>NTFS</Format>
                            <Letter>C</Letter>
                            <Order>1</Order>
                            <PartitionID>1</PartitionID>
                            <Label>{{.OSName}}</Label>
                        </ModifyPartition>
                    </ModifyPartitions>
                    <DiskID>0</DiskID>
                    <WillWipeDisk>true</WillWipeDisk>
                </Disk>
                <WillShowUI>OnError</WillShowUI>
            </DiskConfiguration>
            <UserData>
                <AcceptEula>true</AcceptEula>
                <FullName>{{.Username}}</FullName>
                <Organization></Organization>
                <ProductKey>
										{{ if .ProductKey }}
										<Key>{{.ProductKey}}</Key>
										{{ end }}
                    <WillShowUI>Never</WillShowUI>
                </ProductKey>
            </UserData>
            <ImageInstall>
                <OSImage>
                    <InstallTo>
                        <DiskID>0</DiskID>
                        <PartitionID>1</PartitionID>
                    </InstallTo>
                    <WillShowUI>OnError</WillShowUI>
                    <InstallToAvailablePartition>false</InstallToAvailablePartition>
                    <InstallFrom>
                        <MetaData wcm:action="add">
                            <Key>/IMAGE/NAME</Key>
                            <Value>{{.WindowsImageName}}</Value>
                        </MetaData>
                    </InstallFrom>
                </OSImage>
            </ImageInstall>
        </component>
        <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-International-Core-WinPE" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <SetupUILanguage>
                <UILanguage>en-US</UILanguage>
            </SetupUILanguage>
            <InputLocale>en-US</InputLocale>
            <SystemLocale>en-US</SystemLocale>
            <UILanguage>en-US</UILanguage>
            <UILanguageFallback>en-US</UILanguageFallback>
            <UserLocale>en-US</UserLocale>
        </component>
    </settings>
    <settings pass="offlineServicing">
        <component name="Microsoft-Windows-LUA-Settings" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <EnableLUA>false</EnableLUA>
        </component>
    </settings>
    <settings pass="oobeSystem">
        <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <UserAccounts>
                <AdministratorPassword>
                    <Value>{{.Password}}</Value>
                    <PlainText>true</PlainText>
                </AdministratorPassword>
                <LocalAccounts>
                    <LocalAccount wcm:action="add">
                        <Password>
                            <Value>{{.Password}}</Value>
                            <PlainText>true</PlainText>
                        </Password>
                        <Description>Vagrant User</Description>
                        <DisplayName>vagrant</DisplayName>
                        <Group>administrators</Group>
                        <Name>{{.Username}}</Name>
                    </LocalAccount>
                </LocalAccounts>
            </UserAccounts>
            <OOBE>
                <HideEULAPage>true</HideEULAPage>
                <HideWirelessSetupInOOBE>true</HideWirelessSetupInOOBE>
                <NetworkLocation>Home</NetworkLocation>
                <ProtectYourPC>1</ProtectYourPC>
            </OOBE>
            <AutoLogon>
                <Password>
                    <Value>{{.Password}}</Value>
                    <PlainText>true</PlainText>
                </Password>
                <Username>{{.Username}}</Username>
                <Enabled>true</Enabled>
            </AutoLogon>
            <FirstLogonCommands>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c powershell -Command "Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Force"</CommandLine>
                    <Description>Set Execution Policy 64 Bit</Description>
                    <Order>1</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>C:\Windows\SysWOW64\cmd.exe /c powershell -Command "Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Force"</CommandLine>
                    <Description>Set Execution Policy 32 Bit</Description>
                    <Order>2</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c reg add "HKLM\System\CurrentControlSet\Control\Network\NewNetworkWindowOff"</CommandLine>
                    <Description>Network prompt</Description>
                    <Order>3</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe -File a:\fixnetwork.ps1</CommandLine>
                    <Description>Fix public network</Description>
                    <Order>4</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c winrm quickconfig -q</CommandLine>
                    <Description>winrm quickconfig -q</Description>
                    <Order>5</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c winrm quickconfig -transport:http</CommandLine>
                    <Description>winrm quickconfig -transport:http</Description>
                    <Order>6</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c winrm set winrm/config @{MaxTimeoutms="1800000"}</CommandLine>
                    <Description>Win RM MaxTimoutms</Description>
                    <Order>7</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c winrm set winrm/config/winrs @{MaxMemoryPerShellMB="800"}</CommandLine>
                    <Description>Win RM MaxMemoryPerShellMB</Description>
                    <Order>8</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c winrm set winrm/config/service @{AllowUnencrypted="true"}</CommandLine>
                    <Description>Win RM AllowUnencrypted</Description>
                    <Order>9</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c winrm set winrm/config/service/auth @{Basic="true"}</CommandLine>
                    <Description>Win RM auth Basic</Description>
                    <Order>10</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c winrm set winrm/config/client/auth @{Basic="true"}</CommandLine>
                    <Description>Win RM client auth Basic</Description>
                    <Order>11</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c winrm set winrm/config/listener?Address=*+Transport=HTTP @{Port="5985"} </CommandLine>
                    <Description>Win RM listener Address/Port</Description>
                    <Order>12</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c netsh advfirewall firewall set rule group="remote administration" new enable=yes </CommandLine>
                    <Description>Win RM adv firewall enable</Description>
                    <Order>13</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c netsh firewall add portopening TCP 5985 "Port 5985" </CommandLine>
                    <Description>Win RM port open</Description>
                    <Order>14</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c net stop winrm </CommandLine>
                    <Description>Stop Win RM Service </Description>
                    <Order>15</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c sc config winrm start= auto</CommandLine>
                    <Description>Win RM Autostart</Description>
                    <Order>16</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c net start winrm</CommandLine>
                    <Description>Start Win RM Service</Description>
                    <Order>17</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced\ /v HideFileExt /t REG_DWORD /d 0 /f</CommandLine>
                    <Order>18</Order>
                    <Description>Show file extensions in Explorer</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKCU\Console /v QuickEdit /t REG_DWORD /d 1 /f</CommandLine>
                    <Order>19</Order>
                    <Description>Enable QuickEdit mode</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced\ /v Start_ShowRun /t REG_DWORD /d 1 /f</CommandLine>
                    <Order>20</Order>
                    <Description>Show Run command in Start Menu</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Explorer\Advanced\ /v StartMenuAdminTools /t REG_DWORD /d 1 /f</CommandLine>
                    <Order>21</Order>
                    <Description>Show Administrative Tools in Start Menu</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKLM\SYSTEM\CurrentControlSet\Control\Power\ /v HibernateFileSizePercent /t REG_DWORD /d 0 /f</CommandLine>
                    <Order>22</Order>
                    <Description>Zero Hibernation File</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>%SystemRoot%\System32\reg.exe ADD HKLM\SYSTEM\CurrentControlSet\Control\Power\ /v HibernateEnabled /t REG_DWORD /d 0 /f</CommandLine>
                    <Order>23</Order>
                    <Description>Disable Hibernation Mode</Description>
                </SynchronousCommand>
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c wmic useraccount where "name='vagrant'" set PasswordExpires=FALSE</CommandLine>
                    <Order>24</Order>
                    <Description>Disable password expiration for vagrant user</Description>
                </SynchronousCommand>
                {{ if .WindowsUpdates }}
                <!-- Include non-Windows MS updates -->
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c a:\microsoft-updates.bat</CommandLine>
                    <Order>99</Order>
                    <Description>Enable Microsoft Updates</Description>
                </SynchronousCommand>
                <!-- Install Windows Updates, win-updates.ps1 will start SSH when done -->
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe -File a:\win-updates.ps1</CommandLine>
                    <Description>Install Windows Updates</Description>
                    <Order>100</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                {{ else }}
                <!-- Skipping Windows Updates, directly start SSH -->
                <SynchronousCommand wcm:action="add">
                    <CommandLine>cmd.exe /c C:\Windows\System32\WindowsPowerShell\v1.0\powershell.exe -File a:\openssh.ps1 -AutoStart</CommandLine>
                    <Description>Install OpenSSH</Description>
                    <Order>100</Order>
                    <RequiresUserInput>true</RequiresUserInput>
                </SynchronousCommand>
                {{ end }}
            </FirstLogonCommands>
            <ShowWindowsLive>false</ShowWindowsLive>
        </component>
    </settings>
    <settings pass="specialize">
        <component name="Microsoft-Windows-Shell-Setup" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <OEMInformation>
                <HelpCustomized>false</HelpCustomized>
            </OEMInformation>
            <!-- Rename computer here. -->
            <ComputerName>vagrant-{{Replace .OSName "windows" "win" -1}}</ComputerName>
            <TimeZone>Pacific Standard Time</TimeZone>
            <RegisteredOwner/>
        </component>
        <component xmlns:wcm="http://schemas.microsoft.com/WMIConfig/2002/State" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" name="Microsoft-Windows-Security-SPP-UX" processorArchitecture="amd64" publicKeyToken="31bf3856ad364e35" language="neutral" versionScope="nonSxS">
            <SkipAutoActivation>true</SkipAutoActivation>
        </component>
    </settings>
    <cpi:offlineImage xmlns:cpi="urn:schemas-microsoft-com:cpi" cpi:source="catalog:d:/sources/install_windows 7 ENTERPRISE.clg"/>
</unattend>
`
	pt.VagrantfileTpl = `# -*- mode: ruby -*-
# vi: set ft=ruby :

Vagrant.require_version ">= 1.6.2"

Vagrant.configure("2") do |config|
	config.vm.box = "{{.OSName}}"
	config.vm.communicator = "winrm"

	# Admin user name and password
	config.winrm.username = "{{.Username}}"
	config.winrm.password = "{{.Password}}"

	config.vm.network :forwarded_port, guest: 3389, host: 3389, id: "rdp", auto_correct: true
	config.vm.network :forwarded_port, guest: 22, host: 2222, id: "ssh", auto_correct: true

	config.vm.provider :virtualbox do |v, override|
		v.customize ["modifyvm", :id, "--memory", {{.RAM}}]
		v.customize ["modifyvm", :id, "--cpus", {{.CPU}}]
		v.customize ["setextradata", "global", "GUI/SuppressMessages", "all" ]
	end

  config.vm.provider :vmware_fusion do |v, override|
		v.vmx["memsize"] = "{{.RAM}}"
		v.vmx["numvcpus"] = "{{.CPU}}"
		v.vmx["ethernet0.virtualDev"] = "vmxnet3"
		v.vmx["RemoteDisplay.vnc.enabled"] = "false"
		v.vmx["RemoteDisplay.vnc.port"] = "5900"
		v.vmx["scsi0.virtualDev"] = "lsisas1068"
	end

	config.vm.provider :vmware_workstation do |v, override|
		v.vmx["memsize"] = "{{.RAM}}"
		v.vmx["numvcpus"] = "{{.CPU}}"
		v.vmx["ethernet0.virtualDev"] = "vmxnet3"
		v.vmx["RemoteDisplay.vnc.enabled"] = "false"
		v.vmx["RemoteDisplay.vnc.port"] = "5900"
		v.vmx["scsi0.virtualDev"] = "lsisas1068"
	end
end
`
	return pt
}
