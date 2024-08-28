#
# COPYRIGHT Ericsson 2024
#
#
#
# The copyright to the computer program(s) herein is the property of
#
# Ericsson Inc. The programs may be used and/or copied only with written
#
# permission from Ericsson Inc. or in accordance with the terms and
#
# conditions stipulated in the agreement/contract under which the
#
# program(s) have been supplied.
#

ARG BUILD_IMAGE_TAG="1.20"
ARG BUILD_IMAGE_VERSION=""

FROM armdocker.rnd.ericsson.se/proj-am/sles/sles-golang:${BUILD_IMAGE_TAG} AS build

WORKDIR /src

COPY . .

RUN go mod download && go build -o main cmd/helm-registry-adapter/main.go

CMD ["./main"]
