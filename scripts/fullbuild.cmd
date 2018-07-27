@echo off
echo -- Generators
vgo generate .\...
echo -- Start tests
vgo test -coverprofile=cover.out .\... 

if errorlevel 1 (
	echo -- Tests Failed
	exit /b %errorlevel%
)
echo -- Building
vgo install
vgo install .\cmd\...
echo -- Done.
rem vgo tool cover -func=cover.out
rem vgo tool cover -html=cover.out