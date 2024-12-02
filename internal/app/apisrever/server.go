package apisrever

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
	"web/internal/app/apisrever/utils"
	"web/internal/app/model"
	"web/internal/app/roles"
	"web/internal/app/store"

	"github.com/casbin/casbin/v2"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "unsosik"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router        *mux.Router
	logger        *logrus.Logger
	store         store.Store
	sessionsStore sessions.Store
	enforcer      *casbin.Enforcer
}

func newServer(store store.Store, sessionsStore sessions.Store, enforcer *casbin.Enforcer) *server {
	s := &server{
		router:        mux.NewRouter(),
		logger:        logrus.New(),
		store:         store,
		sessionsStore: sessionsStore,
		enforcer:      enforcer,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))
	//all person can access this router

	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("internal/app/apisrever/static/"))))

	s.router.HandleFunc("/users", s.handleUsersCreate()).Methods("POST", "GET")
	s.router.HandleFunc("/sessions", s.handleSessionsCreate()).Methods("POST", "GET")

	s.router.HandleFunc("/profile", s.handleUnAuthProfile()).Methods("GET")

	s.router.HandleFunc("/logout", s.handleLogout()).Methods("POST", "GET")

	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "internal/app/apisrever/templates/index.html") // сделать так
	}).Methods("GET")

	// only for: /private/***
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
	private.HandleFunc("/profile", s.handleProfile()).Methods("GET", "POST")

	private.HandleFunc("/roles", s.getRoles).Methods("GET")
	private.HandleFunc("/alltopics", s.handleFindAll()).Methods("GET")
	private.HandleFunc("/topic/{id:[0-9]+}", s.handleTopic()).Methods("GET", "POST")
}

func (s *server) handleTopic() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ctxKeyUser).(*model.User)
		if !ok || user == nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		vars := mux.Vars(r)
		topicID, err := strconv.Atoi(vars["id"])

		if err != nil {
			s.error(w, r, http.StatusBadRequest, errors.New("invalid topic ID"))
			return
		}

		switch r.Method {
		case http.MethodGet:
			s.getTopicById(w, r, topicID)
		case http.MethodPost:

			s.updateTopicById(w, r, topicID)
		default:
			s.error(w, r, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		}

	}
}

func (s *server) updateTopicById(w http.ResponseWriter, r *http.Request, topicID int) {
	topic, err := s.store.Topic().FindByID(topicID)

	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			s.error(w, r, http.StatusNotFound, errors.New("topic not found"))
			return
		}
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	user, _ := r.Context().Value(ctxKeyUser).(*model.User)

	fmt.Println("user.ID: ", user.ID)
	fmt.Println("topicID", topicID)
	allowed, err := s.enforcer.Enforce(strconv.Itoa(user.ID), "topic", "edit", strconv.Itoa(topic.UserID))

	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}
	if !allowed {
		s.error(w, r, http.StatusForbidden, errors.New("permission denied on update"))
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&topic); err != nil { // через форму сделать
		s.error(w, r, http.StatusBadRequest, errors.New("invalid request payload"))
		return
	}

	if err := s.store.Topic().UpdateTopic(topic); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	s.respond(w, r, http.StatusOK, topic)

}

func (s *server) getTopicById(w http.ResponseWriter, r *http.Request, topicID int) {
	topic, err := s.store.Topic().FindByID(topicID)

	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			s.error(w, r, http.StatusNotFound, errors.New("topic not found"))
			return
		}
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	tmpl, err := template.ParseFiles("internal/app/apisrever/templates/topic.html")

	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err = tmpl.Execute(w, topic); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	// http.ServeFile(w, r, "internal/app/apisrever/templates/reg.html") // сделать так

	// s.respond(w, r, http.StatusOK, topic)

}

func (s *server) handleUnAuthProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, 123)
	}
}

func (s *server) handleFindAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ctxKeyUser).(*model.User)

		if !ok || user == nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		topics, err := s.store.Topic().FindAll()

		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, topics)
	}
}

func (s *server) getRoles(w http.ResponseWriter, r *http.Request) {
	roles, err := s.enforcer.GetAllRoles()

	if err != nil {
		log.Fatal(err)
		return
	}

	users, err := s.enforcer.GetAllSubjects()

	if err != nil {
		log.Fatal(err)
		return
	}
	data := map[string]interface{}{
		"roles": roles,
		"users": users,
	}

	println(data)

	// s.respond(w, r, http.StatusOK, data) // здесь как то сделать так чтобы были все роли и юзеры их

	s.respond(w, r, http.StatusOK, roles)
}

func (s *server) handleProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ctxKeyUser).(*model.User)
		if !ok || user == nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		switch r.Method {
		case http.MethodGet:
			s.renderProfilePage(w, r, user)
		case http.MethodPost:
			// Если разрешение есть, создаем топик
			if s.enforcer == nil {
				fmt.Println("s.enforcer nil")
			}
			allowed, err := s.enforcer.Enforce(strconv.Itoa(user.ID), "topic", "create", "*")
			if err != nil {
				fmt.Println("ERR ERR")
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
			if !allowed {
				fmt.Println("ALLOWED ALLOWED ALLOWED")

				s.error(w, r, http.StatusForbidden, errors.New("permission denied on create"))
				return
			}
			s.createTopic(w, r, user)
		default:
			s.error(w, r, http.StatusMethodNotAllowed, errors.New("method not allowed"))
		}
	}
}

func (s *server) renderProfilePage(w http.ResponseWriter, r *http.Request, user *model.User) {
	tmpl, err := template.ParseFiles("internal/app/apisrever/templates/auth.html")
	if err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	if err := tmpl.Execute(w, user); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
	}
}

func (s *server) createTopic(w http.ResponseWriter, r *http.Request, user *model.User) {
	m, err := utils.ParseFormFields(r, []string{"topicname", "topicdescription", "isprivate", "topicabout", "topiccategory"})
	if err != nil {
		s.error(w, r, http.StatusBadRequest, err)
		return
	}

	if m["topicname"] == "" || m["topicdescription"] == "" {
		s.error(w, r, http.StatusBadRequest, errors.New("topic name and description cannot be empty"))
		return
	}

	isPrivate := false
	if m["isprivate"] == "on" {
		isPrivate = true
	}

	topic := &model.Topic{
		UserID:      user.ID,
		TopicName:   m["topicname"],
		Description: m["topiccategory"],
		Visibility:  isPrivate,
		Content:     m["topicabout"],
	}

	fmt.Println(*topic)

	if err := s.store.Topic().Create(topic); err != nil {
		s.error(w, r, http.StatusInternalServerError, err)
		return
	}

	// _, err = s.enforcer.AddPolicy(
	// 	strconv.Itoa(user.ID),
	// 	"topic",
	// 	"edit",
	// 	strconv.Itoa(topic.UserID),
	// )
	http.Redirect(w, r, "/private/profile", http.StatusSeeOther)
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remore_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestID),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()

		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start))
	})
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionsStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id.(int))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u))) //чтобы при повторном запросе снова не проверял

	})
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) handleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet { // также примиать post
			// w.Header().Set("Content-Type", "text/html")
			http.ServeFile(w, r, "internal/app/apisrever/templates/reg.html") // сделать так

			return
		}

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			email := r.FormValue("email")
			password := r.FormValue("password")

			if email == "" || password == "" {
				s.error(w, r, http.StatusBadRequest, fmt.Errorf("email or password cannot be empty"))
				return
			}

			// req := &request{}
			// if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			// 	s.error(w, r, http.StatusBadRequest, err)
			// 	return
			// }

			fmt.Println(email)
			fmt.Println(password)

			u := &model.User{
				Email:    email,
				Password: password,
			}

			if err := s.store.User().Create(u); err != nil {
				s.error(w, r, http.StatusUnprocessableEntity, err)
				fmt.Println("Здесь ПРОБЕЛМА")
				return
			}

			u.Sanitaze() // ответ без пароля юзера
			// s.respond(w, r, http.StatusCreated, u) // тут сделать перенапрвление на вход

			fmt.Println(u.Email)
			// if err := addRoleForUser(u.Email, roles.Editor, s.enforcer); err != nil {
			// 	s.error(w, r, http.StatusInternalServerError, err)
			// 	return
			// }

			http.Redirect(w, r, "/sessions", http.StatusSeeOther)
		}
	}
}

func addRoleForUser(name, role string, e *casbin.Enforcer) error {
	_, err := e.AddGroupingPolicy(name, role)

	if err != nil {
		fmt.Println("Error in addRoleForUser")

		return err
	}

	return nil
}

func (s *server) handleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {

		if r.Method == http.MethodGet {
			http.ServeFile(w, r, "internal/app/apisrever/templates/login.html") // сделать так

			return
		} else if r.Method == http.MethodPost {

			requiredFields := []string{"email", "password"}

			formData, err := utils.ParseFormFields(r, requiredFields)

			if err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}

			email := formData["email"]
			password := formData["password"]

			// req := &request{}
			// if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			// 	s.error(w, r, http.StatusBadRequest, err)
			// 	return
			// }

			u, err := s.store.User().FindByEmail(email)
			if err != nil || !u.CompatePassword(password) {
				s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
				return
			}

			session, err := s.sessionsStore.Get(r, sessionName)
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			session.Values["user_id"] = u.ID
			if err := s.sessionsStore.Save(r, w, session); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			// s.respond(w, r, http.StatusOK, nil) // тут можно сделать перенаправление на профиль пользователя

			// err = addRoleForUser("name", s.enforcer)

			// if err != nil {
			// 	log.Fatal(err)
			// 	return
			// }

			if err := addRoleForUser(strconv.Itoa(u.ID), roles.Editor, s.enforcer); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}

			http.Redirect(w, r, "/private/profile", http.StatusSeeOther)
		}

	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

// func (s *server) handleLogout() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		http.SetCookie(w, &http.Cookie{
// 			Name:     "unsosik",
// 			Value:    "",
// 			MaxAge:   -1,
// 			HttpOnly: true,
// 		})
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Successfully logged out"))
// 	}
// }

func (s *server) handleLogout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionsStore.Get(r, sessionName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values = make(map[interface{}]interface{})

		if err := s.sessionsStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     sessionName,
			Value:    "",
			MaxAge:   -1,
			HttpOnly: true,
		})

		http.Redirect(w, r, "/sessions", http.StatusSeeOther)
	}
}
