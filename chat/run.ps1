$ErrorActionPreference = "Stop"

$AppName = "litellm_chat"
$ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$DistDir = Join-Path $ScriptDir "dist"

$os = "windows"
$arch = [System.Runtime.InteropServices.RuntimeInformation]::ProcessArchitecture.ToString().ToLower()

switch ($arch) {
    "x64" { $goarch = "amd64" }
    "arm64" { $goarch = "arm64" }
    default {
        Write-Error "Unsupported architecture: $arch. Please run a matching binary from $DistDir manually."
    }
}

$bin = Join-Path $DistDir "$AppName-$os-$goarch.exe"

if (!(Test-Path $bin)) {
    Write-Error "Binary not found: $bin. Run build_all.sh first."
}

Write-Host "Running $bin ..."
& $bin
