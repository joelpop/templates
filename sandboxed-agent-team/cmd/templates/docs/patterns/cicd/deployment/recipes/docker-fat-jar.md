# Docker Image — Fat JAR

When containerizing a fat-JAR deployment, build from an Eclipse Temurin JRE base
image and inject secrets via environment variables so no credentials are baked into
the image.

Build the production JAR first (see `docs/patterns/cicd/deployment/recipes/fat-jar.md`):

```bash
mvn clean package -Pproduction
```

```dockerfile
FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY {app}-app/target/{app}-app.jar app.jar
EXPOSE 8080
ENTRYPOINT ["java", "-jar", "app.jar"]
```

In production, supply secrets via `--env-file` or orchestrator secrets (Docker Swarm
`--secret`, Kubernetes `Secret`) rather than inline `-e` flags:

```bash
docker build -t {app}:latest .
docker run -p 8080:8080 \
  -e DB_URL=jdbc:postgresql://host/db \
  -e DB_USER=appuser \
  -e DB_PASSWORD=secret \
  {app}:latest
```

## Related

- `docs/patterns/cicd/deployment/recipes/fat-jar.md` — produces the JAR this image packages.
- `docs/patterns/cicd/deployment/recipes/docker-war.md` — alternative: WAR deployed into a containerized servlet container.