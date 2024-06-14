package enrichstorage

import (
	"github.com/gin-gonic/gin"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
)

func getAuthorizationMiddleware(config *Config) gin.HandlerFunc {
	keycloakConfig := ginkeycloak.BuilderConfig{
		Service: "enricher",
		Url:     "<your token url>",
		Realm:   "<your realm to get the public keys>",
	}
	return ginkeycloak.NewAccessBuilder(keycloakConfig).
		RestrictButForRole("user").
		Build()
}
