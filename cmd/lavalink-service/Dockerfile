# syntax=docker/dockerfile:1

FROM eclipse-temurin:17-jdk-jammy

WORKDIR /app

COPY internal/lavalink ./
COPY internal/lavalink/application.yml ./application.yml

EXPOSE 2333

CMD ["java", "-jar", "Lavalink.jar"]