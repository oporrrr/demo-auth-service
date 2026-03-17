import { createContext, useContext, useState, useCallback, useEffect } from 'react'
import { api } from '../api'

const AuthContext = createContext(null)

export function AuthProvider({ children }) {
  const [permissions, setPermissions] = useState(() => {
    try {
      return JSON.parse(localStorage.getItem('permissions') || '[]')
    } catch {
      return []
    }
  })
  const [menus, setMenus] = useState(() => {
    try {
      return JSON.parse(localStorage.getItem('menus') || '[]')
    } catch {
      return []
    }
  })

  const loadPermissions = useCallback(async () => {
    try {
      const [permRes, menuRes] = await Promise.all([
        api.getPermissions(),
        api.getMenus(),
      ])
      if (permRes.data.code === 'SUCCESS') {
        const perms = permRes.data.data ?? []
        setPermissions(perms)
        localStorage.setItem('permissions', JSON.stringify(perms))
      }
      if (menuRes.data.code === 'SUCCESS') {
        const m = menuRes.data.data ?? []
        setMenus(m)
        localStorage.setItem('menus', JSON.stringify(m))
      }
    } catch (e) {
      console.warn('Failed to load permissions/menus:', e)
    }
  }, [])

  // fetch fresh on page reload if token exists
  useEffect(() => {
    if (localStorage.getItem('accessToken')) {
      loadPermissions()
    }
  }, [loadPermissions])

  const clearAuth = useCallback(() => {
    localStorage.removeItem('accessToken')
    localStorage.removeItem('refreshToken')
    localStorage.removeItem('permissions')
    localStorage.removeItem('menus')
    setPermissions([])
    setMenus([])
  }, [])

  const can = useCallback(
    (permission) => {
      if (!permission) return false
      if (permissions.includes(permission)) return true
      const [reqResource, reqAction] = permission.split(':')
      return permissions.some((p) => {
        const [r, a] = p.split(':')
        return (r === '*' || r === reqResource) && (a === '*' || a === reqAction)
      })
    },
    [permissions]
  )

  return (
    <AuthContext.Provider value={{ permissions, menus, can, loadPermissions, clearAuth }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  return useContext(AuthContext)
}
