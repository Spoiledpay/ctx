param (
    [switch]$SkipIncrement
)

$projectRoot = $PSScriptRoot
$versionFile = Join-Path $projectRoot "cmd\ctx\version.bin"
$outputDir = Join-Path $projectRoot "bin"
$mainPath = Join-Path $projectRoot "cmd\ctx"

if (-not (Test-Path $outputDir)) {
    New-Item -ItemType Directory -Path $outputDir -Force | Out-Null
}

if (-not $SkipIncrement) {
    if (Test-Path $versionFile) {
        $currentVersion = (Get-Content $versionFile -Raw).Trim()
        $parts = $currentVersion.Split('.')
        if ($parts.Count -eq 3) {
            $patch = [int]$parts[2] + 1
            $newVersion = "$($parts[0]).$($parts[1]).$patch"
        } else {
            $newVersion = "1.07.1"
        }
    } else {
        $newVersion = "1.07.1"
    }

    Set-Content -Path $versionFile -Value $newVersion -NoNewline
    Write-Host "Version incremented: $newVersion"
} else {
    $newVersion = (Get-Content $versionFile -Raw).Trim()
    Write-Host "Building version: $newVersion (skip increment)"
}

Write-Host "Building ctx.exe..."
Set-Location $projectRoot
go build -o (Join-Path $outputDir "ctx.exe") $mainPath

if ($LASTEXITCODE -eq 0) {
    Write-Host "Build successful: bin\ctx.exe"
} else {
    Write-Host "Build failed!"
    exit 1
}
