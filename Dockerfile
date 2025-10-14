FROM golang:1.25.2-alpine3.22 as build

WORKDIR /sapopinguino-translate

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -tags prod -ldflags="-s -w" -o ./bin/main ./main.go 

FROM alpine:3.22

WORKDIR /sapopinguino-translate

COPY --from=build /sapopinguino-translate/bin/main ./main

COPY --from=build /sapopinguino-translate/assets ./assets

COPY --from=build /sapopinguino-translate/config/config.prod.yml ./config/config.prod.yml

ENTRYPOINT [ "./main" ]

