all: build_controller build_ctl
	
build_controller:
	CGO_ENABLED=1 GOOS=linux GOGC=off go build -o deploy/controller/psi-controller cmd/controller/main.go
build_ctl:
	CGO_ENABLED=1 GOOS=linux GOGC=off go build -o deploy/controller/psictl cmd/psictl/main.go
