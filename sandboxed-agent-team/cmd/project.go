package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

// ProjectInfo holds everything the tool auto-discovers about the target
// project. Values may be empty strings if unknown; callers should treat
// "" as "ask the user".
type ProjectInfo struct {
	Name              string
	JavaVersion       string
	VaadinVersion     string
	SpringBootVersion string
	JUnitVersion      string
	Database          string
}

// pomProject is a minimal partial mapping of pom.xml suitable for
// extracting the fields the kit needs. Unmapped elements are ignored.
type pomProject struct {
	XMLName    xml.Name  `xml:"project"`
	ArtifactID string    `xml:"artifactId"`
	Parent     pomParent `xml:"parent"`
	Properties struct {
		JavaVersion         string `xml:"java.version"`
		MavenCompilerSource string `xml:"maven.compiler.source"`
		VaadinVersion       string `xml:"vaadin.version"`
		SpringBootVersion   string `xml:"spring-boot.version"`
	} `xml:"properties"`
	Dependencies struct {
		Dependency []pomDependency `xml:"dependency"`
	} `xml:"dependencies"`
}

type pomParent struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

type pomDependency struct {
	GroupID    string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version    string `xml:"version"`
}

// DiscoverProject reads pom.xml at the given path (typically "pom.xml"
// in the project root) and returns whatever it can infer. If pom.xml
// doesn't exist, returns a zero-value ProjectInfo and no error — the
// caller will fall back to prompting.
func DiscoverProject(pomPath string) (ProjectInfo, error) {
	data, err := os.ReadFile(pomPath)
	if os.IsNotExist(err) {
		return ProjectInfo{}, nil
	}
	if err != nil {
		return ProjectInfo{}, fmt.Errorf("read %s: %w", pomPath, err)
	}

	var p pomProject
	if err := xml.Unmarshal(data, &p); err != nil {
		return ProjectInfo{}, fmt.Errorf("parse %s: %w", pomPath, err)
	}

	info := ProjectInfo{
		Name:              p.ArtifactID,
		JavaVersion:       firstNonEmpty(p.Properties.JavaVersion, p.Properties.MavenCompilerSource),
		VaadinVersion:     p.Properties.VaadinVersion,
		SpringBootVersion: detectSpringBootVersion(p),
		JUnitVersion:      detectJUnitVersion(p.Dependencies.Dependency),
		Database:          detectDatabase(p.Dependencies.Dependency),
	}
	return info, nil
}

func detectSpringBootVersion(p pomProject) string {
	if p.Properties.SpringBootVersion != "" {
		return p.Properties.SpringBootVersion
	}
	if strings.HasPrefix(p.Parent.ArtifactID, "spring-boot-") {
		return p.Parent.Version
	}
	return ""
}

func detectJUnitVersion(deps []pomDependency) string {
	for _, d := range deps {
		switch {
		case strings.Contains(d.ArtifactID, "junit-jupiter"):
			return "5"
		case strings.Contains(d.ArtifactID, "junit-framework"):
			return "6"
		case d.GroupID == "junit" && d.ArtifactID == "junit":
			return "4"
		// Vaadin testing dependencies embed the JUnit major version
		// in the artifact name (e.g., vaadin-testbench-junit5,
		// browserless-test-junit6). Useful for Spring Boot projects
		// that pull JUnit transitively via a starter and declare
		// only the Vaadin test library directly.
		case strings.Contains(d.ArtifactID, "vaadin-testbench-junit5"),
			strings.Contains(d.ArtifactID, "browserless-test-junit5"):
			return "5"
		case strings.Contains(d.ArtifactID, "vaadin-testbench-junit6"),
			strings.Contains(d.ArtifactID, "browserless-test-junit6"):
			return "6"
		}
	}
	return ""
}

// detectDatabase returns a human-readable database label based on which
// JDBC driver dependency is declared. Returns "" if none found.
func detectDatabase(deps []pomDependency) string {
	for _, d := range deps {
		name := d.ArtifactID
		switch {
		case strings.Contains(name, "postgresql"):
			return "PostgreSQL"
		case strings.Contains(name, "mysql-connector"):
			return "MySQL"
		case strings.Contains(name, "mariadb"):
			return "MariaDB"
		case name == "h2":
			return "H2"
		case strings.Contains(name, "sqlite"):
			return "SQLite"
		case strings.Contains(name, "mssql-jdbc") || strings.Contains(name, "sqlserver"):
			return "SQL Server"
		case strings.Contains(name, "ojdbc"):
			return "Oracle"
		}
	}
	return ""
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}

// ToVariables returns the subset of ProjectInfo that can be plugged
// directly into the Variables map as auto-discovered values.
func (p ProjectInfo) ToVariables() Variables {
	v := Variables{}
	if p.Name != "" {
		v["PROJECT_NAME"] = p.Name
	}
	if p.JavaVersion != "" {
		v["JAVA_VERSION"] = p.JavaVersion
	}
	if p.VaadinVersion != "" {
		v["VAADIN_VERSION"] = p.VaadinVersion
	}
	if p.SpringBootVersion != "" {
		v["SPRING_BOOT_VERSION"] = p.SpringBootVersion
	}
	if p.JUnitVersion != "" {
		v["JUNIT_VERSION"] = p.JUnitVersion
	}
	if p.Database != "" {
		v["DATABASE"] = p.Database
	}
	if s := p.StackSummary(); s != "" {
		v["STACK_SUMMARY"] = s
	}
	return v
}

// StackSummary composes a one-line human-readable stack description
// suitable for the STACK_SUMMARY placeholder in generated docs.
func (p ProjectInfo) StackSummary() string {
	var parts []string
	if p.VaadinVersion != "" {
		parts = append(parts, "Vaadin "+p.VaadinVersion)
	}
	if p.SpringBootVersion != "" {
		parts = append(parts, "Spring Boot "+p.SpringBootVersion)
	}
	if p.JavaVersion != "" {
		parts = append(parts, "Java "+p.JavaVersion)
	}
	if p.Database != "" {
		parts = append(parts, p.Database)
	}
	return strings.Join(parts, " / ")
}
