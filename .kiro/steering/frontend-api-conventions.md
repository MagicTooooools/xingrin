---
inclusion: fileMatch
fileMatchPattern: "frontend/**/*.{ts,tsx}"
---

# Frontend API Conventions

## API Client Usage

### ALWAYS Use `api` Client

```typescript
// ✅ CORRECT - Uses api client, automatically adds JWT token
import { api } from '@/lib/api-client'

const response = await api.get<WebsiteListResponse>('/targets/${id}/websites/', { params })
const response = await api.post('/websites/bulk-delete/', { ids })

// ❌ WRONG - Native fetch doesn't add JWT token, will get 401
const response = await fetch('/api/targets/${id}/websites/')
```

### Why This Matters

The `api` client (axios instance) automatically:
1. Adds `Authorization: Bearer <token>` header
2. Handles token refresh on 401 errors
3. Logs requests/responses in development
4. Provides consistent error handling


## JWT Token Auto-Refresh

The API client automatically refreshes tokens:

1. Request returns 401 (token expired)
2. Client uses Refresh Token to get new Access Token
3. Original request is retried with new token
4. Multiple concurrent 401s are queued and retried together

Token expiration:
- Access Token: 15 minutes
- Refresh Token: 7 days

Users only need to re-login if Refresh Token expires (after 7 days of inactivity).


## Service Layer Pattern

### Organize API Calls in Services

```typescript
// ✅ Good - Centralized in service file
// frontend/services/website.service.ts
export class WebsiteService {
  static async bulkDelete(ids: number[]) {
    const response = await api.post('/websites/bulk-delete/', { ids })
    return response.data
  }
}

// ❌ Avoid - Scattered API calls in components
// frontend/components/SomeComponent.tsx
const response = await api.post('/websites/bulk-delete/', { ids })
```

### Hooks for React Query

```typescript
// frontend/hooks/use-websites.ts
export function useBulkDeleteWebsites() {
  return useMutation({
    mutationFn: (ids: number[]) => WebsiteService.bulkDelete(ids),
    // ... toast messages, cache invalidation
  })
}
```


## API Path Conventions

### Match Backend Routes Exactly

```typescript
// Backend route: GET /api/targets/:id/websites
// Frontend call:
api.get(`/targets/${id}/websites/`)

// Backend route: POST /api/websites/bulk-delete
// Frontend call:
api.post('/websites/bulk-delete/', { ids })
```

### Trailing Slash

Always include trailing slash to match Next.js rewrite rules:

```typescript
// ✅ With trailing slash
api.get('/targets/')
api.post('/websites/bulk-delete/')

// ❌ Without trailing slash (may cause redirect issues)
api.get('/targets')
api.post('/websites/bulk-delete')
```


## Blob Downloads (CSV/Excel Export)

### Handle Error Responses in Blob Mode

When using `responseType: 'blob'`, error responses are also returned as Blob:

```typescript
static async exportByTargetId(targetId: number): Promise<Blob> {
  const response = await api.get<Blob>(`/targets/${targetId}/websites/export/`, {
    responseType: 'blob',
  })
  
  // Check if response is actually an error (JSON instead of CSV)
  if (response.data.type === 'application/json') {
    const text = await response.data.text()
    const error = JSON.parse(text)
    throw new Error(error.error?.message || 'Export failed')
  }
  
  return response.data
}
```
