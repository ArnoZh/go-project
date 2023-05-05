module commonserver

go 1.19

require (
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.7.1
	github.com/gogo/protobuf v1.3.2
	github.com/golang/protobuf v1.5.3
	github.com/sirupsen/logrus v1.9.0
	go.mongodb.org/mongo-driver v1.11.4
	google.golang.org/grpc v1.54.0
)

require (
	github.com/google/go-cmp v0.5.9 // indirect
	github.com/onsi/ginkgo v1.16.5 // indirect
	github.com/onsi/gomega v1.27.6 // indirect
	golang.org/x/net v0.8.0 // indirect
	golang.org/x/sys v0.6.0 // indirect
	golang.org/x/text v0.8.0 // indirect
	google.golang.org/genproto v0.0.0-20200423170343-7949de9c1215 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
