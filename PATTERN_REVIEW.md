# Pattern File Review Checklist

## Procedure

Work through files one at a time in checklist order, skipping `[x]` items:

1. Read the file and give the user the full workspace-relative path (unformatted, not a markdown link — IJ uses it to open the file directly).
2. Give a brief assessment: what's good, what's missing, what's wrong.
3. Wait for user feedback. Make edits directly so the user sees the diff in IntelliJ — don't show a markdown preview in chat. For structural proposals (splits, deletes, large rewrites), describe the approach and wait for approval first.
4. Mark `[x]` only when the user signals done (e.g., "next", "done", "looks good").

Additional rules in memory: "Avoid" first / "Preferred" second in *separate* code blocks; US English spellings; keep ` ```sql ` language tags; "When … use/do … so …" obligation form for INDEX.md entries.

## cicd/
- [x] `cicd/deployment/recipes/fat-jar.md`
- [x] `cicd/deployment/recipes/war.md`
- [x] `cicd/deployment/recipes/docker-fat-jar.md`
- [x] `cicd/deployment/recipes/docker-war.md`
- [x] `cicd/deployment/spring-profiles.md`
- [x] `cicd/deployment/hikari-pool.md`
- [x] `cicd/flyway/versioned-scripts.md`
- [x] `cicd/flyway/idempotent-scripts.md`
- [x] `cicd/flyway/ddl-strategy.md`
- [x] `cicd/flyway/seed-data.md`
- [x] `cicd/flyway/rollback-docs.md`
- [x] `cicd/logging-framework.md`
- [x] `cicd/structured-logging.md`
- [x] `cicd/startup-logging.md`

## design/
### design/figma/recipes/
- [ ] `design/figma/recipes/figma-code-comments.md`
- [ ] `design/figma/recipes/figma-component-mapping.md`
- [ ] `design/figma/recipes/figma-lumo-mapping.md`
- [ ] `design/figma/recipes/figma-quality-standards.md`
- [ ] `design/figma/recipes/figma-requirements.md`
- [ ] `design/figma/recipes/figma-to-lumo-theme.md`
- [ ] `design/figma/recipes/figma-to-vaadin.md`
- [ ] `design/figma/recipes/lumo-styles-location.md`
- [ ] `design/figma/recipes/lumo-value-conversion.md`

## language/
- [x] `language/singular-form.md`

### language/java/
- [x] `language/java/abstraction/base-class.md`
- [x] `language/java/abstraction/component-package.md`
- [x] `language/java/abstraction/placement.md`
- [x] `language/java/abstraction/third-instance-rule.md`
- [x] `language/java/abstraction/value-objects.md`
- [x] `language/java/abstraction/version-isolation.md`
- [x] `language/java/comments.md`
- [x] `language/java/dependency-injection.md`
- [x] `language/java/post-construct.md`
- [x] `language/java/dont-configure-consumer.md`
- [x] `language/java/fix-context-placement.md`
- [x] `language/java/fix-mode-trap.md`
- [x] `language/java/javadoc.md`
- [x] `language/java/lombok/enum.md`
- [x] `language/java/lombok/logging.md`
- [x] `language/java/lombok/pojo.md`
- [x] `language/java/nested-types.md`
- [x] `language/java/no-task-references.md`
- [x] `language/java/no-todos.md`
- [x] `language/java/suppress-warnings.md`

### language/java/lambda/
- [x] `language/java/lambda/lambdas.md`
- [x] `language/java/lambda/null-checks.md`
- [x] `language/java/lambda/unused-lambda-params.md`

### language/java/naming/
- [x] `language/java/naming/event-handler-naming.md`
- [x] `language/java/naming/naming.md`

### language/java/variable/
- [x] `language/java/variable/initialization.md`
- [x] `language/java/variable/local-variables.md`
- [x] `language/java/variable/type-inference.md`

## persistence/
### persistence/spring-data-jpa/
- [x] `persistence/spring-data-jpa/artifact-naming.md`
- [x] `persistence/spring-data-jpa/association.md`
- [x] `persistence/spring-data-jpa/attribute-converter.md`
- [x] `persistence/spring-data-jpa/audit-projection.md`
- [x] `persistence/spring-data-jpa/batch.md`
- [x] `persistence/spring-data-jpa/bulk.md`
- [x] `persistence/spring-data-jpa/composition.md`
- [x] `persistence/spring-data-jpa/dynamic-update.md`
- [ ] `persistence/spring-data-jpa/embeddable.md`
- [x] `persistence/spring-data-jpa/entity-methods.md`
- [ ] `persistence/spring-data-jpa/entity-hierarchy/audited-entity.md`
- [ ] `persistence/spring-data-jpa/entity-hierarchy/overview.md`
- [ ] `persistence/spring-data-jpa/entity-hierarchy/root-entity.md`
- [ ] `persistence/spring-data-jpa/entity-hierarchy/temporal-types.md`
- [ ] `persistence/spring-data-jpa/entity-hierarchy/versioned-entity.md`
- [ ] `persistence/spring-data-jpa/entity-lombok.md`
- [x] `persistence/spring-data-jpa/enum-column-naming.md`
- [ ] `persistence/spring-data-jpa/enum-mapping.md`
- [ ] `persistence/spring-data-jpa/insert-update.md`
- [x] `persistence/spring-data-jpa/jpa-scan-config.md`
- [x] `persistence/spring-data-jpa/key-id-naming.md`
- [x] `persistence/spring-data-jpa/many-to-many.md`
- [ ] `persistence/spring-data-jpa/many-to-one-lazy.md`
- [ ] `persistence/spring-data-jpa/n-plus-one.md`
- [ ] `persistence/spring-data-jpa/named-parameters.md`
- [ ] `persistence/spring-data-jpa/nested-projections.md`
- [x] `persistence/spring-data-jpa/osiv-disabled.md`
- [ ] `persistence/spring-data-jpa/picker-projections.md`
- [ ] `persistence/spring-data-jpa/projection-inheritance.md`
- [x] `persistence/spring-data-jpa/projection-naming.md`
- [ ] `persistence/spring-data-jpa/projections.md`

## security/
### security/authentication/
- [x] `security/authentication/account-lockout-notification.md`
- [ ] `security/authentication/bcrypt.md`
- [x] `security/authentication/login-failure-messages.md`
- [ ] `security/authentication/password-strength.md`
- [ ] `security/authentication/rate-limiting.md`

### security/authorization/
- [ ] `security/authorization/layout-access.md`
- [ ] `security/authorization/role-based-rendering.md`
- [ ] `security/authorization/role-constants.md`
- [ ] `security/authorization/role-enum-placement.md`
- [ ] `security/authorization/role-security-name.md`
- [ ] `security/authorization/self-editing.md`
- [ ] `security/authorization/view-access.md`

### security/data-protection/
- [ ] `security/data-protection/dependency-scanning.md`
- [ ] `security/data-protection/error-leakage.md`
- [ ] `security/data-protection/file-upload.md`
- [ ] `security/data-protection/pii-logging.md`
- [ ] `security/data-protection/response-header-disclosure.md`
- [ ] `security/data-protection/secrets.md`
- [ ] `security/data-protection/tls.md`
- [ ] `security/data-protection/xss.md`

### security/form-login/recipes/
- [ ] `security/form-login/recipes/conditional-auth.md`
- [ ] `security/form-login/recipes/form-login.md`

### security/oidc/recipes/
- [ ] `security/oidc/recipes/oidc-sso.md`

### security/passkey/recipes/
- [ ] `security/passkey/recipes/passkey.md`

### security/spring/
- [ ] `security/spring/csrf-non-vaadin.md`
- [ ] `security/spring/security-headers.md`
- [x] `security/spring/session-cookie-flags.md`
- [x] `security/spring/session-timeout.md`
- [ ] `security/spring/strict-csp.md`
- [ ] `security/spring/vaadin-web-security.md`

### security/spring/recipes/
- [ ] `security/spring/recipes/audited-principal.md`

## structure/
### structure/modules/
- [ ] `structure/modules/allowed-packages.md`
- [x] `structure/modules/app-entry-point.md`
- [ ] `structure/modules/jpa-config.md`
- [ ] `structure/modules/layer-separation.md`
- [ ] `structure/modules/mapstruct-config.md`
- [ ] `structure/modules/overview.md`
- [ ] `structure/modules/package-naming.md`
- [ ] `structure/modules/version-management.md`

### structure/services/
- [ ] `structure/services/service-caching.md`
- [ ] `structure/services/service-errors.md`
- [x] `structure/services/service-grid-loading.md`
- [ ] `structure/services/service-insert.md`
- [ ] `structure/services/service-interface.md`
- [ ] `structure/services/service-mapper.md`
- [ ] `structure/services/service-method-naming.md`
- [ ] `structure/services/service-naming.md`
- [ ] `structure/services/service-osiv.md`
- [ ] `structure/services/service-transactional.md`
- [ ] `structure/services/service-update.md`

## testing/
- [ ] `testing/acceptance-tracing.md`
- [ ] `testing/ac-coverage-gaps.md`
- [ ] `testing/ac-cross-reference.md`
- [ ] `testing/browserless-limitations.md`
- [ ] `testing/browserless-ui.md`
- [ ] `testing/coverage-targets.md`
- [ ] `testing/data-jpa-test.md`
- [ ] `testing/h2-compat-mode.md`
- [ ] `testing/n-plus-one-detection.md`
- [ ] `testing/pyramid.md`
- [ ] `testing/test-naming.md`
- [ ] `testing/test-per-class.md`
- [ ] `testing/tests-as-documentation.md`
- [ ] `testing/transactional-rollback.md`
- [ ] `testing/unit-structure.md`

### testing/page-objects/
- [ ] `testing/page-objects/browserless.md`
- [ ] `testing/page-objects/playwright.md`
- [ ] `testing/page-objects/testbench.md`

### testing/recipes/
- [ ] `testing/recipes/testbench-browserless.md`
- [ ] `testing/recipes/testbench-e2e-parallel.md`
- [ ] `testing/recipes/testbench-e2e-server.md`

## ui/
### ui/vaadin/
- [x] `ui/vaadin/allowed-packages-config.md`
- [ ] `ui/vaadin/async-operation-notification.md`
- [x] `ui/vaadin/component-add.md`
- [x] `ui/vaadin/avatar.md`
- [x] `ui/vaadin/binder.md`
- [x] `ui/vaadin/button-loading-state.md`
- [x] `ui/vaadin/client-details-impl.md`
- [x] `ui/vaadin/client-details-mapstruct.md`
- [x] `ui/vaadin/client-details-service.md`
- [x] `ui/vaadin/composite-base.md`
- [x] `ui/vaadin/conditional-rendering.md`
- [x] `ui/vaadin/confirmation-dialog.md`
- [x] `ui/vaadin/item-button.md`
- [x] `ui/vaadin/item-confirmation-dialog.md`
- [x] `ui/vaadin/constructor-ordering.md`
- [x] `ui/vaadin/datetime.md`
- [x] `ui/vaadin/dialog-delegation.md`
- [x] `ui/vaadin/error-views/access-denied.md`
- [x] `ui/vaadin/error-views/invalid-request.md`
- [x] `ui/vaadin/error-views/logging.md`
- [x] `ui/vaadin/error-views/not-found.md`
- [x] `ui/vaadin/error-views/overview.md`
- [x] `ui/vaadin/error-views/error-endpoint.md`
- [x] `ui/vaadin/error-views/sensitive-info.md`
- [x] `ui/vaadin/error-views/shared-base.md`
- [x] `ui/vaadin/error-views/system-error.md`
- [x] `ui/vaadin/form-layout.md`
- [x] `ui/vaadin/grid-loading-state.md`
- [x] `ui/vaadin/layout-annotation.md`
- [x] `ui/vaadin/layout-diagram.md`
- [x] `ui/vaadin/navigation/active-route.md`
- [x] `ui/vaadin/navigation/after-navigation.md`
- [x] `ui/vaadin/navigation/app-layout.md`
- [x] `ui/vaadin/navigation/body-scrolling.md`
- [x] `ui/vaadin/navigation/conditional-nav.md`
- [x] `ui/vaadin/navigation/drawer-toggle.md`
- [x] `ui/vaadin/navigation/menu-annotation.md`
- [ ] `ui/vaadin/navigation/mobile-nav.md`
- [x] `ui/vaadin/navigation/navigation-guards.md`
- [x] `ui/vaadin/navigation/route-path-grouping.md`
- [x] `ui/vaadin/quick-filter.md`

- [x] `ui/vaadin/responsive/breakpoints.md`
- [x] `ui/vaadin/responsive/server-side-detection.md`
- [x] `ui/vaadin/router-layout-interface.md`
- [x] `ui/vaadin/service-error-handling.md`
- [x] `ui/vaadin/session-scope.md`
- [x] `ui/vaadin/signals.md`
- [x] `ui/vaadin/spring-component.md`
- [x] `ui/vaadin/theming/brand-assets.md`
- [x] `ui/vaadin/theming/stylesheet-and-theme-loading.md`
- [x] `ui/vaadin/theming/brand-customization.md`
- [x] `ui/vaadin/theming/component-variants.md`
- [x] `ui/vaadin/theming/custom-css.md`
- [x] `ui/vaadin/theming/dark-mode.md`
- [x] `ui/vaadin/theming/theme-selection.md`
- [x] `ui/vaadin/theming/tailwind.md`
- [x] `ui/vaadin/theming/lumo-utility.md`
- [x] `ui/vaadin/uimodel/has-caption.md`
- [x] `ui/vaadin/uimodel/naming.md`
- [x] `ui/vaadin/vaadin-dev-dependency.md`
- [x] `ui/vaadin/view-package.md`
- [x] `ui/vaadin/app-icon.md`

### ui/vaadin/recipes/
- [x] `ui/vaadin/recipes/app-icon.md`
- [x] `ui/vaadin/recipes/base-view.md`
- [x] `ui/vaadin/recipes/caching-data-provider.md`
- [ ] `ui/vaadin/recipes/headroom-scrolling.md`
- [x] `ui/vaadin/recipes/item-browser.md`
- [x] `ui/vaadin/recipes/java-lit-bridge.md`
- [x] `ui/vaadin/recipes/loading-button.md`
- [ ] `ui/vaadin/recipes/touch-bottom-bar.md`
- [ ] `ui/vaadin/recipes/touch-secondary-tabs.md`
- [x] `ui/vaadin/recipes/view-icon.md`
