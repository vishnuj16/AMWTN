const ADMIN_TOKEN_KEY = 'admin_token'

async function request(path, options = {}) {
  const res = await fetch(`/api${path}`, {
    ...options,
    headers: {
      'Content-Type': 'application/json',
      ...options.headers,
    },
  })
  const data = await res.json().catch(() => ({}))
  if (!res.ok) {
    throw new Error(data.error || `Request failed (${res.status})`)
  }
  return data
}

export const api = {
  getStory: () => request('/story'),

  getComments: () => request('/comments'),
  addComment: (name, text) =>
    request('/comments', {
      method: 'POST',
      body: JSON.stringify({ name, text }),
    }),

  recordView: () =>
    request('/view', { method: 'POST' }).catch(() => {}),

  adminLogin: (password) =>
    request('/admin/login', {
      method: 'POST',
      body: JSON.stringify({ password }),
    }),

  adminGetComments: (token) =>
    request('/admin/comments', {
      headers: { Authorization: `Bearer ${token}` },
    }),
  adminDeleteComment: (token, id) =>
    request(`/admin/comments?id=${encodeURIComponent(id)}`, {
      method: 'DELETE',
      headers: { Authorization: `Bearer ${token}` },
    }),
  adminGetStats: (token) =>
    request('/admin/stats', {
      headers: { Authorization: `Bearer ${token}` },
    }),
  adminSaveStory: (token, title, body) =>
    request('/admin/story', {
      method: 'POST',
      headers: { Authorization: `Bearer ${token}` },
      body: JSON.stringify({ title, body }),
    }),
}

export function getStoredToken() {
  return sessionStorage.getItem(ADMIN_TOKEN_KEY)
}
export function storeToken(token) {
  sessionStorage.setItem(ADMIN_TOKEN_KEY, token)
}
export function clearToken() {
  sessionStorage.removeItem(ADMIN_TOKEN_KEY)
}
