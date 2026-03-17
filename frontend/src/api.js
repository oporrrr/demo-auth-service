const BASE_URL = 'http://localhost:3000/api/v1'

function authHeader() {
  const token = localStorage.getItem('accessToken')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

async function request(method, path, body) {
  const res = await fetch(`${BASE_URL}${path}`, {
    method,
    headers: { 'Content-Type': 'application/json', ...authHeader() },
    body: body ? JSON.stringify(body) : undefined,
  })
  const data = await res.json()
  return { status: res.status, data }
}

export const api = {
  // Auth
  register: (body) => request('POST', '/auth/register', body),
  login: (body) => request('POST', '/auth/login', body),
  logout: (refreshToken) => request('POST', '/auth/logout', { refreshToken }),

  // Account
  getInformation: () => request('GET', '/account/information'),

  // Permissions + Menus (via role service)
  getPermissions: () => request('GET', '/me/permissions'),
  getMenus: () => request('GET', '/me/menus'),

  // Demo permission-protected endpoints
  demo: {
    view:   (resource)     => request('GET',    `/demo/${resource}`),
    create: (resource)     => request('POST',   `/demo/${resource}`),
    edit:   (resource, id) => request('PUT',    `/demo/${resource}/${id}`),
    delete: (resource, id) => request('DELETE', `/demo/${resource}/${id}`),
  },
}
