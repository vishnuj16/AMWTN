import './Masthead.css'

export default function Masthead({ title }) {
  const today = new Date().toLocaleDateString('en-IN', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })

  return (
    <header className="masthead">
      <div className="masthead__eyebrow-row">
        <span>Kollengode &amp; Elsewhere</span>
        <span>{today}</span>
        <span>Price: One Read</span>
      </div>
      <h1 className="masthead__title">{title}</h1>
      <p className="masthead__deck">
        A short story about a boy, a level crossing, and a stranger who never gave his name —
        only the paper, and one instruction.
      </p>
      <div className="masthead__rule masthead__rule--thick" />
      <div className="masthead__rule masthead__rule--thin" />
    </header>
  )
}
