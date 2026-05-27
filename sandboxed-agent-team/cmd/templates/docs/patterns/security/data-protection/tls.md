# TLS in Production

When deploying the application, require TLS for all client-server communication
so credentials and session tokens cannot be intercepted in transit.

All client-server communication in production uses TLS (HTTPS). HTTP requests
are redirected to HTTPS. The application never serves content over plain HTTP in
production or staging.
