@echo off

SET gspath=%cd%/../


cd %gspath%/pb
del /f/s/q *.pb.go

echo "5. gen pb"
cd %gspath%
protoc -I=./vendor/github.com/gogo/protobuf/protobuf --proto_path=. --gofast_out=plugins=grpc:. proto/*.proto
cd %gspath%/proto
copy *.pb.go ..\pb
del /f/s/q *.pb.go

echo "finsihed..."

pause
