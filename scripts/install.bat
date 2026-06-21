@echo off
setlocal enabledelayedexpansion

echo 📦 CTX Install Script for Windows
echo ==================================

REM Cores
set "GREEN=[92m"
set "YELLOW=[93m"
set "RED=[91m"
set "CYAN=[96m"
set "RESET=[0m"

REM Detecta arquitetura
if "%PROCESSOR_ARCHITECTURE%"=="AMD64" (
    set ARCH=amd64
) else if "%PROCESSOR_ARCHITECTURE%"=="x86" (
    set ARCH=386
) else (
    set ARCH=amd64
)

echo.
echo %CYAN%🔍 Detected architecture: %ARCH%%RESET%

REM Define caminho de instalação
if defined GOBIN (
    set INSTALL_DIR=%GOBIN%
) else if defined GOPATH (
    set INSTALL_DIR=%GOPATH%\bin
) else (
    set INSTALL_DIR=%USERPROFILE%\go\bin
)

echo %CYAN%📁 Install directory: %INSTALL_DIR%%RESET%

REM Cria diretório se não existir
if not exist "%INSTALL_DIR%" (
    echo %CYAN%📁 Creating install directory...%RESET%
    mkdir "%INSTALL_DIR%"
)

REM Build para Windows
echo.
echo %CYAN%🔨 Building CTX for Windows...%RESET%
set GOOS=windows
set GOARCH=%ARCH%
go build -ldflags="-X main.version=dev" -o "%INSTALL_DIR%\ctx.exe" ./cmd/ctx

if %errorlevel% neq 0 (
    echo %RED%✗ Build failed%RESET%
    exit /b 1
)

echo %GREEN%✓ CTX installed successfully to %INSTALL_DIR%\ctx.exe%RESET%

REM Verifica se está no PATH
echo.
echo %CYAN%🔍 Checking PATH...%RESET%
echo %PATH% | findstr /C:"%INSTALL_DIR%" >nul
if %errorlevel% equ 0 (
    echo %GREEN%✓ Install directory is in PATH%RESET%
) else (
    echo %YELLOW%⚠️ Install directory is NOT in PATH%RESET%
    echo.
    echo To add to PATH, run:
    echo   setx PATH "%%PATH%%;%INSTALL_DIR%"
    echo.
    echo Or add manually to your system PATH
)

echo.
echo %CYAN%🎯 Testing installation...%RESET%
"%INSTALL_DIR%\ctx.exe" version >nul 2>nul
if %errorlevel% equ 0 (
    echo %GREEN%✓ CTX is working!%RESET%
    "%INSTALL_DIR%\ctx.exe" version
) else (
    echo %RED%✗ Installation test failed%RESET%
)

echo.
echo %GREEN%✅ Installation complete!%RESET%