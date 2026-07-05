export async function onRequestPost({ request, env }) {
  const json = (body, status = 200) =>
    new Response(JSON.stringify(body), {
      status,
      headers: { 'Content-Type': 'application/json' }
    });

  let data;
  try {
    data = await request.json();
  } catch {
    return json({ ok: false, error: 'invalid json' }, 400);
  }

  // Honeypot: bots fill the hidden field; pretend success.
  if (data.website) return json({ ok: true });

  const email = (data.email || '').trim().toLowerCase();
  if (!/^[^@\s]+@[^@\s]+\.[^@\s]+$/.test(email)) {
    return json({ ok: false, error: 'invalid email' }, 400);
  }

  try {
    await env.DB.prepare(
      `INSERT INTO leads (email, linkedin_url, profile, goal, maturity, challenge, lang, utm, created_at)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?, datetime('now'))`
    )
      .bind(
        email,
        (data.linkedin_url || '').slice(0, 300),
        (data.profile || '').slice(0, 120),
        (data.goal || '').slice(0, 120),
        (data.maturity || '').slice(0, 40),
        (data.challenge || '').slice(0, 40),
        (data.lang || '').slice(0, 10),
        (data.utm || '{}').slice(0, 500)
      )
      .run();
  } catch (e) {
    // Unique-constraint violation = duplicate email: treat as success (idempotent).
    if (String(e).includes('UNIQUE')) return json({ ok: true });
    return json({ ok: false, error: 'storage error' }, 500);
  }

  return json({ ok: true });
}
