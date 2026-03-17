import { usePermission } from '../hooks/usePermission'
import { useAuth } from '../context/AuthContext'

function PermBadge({ perm }) {
  const can = usePermission()
  const allowed = can(perm)
  return (
    <span className={`perm-badge ${allowed ? 'allowed' : 'denied'}`}>
      {allowed ? '✓' : '✗'} {perm}
    </span>
  )
}

function GuardedSection({ title, permission, children }) {
  const can = usePermission()
  return (
    <div className="demo-card">
      <div className="section-header">
        <h3>{title}</h3>
        <PermBadge perm={permission} />
      </div>
      {can(permission)
        ? children
        : <p className="no-access">No access — missing <code>{permission}</code></p>
      }
    </div>
  )
}

function ActionBtn({ permission, label, variant = 'btn-default' }) {
  const can = usePermission()
  if (can(permission)) {
    return <button className={`btn ${variant}`}>{label}</button>
  }
  return <span className="hidden-hint">✗ {label} <small>({permission})</small></span>
}

export default function PermissionDemo() {
  const { permissions } = useAuth()

  return (
    <div className="demo-page">
      <h2>Permission Demo</h2>

      {/* Current permissions */}
      <div className="demo-card">
        <h3>Your Permissions</h3>
        <div className="perm-list">
          {permissions.length === 0
            ? <span className="no-perm">No permissions — please login</span>
            : permissions.map((p) => (
                <span key={p} className="perm-badge allowed">✓ {p}</span>
              ))
          }
        </div>
      </div>

      {/* Setting section */}
      <GuardedSection title="Setting" permission="setting:view">
        <div className="btn-group" style={{ marginTop: 12 }}>
          <ActionBtn permission="setting:create" label="+ New Setting" variant="btn-primary" />
          <ActionBtn permission="setting:edit"   label="Edit"          variant="btn-warning" />
          <ActionBtn permission="setting:delete" label="Delete"        variant="btn-danger" />
        </div>
      </GuardedSection>

      {/* Dashboard section */}
      <GuardedSection title="Dashboard" permission="menu_dashboard:view">
        <p style={{ marginTop: 12, color: '#555', fontSize: '0.9rem' }}>
          You can see this section because you have <code>menu_dashboard:view</code>
        </p>
      </GuardedSection>

      {/* Order section — likely no access */}
      <GuardedSection title="Order" permission="order:view">
        <div className="btn-group" style={{ marginTop: 12 }}>
          <ActionBtn permission="order:create"  label="+ New Order"   variant="btn-primary" />
          <ActionBtn permission="order:edit"    label="Edit"          variant="btn-warning" />
          <ActionBtn permission="order:delete"  label="Delete"        variant="btn-danger" />
          <ActionBtn permission="order:approve" label="Approve"       variant="btn-success" />
        </div>
      </GuardedSection>

      {/* Report section — likely no access */}
      <GuardedSection title="Report" permission="report:view">
        <div className="btn-group" style={{ marginTop: 12 }}>
          <ActionBtn permission="report:export" label="Export CSV" variant="btn-primary" />
        </div>
      </GuardedSection>

      {/* Full matrix */}
      <div className="demo-card">
        <h3>Full Permission Matrix</h3>
        <div className="perm-list" style={{ marginTop: 8 }}>
          {[
            'menu_dashboard:view', 'setting:view', 'setting:create', 'setting:edit', 'setting:delete',
            'order:view', 'order:create', 'order:edit', 'order:delete', 'order:approve',
            'report:view', 'report:export',
            'user:view', 'user:create', 'user:edit', 'user:delete',
          ].map((p) => <PermBadge key={p} perm={p} />)}
        </div>
      </div>
    </div>
  )
}
