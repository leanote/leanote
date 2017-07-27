@echo off

cd..
set SCRIPTPATH=%cd%

: top src directory
set leanotePath="%SCRIPTPATH%\bin\src\github.com\leanote"

if not exist "%leanotePath%" mkdir "%leanotePath%"

: create software link
if exist "%leanotePath%\leanote" del /Q "%leanotePath%\leanote"
mklink /D "%leanotePath%\leanote"  %SCRIPTPATH%

: set GOPATH
set GOPATH="%SCRIPTPATH%\bin"

: run
if %processor_architecture%==x86 (
	"%SCRIPTPATH%\bin\leanote-windows-386.exe" -importPath github.com/leanote/leanote
) else (
	"%SCRIPTPATH%\bin\leanote-windows-amd64.exe" -importPath github.com/leanote/leanote
)
