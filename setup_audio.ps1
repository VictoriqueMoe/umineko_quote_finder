$AudioDir = "internal\quote\data\audio"

$EnvFile = Join-Path $PSScriptRoot ".env"
if (Test-Path $EnvFile) {
    Get-Content $EnvFile | ForEach-Object {
        if ($_ -match '^\s*([^#][^=]+)=(.*)$') {
            [System.Environment]::SetEnvironmentVariable($Matches[1].Trim(), $Matches[2].Trim(), 'Process')
        }
    }
}

$ZipSource = $env:VOICE_ZIP_URL
if (-not $ZipSource) {
    Write-Error "VOICE_ZIP_URL is not set. Create a .env file with VOICE_ZIP_URL=<url or path>"
    exit 1
}

if (Test-Path $AudioDir) {
    Write-Output "Audio directory already exists at $AudioDir, skipping download."
    exit 0
}

$TmpDir = "$env:TEMP\voice"

if (Test-Path $ZipSource) {
    Write-Output "Extracting from local file: $ZipSource"
    New-Item -ItemType Directory -Force -Path "internal\quote\data" | Out-Null
    Expand-Archive -Path $ZipSource -DestinationPath $TmpDir
    Move-Item -Path "$TmpDir\voice" -Destination $AudioDir
    Remove-Item -Recurse -Force $TmpDir
} else {
    $TmpZip = "$env:TEMP\voice.zip"
    Write-Output "Downloading voice files..."
    Invoke-WebRequest -Uri $ZipSource -OutFile $TmpZip
    Write-Output "Extracting..."
    New-Item -ItemType Directory -Force -Path "internal\quote\data" | Out-Null
    Expand-Archive -Path $TmpZip -DestinationPath $TmpDir
    Move-Item -Path "$TmpDir\voice" -Destination $AudioDir
    Remove-Item -Recurse -Force $TmpZip, $TmpDir
}

Write-Output "Done. Audio files extracted to $AudioDir"

$SeDir = "internal\quote\data\se"
$SeSource = $env:SE_ZIP_URL

if (-not $SeSource) {
    Write-Output "SE_ZIP_URL is not set, skipping SE download."
    exit 0
}

if (Test-Path $SeDir) {
    Write-Output "SE directory already exists at $SeDir, skipping download."
    exit 0
}

$TmpSeDir = "$env:TEMP\se"

if (Test-Path $SeSource) {
    Write-Output "Extracting SE from local file: $SeSource"
    New-Item -ItemType Directory -Force -Path "internal\quote\data" | Out-Null
    Expand-Archive -Path $SeSource -DestinationPath $TmpSeDir
    Move-Item -Path "$TmpSeDir\se" -Destination $SeDir
    Remove-Item -Recurse -Force $TmpSeDir
} else {
    $TmpSeZip = "$env:TEMP\se.zip"
    Write-Output "Downloading SE files..."
    Invoke-WebRequest -Uri $SeSource -OutFile $TmpSeZip
    Write-Output "Extracting..."
    New-Item -ItemType Directory -Force -Path "internal\quote\data" | Out-Null
    Expand-Archive -Path $TmpSeZip -DestinationPath $TmpSeDir
    Move-Item -Path "$TmpSeDir\se" -Destination $SeDir
    Remove-Item -Recurse -Force $TmpSeZip, $TmpSeDir
}

Write-Output "Done. SE files extracted to $SeDir"
