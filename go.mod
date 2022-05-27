module github.com/traefik/yaegi

go 1.16

require (
	github.com/jinzhu/copier v0.3.5
	github.com/spf13/cast v1.5.0
)

replace (
	github.com/jinzhu/copier v0.3.5 => github.com/gmemcc/copier v0.3.6-0.20220523071557-031548cb776a
	github.com/spf13/cast v1.5.0 => github.com/gmemcc/cast v1.5.1-0.20220527141822-f9f59eccb45b
)
