package router

import (
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegister_RoutesDoNotConflict(t *testing.T) {
	gin.SetMode(gin.TestMode)
	engine := gin.New()

	// Gin panics at registration time on conflicting paths; handlers are never invoked here.
	(&Router{}).Register(engine)

	want := map[string]bool{
		"GET /api/products":                        true,
		"POST /api/products":                       true,
		"GET /api/products/:public_id":             true,
		"PATCH /api/products/:public_id":           true,
		"DELETE /api/products/:public_id":          true,
		"POST /api/products/:public_id/variants":   true,
		"GET /api/products/:public_id/variants":    true,
		"GET /api/products/variants/:public_id":    true,
		"PATCH /api/products/variants/:public_id":  true,
		"DELETE /api/products/variants/:public_id": true,
	}
	for _, route := range engine.Routes() {
		delete(want, route.Method+" "+route.Path)
	}
	for missing := range want {
		t.Errorf("route not registered: %s", missing)
	}
}
