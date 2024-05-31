protoc --go_out=./pkg/pb --go_opt=module=github.com/uandersonricardo/masterclass-go/pkg/pb \
       --go-grpc_out=./pkg/pb --go-grpc_opt=module=github.com/uandersonricardo/masterclass-go/pkg/pb \
       ./protos/example.proto
