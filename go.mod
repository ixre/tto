module github.com/ixre/tto

go 1.19

require (
	github.com/fsnotify/fsnotify v1.7.0
	github.com/go-sql-driver/mysql v1.6.0
	github.com/ixre/gof v1.15.10
	github.com/pelletier/go-toml v1.9.5
)

require (
	github.com/gomodule/redigo v1.8.9 // indirect
	github.com/lib/pq v1.10.7 // indirect
	github.com/stretchr/testify v1.8.4 // indirect
	golang.org/x/sys v0.14.0 // indirect
)

replace github.com/ixre/gof => ../gof
