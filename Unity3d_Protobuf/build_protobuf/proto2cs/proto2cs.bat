@echo off
:: mode con:cols=120 lines=40

:: 重设当前的CMD路径 

:: 就是进入批处理所在目录了

cd..
cd..
set cs_out_file=%cd%\Assets\Scripts\NetWork\Proto\

cd /d %~dp0 

set currentPath=%cd%

::echo.%currentPath%

:: _protoSrc 是你的proto文件目录的位置
set _protoSrc=proto

:: protoExe 是用于从proto生成cs的protogen.exe程序的位置
set protoExe=protoc

:: cs_out_file 存放生成的cs文件目录的位置
::set cs_out_file=%currentPath%\Assets\Scripts\NetWork\Proto\
echo.%cs_out_file%

for /R "%_protoSrc%" %%i in (*) do (
 if "%%~xi"  == ".proto" (

 :: echo.%protoExe% -i:%%i -o:%cs_out_file%%%~ni.cs -ns:MyProtoBuf
 ::%protoExe% -i:%%i -o:%cs_out_file%%%~ni.cs -ns:MyProtoBuf
 %protoExe% -I=%currentPath%\proto\  --csharp_out=%cs_out_file% %%i
  echo.%%~nxi to csharp ok !!!
 )
)

pause