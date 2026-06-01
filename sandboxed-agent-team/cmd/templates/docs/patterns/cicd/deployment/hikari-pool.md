# HikariCP Keepalive

When deploying to a cloud environment, configure HikariCP keepalive and max-lifetime
so stale connections are detected and replaced before the application uses them —
cloud databases and network firewalls silently drop idle connections, and without
these settings the pool will hold dead connections until a request fails.

```properties
spring.datasource.hikari.keepalive-time=600000
spring.datasource.hikari.max-lifetime=1800000
```

`keepalive-time` (10 min) sends a test query to idle connections, evicting any that
have been dropped. `max-lifetime` (30 min) retires connections before the database's
own idle timeout would close them. Set `max-lifetime` to less than the database's
`wait_timeout` (PostgreSQL: `tcp_keepalives_idle`; RDS: typically 3600 s).