.PHONY: all dep test lint build package

WORKSPACE ?= $$(pwd)

GO_PKG_LIST := $(shell go list ./... | grep -v /vendor/ | grep -v /mock)

export GOFLAGS := -mod=vendor

all: clean package
	@echo "Done"

clean:
	@rm -rf ./bin/
	@mkdir -p ./bin
	@echo "Clean complete"

resolve-dependencies:
	@echo "Resolving go package dependencies"
	@go mod tidy
	@go mod vendor
	@echo "Package dependencies completed"

dep: resolve-dependencies

dep-check:
	@go mod verify

dep-version:
	@export sdk_version=$(sdk) && export da_version=$(da) && make update-sdk && make dep

dep-sdk:
	@make sdk=master da=master dep-version

update-sdk:
	@echo "Updating SDK dependencies"
	@export GOFLAGS="" && go get "git.ecd.axway.org/apigov/apic_agents_sdk@${sdk_version}" "git.ecd.axway.org/apigov/v7_discovery_agent@${da_version}"

test:
	@go vet ${GO_PKG_LIST}
	@go test -short -coverprofile=${WORKSPACE}/gocoverage.out -count=1 ${GO_PKG_LIST}

test-sonar:
	@go vet ${GO_PKG_LIST}
	@go test -short -coverpkg=./... -coverprofile=${WORKSPACE}/gocoverage.out -count=1 ${GO_PKG_LIST} -json > ${WORKSPACE}/goreport.json

error-check:
	./build/scripts/error_check.sh ./pkg ./vendor/git.ecd.axway.org/apigov/apic_agents_sdk ./vendor/git.ecd.axway.org/apigov/v7_discovery_agent

sonar: test-sonar
	sonar-scanner -X \
		-Dsonar.host.url=http://quality1.ecd.axway.int \
		-Dsonar.language=go \
		-Dsonar.projectName=V7_TraceabilityAgent \
		-Dsonar.projectVersion=1.0 \
		-Dsonar.projectKey=V7_TraceabilityAgent \
		-Dsonar.sourceEncoding=UTF-8 \
		-Dsonar.projectBaseDir=${WORKSPACE} \
		-Dsonar.sources=. \
		-Dsonar.tests=. \
		-Dsonar.exclusions=**/mock/**,**/vendor/** \
		-Dsonar.test.inclusions=**/*test*.go \
		-Dsonar.go.tests.reportPaths=goreport.json \
		-Dsonar.go.coverage.reportPaths=gocoverage.out

lint:
	@golint -set_exit_status ${GO_PKG_LIST}

${WORKSPACE}/traceability_agent:
	@export time=`date +%Y%m%d%H%M%S` && \
	export version=`cat version` && \
	export commit_id=`cat commit_id` && \
  export GOOS=linux && \
	export CGO_ENABLED=0 && \
	export GOARCH=amd64 && \
	go build -v -tags static_all \
		-ldflags="-X 'git.ecd.axway.org/apigov/apic_agents_sdk/pkg/cmd.BuildTime=$${time}' \
				-X 'git.ecd.axway.org/apigov/apic_agents_sdk/pkg/cmd.BuildVersion=$${version}' \
				-X 'git.ecd.axway.org/apigov/apic_agents_sdk/pkg/cmd.BuildCommitSha=$${commit_id}' \
				-X 'git.ecd.axway.org/apigov/apic_agents_sdk/pkg/cmd.BuildAgentName=SpringTraceabilityAgent' \
							-X 'git.ecd.axway.org/apigov/apic_agents_sdk/pkg/cmd/service.Name=traceability-agent' \
							-X 'git.ecd.axway.org/apigov/apic_agents_sdk/pkg/cmd/service.Description=Spring Traceability Agent'" \
		-a -o ${WORKSPACE}/bin/traceability_agent ${WORKSPACE}/traceability_agent.go

build: ${WORKSPACE}/traceability_agent
	@echo "Build completed"

docker:
	docker build -t traceability_agent:latest -f ${WORKSPACE}/build/docker/Dockerfile .
	@echo "Docker build complete"