Regenrate the grpc golang files from the proto files using the following command

protoc --go_out=./ --go-grpc_out=./ dag.proto

Use the build.bat to build the web app and then both the server and client exe
Make sure to use the bat files for starting both of them too as they have startup parameters