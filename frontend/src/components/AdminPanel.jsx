import { useEffect, useState } from 'react'
import { api, clearToken } from '../api.js'
import './AdminPanel.css'

export default function AdminPanel({ token, story, onLogout, onClose, onStoryUpdated }) {
  const [comments, setComments] = useState([])
  const [stats, setStats] = useState(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState('')

  const [title, setTitle] = useState(story.title)
  const [body, setBody] = useState(story.body)
  const [saving, setSaving] = useState(false)
  const [saveMessage, setSaveMessage] = useState('')

  useEffect(() => {
    let cancelled = false
    async function load() {
      try {
        const [c, s] = await Promise.all([
          api.adminGetComments(token),
          api.adminGetStats(token),
        ])
        if (!cancelled) {
          setComments(c.comments || [])
          setStats(s)
        }
      } catch (err) {
        if (!cancelled) setError(err.message)
      } finally {
        if (!cancelled) setLoading(false)
      }
    }
    load()
    return () => {
      cancelled = true
    }
  }, [token])

  async function handleDelete(id) {
    if (!window.confirm('Delete this letter? This cannot be undone.')) return
    try {
      await api.adminDeleteComment(token, id)
      setComments((prev) => prev.filter((c) => c.id !== id))
    } catch (err) {
      setError(err.message)
    }
  }

  async function handleSaveStory(e) {
    e.preventDefault()
    setSaving(true)
    setSaveMessage('')
    try {
      await api.adminSaveStory(token, title, body)
      onStoryUpdated({ title, body })
      setSaveMessage('Saved. The published story now reflects this text.')
    } catch (err) {
      setSaveMessage(err.message || 'Could not save.')
    } finally {
      setSaving(false)
    }
  }

  function handleLogout() {
    clearToken()
    onLogout()
  }

  return (
    <div className="admin-panel__backdrop">
      <div className="admin-panel">
        <div className="admin-panel__topbar">
          <h2>Editor's Desk</h2>
          <div className="admin-panel__topbar-actions">
            <button className="admin-panel__close" onClick={onClose}>
              Close
            </button>
            <button className="admin-panel__logout" onClick={handleLogout}>
              Log out
            </button>
          </div>
        </div>

        {error ? <p className="admin-panel__error">{error}</p> : null}

        <section className="admin-panel__section">
          <h3>Readership</h3>
          {loading ? (
            <p>Loading…</p>
          ) : (
            <div className="admin-panel__stats">
              <div>
                <strong>{stats?.views ?? '—'}</strong>
                <span>page views</span>
              </div>
              <div>
                <strong>{stats?.comments ?? comments.length}</strong>
                <span>letters received</span>
              </div>
            </div>
          )}
        </section>

        <section className="admin-panel__section">
          <h3>Letters ({comments.length})</h3>
          {loading ? (
            <p>Loading…</p>
          ) : comments.length === 0 ? (
            <p className="admin-panel__empty">No letters yet.</p>
          ) : (
            <ul className="admin-panel__comments">
              {comments
                .slice()
                .reverse()
                .map((c) => (
                  <li key={c.id}>
                    <div>
                      <strong>{c.name}</strong>
                      <p>{c.text}</p>
                    </div>
                    <button onClick={() => handleDelete(c.id)}>Delete</button>
                  </li>
                ))}
            </ul>
          )}
        </section>

        <section className="admin-panel__section">
          <h3>Manuscript</h3>
          <form onSubmit={handleSaveStory}>
            <label className="admin-panel__field">
              <span>Title</span>
              <input value={title} onChange={(e) => setTitle(e.target.value)} />
            </label>
            <label className="admin-panel__field">
              <span>Story text</span>
              <textarea
                value={body}
                onChange={(e) => setBody(e.target.value)}
                rows={14}
              />
            </label>
            <p className="admin-panel__hint">
              Separate paragraphs with a blank line. This replaces what every visitor sees.
            </p>
            {saveMessage ? <p className="admin-panel__save-msg">{saveMessage}</p> : null}
            <button type="submit" className="admin-panel__save" disabled={saving}>
              {saving ? 'Saving…' : 'Save & publish'}
            </button>
          </form>
        </section>
      </div>
    </div>
  )
}
