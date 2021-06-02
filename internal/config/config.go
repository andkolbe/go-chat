package config

import "github.com/alexedwards/scs/v2"

// exports to all parts of our application, but doesn't import anything from anywhere else
// only uses packages already built into our standard library
// because it is a struct, we can put anything we need sitewide for our app, and it will be available to every package that imports this package
// our session is initialized in main package but we need to use it in the handlers package. Put it here so it can easily be used in both

type AppConfig struct {
	InProduction bool
	Session *scs.SessionManager
}