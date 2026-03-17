import { useState } from 'react'
import { api } from '../api'

export default function Register({ onSuccess }) {
  const [form, setForm] = useState({
    prefixName: '',
    firstName: '',
    lastName: '',
    gender: '',
    dateOfBirth: '',
    email: '',
    countryCode: '',
    phoneNumber: '',
    password: '',
  })
  const [loginType, setLoginType] = useState('EMAIL')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  function handleChange(e) {
    setForm((f) => ({ ...f, [e.target.name]: e.target.value }))
  }

  async function handleSubmit(e) {
    e.preventDefault()
    setError('')
    setLoading(true)

    const body = {
      prefixName: form.prefixName,
      firstName: form.firstName,
      lastName: form.lastName,
      gender: form.gender,
      dateOfBirth: form.dateOfBirth,
      password: form.password,
    }

    if (loginType === 'EMAIL') {
      body.email = form.email
    } else {
      body.countryCode = form.countryCode
      body.phoneNumber = form.phoneNumber
    }

    try {
      const { status, data } = await api.register(body)
      if (data.code === 'SUCCESS') {
        onSuccess('Register successful! Please login.')
      } else {
        setError(data.message || `Error ${status}`)
      }
    } catch (err) {
      setError('Network error')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="card">
      <h2>Register</h2>
      <div className="toggle-group">
        <button
          className={loginType === 'EMAIL' ? 'active' : ''}
          onClick={() => setLoginType('EMAIL')}
          type="button"
        >
          Email
        </button>
        <button
          className={loginType === 'PHONE_NUMBER' ? 'active' : ''}
          onClick={() => setLoginType('PHONE_NUMBER')}
          type="button"
        >
          Phone
        </button>
      </div>
      <form onSubmit={handleSubmit}>
        <select name="prefixName" value={form.prefixName} onChange={handleChange}>
          <option value="">-- Prefix --</option>
          <option value="MR">MR</option>
          <option value="MRS">MRS</option>
          <option value="MISS">MISS</option>
          <option value="MS">MS</option>
        </select>
        <div className="row">
          <input name="firstName" placeholder="First name *" value={form.firstName} onChange={handleChange} required />
          <input name="lastName" placeholder="Last name *" value={form.lastName} onChange={handleChange} required />
        </div>
        <select name="gender" value={form.gender} onChange={handleChange}>
          <option value="">-- Gender --</option>
          <option value="M">Male</option>
          <option value="F">Female</option>
          <option value="NP">Not specified</option>
        </select>
        <input name="dateOfBirth" type="date" placeholder="Date of birth" value={form.dateOfBirth} onChange={handleChange} />
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
        <button type="submit" disabled={loading}>{loading ? 'Registering…' : 'Register'}</button>
      </form>
    </div>
  )
}
