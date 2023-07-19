FROM golang:1.20-bookworm as builder

WORKDIR /go/src
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY components/notification/ /go/src/components/notification/
WORKDIR /go/src/components/notification/cmd
RUN go build -o api -tags osusergo,netgo

FROM gcr.io/distroless/static-debian11:nonroot

COPY --from=builder --chown=nonroot:nonroot /go/src/components/notification/cmd/api /bin/api

ENTRYPOINT [ "/bin/api" ]