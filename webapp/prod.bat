@ ECHO OFF

SET ListOfFrontend=(adduser,login,bmldgkskuhist,bmldgkskumanualmatch,bmldgktable,categorypage,homepage,suppliermatching)
SET FrontendFolderName=frontend
SET ExecutableNameWithOUTExtension=webapp

setlocal ENABLEDELAYEDEXPANSION

    if exist "prod\" rd /S /Q "prod\"
    mkdir prod\
    mkdir prod\%FrontendFolderName%\

    FOR %%A IN %ListOfFrontend% DO mkdir prod\%FrontendFolderName%\%%A\
    FOR %%A IN %ListOfFrontend% DO xcopy %FrontendFolderName%\%%A\build prod\%FrontendFolderName%\%%A\build\
    FOR %%A IN %ListOfFrontend% DO xcopy %FrontendFolderName%\%%A\build\static\js prod\%FrontendFolderName%\%%A\build\static\js\
    FOR %%A IN %ListOfFrontend% DO xcopy %FrontendFolderName%\%%A\build\static\css prod\%FrontendFolderName%\%%A\build\static\css\

    xcopy %ExecutableNameWithOUTExtension%.exe prod\
    xcopy %ExecutableNameWithOUTExtension%.bat prod\
    xcopy envConfig.json prod\
    xcopy creds.json prod\

        