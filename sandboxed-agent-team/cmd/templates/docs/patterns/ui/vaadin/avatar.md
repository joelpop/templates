# Avatar Component

When displaying a user profile photo or entity logo, use the Vaadin `Avatar`
component so the initials fall-back is shown automatically when no image is set
and display is consistent across views.

```java
// User avatar
var userAvatar = new Avatar(user.getFullName());
userAvatar.setImage(user.getPhotoUrl());  // displays initials fallback when null

// Entity logo
var entityLogoAvatar = new Avatar(entity.getName());
entityLogoAvatar.setImage(entity.getLogoUrl());
```

Do not use `<img>` tags or custom image components for user photos or entity logos.
