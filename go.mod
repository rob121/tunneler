module github.com/rob121/tunneler

go 1.24

require (
	github.com/caseymrm/menuet v1.0.0
	gopkg.in/yaml.v3 v3.0.1
)

require github.com/caseymrm/askm v1.0.0 // indirect

replace github.com/caseymrm/menuet => ./third_party/menuet
