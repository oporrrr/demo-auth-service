import { useEffect, useState } from 'react'
import { api } from '../api'
import { useAuth } from '../context/AuthContext'

export default function Profile({ onLogout }) {
  const { clearAuth } = useAuth()
  const [info, setInfo] = useState(null)
  const [error, setError] = useState('')

  useEffect(() => {
    api.getInformation().then(({ data }) => {
      if (data.code === 'SUCCESS') {
        setInfo(data.data)
      } else {
        setError(data.message || 'Failed to load profile')
      }
    }).catch(() => setError('Network error'))
  }, [])

  async function handleLogout() {
    const refreshToken = localStorage.getItem('refreshToken')
    if (refreshToken) {
      await api.logout(refreshToken).catch(() => {})
    }
    clearAuth()
    onLogout()
  }

  if (error) return (
    <div className="card">
      <p className="error">{error}</p>
      <button onClick={handleLogout}>Back to Login</button>
    </div>
  )

  if (!info) return <div className="card"><p>Loading…</p></div>

  const rows = [
    ['ID', info.id],
    ['Prefix', info.prefixName],
    ['First name', info.firstName],
    ['Last name', info.lastName],
    ['Gender', info.gender],
    ['Date of birth', info.dateOfBirth],
    ['Email', info.email],
    ['Country code', info.countryCode],
    ['Phone', info.phoneNumber],
    ['Username', info.username],
    ['Account status', info.accountStatus],
    ['CIS number', info.cisNumber ?? '—'],
  ]

  return (
    <div className="card">
      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
        <h2>Profile</h2>
        <button onClick={handleLogout} style={{ background: '#e74c3c' }}>Logout</button>
      </div>
      <table className="profile-table">
        <tbody>
          {rows.map(([label, value]) => (
            <tr key={label}>
              <td className="label">{label}</td>
              <td>{value || '—'}</td>
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  )
}
