import './Footer.css'

// The small dot after the year is the only way into the admin panel. It has
// no label, no hover tooltip, and no visual weight beyond a punctuation
// mark - by design, its existence isn't meant to be obvious from the page.
export default function Footer({ onSecretTrigger }) {
  const year = new Date().getFullYear()
  return (
    <footer className="site-footer">
      <div className="site-footer__rule" />
      <p className="site-footer__line">
        &copy; {year} Vishnu. All rights reserved.
        <button
          type="button"
          aria-hidden="true"
          tabIndex={-1}
          className="site-footer__mark"
          onClick={onSecretTrigger}
        >
          &middot;
        </button>
      </p>
    </footer>
  )
}
