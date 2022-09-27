---
reviewer: {{ .Comment.GetUser.GetLogin }}
created_at: {{ .Comment.GetCreatedAt }}
updated_at: {{ .Comment.GetUpdatedAt }}
---

{{ .Comment.GetBody }}