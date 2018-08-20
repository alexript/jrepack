@echo off

rem Copyright (C) 2018  Alexander Malyshev

rem Permission is hereby granted, free of charge, to any person obtaining a copy
rem of this software and associated documentation files (the "Software"), to deal
rem in the Software without restriction, including without limitation the rights
rem to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
rem copies of the Software, and to permit persons to whom the Software is
rem furnished to do so, subject to the following conditions:
rem The above copyright notice and this permission notice shall be included in
rem all copies or substantial portions of the Software.

rem THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
rem IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
rem FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
rem AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
rem LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
rem OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
rem THE SOFTWARE.


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
