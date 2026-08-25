# Error View Overview

When implementing the application's error views, create one view per error
type from the table below so every error condition shows a consistent,
friendly view and attackers cannot probe the system by triggering errors.

This document describes the *shape* each view must conform to; the project
picks its own copy, layout primitives, icons, class names, and home-view target.

| Error type      | HTTP equivalent        | Exception type             |
|:----------------|:-----------------------|:---------------------------|
| Not Found       | 404                    | `NotFoundException`        |
| Access Denied   | 403 (displayed as 404) | `AccessDeniedException`    |
| System Error    | 500                    | `Exception`                |
| Invalid Request | 400                    | `IllegalArgumentException` |

No error condition results in a raw stack trace or framework error overlay
visible to the user.
