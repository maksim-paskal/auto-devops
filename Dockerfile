FROM alpine:latest

COPY ./auto-devops /app/auto-devops

ENTRYPOINT [ "/app/auto-devops" ]
