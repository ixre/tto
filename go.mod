module github.com/ixre/tto

go 1.23.0

toolchain go1.24.2

require (
	github.com/fsnotify/fsnotify v1.9.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/ixre/gof v1.17.11
	github.com/pelletier/go-toml v1.9.5
)

require (
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/sys v0.32.0 // indirect
	golang.org/x/text v0.21.0 // indirect
)

replace github.com/ixre/gof => ../gof
