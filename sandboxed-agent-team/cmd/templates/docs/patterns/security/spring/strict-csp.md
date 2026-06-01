# Strict CSP (Nonce-Based)

When a strict Content Security Policy beyond Vaadin's default relaxations is required, use Vaadin's nonce-based strict CSP mechanism — not Spring Security's `headers().contentSecurityPolicy(...)`.

Vaadin's [nonce-based strict CSP](https://vaadin.com/docs/latest/flow/security/advanced-topics/strict-csp)
requires:

- An `IndexHtmlRequestListener` that generates a random nonce per request and
  sets the `Content-Security-Policy` header with `script-src 'nonce-...'`.
- Nonce injection into script tags in the index file.
- JavaScript overrides that replace `Function` / `eval()` with CSP-compliant
  versions.

Production-mode only with prerequisites — follow the linked Vaadin docs rather
than reproducing the recipe here.
