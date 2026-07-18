import { useState } from 'react'
import { api, storeToken } from '../api.js'
import './AdminLoginModal.css'

export default function AdminLoginModal({ onClose, onLoggedIn }) {
  const [password, setPassword] = useState('')
  const [error, setError] = useState('')
  const [loading, setLoading] = useState(false)

  async function handleSubmit(e) {
    e.preventDefault()
    setLoading(true)
    setError('')
    try {
      const { token } = await api.adminLogin(password)
      storeToken(token)
      onLoggedIn(token)
    } catch (err) {
      setError(err.message || 'Incorrect password')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="admin-modal__backdrop" onClick={onClose}>
      <div className="admin-modal" onClick={(e) => e.stopPropagation()}>
        <form onSubmit={handleSubmit}>
          <label className="admin-modal__label">
            <span>Editor access</span>
            <input
              type="password"
              autoFocus
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              placeholder="Password"
            />
          </label>
          {error ? <p className="admin-modal__error">{error}</p> : null}
          <div className="admin-modal__actions">
            <button type="submit" disabled={loading}>
              {loading ? 'Checking…' : 'Enter'}
            </button>
            <button type="button" className="admin-modal__cancel" onClick={onClose}>
              Cancel
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}
