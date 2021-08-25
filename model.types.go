package gou

import (
	"github.com/yaoapp/kun/maps"
)

// Relationship types
const (
	RelHasOne         = "hasOne"         // 1 v 1
	RelHasMany        = "hasMany"        // 1 v n
	RelBelongsTo      = "belongsTo"      // inverse  1 v 1 / 1 v n / n v n
	RelHasOneThrough  = "hasOneThrough"  // 1 v 1 ( t1 <-> t2 <-> t3)
	RelHasManyThrough = "hasManyThrough" // 1 v n ( t1 <-> t2 <-> t3)
	RelBelongsToMany  = "belongsToMany"  // 1 v1 / 1 v n / n v n
	RelMorphOne       = "morphOne"
	RelMorphMany      = "morphMany"
	RelMorphToMany    = "morphToMany"
	RelMorphByMany    = "morphByMany"
	RelMorphMap       = "morphMap"
)

// Model 数据模型
type Model struct {
	Name     string
	Source   string
	MetaData MetaData
	Columns  map[string]*Column
}

// MetaData 元数据
type MetaData struct {
	Name      string              `json:"name,omitempty"`      // 元数据名称
	Table     Table               `json:"table,omitempty"`     // 数据表选项
	Columns   []Column            `json:"columns,omitempty"`   // 字段定义
	Indexes   []Index             `json:"indexes,omitempty"`   // 索引定义
	Relations map[string]Relation `json:"relations,omitempty"` // 映射关系定义
	Values    []maps.MapStrAny    `json:"values,omitempty"`    // 初始数值
	Option    Option              `json:"option,omitempty"`    // 元数据配置
}

// Column the field description struct
type Column struct {
	Name        string       `json:"name"`
	Type        string       `json:"type,omitempty"`
	Title       string       `json:"title,omitempty"`
	Description string       `json:"description,omitempty"`
	Comment     string       `json:"comment,omitempty"`
	Length      int          `json:"length,omitempty"`
	Precision   int          `json:"precision,omitempty"`
	Scale       int          `json:"scale,omitempty"`
	Nullable    bool         `json:"nullable,omitempty"`
	Option      []string     `json:"option,omitempty"`
	Default     interface{}  `json:"default,omitempty"`
	DefaultRaw  string       `json:"default_raw,omitempty"`
	Example     interface{}  `json:"example,omitempty"`
	Generate    string       `json:"generate,omitempty"` // Increment, UUID,...
	Crypt       string       `json:"crypt,omitempty"`    // AES, PASSWORD, AES-256, AES-128, PASSWORD-HASH, ...
	Validations []Validation `json:"validations,omitempty"`
	Index       bool         `json:"index,omitempty"`
	Unique      bool         `json:"unique,omitempty"`
	Primary     bool         `json:"primary,omitempty"`
}

// Validation the field validation struct
type Validation struct {
	Title       string `json:"title,omitempty"`
	Description string `json:"description,omitempty"`
	Type        string `json:"type,omitempty"`
	Pattern     string `json:"pattern,omitempty"`
}

// Index the search index struct
type Index struct {
	Comment string   `json:"comment,omitempty"`
	Name    string   `json:"name,omitempty"`
	Columns []string `json:"columns,omitempty"`
	Type    string   `json:"type,omitempty"` // primary,unique,index,match
}

// Table the model mapping table in DB
type Table struct {
	Name        string   `json:"name"`
	Comment     string   `json:"comment,omitempty"`
	Engine      string   `json:"engine,omitempty"` // InnoDB,MyISAM ( MySQL Only )
	Collation   string   `json:"collation"`
	Charset     string   `json:"charset"`
	PrimaryKeys []string `json:"primarykeys"`
}

// Relation the new xun model relation
type Relation struct {
	Type     string `json:"type"`
	Key      string `json:"key,omitempty"`
	Model    string `json:"model"`
	Foreign  string `json:"foreign,omitempty"`
	Parent   string `json:"parent,omitempty"`
	MaxLevel int    `json:"max_level,omitempty"`
}

// Option 模型配置选项
type Option struct {
	Timestamps  bool `json:"timestamps,omitempty"`   // + created_at, updated_at 字段
	SoftDeletes bool `json:"soft_deletes,omitempty"` // + deleted_at 字段
	Trackings   bool `json:"trackings,omitempty"`    // + created_by, updated_by, deleted_by 字段
	Constraints bool `json:"constraints,omitempty"`  // + 约束定义
	Permission  bool `json:"permission,omitempty"`   // + __permission 字段
	Logging     bool `json:"logging,omitempty"`      // + __logging_id 字段
}