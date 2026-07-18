import BarrierDivider from './BarrierDivider.jsx'
import './StoryContent.css'

// The line the man on the train repeats across the decades - used once as a
// pulled clipping, independent of exact wording in the (editable) body text.
const PULL_QUOTE = '"Read the paper. It\u2019s a good habit."'

function toParagraphs(body) {
  return body
    .split(/\n\s*\n/)
    .map((p) => p.trim())
    .filter(Boolean)
}

export default function StoryContent({ body }) {
  const paragraphs = toParagraphs(body)
  const midpoint = Math.floor(paragraphs.length * 0.42)

  return (
    <article className="story">
      <p className="story__byline">Words by Vishnu &middot; a Chennai&ndash;Kollengode story</p>

      <BarrierDivider label="Boom barrier down &middot; the 7:45 is waiting" />

      <div className="story__columns">
        {paragraphs.map((para, i) => {
          const lines = para.split('\n')
          const isFirst = i === 0
          return (
            <div key={i}>
              <p className={isFirst ? 'story__para story__para--first' : 'story__para'}>
                {lines.map((line, li) => (
                  <span key={li}>
                    {line}
                    {li < lines.length - 1 ? <br /> : null}
                  </span>
                ))}
              </p>
              {i === midpoint && paragraphs.length > 10 ? (
                <blockquote className="story__pullquote">
                  <span>{PULL_QUOTE}</span>
                </blockquote>
              ) : null}
            </div>
          )
        })}
      </div>

      <BarrierDivider label="Barrier lifting &middot; letters to the editor ahead" />
    </article>
  )
}
