where /q curl
IF ERRORLEVEL 1 (
    choco install curl -y --no-progress --skipdownloadcache
    IF ERRORLEVEL 1 (
        REM curl installation has failed.
        exit /b 1
    )
)
mkdir %WORKSPACE%\bin

IF EXIST "%PROGRAMFILES(X86)%" (
    REM Force the gvm installation.
    SET GVM_BIN=gvm.exe
    curl -L -o %WORKSPACE%\bin\gvm.exe https://github.com/andrewkroh/gvm/releases/download/v0.2.4/gvm-windows-amd64.exe
    IF ERRORLEVEL 1 (
        REM gvm installation has failed.
        exit /b 1
    )
) ELSE (
    REM Windows 7 workers got a broken gvm installation.
    curl -L -o %WORKSPACE%\bin\gvm.exe https://github.com/andrewkroh/gvm/releases/download/v0.2.4/gvm-windows-386.exe
    IF ERRORLEVEL 1 (
        REM gvm installation has failed.
        exit /b 1
    )
)

SET GVM_BIN=gvm.exe
WHERE /q %GVM_BIN%
%GVM_BIN% version

REM Install the given go version
%GVM_BIN% --debug install %GO_VERSION%

REM Configure the given go version
FOR /f "tokens=*" %%i IN ('"%GVM_BIN%" use %GO_VERSION% --format=batch') DO %%i

go env
IF ERRORLEVEL 1 (
    REM go is not configured correctly.
    exit /b 1
)

where /q gcc
IF ERRORLEVEL 1 (
    REM Install mingw 5.3.0
    choco install mingw -y -r --no-progress --version 5.3.0
    IF NOT ERRORLEVEL 0 (
        exit /b 1
    )
)
gcc --version
where gcc
