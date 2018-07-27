@echo off
echo -- Generators
vgo generate .\...
echo -- Start tests
vgo test .\...
if errorlevel 1 (
	echo -- Tests Failed
	exit /b %errorlevel%
)
echo -- Building
vgo install
vgo install .\cmd\...
echo -- Done.