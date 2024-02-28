FROM eclipse-temurin:21-jdk-jammy

WORKDIR /app

COPY internal/lavalink ./
COPY internal/lavalink/application.yml ./application.yml

EXPOSE 2333

CMD ["java", "-jar", "Lavalink.jar"]