package apisrever

import (
	"database/sql"
	"net/http"
	"web/casbin"
	"web/internal/app/store/sqlstore"

	"github.com/gorilla/sessions"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)

	if err != nil {
		return err
	}

	defer db.Close()

	store := sqlstore.New(db)
	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	enforcer, _ := casbin.NewCasbin()
	srv := newServer(store, sessionStore, enforcer)

	return http.ListenAndServe(config.BindAddr, srv)
	// m = g(r.sub, p.sub) && (r.obj == p.obj || p.obj == "*") && (r.act == p.act || p.act == "*") && (r.is_premium == p.is_premium || r.is_premium == "*") && timeFunc(r.updated_at) || r.updated_at == "*" || (r.sub == r.creator && r.act == "edit") || (r.sub == r.creator && r.act == "delete") || (g(r.sub, "admin") && r.act == "delete")
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseURL)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
