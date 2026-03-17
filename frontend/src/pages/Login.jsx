import { useState } from 'react'
import { api } from '../api'
import { useAuth } from '../context/AuthContext'

export default function Login({ onLogin }) {
  const { loadPermissions } = useAuth()
  const [loginType, setLoginType] = useState('EMAIL')
  const [form, setForm] = useState({ email: '', countryCode: '', phoneNumber: '', password: '' })
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  function handleChange(e) {
    setForm((f) => ({ ...f, [e.target.name]: e.target.value }))
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    setLoading(true)

    const body = { loginType, password: form.password }
    if (loginType === 'EMAIL') {
      body.email = form.email
    } else {
      body.countryCode = form.countryCode
      body.phoneNumber = form.phoneNumber
    }

    try {
      const { status, data } = await api.login(body)
      if (data.code === 'SUCCESS' && data.data?.accessToken) {
        localStorage.setItem('accessToken', data.data.accessToken)
        localStorage.setItem('refreshToken', data.data.refreshToken)
        await loadPermissions()
        onLogin()
      } else {
        setError(data.message || `Login failed (${status})`)
      }
    } catch (err) {
      setError('Network error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="card">
      <h2>Login</h2>
      <div className="toggle-group">
        <button className={loginType === 'EMAIL' ? 'active' : ''} onClick={() => setLoginType('EMAIL')} type="button">
          Email
        </button>
        <button className={loginType === 'PHONE_NUMBER' ? 'active' : ''} onClick={() => setLoginType('PHONE_NUMBER')} type="button">
          Phone
        </button>
      </div>
      <form onSubmit={handleSubmit}>
        {loginType === 'EMAIL' ? (
          <input name="email" type="email" placeholder="Email" value={form.email} onChange={handleChange} required />
        ) : (
          <div className="row">
            <input name="countryCode" placeholder="+66" value={form.countryCode} onChange={handleChange} style={{ width: 80 }} required />
            <input name="phoneNumber" placeholder="Phone number" value={form.phoneNumber} onChange={handleChange} required />
          </div>
        )}
        <input name="password" type="password" placeholder="Password" value={form.password} onChange={handleChange} required />
        {error && <p className="error">{error}</p>}
        <button type="submit" disabled={loading}>{loading ? 'Logging in…' : 'Login'}</button>
      </form>
    </div>
  )
}
