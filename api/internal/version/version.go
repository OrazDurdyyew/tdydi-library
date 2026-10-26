package version

// Version is the application version. It is overridden at build time via
// -ldflags "-X library/internal/version.Version=<value>" (see api/Dockerfile),
// which reads the root VERSION file. Defaults to "dev" for local builds.
var Version = "dev"
