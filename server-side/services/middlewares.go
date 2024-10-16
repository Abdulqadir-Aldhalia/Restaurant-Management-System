package services

import (
	"context"
	"fmt"
	"net/http"
	"server-side/model"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("032440aef15d1bab49241d688d2dcb61a4c9072384da9de349baee0eea00d49ca77ca386051ecbe1ac2d3bc24af10ca72ad2256ea49aae28c398c04af868935a0c9f10c63129a602225643d203d5b0ff2166fef2ac56c5e694c1d7584bfc746cac47d4e7e01d2c838ba1e90cead907d3626aee30e887af50b5be1870017fc9bac192d4f5e6eeaa55474bc7ce0b063ce8d9a7bed88a1098c0426ec6d33e39d36115f45c4911fc3a792c8f2f6b2bc828c214724d0e34840050822b7ca732d3eaec0afa535ca62cee6e743b444d18eb6b237bc515849ce8fbee6d1a40a67221b7a1a4ccb16e257287131d1f32161d543b7d5451b4a84100a7451498a390ec50f80413ad3bd281194592e28fee460158ea76f77b0830584506b7d3a68dc3dd3b153126522b533592fd2a2736fa0f2657bc6f1f01aabd1aae44dd8a782baee9d86a9602660788b72e2d29ef9a9a7f0f564a4b5c180c7ba54c1b14c1a8654e26d4fb5c18943542a05785e2d4bca7189a6295b9fe39828acb1545d7b6b227bd8cc915749f6cb43bc784f67176b16adcdc704188284f0442df881c4cd8f9aa7e2101c7eb8600ae5de4683a1e89b502a3e360d19bddc344223a6c2c24a2274122ec9735dba485407e105fe4a00a1e8ca576765539dbe04b504edac20af5f4bd6bce2434695b7f3af82a66e33b3ef2e4e25ca3bc30ef403e58e6d28605d9607e42e30c707c")

type contextKey string

const userContextKey contextKey = "user"

type UserDetails struct {
	UserData  model.User
	UserRoles []string
}

type Claims struct {
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

func GenerateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	var userRolesIds []string
	var userId string

	query, args, err := statement.Select("id").
		From("users").
		Where("email = ?", username).
		ToSql()
	if err != nil {
		return "", fmt.Errorf("error while generating query: %s", err)
	}
	err = db.Get(&userId, query, args...)
	if err != nil {
		return "", fmt.Errorf("error while executing query: %s", err)
	}

	query2, args2, err2 := statement.Select("role_id").
		From("user_roles").
		Where("user_id = ?", userId).
		ToSql()

	if err2 != nil {
		return "", fmt.Errorf("error while creating query: %s", err2)
	}
	err = db.Select(&userRolesIds, query2, args2...)
	if err != nil {
		return "", fmt.Errorf("error while executing query: %s", err)
	}

	roleMap := map[string]string{
		"1": "admin",
		"2": "vendor",
		"3": "customer",
	}

	var userRoles []string
	for _, roleId := range userRolesIds {
		if roleName, exists := roleMap[roleId]; exists {
			userRoles = append(userRoles, roleName)
		}
	}

	fmt.Println(userRoles)

	claims.Roles = userRoles

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Authorize(allowedRoles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			roles := r.Header.Get("roles")
			if roles == "" {
				http.Error(w, "Missing roles", http.StatusForbidden)
				return
			}

			hasRole := false
			for _, role := range allowedRoles {
				if strings.Contains(roles, role) {
					hasRole = true
					break
				}
			}

			if !hasRole {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func AuthenticateJWT(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				http.Error(w, "Invalid token signature", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		var user model.User
		if err = db.Get(&user, "SELECT * FROM users WHERE email = $1", claims.Username); err != nil {
			SendErrorResponse(w, err)
			return
		}

		userDetails := UserDetails{
			UserData:  user,
			UserRoles: claims.Roles,
		}

		r.Header.Set("username", claims.Username)
		r.Header.Set("roles", strings.Join(claims.Roles, ","))

		ctx := context.WithValue(r.Context(), userContextKey, userDetails)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// this one is for Embedded Rust !
func AuthenticateAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		if apiKey == "" || apiKey != "[Q<-(C*V{u/AJim+<qwJ0|~Jus{u',pYJ]vEflDl~sb5LiLx2JA}F,.&cJB'a{u" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
