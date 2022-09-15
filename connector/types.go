package connector

const (
	// DATABASE the database connector (mysql, pgsql, oracle, sqlite ... )
	DATABASE = iota + 1

	// REDIS the redis connector
	REDIS

	// MONGO the mongodb connector
	MONGO

	// ELASTICSEARCH the elasticsearch connector
	ELASTICSEARCH

	// KAFKA the kafka connector
	KAFKA

	// SCRIPT ? the script connector ( difference with widget ?)
	SCRIPT
)

var types = map[string]int{
	"mysql":         DATABASE,
	"sqlite":        DATABASE,
	"sqlite3":       DATABASE,
	"postgres":      DATABASE,
	"oracle":        DATABASE,
	"redis":         REDIS,
	"mongo":         MONGO,
	"elasticsearch": ELASTICSEARCH,
	"es":            ELASTICSEARCH,
	"kafka":         KAFKA,
	"script":        SCRIPT, // ?
}

// Connector the connector interface
type Connector interface {
	Register(name string, dsl []byte) error
	Is(int) bool
}

// DSL the connector DSL
type DSL struct {
	Name    string                 `json:"-"`
	Type    string                 `json:"type"`
	Label   string                 `json:"label,omitempty"`
	Version string                 `json:"version,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}
