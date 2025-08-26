@echo off
setlocal enabledelayedexpansion
cd /d "%~dp0"
git init -b main
if not defined GIT_USER_NAME (
  set /p GIT_USER_NAME=Enter git user.name: 
)
if not defined GIT_USER_EMAIL (
  set /p GIT_USER_EMAIL=Enter git user.email: 
)
git config user.name "%GIT_USER_NAME%"
git config user.email "%GIT_USER_EMAIL%"
git config core.filemode false
git config core.autocrlf false
git config --global push.autoSetupRemote true
git branch -M main
rem Add origin if missing; otherwise update it
git remote get-url origin 1>nul 2>nul && (
  git remote set-url origin git@github.com:mechiko/ktrade_xlsx_agregation.git
) || (
  git remote add origin git@github.com:mechiko/ktrade_xlsx_agregation.git
)