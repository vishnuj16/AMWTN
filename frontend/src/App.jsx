import { useEffect, useState } from 'react'
import { api, getStoredToken } from './api.js'
import Masthead from './components/Masthead.jsx'
import StoryContent from './components/StoryContent.jsx'
import CommentSection from './components/CommentSection.jsx'
import Footer from './components/Footer.jsx'
import AdminLoginModal from './components/AdminLoginModal.jsx'
import AdminPanel from './components/AdminPanel.jsx'
import './App.css'

export default function App() {
  const [story, setStory] = useState(null)
  const [comments, setComments] = useState([])
  const [loading, setLoading] = useState(true)
  const [loadError, setLoadError] = useState('')

  const [showLogin, setShowLogin] = useState(false)
  const [showAdminPanel, setShowAdminPanel] = useState(false)
  const [adminToken, setAdminToken] = useState(getStoredToken())

  function handleSecretTrigger() {
    if (adminToken) {
      setShowAdminPanel(true)
    } else {
      setShowLogin(true)
    }
  }

  useEffect(() => {
    let cancelled = false
    async function load() {
      try {
        const [storyData, commentsData] = await Promise.all([
          api.getStory(),
          api.getComments(),
        ])
        if (!cancelled) {
          setStory(storyData)
          setComments(commentsData.comments || [])
        }
      } catch (err) {
        if (!cancelled) setLoadError(err.message || 'Could not load the story right now.')
      } finally {
        if (!cancelled) setLoading(false)
      }
    }
    load()
    api.recordView()
    return () => {
      cancelled = true
    }
  }, [])

  if (loading) {
    return (
      <div className="page-status">
        <p>Waiting for the boom barrier to lift…</p>
      </div>
    )
  }

  if (loadError) {
    return (
      <div className="page-status">
        <p>{loadError}</p>
      </div>
    )
  }

  return (
    <>
      <Masthead title={story.title} />
      <StoryContent body={story.body} />
      <CommentSection
        comments={comments}
        onCommentAdded={(c) => setComments((prev) => [...prev, c])}
      />
      <Footer onSecretTrigger={handleSecretTrigger} />

      {showLogin && !adminToken ? (
        <AdminLoginModal
          onClose={() => setShowLogin(false)}
          onLoggedIn={(token) => {
            setAdminToken(token)
            setShowLogin(false)
            setShowAdminPanel(true)
          }}
        />
      ) : null}

      {adminToken && showAdminPanel ? (
        <AdminPanel
          token={adminToken}
          story={story}
          onLogout={() => {
            setAdminToken(null)
            setShowAdminPanel(false)
          }}
          onClose={() => setShowAdminPanel(false)}
          onStoryUpdated={(updated) => setStory((prev) => ({ ...prev, ...updated }))}
        />
      ) : null}
    </>
  )
}
