import './BarrierDivider.css'

/**
 * The recurring visual signature of the page: a level-crossing boom barrier,
 * lifting to let the reader through. Used sparingly at two thresholds -
 * entering the story, and crossing from the story into the letters section.
 */
export default function BarrierDivider({ label }) {
  return (
    <div className="barrier" role="presentation">
      <div className="barrier__post barrier__post--left" />
      <div className="barrier__post barrier__post--right" />
      <div className="barrier__arm-track">
        <div className="barrier__arm">
          <svg viewBox="0 0 400 26" preserveAspectRatio="none" className="barrier__stripes">
            <rect width="400" height="26" fill="#ede6d6" />
            {Array.from({ length: 10 }).map((_, i) => (
              <polygon
                key={i}
                points={`${i * 44},0 ${i * 44 + 22},0 ${i * 44 + 2},26 ${i * 44 - 20},26`}
                fill={i % 2 === 0 ? '#a63b2e' : '#1c1a17'}
              />
            ))}
          </svg>
        </div>
      </div>
      {label ? <span className="barrier__label">{label}</span> : null}
    </div>
  )
}
