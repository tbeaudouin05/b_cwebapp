@ ECHO OFF

::192.168.100.160.xip.io
::127.0.0.1
SET APIHandlerIP=192.168.100.160.xip.io
::127.0.0.1
::192.168.100.160
SET WebsocketIP=192.168.100.160
SET APIHandlerPort=8082
SET WebsocketPort=4000
SET ListOfFrontend=(adduser,login,bmldgkskuhist,bmldgkskumanualmatch,bmldgktable,categorypage,homepage,suppliermatching)
SET ListOfFrontendB=(doNotErase,adduser,login,bmldgkskuhist,bmldgkskumanualmatch,bmldgktable,categorypage,homepage,suppliermatching)
SET FrontendFolderName=frontend

setlocal ENABLEDELAYEDEXPANSION
set slash=/
set homeDir=%~dp0
set homeDir=%homeDir:\=!slash!%

FOR %%A IN %ListOfFrontend% DO echo {"APIHandlerIP":"%APIHandlerIP%","WebsocketIP":"%WebsocketIP%","APIHandlerPort":"%APIHandlerPort%","WebsocketPort":"%WebsocketPort%"} > %homeDir%%FrontendFolderName%/%%A/src/envConfig.json
echo {"APIHandlerIP":"%APIHandlerIP%","WebsocketIP":"%WebsocketIP%","APIHandlerPort":"%APIHandlerPort%","WebsocketPort":"%WebsocketPort%"} > envConfig.json
go build

FOR %%A IN %ListOfFrontendB% DO (
    cd %homeDir%%FrontendFolderName%/%%A
    npm install
    yarn build
)