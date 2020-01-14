package mongo

// Config contains the environment variables requirements to initialize mongodb.
type Config struct {
	URL      string
	Database string
}
