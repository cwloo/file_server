SET BINARY_NAME=loader.exe
SET BIN_DIR=..\..\deploy\bin

@REM windows
SET CGO_ENABLED=0 GOOS=windows GOARCH=amd64

go build -ldflags="-w -s"

move /Y %BINARY_NAME% %BIN_DIR%