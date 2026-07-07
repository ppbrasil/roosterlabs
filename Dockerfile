# Build
FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# Runtime — imagem mínima, binário estático.
# O Lambda Web Adapter entra como extensão: dentro do Lambda ele traduz
# invocações em HTTP para o app; fora do Lambda é ignorado. Mesmo container
# roda em Lambda hoje e em ECS/Fargate (VPC própria) amanhã.
# Variante root (não :nonroot): o Lambda executa a imagem com seu próprio
# usuário sandbox e falha com "permission denied" no entrypoint quando a
# imagem usa o nonroot da distroless (workdir /home/nonroot inacessível).
# Ver https://github.com/ko-build/ko/issues/669. Revisitar :nonroot na
# migração futura para ECS/Fargate.
FROM gcr.io/distroless/static-debian12:latest
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.9.1 /lambda-adapter /opt/extensions/lambda-adapter
# --chmod garante bit de execução para o usuário não-root do Lambda;
# sem isso o modo do arquivo depende do umask do ambiente de build.
COPY --from=build --chmod=0755 /out/server /server
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/server"]
