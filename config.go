package main

type Config struct {
	Dir  string
	Host string
	Port string
}

func NewConfig() Config {
	return Config{
		Host: "localhost",
		Port: "8080",
	}
}

func (c Config) Valid() map[string]string {
	problems := make(map[string]string, 0)
	if c.Dir == "" {
		problems["ErrDirEmpty"] = "dir must be given using the FILEVAULT_DIR env or the -dir flag"
		return problems
	}
	if c.Port == "" {
		problems["ErrPortEmpty"] = "specify the port on which the http server should listen"
		return problems
	}
	if c.Host == "" {
		problems["ErrHostEmpty"] = "host is not specified"
		return problems
	}
	return problems
}
