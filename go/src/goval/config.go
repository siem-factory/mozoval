package main

type Config struct {
	flag_debug	bool
	flag_list	string
	flag_run	string
}

func default_config() Config {
	cfg := Config{
		flag_debug: false,
	}
	return cfg
}
