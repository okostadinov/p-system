package main

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
)

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Security-Policy", "default-src https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/js/bootstrap.min.js 'self'; style-src https://cdn.jsdelivr.net/npm/bootstrap@5.3.2/dist/css/bootstrap.min.css 'self' 'unsafe-inline'; img-src 'self' data:")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.infoLog.Printf("%s - %s %s %s", r.RemoteAddr, r.Proto, r.Method, r.URL.RequestURI())

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !app.isAuthenticated(w, r) {
			http.Redirect(w, r, "/users/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := app.getUserId(w, r)
		if userId == 0 {
			next.ServeHTTP(w, r)
			return
		}

		exists, err := app.users.Exists(userId)
		if err != nil {
			app.serverError(w, err)
			return
		}

		if exists {
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			ctx = context.WithValue(ctx, userIdContextKey, userId)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) verifyAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userId := app.getUserIdFromContext(w, r)
		err := r.ParseForm()
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		targetUserId, err := strconv.Atoi(r.PostForm.Get("user_id"))
		if err != nil {
			app.clientError(w, http.StatusBadRequest)
			return
		}

		if targetUserId != userId {
			app.setFlash(w, r, "Unauthorized action!", FlashTypeDanger)
			http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// options setup for gorilla csrf middleware
func csrfProtect(key string) func(http.Handler) http.Handler {
	return csrf.Protect(
		[]byte(key),
		csrf.Path("/"),
		csrf.SameSite(csrf.SameSiteLaxMode),
	)
}
