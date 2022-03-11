# Sturdy + Docker Compose + TLS = ❤️

To test Sturdy with TLS locally, you can use [mkcert](https://github.com/FiloSottile/mkcert) to create certificates for development: `mkcert -install && mkcert mysturdy.example.com localhost 127.0.0.1`

1. Get your certificates
2. Update `docker-compose.yaml` to fit your needs (paths to certificates, ports, etc)
3. `docker compose up`
4. Sturdy is now running with TLS on port 443!