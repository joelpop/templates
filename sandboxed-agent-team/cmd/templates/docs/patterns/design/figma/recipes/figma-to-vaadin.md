---
name: figma-to-vaadin
description: Translate Figma designs to well-structured Vaadin Flow code using Figma MCP and Vaadin MCP. Always extract design context first, check annotations, review documentation, then implement proper Vaadin components with correct themes and semantic structure.
---

# Figma to Vaadin Implementation

When translating a Figma frame into a Vaadin Flow view, extract design context
via the Figma MCP first, resolve any Figma/`docs/reqs` discrepancies before
coding, and map design primitives to Vaadin components plus Lumo tokens so the
resulting view matches the design's semantic structure rather than reproducing
absolute-positioned layout.


## Required Implementation Workflow
Create TODOs based on these steps.

### Step 1. ALWAYS Start with `get_design_context` tool
- Contains the most detailed component information
- Check `data-name` attribute to get the type of the component
- Review component description for identification of correct Vaadin component
- Identify theme/variant hints
- Text styles and typography information (font size, weight, line height)

### Step 2. Figma Component Instance Annotation Checker

When you receive design context from Figma MCP that contains a component instance:

1. **Detect component instances** by checking for node IDs in the format `I[instance-id];[master-component-id]` (e.g., `I1:6;1:3`)

2. **Check instance annotations FIRST** - Examine the component instance for annotations that provide specific implementation guidance.

3. **Extract the master component ID** from the instance node ID:
   - Split on semicolon `;`
   - Take the second part as the master component node ID

4. **Fetch the master component** using the master component ID

5. **Check master component annotations** - These provide default/fallback guidance:
   - Accessibility requirements
   - Recommended Vaadin components
   - Implementation notes
   - Additional content or behavior details

6. **Merge annotations with priority**:
   - **Instance annotations override master component annotations** when both exist
   - Instance-specific annotations are assumed to be more accurate and contextual
   - Use master component annotations only when instance lacks specific guidance

7. **Extract component documentation** from both responses:
   - Component descriptions
   - Documentation links
   - Usage guidelines

Always check both sources; instance annotations take priority.


### Step 3. Use `get_metadata` tool for Structure and identification of components
- Component `name` is the name of the layer and might not correspond to right Vaadin component.
- Plan component hierarchy and relationships
- Analyze node IDs and relationships
- Identify layout patterns and nesting

### Step 4: Component Research (MANDATORY - No Implementation Without This)
**For EACH component identified in Steps 1-2:**

#### 4.1 Component Discovery
- Use `search_vaadin_docs` tool to find relevant components
- Record `file_path` for each component found
- Search results are previews only

#### 4.2 Complete Documentation Review (MANDATORY)
**For each component, call `get_full_document` with the file_path** before implementation:
- TextField → `get_full_document("components/text-field/index-flow.md")`
- DatePicker → `get_full_document("components/date-picker/index-flow.md")`
- Button → `get_full_document("components/button/index-flow.md")`
- etc.

#### 4.3 Implementation Planning
- Document available theme variants
- Note component-specific features
- Identify any limitations or gaps
- Plan component configuration approach

❌ NEVER implement without completing full documentation review
❌ WORKFLOW VIOLATIONS = REJECTION
⚠️ Search results are previews only - not sufficient for implementation

### Step 5: Implement user interface
- Use proper Vaadin components (and project-specific custom elements), not generic HTML
- Apply correct themes and variants
- Ignore visual styling of elements
- Use Lumo Utilities for layouts, padding, borders, background colors
- Don't add spacing or gap to layouts with input fields
- Ensure semantic correctness and correct heading levels based on text styles
- Include accessibility attributes where needed

### Step 6: Don't run tests
- No terminal commands
- No browser or screenshots
