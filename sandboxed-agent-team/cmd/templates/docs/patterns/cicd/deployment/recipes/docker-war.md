# Docker Image — WAR

When containerizing a WAR deployment, build from a Tomcat base image and copy the
WAR into `webapps/` so the container manages the servlet lifecycle.

Build the production WAR first (see `docs/patterns/cicd/deployment/recipes/war.md`):

```bash
mvn clean package -Pproduction
```

```dockerfile
FROM tomcat:10-jre21-temurin-alpine
RUN rm -rf /usr/local/tomcat/webapps/ROOT
COPY {app}-app/target/{app}-app.war /usr/local/tomcat/webapps/ROOT.war
EXPOSE 8080
CMD ["catalina.sh", "run"]
```

Deploying as `ROOT.war` serves the application at the root context path (`/`). To
serve at a sub-path, name the WAR after the desired path segment instead.

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

- `docs/patterns/cicd/deployment/recipes/war.md` — produces the WAR this image packages.
- `docs/patterns/cicd/deployment/recipes/docker-fat-jar.md` — alternative: fat JAR without an external servlet container.