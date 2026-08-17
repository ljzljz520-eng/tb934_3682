package web

const templates = `
{{define "shell-head"}}
<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>{{if .Guide}}{{.Guide.Title}}{{else}}Wedding guide{{end}}</title>
  <style>
    :root { color-scheme: light; --ink:#17212b; --muted:#5d6b75; --paper:#fffdfa; --line:#d9e0dc; --jade:#176b5b; --jade-dark:#0e4e43; --coral:#b55245; --gold:#a66a22; --wash:#eef5f0; }
    * { box-sizing:border-box; }
    body { margin:0; font-family:system-ui,-apple-system,BlinkMacSystemFont,"Segoe UI",sans-serif; color:var(--ink); background:var(--paper); line-height:1.55; }
    a { color:var(--jade-dark); }
    .topbar { border-bottom:1px solid var(--line); background:#f4f8f4; }
    .topbar-inner { max-width:980px; margin:0 auto; display:flex; align-items:center; justify-content:space-between; gap:20px; padding:14px 22px; }
    .brand { font-weight:750; letter-spacing:.02em; text-decoration:none; color:var(--ink); }
    .nav { display:flex; gap:16px; font-size:.9rem; }
    .nav a { text-decoration:none; }
    main { max-width:980px; margin:0 auto; padding:34px 22px 72px; }
    .hero { display:grid; grid-template-columns:1.4fr .9fr; gap:34px; align-items:end; padding:26px 0 42px; }
    .eyebrow { color:var(--coral); font-size:.76rem; text-transform:uppercase; letter-spacing:.15em; font-weight:700; }
    h1 { font-family:Georgia,serif; font-size:clamp(2.2rem,6vw,4.4rem); line-height:1.04; font-weight:500; margin:9px 0 18px; letter-spacing:0; }
    h2 { font-family:Georgia,serif; font-size:1.8rem; font-weight:500; margin:0 0 15px; letter-spacing:0; }
    h3 { margin:0 0 5px; font-size:1.05rem; }
    .lead { color:var(--muted); font-size:1.1rem; max-width:52ch; }
    .hero-note { border-left:4px solid var(--gold); padding:14px 0 14px 18px; color:var(--muted); }
    .section { border-top:1px solid var(--line); padding:34px 0; }
    .section-grid { display:grid; grid-template-columns:1fr 1fr; gap:28px; }
    .schedule { display:grid; gap:12px; }
    .schedule-item { display:grid; grid-template-columns:120px 1fr; gap:16px; padding:15px 0; border-bottom:1px solid var(--line); }
    .time { color:var(--jade-dark); font-weight:700; font-size:.9rem; }
    .meta { color:var(--muted); font-size:.9rem; }
    .address { font-style:normal; color:var(--muted); }
    .tip { background:var(--wash); padding:22px; border-radius:6px; }
    .tip strong { color:var(--jade-dark); }
    .actions { display:flex; flex-wrap:wrap; gap:10px; margin-top:20px; }
    .button { display:inline-flex; align-items:center; justify-content:center; min-height:42px; padding:9px 16px; border-radius:4px; border:1px solid var(--jade); background:var(--jade); color:white; text-decoration:none; font-weight:700; }
    .button.secondary { background:transparent; color:var(--jade-dark); }
    .hint { margin:20px 0 0; padding:14px 16px; border-left:4px solid var(--coral); background:#fff1ed; color:#6e3029; }
    .form { display:grid; gap:14px; max-width:640px; }
    label { display:grid; gap:6px; font-weight:650; }
    input, textarea { font:inherit; border:1px solid #b9c6c0; border-radius:4px; padding:10px 12px; color:var(--ink); background:#fff; }
    textarea { min-height:130px; resize:vertical; }
    .admin-header { display:flex; justify-content:space-between; align-items:start; gap:20px; }
    .badge { font-size:.76rem; text-transform:uppercase; letter-spacing:.12em; color:var(--jade-dark); border:1px solid #9ab8ad; padding:5px 8px; border-radius:3px; white-space:nowrap; }
    .audit-list { display:grid; gap:0; border-top:1px solid var(--line); }
    .audit-row { display:grid; grid-template-columns:190px 180px 1fr; gap:14px; padding:13px 0; border-bottom:1px solid var(--line); font-size:.9rem; }
    .footer { max-width:980px; margin:0 auto; padding:24px 22px 40px; color:var(--muted); font-size:.82rem; }
    @media (max-width:700px) { .topbar-inner { align-items:flex-start; flex-direction:column; gap:8px; } .hero,.section-grid { grid-template-columns:1fr; gap:18px; } main { padding-top:20px; } .schedule-item { grid-template-columns:1fr; gap:3px; } .admin-header { flex-direction:column; } .audit-row { grid-template-columns:1fr; gap:3px; } .actions { flex-direction:column; align-items:stretch; } .button { width:100%; } }
  </style>
</head>
<body>
  <header class="topbar"><div class="topbar-inner"><a class="brand" href="/">Willow House / Lin &amp; Morgan</a><nav class="nav"><a href="/admin/{{if .Guide}}{{.Guide.ID}}{{else}}{{.GuideID}}{{end}}">Editor preview</a><a href="/audits/{{if .Guide}}{{.Guide.ID}}{{else}}{{.GuideID}}{{end}}">Activity</a></nav></div></header>
{{end}}
{{define "shell-foot"}}
  <footer class="footer">Made for a calm, well-informed wedding day.</footer>
</body>
</html>
{{end}}

{{define "guide"}}{{template "shell-head" .}}
<main>
  <section class="hero"><div><div class="eyebrow">The wedding of</div><h1>{{.Guide.Couple}}</h1><p class="lead">{{.Guide.Welcome}}</p><div class="actions">{{range .Guide.Links}}<a class="button {{if eq .Kind "seating"}}secondary{{end}}" href="/action/{{$.Guide.ID}}/{{.ID}}?visitor={{$.VisitorKey}}">{{.Label}}</a>{{end}}</div>{{if .AtLimit}}<p class="hint">{{.Hint}}</p>{{end}}</div><div class="hero-note"><strong>Today at a glance</strong><br>{{.Guide.Venue.Name}}<br>{{.Guide.Venue.City}}, {{.Guide.Venue.Country}}<br><span class="meta">Visit {{.Visitor.VisitCount}} of the welcome guide</span></div></section>
  <section class="section"><h2>Schedule</h2><div class="schedule">{{range .Guide.Schedule}}<article class="schedule-item"><div class="time">{{.StartsAt}}<br><span class="meta">{{.EndsAt}}</span></div><div><h3>{{.Title}}</h3><div class="meta">{{.Location}}</div><p>{{.Details}}</p></div></article>{{end}}</div></section>
  <section class="section section-grid"><div><h2>Getting there</h2><address class="address"><strong>{{.Guide.Venue.Name}}</strong><br>{{.Guide.Venue.Line1}}<br>{{.Guide.Venue.Line2}}<br>{{.Guide.Venue.City}}, {{.Guide.Venue.Region}} {{.Guide.Venue.PostalCode}}<br>{{.Guide.Venue.Country}}</address></div><div class="tip"><h2>What to wear</h2><p><strong>{{.Guide.Attire.Summary}}</strong></p><p>{{.Guide.Attire.Description}}</p><p class="meta">Palette: {{.Guide.Attire.ColorHint}}<br>{{.Guide.Attire.WeatherNote}}</p></div></section>
  <section class="section"><h2>A note for you</h2><p class="lead">Use the buttons above for the practical details, and leave a few words for the couple when you have a quiet moment.</p></section>
</main>
{{template "shell-foot" .}}{{end}}

{{define "admin"}}{{template "shell-head" .}}
<main><div class="admin-header"><div><div class="eyebrow">Private editor preview</div><h1>{{.Guide.Title}}</h1><p class="lead">Revision {{.Guide.Revision}} · {{.ScheduleCount}} schedule items · {{.LinkCount}} buttons</p></div><span class="badge">{{if .Guide.Published}}Published{{else}}Draft{{end}}</span></div>
  <section class="section"><form class="form" method="post"><label>Page title<input name="title" value="{{.Guide.Title}}" required></label><label>Welcome message<textarea name="welcome" required>{{.Guide.Welcome}}</textarea></label><label>Attire summary<input name="attire" value="{{.Guide.Attire.Summary}}" required></label><label>Editor name<input name="actor" value="planner" required></label><label><input type="checkbox" name="publish" value="yes"> Publish after saving</label><button class="button" type="submit">Save changes</button></form></section>
  <section class="section"><h2>Preview stays private</h2><p class="meta">This screen reads the draft directly and does not increment a visitor record.</p><a class="button secondary" href="/guide/{{.Guide.ID}}?visitor=preview-reader">Open guest page</a></section></main>
{{template "shell-foot" .}}{{end}}

{{define "blessing"}}{{template "shell-head" .}}
<main><div class="eyebrow">A few words</div><h1>Leave a blessing</h1><p class="lead">Your note will be saved with the wedding memories.</p><section class="section"><form class="form" method="post"><input type="hidden" name="visitor" value="guest"><label>Your name<input name="name" required></label><label>Your message<textarea name="message" minlength="8" required></textarea></label><button class="button" type="submit">Send blessing</button></form></section></main>
{{template "shell-foot" .}}{{end}}

{{define "audits"}}{{template "shell-head" .}}
<main><div class="eyebrow">Editor activity</div><h1>Guide timeline</h1><form class="form" method="get"><label>Filter activity<input name="q" value="{{.Term}}" placeholder="publish, guest, import"></label><button class="button secondary" type="submit">Filter</button></form><section class="section"><div class="audit-list">{{range .Entries}}<div class="audit-row"><span class="meta">{{.CreatedAt}}</span><strong>{{.Action}}</strong><span>{{.Entity}} · {{.EntityID}} · {{.Detail}}</span></div>{{else}}<p class="meta">No activity yet.</p>{{end}}</div></section></main>
{{template "shell-foot" .}}{{end}}
`
