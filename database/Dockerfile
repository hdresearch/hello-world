FROM golang:1.23.5-alpine3.21 AS gateway-build
WORKDIR /src
COPY gateway/go.mod gateway/go.sum ./
RUN go mod download
COPY gateway/main.go ./
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/hello-world-gateway .

FROM postgres:16.6-alpine3.21
ENV POSTGRES_DB=webstack \
    POSTGRES_USER=postgres
COPY --from=gateway-build /out/hello-world-gateway /usr/local/bin/hello-world-gateway
COPY migrations /opt/hello-world/migrations
COPY entrypoint.sh /usr/local/bin/hello-world-entrypoint
RUN chmod +x /usr/local/bin/hello-world-entrypoint
EXPOSE 80
HEALTHCHECK --interval=10s --timeout=2s --retries=6 \
  CMD wget -qO- http://127.0.0.1/healthz >/dev/null || exit 1
ENTRYPOINT ["/usr/local/bin/hello-world-entrypoint"]
