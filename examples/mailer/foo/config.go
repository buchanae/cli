package foo

// Server config docs.
type ServerConfig struct {
  private string
  // Server port to listen on. 
  Port int
  // Server host to listen on.
  Host string
}

// User Config docs.
type UserConfig struct {
  // User name for login.
  Username string
  // Password for login.
  Password string
}

// Config docs.
type Config struct {
  // Server Config 2
  ServerConfig
  // User doc2
  User UserConfig
}

func DefaultConfig() Config {
  return Config{}
}
