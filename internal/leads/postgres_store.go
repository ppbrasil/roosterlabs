package leads

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type postgresStore struct {
	db *sql.DB
}

// NewPostgresStore conecta no Postgres usando DATABASE_URL.
func NewPostgresStore(ctx context.Context, databaseURL string) (Store, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("open postgres: %w", err)
	}
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(4)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping postgres: %w", err)
	}

	return &postgresStore{db: db}, nil
}

func (p *postgresStore) RecordViewEvent(ctx context.Context, ev ViewEvent) error {
	_, err := p.db.ExecContext(ctx, `
		INSERT INTO funnel_events (
			token, event_type, step, language, page_path, payload,
			utm_source, utm_medium, utm_campaign, utm_term, utm_content
		) VALUES ($1, 'view', NULL, $2, $3, NULL, $4, $5, $6, $7, $8)
	`, ev.Token, NormalizeLanguage(ev.Language), ev.PagePath,
		nullIfEmpty(ev.UTM.Source),
		nullIfEmpty(ev.UTM.Medium),
		nullIfEmpty(ev.UTM.Campaign),
		nullIfEmpty(ev.UTM.Term),
		nullIfEmpty(ev.UTM.Content),
	)
	if err != nil {
		return fmt.Errorf("insert view event: %w", err)
	}
	return nil
}

func (p *postgresStore) RecordFormAnswer(ctx context.Context, ev FormAnswerEvent) error {
	payload, err := json.Marshal(ev.Payload)
	if err != nil {
		return fmt.Errorf("marshal answer payload: %w", err)
	}
	_, err = p.db.ExecContext(ctx, `
		INSERT INTO funnel_events (
			token, event_type, step, language, page_path, payload,
			utm_source, utm_medium, utm_campaign, utm_term, utm_content
		) VALUES ($1, 'answer', $2, $3, $4, $5::jsonb, $6, $7, $8, $9, $10)
	`, ev.Token, ev.Step, NormalizeLanguage(ev.Language), ev.PagePath, string(payload),
		nullIfEmpty(ev.UTM.Source),
		nullIfEmpty(ev.UTM.Medium),
		nullIfEmpty(ev.UTM.Campaign),
		nullIfEmpty(ev.UTM.Term),
		nullIfEmpty(ev.UTM.Content),
	)
	if err != nil {
		return fmt.Errorf("insert answer event: %w", err)
	}
	return nil
}

func (p *postgresStore) FinalizeLead(ctx context.Context, in FinalizeInput) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	submitPayload, err := json.Marshal(map[string]string{
		"email":    strings.TrimSpace(in.Email),
		"linkedin": strings.TrimSpace(in.LinkedIn),
	})
	if err != nil {
		return fmt.Errorf("marshal submit payload: %w", err)
	}

	_, err = tx.ExecContext(ctx, `
		INSERT INTO funnel_events (
			token, event_type, step, language, page_path, payload,
			utm_source, utm_medium, utm_campaign, utm_term, utm_content
		) VALUES ($1, 'submit', 5, $2, $3, $4::jsonb, $5, $6, $7, $8, $9)
	`, in.Token, NormalizeLanguage(in.Language), in.PagePath, string(submitPayload),
		nullIfEmpty(in.UTM.Source),
		nullIfEmpty(in.UTM.Medium),
		nullIfEmpty(in.UTM.Campaign),
		nullIfEmpty(in.UTM.Term),
		nullIfEmpty(in.UTM.Content),
	)
	if err != nil {
		return fmt.Errorf("insert submit event: %w", err)
	}

	rows, err := tx.QueryContext(ctx, `
		SELECT DISTINCT ON (step) step, payload
		FROM funnel_events
		WHERE token = $1 AND event_type = 'answer' AND step BETWEEN 1 AND 4
		ORDER BY step, created_at DESC
	`, in.Token)
	if err != nil {
		return fmt.Errorf("query answers: %w", err)
	}
	// rows.Err() (checado abaixo) captura erros de iteração; o erro de Close
	// num result set já consumido não tem ação útil — descartado de propósito.
	defer func() { _ = rows.Close() }()

	answers := map[int]map[string]string{}
	for rows.Next() {
		var step int
		var raw []byte
		if err := rows.Scan(&step, &raw); err != nil {
			return fmt.Errorf("scan answer: %w", err)
		}
		payload := map[string]string{}
		if err := json.Unmarshal(raw, &payload); err != nil {
			return fmt.Errorf("decode answer payload: %w", err)
		}
		answers[step] = payload
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("iterate answers: %w", err)
	}

	for i := 1; i <= 4; i++ {
		if _, ok := answers[i]; !ok {
			return fmt.Errorf("missing answer for step %d", i)
		}
	}

	profile := firstNonEmpty(answers[1]["choice"], answers[1]["other"])
	goal := firstNonEmpty(answers[2]["choice"], answers[2]["other"])
	maturity := answers[3]["choice"]
	challenge := answers[4]["choice"]

	_, err = tx.ExecContext(ctx, `
		INSERT INTO leads (
			token, language, profile, goal, maturity, challenge,
			email, linkedin_url, utm_source, utm_medium, utm_campaign, utm_term, utm_content,
			created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6,
			$7, $8, $9, $10, $11, $12, $13,
			NOW(), NOW()
		)
		ON CONFLICT (token) DO UPDATE SET
			language = EXCLUDED.language,
			profile = EXCLUDED.profile,
			goal = EXCLUDED.goal,
			maturity = EXCLUDED.maturity,
			challenge = EXCLUDED.challenge,
			email = EXCLUDED.email,
			linkedin_url = EXCLUDED.linkedin_url,
			utm_source = EXCLUDED.utm_source,
			utm_medium = EXCLUDED.utm_medium,
			utm_campaign = EXCLUDED.utm_campaign,
			utm_term = EXCLUDED.utm_term,
			utm_content = EXCLUDED.utm_content,
			updated_at = NOW()
	`, in.Token, NormalizeLanguage(in.Language), strings.TrimSpace(profile), strings.TrimSpace(goal), strings.TrimSpace(maturity), strings.TrimSpace(challenge),
		strings.TrimSpace(in.Email), strings.TrimSpace(in.LinkedIn),
		nullIfEmpty(in.UTM.Source),
		nullIfEmpty(in.UTM.Medium),
		nullIfEmpty(in.UTM.Campaign),
		nullIfEmpty(in.UTM.Term),
		nullIfEmpty(in.UTM.Content),
	)
	if err != nil {
		return fmt.Errorf("upsert lead: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}

func (p *postgresStore) Close() error {
	return p.db.Close()
}

func nullIfEmpty(v string) any {
	v = strings.TrimSpace(v)
	if v == "" {
		return nil
	}
	return v
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if s := strings.TrimSpace(v); s != "" {
			return s
		}
	}
	return ""
}
