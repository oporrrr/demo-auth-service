import { useAuth } from '../context/AuthContext'

/**
 * Returns a `can(permission)` function for the current user.
 *
 * Usage:
 *   const can = usePermission()
 *   {can('order:edit') && <button>Edit</button>}
 *   {can('btn_approve:approve') && <button>Approve</button>}
 */
export function usePermission() {
  const { can } = useAuth()
  return can
}
