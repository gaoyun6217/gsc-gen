@echo off
REM GFRD Framework Test Script (Windows)
REM 测试 gen 模块集成

echo ========================================
echo GFRD Framework - Gen Module Test
echo ========================================
echo.

cd /d "%~dp0"

echo Step 1: 检查 gen 模块...
if exist "gen\main.go" (
    echo   [OK] gen 模块存在
) else (
    echo   [FAIL] gen 模块不存在
    exit /b 1
)

echo.
echo Step 2: 构建 gen 模块...
cd gen
go build -o bin\gfrd-gen.exe main.go 2>&1
if errorlevel 1 (
    echo   [FAIL] gen 模块构建失败
    exit /b 1
)
echo   [OK] gen 模块构建成功
cd ..

echo.
echo Step 3: 测试 gen CLI 帮助...
gen\bin\gfrd-gen.exe --help >nul 2>&1
if errorlevel 1 (
    echo   [FAIL] gen CLI 执行失败
    exit /b 1
)
echo   [OK] gen CLI 执行成功

echo.
echo Step 4: 检查 server 目录...
if exist "server\go.mod" (
    echo   [OK] server 目录存在
) else (
    echo   [FAIL] server 目录不存在
    exit /b 1
)

echo.
echo Step 5: 检查 web 目录...
if exist "web\package.json" (
    echo   [OK] web 目录存在
) else (
    echo   [FAIL] web 目录不存在
    exit /b 1
)

echo.
echo Step 6: 测试从根目录运行 gen...
go run ./gen preview --help >nul 2>&1
if errorlevel 1 (
    echo   [FAIL] 从根目录运行 gen 失败
    exit /b 1
)
echo   [OK] 从根目录运行 gen 成功

echo.
echo ========================================
echo 所有测试通过！[OK]
echo ========================================
echo.
echo 使用方法:
echo   # 生成 CRUD 代码
echo   go run ./gen crud ^
echo     --table="sys_user" ^
echo     --db="mysql:root:123456@tcp(127.0.0.1:3306)/gfrd" ^
echo     --output="./server" ^
echo     --web-output="./web/src" ^
echo     --module="sys"
echo.
