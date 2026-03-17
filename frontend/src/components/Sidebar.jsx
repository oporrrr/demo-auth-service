import { useAuth } from '../context/AuthContext'

export default function Sidebar({ activePath, onNavigate }) {
  const { menus } = useAuth()

  // role service already returns only menus the user is allowed to see
  const visible = menus.filter((m) => m.isActive)

  if (visible.length === 0) return null

  return (
    <aside className="sidebar">
      <ul>
        {visible.map((menu) => {
          const children = (menu.children ?? []).filter((c) => c.isActive)
          const isActive =
            activePath === menu.path ||
            children.some((c) => c.path === activePath)

          return (
            <li key={menu.id} className={isActive ? 'active' : ''}>
              <button onClick={() => onNavigate(menu.path)}>
                {menu.icon && <span className="menu-icon">{menu.icon}</span>}
                {menu.name}
              </button>
              {children.length > 0 && (
                <ul className="submenu">
                  {children.map((child) => (
                    <li key={child.id} className={activePath === child.path ? 'active' : ''}>
                      <button onClick={() => onNavigate(child.path)}>
                        {child.icon && <span className="menu-icon">{child.icon}</span>}
                        {child.name}
                      </button>
                    </li>
                  ))}
                </ul>
              )}
            </li>
          )
        })}
      </ul>
    </aside>
  )
}
