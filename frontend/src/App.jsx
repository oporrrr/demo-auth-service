import { useState } from 'react'
import { useAuth } from './context/AuthContext'
import Login from './pages/Login'
import Register from './pages/Register'
import Profile from './pages/Profile'
import MenuPage from './pages/MenuPage'
import Sidebar from './components/Sidebar'
import './App.css'

export default function App() {
  const [page, setPage] = useState(localStorage.getItem('accessToken') ? '/profile' : 'login')
  const [flash, setFlash] = useState('')
  const { menus } = useAuth()

  const loggedIn = !!localStorage.getItem('accessToken')

  const allMenuPaths = menus.flatMap((m) => {
    const paths = [m.path]
    if (m.children) paths.push(...m.children.map((c) => c.path))
    return paths
  })

  function showFlash(msg) {
    setFlash(msg)
    setTimeout(() => setFlash(''), 4000)
  }

  function renderPage() {
    if (page === 'login')    return <Login onLogin={() => setPage('/profile')} />
    if (page === 'register') return <Register onSuccess={(msg) => { showFlash(msg); setPage('login') }} />
    if (page === '/profile') return <Profile onLogout={() => setPage('login')} />
    if (allMenuPaths.includes(page)) return <MenuPage path={page} />
    return <Profile onLogout={() => setPage('login')} />
  }

  return (
    <div className="app">
      <nav>
        <span className="brand">Demo Auth Center</span>
        <div className="nav-links">
          {!loggedIn && page !== 'login'    && <button onClick={() => setPage('login')}>Login</button>}
          {!loggedIn && page !== 'register' && <button onClick={() => setPage('register')}>Register</button>}
          {loggedIn  && page !== '/profile' && <button onClick={() => setPage('/profile')}>Profile</button>}
        </div>
      </nav>

      {flash && <div className="flash">{flash}</div>}

      <div className={loggedIn && menus.length > 0 ? 'layout' : ''}>
        {loggedIn && <Sidebar activePath={page} onNavigate={setPage} />}
        <main className="content">
          {renderPage()}
        </main>
      </div>
    </div>
  )
}
