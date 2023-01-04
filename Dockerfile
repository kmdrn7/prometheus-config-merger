##################################################################
###
##################################################################
FROM golang:1.19 as build

WORKDIR /go/src/prometheus-config-merger

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/prometheus-config-merger

##################################################################
###
##################################################################
FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/prometheus-config-merger /

ENTRYPOINT ["/prometheus-config-merger"]