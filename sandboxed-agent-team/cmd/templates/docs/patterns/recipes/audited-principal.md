# Recipe: Audited Principal — `AuditedPrincipal`, `CurrentUser`, and JPA Auditing

The shared foundation that connects authentication to audit
columns. Every auth method (form login, passkey, OIDC) wraps its
own framework-supplied principal type, but they all expose two
things via a common interface — the **user key** (surrogate FK
into the user table) and the **username** (whatever
`Authentication.getName()` returns for that flow). With those two
methods the audit pipeline populates `@CreatedBy` /
`@LastModifiedBy` columns without a per-write database lookup, and
non-JPA code reads "who is logged in" from a single helper.

## What this produces

- An `AuditedPrincipal` interface exposing `getKey()` (used by
  audit) and `getUsername()` (used by UI / logging).
- A `CurrentUser` static helper returning the
  `Optional<AuditedPrincipal>` from the `SecurityContextHolder`.
  Callers ask for whatever property they need —
  `CurrentUser.get().map(AuditedPrincipal::getKey)`,
  `.map(AuditedPrincipal::getUsername)`, etc.
- An `AuditedEntity<K>` `@MappedSuperclass` carrying
  `@CreatedDate` / `@LastModifiedDate` and `@CreatedBy` /
  `@LastModifiedBy` FKs to your user entity.
- An `AuditorAware<UserEntity>` bean wired through
  `@EnableJpaAuditing` that pulls the key off the principal and
  resolves a `UserEntity` reference via
  `EntityManager.getReference()` — **no SELECT on every audited
  write.**

## Dependencies

- Spring Boot 3+ (Spring Data JPA's auditing API).
- Spring Security 6+ (`SecurityContextHolder`,
  `Authentication.getPrincipal()`).
- A `UserEntity` (or whatever you call it) with a surrogate key.

## Step 1 — Define the `AuditedPrincipal` interface

The contract every auth flow's principal must satisfy. Each
method's purpose is documented inline so adapters (form login,
OIDC, passkey) implement them correctly when they wrap their
framework's principal type.

```java
package {base_package}.security;

/**
 * Implemented by every auth flow's principal so audit and
 * "who is logged in?" queries have a single contract.
 *
 * <ul>
 *   <li>{@link #getKey()} — the user's surrogate key, captured at
 *       login time from whichever lookup the auth flow already
 *       performs. {@code AuditorAware} reads this so audit
 *       columns get a real FK without a DB roundtrip per write.
 *   <li>{@link #getUsername()} — whatever string
 *       {@code Authentication.getName()} returns for this flow:
 *       the form-login username column, the OIDC name claim, the
 *       WebAuthn credential user-entity name. The project picks
 *       what the username column actually holds (email, employee
 *       ID, login handle); this method is *not* an email accessor.
 * </ul>
 */
public interface AuditedPrincipal {

    Long getKey();

    String getUsername();
}
```

> **Note on the key type.** `Long` is the default surrogate-key
> shape this kit assumes. If your user entity uses `UUID` or a
> string key, change the return type — `AuditorAware` and
> `EntityManager.getReference()` accept any key type.

## Step 2 — Define the `CurrentUser` helper

A single static reader for "who is logged in?" — used by JPA
config, services, and UI code that doesn't want to import Spring
Security types directly. Returns the principal so callers compose
whichever properties they need; do **not** narrow this to one
property (e.g., `CurrentUserEmail`), because the moment a caller
needs a different property the helper has to grow a sibling.

```java
package {base_package}.security;

import java.util.Optional;
import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;

public final class CurrentUser {

    private CurrentUser() { }

    /**
     * The currently-authenticated {@link AuditedPrincipal}, or
     * {@link Optional#empty()} when no authenticated principal is
     * present (anonymous request, system-originated work like a
     * dev-seed migration, etc.).
     */
    public static Optional<AuditedPrincipal> get() {
        return Optional.ofNullable(SecurityContextHolder.getContext().getAuthentication())
                .filter(Authentication::isAuthenticated)
                .map(Authentication::getPrincipal)
                .filter(AuditedPrincipal.class::isInstance)
                .map(AuditedPrincipal.class::cast);
    }
}
```

> Why not `Authentication.getName()`-only? Some auth types return
> `Object`'s default `ClassName@hashHex` from `getPrincipal()`'s
> `toString()` if you reach for the principal carelessly; routing
> through the typed `AuditedPrincipal` interface eliminates the
> footgun. `getName()` works for one specific use case (the
> username); `CurrentUser.get()` works for all of them.

## Step 3 — Place the interface and helper alongside the user domain type

`AuditedPrincipal` and `CurrentUser` are about *identifying* the
authenticated user — they're not security-policy code (filter
chains, authorization rules) and they're not framework adapters
(`UserDetailsService` impls). The right seam is the package /
module that holds your **user domain type** — the one your
services accept and your UI displays. Group these three:

- `User` (or whatever your user value object is called).
- `AuditedPrincipal` (interface).
- `CurrentUser` (helper).

Putting them together means callers import one package to get the
"who is the user, and what can I ask about them?" surface. The
auth-method adapters (next recipes) import this seam, then add
their framework-specific wrapping.

If your project is multi-module, this seam often lives in a
`domain` / `uimodel` / `core` module. Single-module projects can
put it directly under `{base_package}.security` or
`{base_package}.user` — the rule is *together with the user
type*, not "in a `common` bucket."

## Step 4 — Add the `AuditedEntity<K>` mapped superclass

Entities that need audit columns extend this. The Spring Data
annotations (`@CreatedBy`, `@CreatedDate`, etc.) drive the
`AuditingEntityListener`, which calls into the `AuditorAware` bean
defined in Step 5.

```java
package {base_package}.entity;

import jakarta.persistence.Column;
import jakarta.persistence.EntityListeners;
import jakarta.persistence.FetchType;
import jakarta.persistence.JoinColumn;
import jakarta.persistence.ManyToOne;
import jakarta.persistence.MappedSuperclass;
import java.time.Instant;
import org.springframework.data.annotation.CreatedBy;
import org.springframework.data.annotation.CreatedDate;
import org.springframework.data.annotation.LastModifiedBy;
import org.springframework.data.annotation.LastModifiedDate;
import org.springframework.data.jpa.domain.support.AuditingEntityListener;

@MappedSuperclass
@EntityListeners(AuditingEntityListener.class)
public abstract class AuditedEntity<K> {

    @CreatedDate
    @Column(name = "created_at", nullable = false, updatable = false)
    private Instant createdAt;

    @LastModifiedDate
    @Column(name = "updated_at", nullable = false)
    private Instant updatedAt;

    @CreatedBy
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "created_by_key", updatable = false)
    private UserEntity createdBy;

    @LastModifiedBy
    @ManyToOne(fetch = FetchType.LAZY)
    @JoinColumn(name = "updated_by_key")
    private UserEntity updatedBy;

    // getters omitted — Lombok @Getter or hand-written
}
```

The `@JoinColumn` names (`created_by_key`, `updated_by_key`) are
the column names on every audited table — keep them consistent so
your DDL and reporting queries can rely on the convention.

## Step 5 — Wire `AuditorAware` with the FK-only optimisation

The `AuditorAware<UserEntity>` bean is the seam between Spring
Security and JPA auditing. The trick is to avoid loading the full
`UserEntity` on every audited write — `EntityManager.getReference()`
returns a Hibernate proxy holding just the key, which JPA then
persists into the FK column directly.

```java
package {base_package}.config;

import jakarta.persistence.EntityManager;
import java.util.Optional;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.domain.AuditorAware;
import org.springframework.data.jpa.repository.config.EnableJpaAuditing;
import {base_package}.entity.UserEntity;
import {base_package}.security.AuditedPrincipal;
import {base_package}.security.CurrentUser;

@Configuration
@EnableJpaAuditing(auditorAwareRef = "auditorAware")
public class JpaConfig {

    /**
     * Resolves the current auditor as a {@link UserEntity} reference
     * so audit columns are real FKs to {@code usr.key}. Pulls the key
     * off the principal — every auth flow's principal implements
     * {@link AuditedPrincipal} and carries the key from its login-time
     * validation step. {@link EntityManager#getReference} returns a
     * Hibernate proxy holding just that key, so JPA persists the FK
     * without issuing a SELECT. Empty when there's no authenticated
     * principal (system-originated writes like dev-seed inserts) — the
     * audit columns are nullable, so {@code NULL} is preserved.
     */
    @Bean
    public AuditorAware<UserEntity> auditorAware(EntityManager entityManager) {
        return () -> CurrentUser.get()
                .map(AuditedPrincipal::getKey)
                .map(key -> entityManager.getReference(UserEntity.class, key));
    }
}
```

This bean reads through `CurrentUser.get()` (Step 2) so the
SecurityContext-extraction logic stays in one place. If you ever
need to handle a special case (e.g., a "system" principal for
background jobs), implement that in `CurrentUser` and the audit
pipeline picks it up automatically.

## Step 6 — Implement `AuditedPrincipal` in each auth flow's principal

The auth-method recipes ([form-login.md](form-login.md),
[passkey.md](passkey.md), [oidc-sso.md](oidc-sso.md)) each define a
principal class that wraps the framework's auth type
(`UserDetails`, `OidcUser`, `PublicKeyCredentialUserEntity`). Each
of those principals implements `AuditedPrincipal` and populates the
key from whatever lookup the flow already performs at login.

Sketch — the recipes have the full implementations:

```java
public class AuditedFormLoginUser implements UserDetails, AuditedPrincipal {
    private final Long key;
    private final String username;
    // password (transient), authorities, etc.

    @Override public Long getKey()     { return key; }
    @Override public String getUsername() { return username; }
}
```

## Decisions this recipe imposes

- **Audit columns persist the user key, not the username.** The
  username is a mutable identity field; the key is stable. Spring
  Data JPA's auditing docs are silent on this; getting it wrong
  couples your audit history to identity changes.
- **`CurrentUser` returns the principal, not a property.** Adding
  one helper per property (`CurrentUserEmail`, `CurrentUserId`,
  `CurrentUserDisplayName`, …) creates a maintenance trap.
  Returning the principal lets callers compose what they need.
- **`Authentication.getName()` is the username, not the email.**
  Some projects use email as the username column; others use
  employee ID, login handle, or an opaque identifier. The
  `AuditedPrincipal.getUsername()` contract reflects what Spring
  Security calls it: a *username*. Project-specific column choice
  belongs in `docs/solutions/`.
- **`EntityManager.getReference()` over `findById()` in
  `AuditorAware`.** A SELECT on every audited write is invisible
  cost; the proxy is the same FK persistence with none of the
  load.
- **Place the interface and helper alongside the user domain
  type.** Not in a generic `common` bucket. The auth-method
  adapters import the user-domain seam to wrap their principals;
  scattering the contract makes that import path opaque.

## What to verify

- An audited entity persists with `created_by_key` populated to
  the logged-in user's `usr.key`. `SELECT created_by_key FROM
  audited_table` returns real FKs, not null, after a write
  performed under an authenticated session.
- A `SELECT` count on the audit-write transaction shows **no
  extra `SELECT` against `usr`** — the `EntityManager.getReference()`
  proxy avoided the load. Use `spring.jpa.show-sql=true` and
  watch the log during a write.
- A system-originated write (e.g., a Flyway/dev-seed migration
  running outside an authenticated context) produces `NULL` audit
  FKs and does not throw — `CurrentUser.get()` returns empty,
  `AuditorAware` returns empty, JPA persists `NULL`.
- `CurrentUser.get().map(AuditedPrincipal::getUsername)` returns
  the same string as `Authentication.getName()` for the current
  request — the contract holds across all auth flows.

## Related

- [conditional-auth.md](conditional-auth.md) — the configuration
  layer that decides which auth methods are active. The
  audited-principal contract is what each method's adapter
  satisfies.
- [form-login.md](form-login.md) — the form-login adapter that
  produces an `AuditedFormLoginUser`.
- [passkey.md](passkey.md) — the passkey adapter that produces an
  `AuditedPasskeyPrincipal`.
- [oidc-sso.md](oidc-sso.md) — the OIDC adapter that produces an
  `AuditedOidcUser`.
- `docs/patterns/architecture/persistence.md` — surrounding JPA
  conventions (`@MappedSuperclass` hierarchy, lazy fetch defaults).
- `docs/solutions/` — your project's specific user-key type
  (`Long` / `UUID` / `String`), username column choice (email /
  employee ID / handle), and any system-principal handling.
