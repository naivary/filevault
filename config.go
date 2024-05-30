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
		Dir:  "/mnt/filevault",
	}
}
