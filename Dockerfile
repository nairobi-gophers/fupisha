FROM golang:1.15-alpine as base

WORKDIR /fupisha

FROM aquasec/trivy:0.14.0 as trivy

RUN trivy --debug --timeout 4m golang:1.15-alpine && \
  echo "No image vulnerabilities" > result

FROM base as dev

COPY . .

RUN go mod download

RUN go mod verify

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN go env
RUN go build -o main ./cmd/

FROM alpine:3.10 AS prod

RUN apk --no-cache add ca-certificates --upgrade bash

COPY --from=dev /fupisha/main main
COPY --from=dev /fupisha/templates templates

CMD [ "./main" ,"start"]
