# Avatar Component

When displaying a user profile photo or entity logo, use the Vaadin `Avatar`
component so the initials fallback is shown automatically when no image is set
and display is consistent across views.

```java
// User avatar
var avatar = new Avatar(user.getFullName());
avatar.setImage(user.getPhotoUrl());  // displays initials fallback when null

// Entity logo
var logo = new Avatar(entity.getName());
logo.setImage(entity.getLogoUrl());
```

Do not use `<img>` tags or custom image components for user photos or entity logos.

Application-specific avatar ring patterns (e.g., role indicator rings) are defined
in the application-specific UI requirements.
