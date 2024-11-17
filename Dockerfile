FROM golang:1.23-alpine as builder

WORKDIR /go/src/app

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o cscheck cmd/cscheck/main.go

FROM 3ayazaya/agscript:4.9

COPY --from=builder /go/src/app/cscheck /opt/cobaltstrike
COPY ./scripts /opt/cobaltstrike/

ENTRYPOINT [ "./cscheck" ]
