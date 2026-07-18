import { useState } from 'react'
import { api } from '../api.js'
import './CommentSection.css'

function formatDate(iso) {
  try {
    return new Date(iso).toLocaleDateString('en-IN', {
      day: 'numeric',
      month: 'short',
      year: 'numeric',
    })
  } catch {
    return ''
  }
}

export default function CommentSection({ comments, onCommentAdded }) {
  const [showForm, setShowForm] = useState(false)
  const [name, setName] = useState('')
  const [text, setText] = useState('')
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

  async function handleSubmit(e) {
    e.preventDefault()
    if (!text.trim()) {
      setError('A letter needs a few words, at least.')
      return
    }
    setSubmitting(true)
    setError('')
    try {
      const saved = await api.addComment(name, text)
      onCommentAdded(saved)
      setName('')
      setText('')
      setShowForm(false)
    } catch (err) {
      setError(err.message || 'Could not send that just now. Try again in a moment.')
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <section className="letters" id="letters">
      <div className="letters__head">
        <h2 className="letters__title">Letters to the Editor</h2>
        <p className="letters__sub">
          Thoughts on the story, from readers who made it to the last page.
        </p>
      </div>

      {!showForm ? (
        <button className="letters__cta" onClick={() => setShowForm(true)}>
          Write a letter
        </button>
      ) : (
        <form className="letters__form" onSubmit={handleSubmit}>
          <label className="letters__field">
            <span>Your name</span>
            <input
              value={name}
              onChange={(e) => setName(e.target.value)}
              placeholder="A Reader"
              maxLength={60}
            />
          </label>
          <label className="letters__field">
            <span>Your letter</span>
            <textarea
              value={text}
              onChange={(e) => setText(e.target.value)}
              placeholder="What stayed with you?"
              rows={4}
              maxLength={2000}
              required
            />
          </label>
          {error ? <p className="letters__error">{error}</p> : null}
          <div className="letters__form-actions">
            <button type="submit" disabled={submitting} className="letters__cta">
              {submitting ? 'Sending…' : 'Send letter'}
            </button>
            <button
              type="button"
              className="letters__cancel"
              onClick={() => {
                setShowForm(false)
                setError('')
              }}
            >
              Cancel
            </button>
          </div>
        </form>
      )}

      <ul className="letters__list">
        {comments.length === 0 ? (
          <li className="letters__empty">No letters yet &mdash; be the first to write one in.</li>
        ) : (
          comments
            .slice()
            .reverse()
            .map((c) => (
              <li key={c.id} className="letters__item">
                <div className="letters__meta">
                  <span className="letters__name">{c.name}</span>
                  <span className="letters__date">{formatDate(c.createdAt)}</span>
                </div>
                <p className="letters__text">{c.text}</p>
              </li>
            ))
        )}
      </ul>
    </section>
  )
}
