# MapStruct Configuration

When configuring MapStruct, declare the annotation processor in the parent POM's
`<pluginManagement>` section so all child modules use the same processor version, and
use Spring component model on all mappers.

```xml
<!-- In root {app}/pom.xml -->
<build>
    <pluginManagement>
        <plugins>
            <plugin>
                <groupId>org.apache.maven.plugins</groupId>
                <artifactId>maven-compiler-plugin</artifactId>
                <configuration>
                    <annotationProcessorPaths>
                        <path>
                            <groupId>org.mapstruct</groupId>
                            <artifactId>mapstruct-processor</artifactId>
                            <version>${mapstruct.version}</version>
                        </path>
                    </annotationProcessorPaths>
                </configuration>
            </plugin>
        </plugins>
    </pluginManagement>
</build>
```

All MapStruct mappers use Spring component model:

```java
@Mapper(componentModel = MappingConstants.ComponentModel.SPRING)
public interface EmployeeMapper {
    EmployeeDetail toDetail(EmployeeDetailProjection projection);
    List<EmployeeListItem> toListItems(List<EmployeeListItemProjection> projections);

    // Updates existing entity from UI model — leaves unrelated fields untouched
    EmployeeEntity toEntity(EmployeeDetail detail, @MappingTarget EmployeeEntity entity);
}
```

Generated mapper sources appear under `target/generated-sources/annotations` and are
injected as Spring beans.
