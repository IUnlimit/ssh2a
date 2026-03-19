FROM debian:bookworm-slim

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY output/ssh2a_linux /app/ssh2a

RUN chmod +x /app/ssh2a

EXPOSE 9080 9022

ENTRYPOINT ["/app/ssh2a"]
