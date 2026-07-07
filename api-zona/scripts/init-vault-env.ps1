param(
    [string]$VaultToken = $env:VAULT_TOKEN
)

if (-not $VaultToken) {
    $VaultToken = Read-Host -Prompt "Ingresa tu Vault Token (o define la variable de entorno VAULT_TOKEN)"
}

if (-not $VaultToken) {
    Write-Error "El Vault Token es requerido para continuar."
    exit 1
}

$VaultAddr = "http://192.168.1.20:8200"
$VaultSecretPath = "v1/retromarket/data/dev/api-auth"
$EnvFile = "..\.env"
$CurrentDir = Split-Path -Parent $MyInvocation.MyCommand.Path

$Url = "$VaultAddr/$VaultSecretPath"
Write-Host "Obteniendo secretos desde Vault: $Url"

$Headers = @{
    "X-Vault-Token" = $VaultToken
}

try {
    $Response = Invoke-RestMethod -Uri $Url -Headers $Headers -Method Get
    
    if ($Response.data -and $Response.data.data) {
        $Secrets = $Response.data.data
        $EnvContent = "# ==============================================================================`n"
        $EnvContent += "# ARCHIVO .ENV GENERADO AUTOMÁTICAMENTE DESDE HASHICORP VAULT`n"
        $EnvContent += "# Generado el: $(Get-Date -Format 'yyyy-MM-dd HH:mm:ss')`n"
        $EnvContent += "# ==============================================================================`n`n"

        # Update environment values to localhost if needed for local dev
        foreach ($Key in $Secrets.PSObject.Properties) {
            $Value = $Key.Value
            # Para entorno local, podríamos querer sobreescribir algunos valores (por ejemplo localhost en lugar de otra IP)
            # si es necesario:
            # if ($Key.Name -eq "DB_HOST") { $Value = "localhost" }
            $EnvContent += "$($Key.Name)=$Value`n"
        }

        $EnvFilePath = Join-Path -Path $CurrentDir -ChildPath "..\.env"
        Set-Content -Path $EnvFilePath -Value $EnvContent -Encoding UTF8
        Write-Host "Archivo .env generado exitosamente en: $(Resolve-Path $EnvFilePath)" -ForegroundColor Green
    } else {
        Write-Error "No se encontraron secretos en la ruta especificada."
    }
} catch {
    Write-Error "Error al conectar a Vault o al procesar la respuesta:"
    Write-Error $_.Exception.Message
}
