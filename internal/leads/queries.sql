-- Leads recentes (com segmentacao principal)
SELECT
  created_at,
  language,
  profile,
  goal,
  maturity,
  challenge,
  email,
  linkedin_url,
  utm_source,
  utm_medium,
  utm_campaign
FROM leads
ORDER BY created_at DESC
LIMIT 200;

-- Conversao visita -> submit (total e por idioma)
WITH visits AS (
  SELECT language, COUNT(*)::numeric AS n
  FROM funnel_events
  WHERE event_type = 'view'
  GROUP BY language
), submits AS (
  SELECT language, COUNT(DISTINCT token)::numeric AS n
  FROM funnel_events
  WHERE event_type = 'submit'
  GROUP BY language
)
SELECT
  COALESCE(v.language, s.language) AS language,
  COALESCE(v.n, 0) AS visits,
  COALESCE(s.n, 0) AS submits,
  CASE
    WHEN COALESCE(v.n, 0) = 0 THEN 0
    ELSE ROUND((COALESCE(s.n, 0) / v.n) * 100, 2)
  END AS conversion_pct
FROM visits v
FULL OUTER JOIN submits s ON s.language = v.language
ORDER BY language;

-- Abandono por etapa do formulario
WITH base AS (
  SELECT step, COUNT(DISTINCT token)::numeric AS users
  FROM funnel_events
  WHERE event_type = 'answer'
  GROUP BY step
), first_step AS (
  SELECT COALESCE(users, 0) AS users FROM base WHERE step = 1
)
SELECT
  b.step,
  b.users,
  CASE
    WHEN fs.users = 0 THEN 0
    ELSE ROUND((b.users / fs.users) * 100, 2)
  END AS retained_pct_from_step1
FROM base b
CROSS JOIN first_step fs
ORDER BY b.step;
