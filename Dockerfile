FROM golang:alpine as builder

RUN apk update && apk add --no-cache git ca-certificates && update-ca-certificates

ENV USER=appuser
ENV UID=10001
RUN adduser \    
    --disabled-password \    
    --gecos "" \    
    --home "/nonexistent" \    
    --shell "/sbin/nologin" \    
    --no-create-home \    
    --uid "${UID}" \    
    "${USER}"

WORKDIR /
COPY . .
RUN mkdir /data
RUN chown -R appuser:appuser /data

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /rescope

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/group /etc/group
COPY --from=builder /rescope /rescope
COPY --from=builder --chown=appuser:appuser /configs /configs
COPY --from=builder --chown=appuser:appuser /data /data

USER appuser:appuser
WORKDIR /data
ENTRYPOINT ["/rescope"]
