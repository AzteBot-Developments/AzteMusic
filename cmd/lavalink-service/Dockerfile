# syntax=docker/dockerfile:1

FROM eclipse-temurin:17-jdk-jammy

WORKDIR /app

COPY internal/lavalink ./

EXPOSE 2333

CMD ["java", "-jar", "Lavalink.jar"]