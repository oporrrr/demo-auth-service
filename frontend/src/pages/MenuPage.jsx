import { useState } from 'react'
import { usePermission } from '../hooks/usePermission'
import { useAuth } from '../context/AuthContext'
import { api } from '../api'

const ACTIONS = [
  { key: 'view',    label: 'View',    variant: 'btn-default' },
  { key: 'create',  label: '+ Create', variant: 'btn-primary' },
  { key: 'edit',    label: 'Edit',    variant: 'btn-warning' },
  { key: 'delete',  label: 'Delete',  variant: 'btn-danger'  },
  { key: 'approve', label: 'Approve', variant: 'btn-success' },
  { key: 'export',  label: 'Export',  variant: 'btn-default' },
  { key: 'import',  label: 'Import',  variant: 'btn-default' },
  { key: 'print',   label: 'Print',   variant: 'btn-default' },
]

function extractResource(menu) {
  // prefer explicit permission field, fallback to menu code
  const raw = menu.permission ?? menu.code ?? ''
  return raw.replace(/^menu_/, '')
}

async function callDemoAPI(resource, actionKey) {
  switch (actionKey) {
    case 'view':   return api.demo.view(resource)
    case 'create': return api.demo.create(resource)
    case 'edit':   return api.demo.edit(resource, 1)
    case 'delete': return api.demo.delete(resource, 1)
    default:       return { status: 200, data: { code: 'SUCCESS', data: `${actionKey} (no API mapped)` } }
  }
}

export default function MenuPage({ path }) {
  const can = usePermission()
  const { menus } = useAuth()
  const [result, setResult] = useState(null)
  const [loading, setLoading] = useState(false)

  const allMenus = menus.flatMap((m) => [m, ...(m.children ?? [])])
  const menu = allMenus.find((m) => m.path === path)

  if (!menu) {
    return (
      <div className="menu-page">
        <div className="menu-page-header">
          <span className="menu-page-path">{path}</span>
          <h2>Page not found</h2>
        </div>
      </div>
    )
  }

  const resource = extractResource(menu)
  const allowedActions = ACTIONS.filter((a) => can(`${resource}:${a.key}`))
  const deniedActions  = ACTIONS.filter((a) => !can(`${resource}:${a.key}`))

  async function handleAction(actionKey) {
    setLoading(true)
    setResult(null)
    try {
      const res = await callDemoAPI(resource, actionKey)
      setResult({ action: actionKey, ...res })
    } catch (e) {
      setResult({ action: actionKey, status: 0, data: { code: 'ERROR', message: e.message } })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="menu-page">
      <div className="menu-page-header">
        <span className="menu-page-path">{menu.path}</span>
        <h2>{menu.icon && <span style={{ marginRight: 8 }}>{menu.icon}</span>}{menu.name}</h2>
        <span className="menu-page-perm">{menu.permission ?? menu.code}</span>
      </div>

      <div className="menu-page-body">
        <div className="demo-card">
          <h3>Actions</h3>
          {allowedActions.length > 0 ? (
            <div className="btn-group" style={{ marginTop: 12 }}>
              {allowedActions.map((a) => (
                <button
                  key={a.key}
                  className={`btn ${a.variant}`}
                  disabled={loading}
                  onClick={() => handleAction(a.key)}
                >
                  {a.label}
                </button>
              ))}
            </div>
          ) : (
            <p style={{ color: '#aaa', fontSize: '0.875rem', marginTop: 8 }}>
              No action permissions for <code>{resource}:*</code>
            </p>
          )}

          {deniedActions.length > 0 && (
            <div className="btn-group" style={{ marginTop: 10 }}>
              {deniedActions.map((a) => (
                <span key={a.key} className="hidden-hint">
                  ✗ {a.label} <small>({resource}:{a.key})</small>
                </span>
              ))}
            </div>
          )}
        </div>

        {result && (
          <div className="demo-card">
            <h3>
              API Response
              <span style={{ marginLeft: 8, fontSize: '0.8rem', fontWeight: 400, color: result.status === 200 ? '#4caf50' : '#f44336' }}>
                {result.action} → HTTP {result.status}
              </span>
            </h3>
            <pre style={{ marginTop: 8, fontSize: '0.85rem', background: '#1a1a1a', padding: 12, borderRadius: 6, overflowX: 'auto' }}>
              {JSON.stringify(result.data, null, 2)}
            </pre>
          </div>
        )}

        <div className="demo-card">
          <h3>Content Area</h3>
          <p style={{ color: '#999', fontSize: '0.9rem', marginTop: 8 }}>
            This is where the <strong>{menu.name}</strong> content would go.
          </p>
        </div>
      </div>
    </div>
  )
}
