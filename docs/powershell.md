# PowerShell

## Server

- Install PowerShell Core 7

    ```
    Invoke-WebRequest -Uri https://github.com/PowerShell/PowerShell/releases/download/v7.0.1/PowerShell-7.0.1-win-x64.msi -OutFile PowerShell.msi
    msiexec.exe /package PowerShell.msi /quiet ENABLE_PSREMOTING=1 REGISTER_MANIFEST=1
    ```
- Enable PowerShell Remoting
    ```
    Enable-PSRemoting -Force
    ```
- Enable WinRM Basic Authentication

    ```
    winrm set winrm/config/Service/Auth @{Basic="true"}
    ```
- Ensure port 5986 (TCP) is open for ingress traffic.

## Client

- Use a Microsoft provided Docker image:

    ```
    docker run -it docker pull mcr.microsoft.com/powershell
    ```
- Connect using HTTPS and Basic Authentication:
    
    ```
    Enter-PSSession -ComputerName ${IP} -Authentication Basic -Credential ${USERNAME} -ConfigurationName "PowerShell.7" -UseSSL -SessionOption (New-PSSessionOption -SkipCACheck -SkipCNCheck)
    ```

Alternatively you may not want to have interactive sign in using the password.
In that case you will need to manually create a `PSCredential` object.

- `$pw = ConvertTo-SecureString "${PW}" -AsPlainText -Force`
- `$cred = New-Object System.Management.Automation.PSCredential ("${USERNAME}", $pw)`
- `Enter-PSSession -ComputerName ${IP} -Authentication Basic -Credential $cred -ConfigurationName "PowerShell.7" -UseSSL -SessionOption (New-PSSessionOption -SkipCACheck -SkipCNCheck)`
