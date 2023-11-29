package main

type contextKey string

const (
	isAuthenticatedContextKey = contextKey("isAuthenticated")
	userIdContextKey          = contextKey("userId")
)
