# GENERATED FILE — do not edit directly.
# Edits here will be lost the next time this file is regenerated.
# To change this file, edit its template in the kit source and
# re-run `agent-team-install`.

FROM docker/sandbox-templates:claude-code

USER root

# ── Suppress interactive prompts during apt-get ───────────────────────────
# Without this, packages like x11-common (a Playwright dependency) fail to
# configure because there is no terminal during docker build.
ENV DEBIAN_FRONTEND=noninteractive

# ── Java / Maven ───────────────────────────────────────────────────────────
# Replace {{JAVA_VERSION}} with the project's Java version (e.g., 21, 17)
# The base image already includes make, curl, unzip, jq, socat, node, npm.
# openssh-client is required for Git SSH remote access (the base image does
# not include it). Remove it if the project uses HTTPS remotes only.
RUN apt-get update \
    && apt-get install -y openjdk-{{JAVA_VERSION}}-jdk maven openssh-client \
    && rm -rf /var/lib/apt/lists/*

# ── Playwright browser (Chromium only) ────────────────────────────────────
# Required for E2E tests (Node.js Playwright) and for visual debugging
# via the playwright MCP server. Install Chromium only — other browsers
# are not needed. This layer rarely changes, so it caches well.
# Typically takes 1-3 minutes.
# Install the Playwright npm package first so `playwright install` does
# not warn about missing project dependencies.
# Playwright may print a "BEWARE: your OS is not officially supported"
# warning on some platforms and use a fallback. This is cosmetic —
# Chromium works fine via the fallback. The warning can be ignored.
ENV PLAYWRIGHT_BROWSERS_PATH=/opt/browsers
RUN npm install -g playwright \
    && timeout 300 npx playwright install --with-deps chromium

# ── MCP server pre-install ────────────────────────────────────────────────
# Command-type MCP servers (spring-docs, playwright, fetch) use npx/npm.
# Install globally so they are immediately available at runtime without
# npx re-downloading them on each invocation. Rebuild the image
# periodically to refresh versions.
RUN timeout 120 npm install -g @enokdev/springdocs-mcp@latest 2>/dev/null || true \
    && timeout 120 npm install -g @playwright/mcp@latest 2>/dev/null || true \
    && timeout 120 npm install -g fetch-mcp@latest 2>/dev/null || true

# ── Cost tracking (ccusage) ───────────────────────────────────────────────
# ccusage reads Claude Code's local JSONL session logs and reports token
# usage and USD cost estimates. Works for both API-key and subscription
# users (unlike Claude Code's built-in /cost, which gives no numbers to
# subscription users). The Integrator invokes it at task conclusion to
# build the cost report included in the final squash-merge commit
# message. See https://github.com/ryoppippi/ccusage.
RUN timeout 120 npm install -g ccusage@latest 2>/dev/null || true

# ── Project-specific extras ────────────────────────────────────────────────
# Pre-install project dependencies so agents don't waste tokens
# COPY pom.xml /tmp/pom.xml
# RUN cd /tmp && mvn dependency:resolve

USER agent

# ── Git identity ──────────────────────────────────────────────────────────
# Replace {{GIT_USER_NAME}} and {{GIT_USER_EMAIL}} with the developer's git
# identity (discovered from host `git config` during setup).
RUN git config --global user.name "{{GIT_USER_NAME}}" \
    && git config --global user.email "{{GIT_USER_EMAIL}}"

