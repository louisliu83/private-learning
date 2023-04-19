all: build_controller build_ctl build_tools
	
build_controller:
	CGO_ENABLED=1 GOOS=linux GOGC=off go build -o deploy/psi-controller cmd/controller/main.go
build_ctl:
	CGO_ENABLED=1 GOOS=linux GOGC=off go build -o deploy/psictl cmd/tools/ctl/psictl.go
build_tools:
	CGO_ENABLED=1 GOOS=linux GOGC=off go build -o deploy/tokengen cmd/tools/register/token_gen.go
