# Dependency Version Management

When declaring dependency or plugin versions, put all versions in the parent POM's
`<dependencyManagement>` and `<pluginManagement>` sections so child module POMs omit
version numbers and upgrade commits touch only one file.

- All dependency versions declared in parent POM `<dependencyManagement>`; child module POMs do not specify version numbers for managed dependencies.
- All plugin versions declared in parent POM `<pluginManagement>`; child module POMs do not specify plugin versions inline.

This ensures version conflicts are resolved centrally and upgrade commits are minimal.
