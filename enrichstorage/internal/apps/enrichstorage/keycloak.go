package enrichstorage

import (
	"github.com/gin-gonic/gin"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
)

func getAuthorizationMiddleware(config *Config) gin.HandlerFunc {
	keycloakConfig := ginkeycloak.BuilderConfig{
		Service: config.Keycloak.ClientId,
		Url:     config.Keycloak.Url,
		Realm:   config.Keycloak.Realm,
	}
	builder := ginkeycloak.NewAccessBuilder(keycloakConfig)
	for _, role := range config.Keycloak.RolesToPermit {
		builder = builder.RestrictButForRole(role)
	}
	return builder.Build()
}
