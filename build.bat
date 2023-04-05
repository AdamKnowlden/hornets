cd web/panel
cmd /C npm i
cd ../..
cd services/server
cmd /C go build -o ../../
cd ../client
cmd /C go build -o ../../
pause