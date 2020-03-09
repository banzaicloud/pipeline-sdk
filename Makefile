OS = $(shell uname | tr A-Z a-z)

PIPELINE_VERSION = process-log
PROTOC_VERSION = 3.11.4
OPENAPI_GENERATOR_VERSION = v4.2.3
OPENAPI_DESCRIPTOR = apis/pipeline/pipeline.yaml

bin/protoc: bin/protoc-${PROTOC_VERSION}
	@ln -sf protoc-${PROTOC_VERSION}/bin/protoc bin/protoc
bin/protoc-${PROTOC_VERSION}:
	@mkdir -p bin/protoc-${PROTOC_VERSION}
ifeq ($(OS), darwin)
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-osx-x86_64.zip > bin/protoc.zip
endif
ifeq ($(OS), linux)
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip > bin/protoc.zip
endif
	unzip bin/protoc.zip -d bin/protoc-${PROTOC_VERSION}
	rm bin/protoc.zip

bin/protoc-gen-go: go.mod
	@mkdir -p bin
	go build -o bin/protoc-gen-go github.com/golang/protobuf/protoc-gen-go

apis/pipeline/api.proto:
	@mkdir -p apis/pipeline
	curl https://raw.githubusercontent.com/banzaicloud/pipeline/${PIPELINE_VERSION}/apis/pipeline/api.proto > apis/pipeline/api.proto

.PHONY: _download-protos
_download-protos: apis/pipeline/api.proto

.PHONY: proto
proto: bin/protoc bin/protoc-gen-go _download-protos ## Generate client and server stubs from the protobuf definition
	mkdir -p .gen/pipeline
	bin/protoc -I bin/protoc-${PROTOC_VERSION} -I apis/pipeline --go_out=plugins=grpc,import_path=pipeline:.gen/pipeline $(shell find apis/pipeline -name '*.proto')

apis/pipeline/pipeline.yaml:
	@mkdir -p apis/pipeline
	curl https://raw.githubusercontent.com/banzaicloud/pipeline/${PIPELINE_VERSION}/apis/pipeline/pipeline.yaml > ${OPENAPI_DESCRIPTOR}

.PHONY: _download-openapis
_download-openapis: apis/pipeline/pipeline.yaml

.PHONY: validate-openapi
validate-openapi: _download-openapis ## Validate the openapi description
	docker run --rm -v $${PWD}:/local openapitools/openapi-generator-cli:${OPENAPI_GENERATOR_VERSION} validate --recommend -i /local/${OPENAPI_DESCRIPTOR}

.PHONY: generate-openapi
generate-openapi: validate-openapi ## Generate go server based on openapi description
	@ if [[ "$$OSTYPE" == "linux-gnu" ]]; then sudo rm -rf ./.gen/pipeline; else rm -rf ./.gen/pipeline/; fi
	docker run --rm -v $${PWD}:/local openapitools/openapi-generator-cli:${OPENAPI_GENERATOR_VERSION} generate \
	--additional-properties packageName=pipeline \
	--additional-properties withGoCodegenComment=true \
	-i /local/${OPENAPI_DESCRIPTOR} \
	-g go \
	-o /local/.gen/pipeline/pipeline
	@ if [[ "$$OSTYPE" == "linux-gnu" ]]; then sudo chown -R $(shell id -u):$(shell id -g) .gen/pipeline/; fi
	rm -r .gen/pipeline/pipeline/{docs,go.*}
