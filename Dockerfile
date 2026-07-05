# Build
FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

# Runtime — imagem mínima, binário estático.
# O Lambda Web Adapter entra como extensão: dentro do Lambda ele traduz
# invocações em HTTP para o app; fora do Lambda é ignorado. Mesmo container
# roda em Lambda hoje e em ECS/Fargate (VPC própria) amanhã.
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=public.ecr.aws/awsguru/aws-lambda-adapter:0.9.1 /lambda-adapter /opt/extensions/lambda-adapter
COPY --from=build /out/server /server
ENV PORT=8080
EXPOSE 8080
ENTRYPOINT ["/server"]
