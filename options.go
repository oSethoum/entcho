package entcho

import "log"

func WithDB(config *DBConfig) option {
	return func(e *extension) {
		if config == nil {
			config = new(DBConfig)
		}

		if config.Path == "" {
			config.Path = "db"
		}

		if config.Driver == "" {
			config.Driver = SQLite
		} else {
			if !in(config.Driver, []string{MySQL, SQLite, PostgreSQL}) {
				log.Fatalln("driver", config.Driver, "is not supported")
			}
		}

		switch config.Driver {
		case SQLite:
			if config.Dsn == "" {
				config.Dsn = "file:entcho.sqlite?_fk=1&cache=shared"
			}
		case MySQL:
			if config.Dsn == "" {
				config.Dsn = "<user>:<pass>@tcp(<host>:<port>)/<database>?parseTime=True"
			}
		case PostgreSQL:
			if config.Dsn == "" {
				config.Dsn = "host=<host> port=<port> user=<user> dbname=<database> password=<pass>"
			}
		}

		e.data.DBConfig = config
	}
}

func WithFiber(config *FiberConfig) option {
	return func(e *extension) {
		if config == nil {
			config = new(FiberConfig)
		}
		if config.HandlersPath == "" {
			config.HandlersPath = "handlers"
		}
		if config.RoutesPath == "" {
			config.RoutesPath = "routes"
		}

		e.data.FiberConfig = config
	}
}

func WithTS(config *TSConfig) option {
	return func(e *extension) {
		if config == nil {
			config = new(TSConfig)
		}
		if config.TypesPath == "" {
			config.TypesPath = "ts/types"
		}

		if config.ApiPath == "" {
			config.ApiPath = "ts/api"
		}

		e.data.TSConfig = config
	}
}
