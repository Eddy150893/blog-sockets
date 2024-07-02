package middleware

import (
	"net/http"
	"strings"

	"github.com/Eddy150893/blog-sockets/models"
	"github.com/Eddy150893/blog-sockets/server"
	"github.com/golang-jwt/jwt"
)

var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

func shouldCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false
		}
	}
	return true
}

// La funcion recibe un s de tipo server.Server y devuelve una funcion
// La funcion que va a devolver tiene como parametro y valor de retorno un handler
// ya que al ser un middleware recibe el handler, hace la verificacion y de estar todo correcto
// devuelve el handler correspondiente, recordemos que pueden haber mas de un middleware
func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !shouldCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r)
				return
			}
			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
