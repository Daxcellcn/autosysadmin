
```markdown
<!-- frontend/docs/architecture.md -->
# Autosysadmin Architecture Overview

## System Architecture

```mermaid
graph TD
  A[Frontend] -->|API Calls| B[Backend API]
  B --> C[Database]
  B --> D[Redis]
  B --> E[Agent Connections]
  E --> F[Managed Servers]