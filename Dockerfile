# syntax=docker/dockerfile:1
FROM golang:1.21-alpine as build

WORKDIR /app

# for Makefile
RUN apk add build-base
# for golang wire
RUN go install github.com/google/wire/cmd/wire@latest
# for golang swaggo
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY . .
RUN make generator
RUN go build -mod vendor -o ./.kanthor/kanthor -buildvcs=false cmd/server/main.go
RUN go build -mod vendor -o ./.kanthor/kanthorhealthz -buildvcs=false cmd/healthz/main.go
RUN go build -mod vendor -o ./.kanthor/kanthordata -buildvcs=false cmd/data/main.go

FROM alpine:3
WORKDIR /app

COPY --from=build /app/data ./data
COPY --from=build /app/configs.yaml ./configs.yaml
COPY --from=build /app/.kanthor/kanthor /usr/bin/kanthor
COPY --from=build /app/.kanthor/kanthorhealthz /usr/bin/kanthorhealthz
COPY --from=build /app/.kanthor/kanthordata /usr/bin/kanthordata

# portal
EXPOSE 8180
# sdk
EXPOSE 8280