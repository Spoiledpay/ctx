@echo off
setlocal enabledelayedexpansion

echo 🧪 CTX Test Script for Windows
echo ===============================

REM Cores
set "GREEN=[92m"
set "YELLOW=[93m"
set "RED=[91m"
set "CYAN=[96m"
set "RESET=[0m"

REM Configurações
set COVERAGE_FILE=coverage.out
set TEST_TIMEOUT=5m

echo.
echo %CYAN%📦 Running tests...%RESET%

REM Testes com race detection
echo.
echo %CYAN%🔍 Running tests with race detection...%RESET%
go test -race -timeout=%TEST_TIMEOUT% ./...
if %errorlevel% neq 0 (
    echo %RED%  ✗ Race detection tests failed%RESET%
    exit /b 1
) else (
    echo %GREEN%  ✓ Race detection tests passed%RESET%
)

REM Testes com cobertura
echo.
echo %CYAN%📊 Running tests with coverage...%RESET%
go test -coverprofile=%COVERAGE_FILE% -covermode=atomic -timeout=%TEST_TIMEOUT% ./...
if %errorlevel% neq 0 (
    echo %RED%  ✗ Coverage tests failed%RESET%
    exit /b 1
) else (
    echo %GREEN%  ✓ Coverage tests passed%RESET%
)

REM Mostra cobertura
echo.
echo %CYAN%📈 Test Coverage:%RESET%
go tool cover -func=%COVERAGE_FILE%

REM Gera relatório HTML
echo.
echo %CYAN%🌐 Generating HTML coverage report...%RESET%
go tool cover -html=%COVERAGE_FILE% -o coverage.html
if %errorlevel% equ 0 (
    echo %GREEN%  ✓ coverage.html generated%RESET%
)

REM Testes de benchmark
echo.
echo %CYAN%⏱️  Running benchmarks...%RESET%
go test -bench=. -benchmem -timeout=%TEST_TIMEOUT% ./... > benchmark.txt
if %errorlevel% equ 0 (
    echo %GREEN%  ✓ Benchmarks completed (see benchmark.txt)%RESET%
)

REM Testes de lint
echo.
echo %CYAN%🔍 Running go vet...%RESET%
go vet ./...
if %errorlevel% equ 0 (
    echo %GREEN%  ✓ go vet passed%RESET%
) else (
    echo %RED%  ✗ go vet failed%RESET%
)

REM Verifica formatação
echo.
echo %CYAN%📝 Checking code formatting...%RESET%
set UNFORMATTED_FILES=0
for /f "tokens=*" %%a in ('gofmt -l .') do (
    if not "%%a"=="" (
        echo %YELLOW%  ! %%a is not formatted%RESET%
        set UNFORMATTED_FILES=1
    )
)
if !UNFORMATTED_FILES! equ 0 (
    echo %GREEN%  ✓ All files are properly formatted%RESET%
) else (
    echo %YELLOW%  ⚠️ Run 'gofmt -w .' to fix formatting%RESET%
)

echo.
echo %GREEN%✅ All tests completed!%RESET%