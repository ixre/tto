module github.com/ixre/tto

go 1.19

require (
	github.com/fsnotify/fsnotify v1.7.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/ixre/gof v1.15.21
	github.com/pelletier/go-toml v1.9.5
)

require (
	github.com/gomodule/redigo v1.9.2 // indirect
	github.com/lib/pq v1.10.9 // indirect
	golang.org/x/sys v0.21.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)

replace github.com/ixre/gof => ../gof
