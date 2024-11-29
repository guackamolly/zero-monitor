@ECHO OFF

rem Save program arguments to later pass on init binary.
set init_args=%*

rem Common urls.
set new_issue_url=https://github.com/guackamolly/zero-monitor/issues/new
set latest_release_url=https://api.github.com/repos/guackamolly/zero-monitor/releases/latest
set jq_release_url=https://github.com/jqlang/jq/releases/download/jq-1.7.1

rem Installation directory and program paths.
set install_dir=%~dp0\config\zero-monitor
set bin_path=%install_dir%\node.exe
set init_bin_path=%install_dir%\init.exe
set jq_bin_path=%install_dir%\jq.exe

if not exist "%install_dir%" (
  mkdir "%install_dir%"
)

:fatal
  echo %1
  exit /b 1

:jq
  "%jq_bin_path%" %*

:exec_bin
  if "%init_args%" NEQ "" (
    "%init_bin_path%" %init_args%
  )
  call "%bin_path%"

:download
  set url=%1
  if "%url%" EQU "" (
    echo Failed to extract url, please raise an issue to alert maintainers about this bug. ^&%new_issue_url%
    goto :fatal
  )
  set bin_name=%url:~-10%

  echo Downloading %bin_name% ...
  for /f "usebackq delims=\\" %%a in (`where wget`) do (
    %%a -O "%install_dir%\%bin_name%" "%url%"
  )
  attrib +x "%install_dir%\%bin_name%"

rem Query host OS
set os=windos

rem Query CPU architecture
set arch=
for /f "tokens=2 delims=:" %%a in ('wmic cpu get architecture /format:List') do set "arch=%%a"

if "%arch%" EQU "x86_64" (
  set arch=amd64
) else if "%arch%" EQU "i686" (
  set arch=386
) else (
  echo "%arch%" is not supported right now, please raise an issue to get support on this architecture. ^&%new_issue_url%
  goto :fatal
)

rem Download jq if not available.
if not exist "%jq_bin_path%" (
  call :download "%jq_release_url%\jq-%os%-%arch%"
)

rem Head latest release
for /f "usebackq delims=:" %%a in (`curl -s ^&%latest_release_url%`) do set "response=%%a"
if not "!errorlevel!" EQU "0" (
  echo Failed to head release, please raise an issue to alert maintainers about this bug. ^&%new_issue_url%
  goto :fatal
)

set latest_release_version=%response:~10,1%

if "%latest_release_version%" EQU "" (
  echo Failed to extract release version, please raise an issue to alert maintainers about this bug. ^&%new_issue_url%
  goto :fatal
)

rem If local init binary does not exist and no arguments have been passed, then bootstrap node.
if not exist "%init_bin_path%" AND "%init_args%" EQU "" (
  set init_args="--node="
)

rem If local target binary version is different than the latest release version, download it again.
if exist "%bin_path%" AND "%latest_release_version%" NEQ "%bin_path% -version%" (
  for /f "usebackq delims=:" %%a in (`echo ^&response% | jq -r '.assets[] | select(.name == "node_%os%_%arch%") | .browser_download_url'`)