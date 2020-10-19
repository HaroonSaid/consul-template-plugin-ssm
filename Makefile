# Other config
NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m

DIR_OUT=$(CURDIR)/dist

.PHONY: all clean deps build

all: clean deps build

deps:
	@echo "$(OK_COLOR)==> Installing AWS SDK dependencies$(NO_COLOR)"
	@go get -u github.com/aws/aws-sdk-go/aws/session
	@go get -u github.com/aws/aws-sdk-go/service/ssm

# Builds the project
build:
	@echo "$(OK_COLOR)==> Building... $(NO_COLOR)"
	/bin/sh -c "VERSION=${VERSION} ./build/build.sh"

# Run unit-tests
test:
	@/bin/sh -c "./build/test.sh"

# Cleans our project: deletes binaries
clean:
	@echo "$(OK_COLOR)==> Cleaning project$(NO_COLOR)"
	if [ -d ${DIR_OUT} ] ; then rm -rf ${DIR_OUT}/* ; fi
