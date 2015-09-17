@echo off

cd..
set SCRIPTPATH=%cd%

: top src directory
set path="%SCRIPTPATH%\bin\src\github.com\leanote"

if not exist "%path%" mkdir "%path%"

: create software link
if exist "%path%\leanote" del /Q "%path%\leanote"
mklink /D "%path%\leanote"  %SCRIPTPATH%

: set GOPATH
set GOPATH="%SCRIPTPATH%\bin"

: run
if %processor_architecture%==x86 (
	"%SCRIPTPATH%\bin\leanote-windows-386.exe" -importPath github.com/leanote/leanote
) else (
	"%SCRIPTPATH%\bin\leanote-windows-amd64.exe" -importPath github.com/leanote/leanote
)
