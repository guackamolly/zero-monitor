@ECHO OFF & setlocal enabledelayedexpansion

rem Save program arguments to later pass on init binary.
set args=%*

rem Common urls.
set new_issue_url=https://github.com/guackamolly/zero-monitor/issues/new
set latest_release_url=https://api.github.com/repos/guackamolly/zero-monitor/releases/latest
set jq_release_url=https://github.com/jqlang/jq/releases/download/jq-1.7.1

rem Installation directory and program paths.
set install_dir=%APPDATA%\zero-monitor
set temp_input_dir=%TEMP%\.zero-monitor
set bin_path=%install_dir%\master.exe
set jq_bin_path=%install_dir%\jq.exe


if not exist "%install_dir%" (
  mkdir "%install_dir%"
)

if not exist "%temp_input_dir%" (
  mkdir "%temp_input_dir%"
)


rem Query host OS
set os=windows

rem Query CPU architecture
set arch=%PROCESSOR_ARCHITECTURE%
if "%arch%" EQU "AMD64" (
  set arch=amd64
) else if "%arch%" EQU "X86" (
  set arch=386
) else if "%arch%" EQU "ARM64" (
  set arch=arm64
) else (
  echo "%arch%" is not supported right now, please raise an issue to get support on this architecture. ^%new_issue_url%
  goto :fatal
)

rem Download jq if not available.
if not exist "%jq_bin_path%" (
  call :download "%jq_release_url%/jq-%os%-%arch%.exe" jq.exe
)

curl -s ^%latest_release_url% > %temp_input_dir%\latest-release

if not %ERRORLEVEL% EQU 0 (
  echo Failed to head release, please raise an issue to alert maintainers about this bug. ^&%new_issue_url%
  goto :fatal
)

rem Head latest release
call:jq .tag_name %temp_input_dir%\latest-release>%temp_input_dir%\version
set /P latest_release_version=<%temp_input_dir%\version

if %latest_release_version% EQU "" (
  echo Failed to extract release version, please raise an issue to alert maintainers about this bug. ^&%new_issue_url%
  goto :fatal
)

rem If local target binary version is different than the latest release version, download it again.
"%bin_path%" "-version">%temp_input_dir%\bin_version
set /P bin_version=<%temp_input_dir%\bin_version

if %latest_release_version% NEQ "%bin_version%" (
  call:jq -r ".assets[] | select(.name == \"master_%os%_%arch%\") | .browser_download_url" %temp_input_dir%\latest-release>%temp_input_dir%\download_url
  set /P download_url=<%temp_input_dir%\download_url
  call :download !download_url! master.exe
)

rem Run the binary.
call:exec_bin

REM %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%
REM %%%%%%%%%%% FUNCTION DEFINITIONS %%%%%%%%%%%%
REM %%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%

:fatal
    echo %1
    exit /b 1

:jq
  "%jq_bin_path%" %*
  exit /b 0

:exec_bin
  call "%bin_path%" "%args%"
  exit /b 0

:download
  set url=%1
  set bin_name=%2
  if "%url%" EQU "" (
    echo Failed to extract url, please raise an issue to alert maintainers about this bug. ^&%new_issue_url%
    goto :fatal
  )

  echo Downloading %bin_name% ...
  curl -L "%url%" -o "%install_dir%\%bin_name%"
  exit /b 0