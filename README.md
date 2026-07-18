# The Man Who Threw Newspapers — story website

A single-page, mobile-responsive site for your short story, with reader
comments ("Letters to the Editor"), a hidden editor login, and page-view
stats — built to run entirely on free tiers.

- **Frontend:** React (Vite), no heavy libraries, ~50KB gzipped
- **Backend:** Go, deployed as Vercel serverless functions
- **Storage:** Upstash Redis (free tier) — stores comments, the manuscript, view counts, and admin sessions
- **Hosting:** Vercel (free tier)

Total monthly cost: **$0**, at the traffic levels a short-story page normally gets.

---

## 1. How it's put together

```
newspaper-site/
├── frontend/          React app (everything a reader sees)
├── api/                One folder per Go serverless function:
│   ├── story/           GET  /api/story          - the manuscript
│   ├── comments/         GET/POST /api/comments   - reader letters
│   ├── view/              POST /api/view          - view counter
│   └── admin/
│       ├── login/         POST /api/admin/login   - password check
│       ├── comments/      GET/DELETE (moderation)
│       ├── story/         POST (replace manuscript)
│       └── stats/         GET (views + comment count)
├── lib/shared/         Shared Go helpers (Redis client, auth, sanitizing)
└── vercel.json
```

Nothing is stored on disk (Vercel's functions don't have persistent disk
storage) — everything the admin can change lives in Upstash Redis. The
original manuscript is also baked into the Go binary as a fallback, so the
story is always readable even before you've touched the admin panel.

## 2. One-time setup

### a) Create a free Upstash Redis database
1. Go to [upstash.com](https://upstash.com) and sign up (free, no card needed).
2. Create a new **Redis** database (any nearby region).
3. On the database page, open the **REST API** section and copy the
   `UPSTASH_REDIS_REST_URL` and `UPSTASH_REDIS_REST_TOKEN` values.

### b) Push this project to GitHub
```bash
cd newspaper-site
git init
git add .
git commit -m "Initial site"
git branch -M main
git remote add origin https://github.com/<your-username>/<repo>.git
git push -u origin main
```

### c) Deploy on Vercel
1. Go to [vercel.com](https://vercel.com) and sign up with your GitHub account (free).
2. **Add New → Project**, pick the repo you just pushed.
3. Vercel should detect the `vercel.json` and use it automatically — you
   don't need to change the framework preset.
4. Before the first deploy, open **Environment Variables** and add:
   - `UPSTASH_REDIS_REST_URL`
   - `UPSTASH_REDIS_REST_TOKEN`
   - `ADMIN_PASSWORD` — a password only you know
5. Click **Deploy**. In a minute or two you'll get a live URL like
   `https://your-story.vercel.app`.
6. (Optional) Add a custom domain under **Project Settings → Domains** —
   still free on Vercel's plan, you'd just pay your domain registrar.

That's it — the site is live, comments work, and view counts are ticking up.

## 3. Using the hidden editor login

There is **no visible "Admin" or "Login" button anywhere on the page.**
The way in:

> At the very bottom of the page, in the copyright line, there's a small
> `·` (a period-like dot) right after "All rights reserved." It looks like
> punctuation. Click it. That opens the password prompt.

Only you know it's there. Log in with the `ADMIN_PASSWORD` you set in
Vercel, and you'll see the **Editor's Desk** panel:

- **Readership** — total page views and letters received.
- **Letters** — every comment, with a **Delete** button for anything you
  decide crosses a line. Deletions are permanent.
- **Manuscript** — a text box with the current title and full story text.
  Edit it and press **Save & publish** to replace what every visitor sees,
  immediately, no redeploy needed. Leave a blank line between paragraphs,
  same as the original file.

Sessions last 12 hours, then you're logged out automatically. "Log out"
ends it immediately; "Close" just hides the panel without logging out.

## 4. Local development (optional)

You'll need the [Vercel CLI](https://vercel.com/docs/cli) to run the Go
functions locally:

```bash
npm i -g vercel
vercel link          # connect this folder to your Vercel project
vercel env pull .env # pulls your env vars into a local .env file
vercel dev           # runs the Go API on http://localhost:3000
```

In a second terminal:
```bash
cd frontend
npm install
npm run dev           # runs React on http://localhost:5173, proxying /api to :3000
```

## 5. Changing the password later

Update `ADMIN_PASSWORD` in Vercel's Environment Variables and redeploy
(Vercel → Deployments → ⋯ → Redeploy). Existing logged-in sessions stay
valid until they expire naturally (up to 12 hours) — the old password
stops working for new logins immediately.

## 6. A note on how "hidden" the login really is

The password check happens on the server, so no one can log in without it
— that part is real security. The *location* of the login button, though,
is just obscurity: anyone who reads the page's source code could find it.
That's an acceptable trade-off for a personal story site, but it isn't the
same as it being truly undiscoverable to a determined visitor.
