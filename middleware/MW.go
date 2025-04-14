package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/untrik/FromSkateToZOH/database"
	"github.com/untrik/FromSkateToZOH/models"
)

var secretKey []byte

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	roleKey   contextKey = "role"
)

func InitSecretKey() {
	secretKey = []byte(os.Getenv("JWT_SECRET"))
	if len(secretKey) == 0 {
		panic("JWT_SECRET is not set")
	}
}
func getUserRole(userID uint) (string, error) {
	role := "unknown"
	var admin models.Admin
	var student models.Student

	errAdmin := database.DB.Where("user_id = ?", userID).First(&admin).Error
	if errAdmin == nil {
		role = "admin"
		log.Printf("User %d is admin", userID)
	} else {
		log.Printf("Admin check error: %v", errAdmin)
	}
	errStudent := database.DB.Where("user_id = ?", userID).First(&student).Error
	if errStudent == nil {
		role = "student"
		log.Printf("User %d is student", userID)
	} else {
		log.Printf("Admin check error: %v", errAdmin)
	}
	return role, nil

}
func GenerateJWT(userID uint) (string, error) {
	role := "unknown"
	role, err := getUserRole(userID)
	if err != nil {
		return "", fmt.Errorf("role detection failed: %w", err)
	}
	claims := jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     jwt.NewNumericDate(time.Now().Add(12 * time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secretKey)
}
func JWTMiddlewareAdmin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		exp, err := claims.GetExpirationTime()
		if err != nil || exp == nil || exp.Before(time.Now()) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		if claims["role"] != "admin" {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, uint(userID))
		ctx = context.WithValue(ctx, roleKey, claims["role"])
		r = r.WithContext(ctx)

		next(w, r)
	}
}
func JWTMiddlewareStudent(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		exp, err := claims.GetExpirationTime()
		if err != nil || exp == nil || exp.Before(time.Now()) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		if claims["role"] != "student" {
			http.Error(w, "Access denied", http.StatusForbidden)
			return
		}

		userID, ok := claims["user_id"].(float64)
		if !ok {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), UserIDKey, uint(userID))
		ctx = context.WithValue(ctx, roleKey, claims["role"])
		r = r.WithContext(ctx)

		next(w, r)
	}
}
