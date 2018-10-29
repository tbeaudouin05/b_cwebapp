@ ECHO OFF

if not exist "log\" mkdir log\

CALL "C:\Program Files\R\R-3.4.2\bin\x64\RScript.exe" "pricerecommendation.R" > log\log_%DATE:~-4%-%DATE:~4,2%-%DATE:~7,2%-%TIME:~0,2%h%TIME:~3,2%m.txt 2>&1