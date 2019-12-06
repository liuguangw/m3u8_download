@echo off
SET CGO_ENABLED=0
SET GOARCH=amd64
SET GO111MODULE=on
SET GOPROXY=https://goproxy.cn
SET GOOS=windows
SET DIST_FILE_NAME=m3u8_download.exe
goto :do_build

Rem for linux.
:os_linux
  SET GOOS=linux
  SET DIST_FILE_NAME=m3u8_download

Rem do build task.
:do_build
  echo build for %GOOS%^<%GOARCH%^>
  go build -o %DIST_FILE_NAME%
  if %ERRORLEVEL% NEQ 0 (
    pause
    exit
  )
  if "%GOOS%" == "windows" goto :os_linux
  echo build complete
  pause