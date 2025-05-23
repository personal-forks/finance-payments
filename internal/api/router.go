package api

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/formancehq/go-libs/v3/api"
	"github.com/formancehq/go-libs/v3/auth"
	"github.com/formancehq/go-libs/v3/health"
	"github.com/formancehq/payments/internal/api/backend"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Version struct {
	Version int
	Builder func(backend backend.Backend, a auth.Authenticator, debug bool) *chi.Mux
}

type versionsSlice []Version

func (v versionsSlice) Len() int {
	return len(v)
}

func (v versionsSlice) Less(i, j int) bool {
	return v[i].Version < v[j].Version
}

func (v versionsSlice) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func NewRouter(
	backend backend.Backend,
	info api.ServiceInfo,
	healthController *health.HealthController,
	a auth.Authenticator,
	debug bool,
	versions ...Version) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			handler.ServeHTTP(w, r)
		})
	})
	r.Get("/_healthcheck", healthController.Check)
	r.Get("/_info", api.InfoHandler(info))

	sortedVersions := versionsSlice(versions)
	sort.Stable(sortedVersions)

	for _, version := range sortedVersions[1:] {
		prefix := fmt.Sprintf("/v%d", version.Version)
		r.Handle(prefix+"/*", http.StripPrefix(prefix, version.Builder(backend, a, debug)))
	}

	r.Handle("/*", versions[0].Builder(backend, a, debug)) // V1 and V2 have no prefix

	return r
}
