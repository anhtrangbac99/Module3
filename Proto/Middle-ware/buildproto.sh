 protoc -I/usr/local/include -I. \
   -I$GOPATH/src \
   -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis \
   -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway \
   --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true:. \
   --swagger_out=logtostderr=true:. \
    protobuf.proto
   #mv *.json ../../../api-gateway/swagger/merchant.swagger.json