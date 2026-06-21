@echo off
setlocal enabledelayedexpansion

echo 🔨 CTX Build Script for Windows
echo =================================

REM Configurações
set VERSION=0.1.1
set OUTPUT_DIR=bin
set MAIN_PATH=./cmd/ctx

REM Cores (ANSI no Windows 10+)
set "GREEN=[92m"
set "YELLOW=[93m"
set "RED=[91m"
set "CYAN=[96m"
set "RESET=[0m"

echo.
echo %CYAN%📦 Building CTX version %VERSION%%RESET%

REM Cria diretório de output
if not exist %OUTPUT_DIR% mkdir %OUTPUT_DIR%

REM Detecta Git para versão
where git >nul 2>nul
if %errorlevel% equ 0 (
    for /f "tokens=*" %%a in ('git describe --tags --always --dirty 2^>nul') do set GIT_VERSION=%%a
    if defined GIT_VERSION (
        set VERSION=%GIT_VERSION%
        echo %GREEN%✓ Git version detected: %VERSION%%RESET%
    )
)

REM Build para Windows (amd64)
echo.
echo %CYAN%🔨 Building for Windows (amd64)...%RESET%
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-X main.version=%VERSION%" -o %OUTPUT_DIR%\ctx-windows-amd64.exe %MAIN_PATH%
if %errorlevel% equ 0 (
    echo %GREEN%  ✓ ctx-windows-amd64.exe built successfully%RESET%
) else (
    echo %RED%  ✗ Failed to build for Windows amd64%RESET%
    exit /b 1
)

REM Build para Windows (386)
echo.
echo %CYAN%🔨 Building for Windows (386)...%RESET%
set GOARCH=386
go build -ldflags="-X main.version=%VERSION%" -o %OUTPUT_DIR%\ctx-windows-386.exe %MAIN_PATH%
if %errorlevel% equ 0 (
    echo %GREEN%  ✓ ctx-windows-386.exe built successfully%RESET%
) else (
    echo %RED%  ✗ Failed to build for Windows 386%RESET%
)

REM Build para Linux (amd64)
echo.
echo %CYAN%🔨 Building for Linux (amd64)...%RESET%
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-X main.version=%VERSION%" -o %OUTPUT_DIR%\ctx-linux-amd64 %MAIN_PATH%
if %errorlevel% equ 0 (
    echo %GREEN%  ✓ ctx-linux-amd64 built successfully%RESET%
) else (
    echo %RED%  ✗ Failed to build for Linux amd64%RESET%
)

REM Build para Linux (arm64)
echo.
echo %CYAN%🔨 Building for Linux (arm64)...%RESET%
set GOARCH=arm64
go build -ldflags="-X main.version=%VERSION%" -o %OUTPUT_DIR%\ctx-linux-arm64 %MAIN_PATH%
if %errorlevel% equ 0 (
    echo %GREEN%  ✓ ctx-linux-arm64 built successfully%RESET%
) else (
    echo %RED%  ✗ Failed to build for Linux arm64%RESET%
)

REM Build para macOS (amd64)
echo.
echo %CYAN%🔨 Building for macOS (amd64)...%RESET%
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-X main.version=%VERSION%" -o %OUTPUT_DIR%\ctx-darwin-amd64 %MAIN_PATH%
if %errorlevel% equ 0 (
    echo %GREEN%  ✓ ctx-darwin-amd64 built successfully%RESET%
) else (
    echo %RED%  ✗ Failed to build for macOS amd64%RESET%
)

REM Build para macOS (arm64)
echo.
echo %CYAN%🔨 Building for macOS (arm64)...%RESET%
set GOARCH=arm64
go build -ldflags="-X main.version=%VERSION%" -o %OUTPUT_DIR%\ctx-darwin-arm64 %MAIN_PATH%
if %errorlevel% equ 0 (
    echo %GREEN%  ✓ ctx-darwin-arm64 built successfully%RESET%
) else (
    echo %RED%  ✗ Failed to build for macOS arm64%RESET%
)

REM Mostra resultados
echo.
echo %CYAN%📊 Build Summary:%RESET%
echo ====================
dir %OUTPUT_DIR%

echo.
echo %GREEN%✅ Build completed!%RESET%
echo Binaries available in %OUTPUT_DIR%\