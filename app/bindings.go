package app

// Bindings is an array of keys used to generate
// the environment variable names when loading the configuration.
// For example: "service.name" will be -> "PREFIX_SERVICE_NAME"
//
// Furthermore, we use these bindings for marshalling into the app.Config.
var Bindings = []string{
	"service.name",
	"service.port",
	"store.client",
	"store.host",
	"store.port",
	"store.timeout",
	"store.name",
	"server.timeout.read",
	"server.timeout.write",
	"server.timeout.idle",
	"server.timeout.shutdown",
	"security.allowed.origins",
	"security.allowed.methods",
	"security.allowed.headers",
	"security.allowCredentials",
	"auth.csrf",
}
