module server

go 1.15

require shenguanchu.com/common v0.0.0-incompatible

replace shenguanchu.com/common => ../common

require (
	github.com/garyburd/redigo v1.6.2
	shenguanchu.com/client v0.0.0-incompatible
)

replace shenguanchu.com/client => ../client
