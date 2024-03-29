module github.com/icinga/icingadb/tests

go 1.16

replace github.com/icinga/icingadb => ../

require (
	github.com/containerd/containerd v1.5.6 // indirect
	github.com/go-redis/redis/v8 v8.11.4
	github.com/google/uuid v1.3.0
	github.com/icinga/icinga-testing v0.0.0-20220315141514-b316c93f8d77
	github.com/icinga/icingadb v1.0.0-rc2
	github.com/jmoiron/sqlx v1.3.4
	github.com/pkg/errors v0.9.1
	github.com/stretchr/testify v1.7.0
	go.uber.org/zap v1.21.0
	golang.org/x/net v0.0.0-20211020060615-d418f374d309 // indirect
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20211025201205-69cdffdb9359 // indirect
	google.golang.org/genproto v0.0.0-20211027162914-98a5263abeca // indirect
	google.golang.org/grpc v1.41.0 // indirect
)
