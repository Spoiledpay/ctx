@echo off
setlocal enabledelayedexpansion

echo 🚀 CTX Release Script for Windows
echo ==================================

REM Cores
set "GREEN=[92m"
set "YELLOW=[93m"
set "RED=[91m"
set "CYAN=[96m"
set "RESET=[0m"

REM Configurações
set RELEASE_DIR=release
set VERSION_FILE=VERSION

REM Lê versão atual
if exist %VERSION_FILE% (
    set /p VERSION=<%VERSION_FILE%
) else (
    set VERSION=0.1.0
)

echo.
echo %CYAN%📦 Current version: %VERSION%%RESET%

REM Pergunta nova versão
set /p NEW_VERSION="%YELLOW%Enter new version (or press Enter to keep %VERSION%): %RESET%"
if defined NEW_VERSION (
    set VERSION=%NEW_VERSION%
    echo %VERSION% > %VERSION_FILE%
    echo %GREEN%✓ Version updated to %VERSION%%RESET%
)

REM Verifica se há mudanças não commitadas
echo.
echo %CYAN%🔍 Checking git status...%RESET%
git status --porcelain > git_status.tmp
set /p GIT_STATUS=<git_status.tmp
if defined GIT_STATUS (
    echo %YELLOW%⚠️ There are uncommitted changes:%RESET%
    type git_status.tmp
    echo.
    set /p CONFIRM="%YELLOW%Continue anyway? (y/N): %RESET%"
    if /i not "!CONFIRM!"=="y" (
        echo %RED%✗ Release cancelled%RESET%
        del git_status.tmp
        exit /b 1
    )
)
del git_status.tmp 2>nul

REM Cria diretório de release
echo.
echo %CYAN%📁 Creating release directory...%RESET%
if exist %RELEASE_DIR% rmdir /s /q %RELEASE_DIR%
mkdir %RELEASE_DIR%

REM Roda testes
echo.
echo %CYAN%🧪 Running tests...%RESET%
call scripts\test.bat
if %errorlevel% neq 0 (
    echo %RED%✗ Tests failed, aborting release%RESET%
    exit /b 1
)

REM Build para todas as plataformas
echo.
echo %CYAN%🔨 Building for all platforms...%RESET%
call scripts\build.bat
if %errorlevel% neq 0 (
    echo %RED%✗ Build failed%RESET%
    exit /b 1
)

REM Gera checksums
echo.
echo %CYAN%🔐 Generating checksums...%RESET%
cd %RELEASE_DIR%
for %%f in (*) do (
    certutil -hashfile "%%f" SHA256 > "%%f.sha256"
    echo %GREEN%  ✓ %%f.sha256 generated%RESET%
)
cd ..

REM Cria arquivo de release notes
echo.
echo %CYAN%📝 Creating release notes...%RESET%
set RELEASE_NOTES=%RELEASE_DIR%\RELEASE_NOTES.md

echo # CTX Release %VERSION% > %RELEASE_NOTES%
echo. >> %RELEASE_NOTES%
echo Release date: %date% >> %RELEASE_NOTES%
echo. >> %RELEASE_NOTES%
echo ## Changes >> %RELEASE_NOTES%
echo. >> %RELEASE_NOTES%

REM Pega últimas commits
git log --pretty=format:"- %s" $(git describe --tags --abbrev=0)..HEAD 2>nul >> %RELEASE_NOTES%
if %errorlevel% neq 0 (
    echo - Initial release >> %RELEASE_NOTES%
)

echo. >> %RELEASE_NOTES%
echo ## Binaries >> %RELEASE_NOTES%
echo. >> %RELEASE_NOTES%

for %%f in (%RELEASE_DIR%\*) do (
    if not "%%~xf"==".sha256" if not "%%~xf"==".md" (
        echo - [%%~nxf](%%~nxf) >> %RELEASE_NOTES%
    )
)

echo %GREEN%✓ Release notes created%RESET%

echo.
echo %CYAN%📊 Release summary:%RESET%
echo ====================
dir %RELEASE_DIR%

echo.
echo %GREEN%✅ Release %VERSION% prepared!%RESET%
echo Release files available in %RELEASE_DIR%\