package gou

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/yaoapp/kun/any"
	"github.com/yaoapp/kun/maps"
	"github.com/yaoapp/kun/utils"
	"github.com/yaoapp/xun/capsule"
)

func TestLoadModel(t *testing.T) {
	source := "file://" + path.Join(TestModRoot, "user.json")
	user := LoadModel(source, "user")
	assert.Equal(t, user.MetaData.Name, "用户")
	assert.Equal(t, user.Name, "user")
	assert.Equal(t, user.Source, source)
}

func TestModelReload(t *testing.T) {
	user := Select("user")
	user.Reload()
	assert.Equal(t, user.MetaData.Name, "用户")
	assert.Equal(t, user.Name, "user")
}

func TestModelMigrate(t *testing.T) {
	for name, mod := range Models {
		utils.Dump(name)
		mod.Migrate(true)
	}
}

func TestModelMustFind(t *testing.T) {
	user := Select("user").MustFind(1, QueryParam{})
	assert.Equal(t, user.Get("mobile"), "13900001111")
	assert.Equal(t, user.Dot().Get("extra.sex"), "男")
}

func TestModelMustFindWiths(t *testing.T) {
	user := Select("user").MustFind(1,
		QueryParam{
			Withs: map[string]With{
				"manu":      {},
				"addresses": {},
				"roles":     {}, // 暂未实现（ 下一版支持 )
				"mother": {
					Query: QueryParam{ // 数据归集存在BUG（ 下一版修复 )
						Withs: map[string]With{
							// "addresses": {},
							// "manu": {},
						},
					},
				},
			},
		})

	userDot := user.Dot()
	assert.Equal(t, userDot.Get("mobile"), "13900001111")
	assert.Equal(t, userDot.Get("extra.sex"), "男")
	assert.Equal(t, userDot.Get("manu.name"), "北京云道天成科技有限公司")
	assert.Equal(t, userDot.Get("addresses.0.location"), "北京国家数字出版基地A103")
	assert.Equal(t, userDot.Get("mother.extra.sex"), "女")
	assert.Equal(t, userDot.Get("mother.friends.friend_id"), int64(2))
	assert.Equal(t, userDot.Get("mother.friends.type"), "monther")
}

func TestModelMustSearch(t *testing.T) {
	user := Select("user").MustSearch(QueryParam{}, 1, 2)
	userDot := user.Dot()
	assert.Equal(t, userDot.Get("total"), 3)
	assert.Equal(t, userDot.Get("next"), 2)
	assert.Equal(t, userDot.Get("page"), 1)
	assert.Equal(t, userDot.Get("data.0.id"), int64(1))
	assert.Equal(t, userDot.Get("data.1.id"), int64(2))
}

func TestModelMustSearchWiths(t *testing.T) {
	user := Select("user").MustSearch(QueryParam{
		Withs: map[string]With{
			"manu":      {},
			"addresses": {},
			"mother":    {},
		},
	}, 1, 2)
	userDot := user.Dot()
	assert.Equal(t, userDot.Get("total"), 3)
	assert.Equal(t, userDot.Get("next"), 2)
	assert.Equal(t, userDot.Get("page"), 1)
	assert.Equal(t, userDot.Get("data.0.id"), int64(1))
	assert.Equal(t, userDot.Get("data.0.manu.name"), "北京云道天成科技有限公司")
	assert.Equal(t, userDot.Get("data.0.mother.extra.sex"), "女")
	assert.Equal(t, userDot.Get("data.0.extra.sex"), "男")
	assert.Equal(t, userDot.Get("data.0.addresses.0.location"), "北京国家数字出版基地A103")
	assert.Equal(t, userDot.Get("data.1.id"), int64(2))
}

func TestModelMustSearchWithsWhere(t *testing.T) {
	user := Select("user").MustSearch(QueryParam{
		Wheres: []QueryWhere{
			{
				Column: "mobile",
				Value:  "13900001111",
			},
		},
		Withs: map[string]With{
			"manu":      {},
			"addresses": {},
			"mother":    {},
		},
	}, 1, 2)
	userDot := user.Dot()
	assert.Equal(t, userDot.Get("total"), 1)
	assert.Equal(t, userDot.Get("next"), -1)
	assert.Equal(t, userDot.Get("page"), 1)
	assert.Equal(t, userDot.Get("data.0.id"), int64(1))
	assert.Equal(t, userDot.Get("data.0.manu.name"), "北京云道天成科技有限公司")
	assert.Equal(t, userDot.Get("data.0.mother.extra.sex"), "女")
	assert.Equal(t, userDot.Get("data.0.extra.sex"), "男")
	assert.Equal(t, userDot.Get("data.0.addresses.0.location"), "北京国家数字出版基地A103")

}

func TestModelMustSearchWithsWheresOrder(t *testing.T) {
	user := Select("user").MustSearch(QueryParam{
		Orders: []QueryOrder{
			{
				Column: "id",
				Option: "desc",
			},
		},
		Wheres: []QueryWhere{
			{
				Wheres: []QueryWhere{
					{
						Column: "mobile",
						Value:  "13900002222",
					}, {
						Column: "mobile",
						Method: "orwhere",
						Value:  "13900001111",
					},
				},
			},
		},
		Withs: map[string]With{
			"manu":      {},
			"addresses": {},
			"mother":    {},
		},
	}, 1, 2)
	userDot := user.Dot()
	assert.Equal(t, userDot.Get("total"), 2)
	assert.Equal(t, userDot.Get("next"), -1)
	assert.Equal(t, userDot.Get("page"), 1)
	assert.Equal(t, userDot.Get("data.1.id"), int64(1))
	assert.Equal(t, userDot.Get("data.1.manu.name"), "北京云道天成科技有限公司")
	assert.Equal(t, userDot.Get("data.1.mother.extra.sex"), "女")
	assert.Equal(t, userDot.Get("data.1.extra.sex"), "男")
	assert.Equal(t, userDot.Get("data.1.addresses.0.location"), "北京国家数字出版基地A103")

}

func TestModelMustSaveNew(t *testing.T) {
	user := Select("user")
	id := user.MustSave(maps.MapStr{
		"name":     "用户创建",
		"manu_id":  2,
		"type":     "user",
		"idcard":   "23082619820207006X",
		"mobile":   "13900004444",
		"password": "qV@uT1DI",
		"key":      "XZ12MiPp",
		"secret":   "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN",
		"status":   "enabled",
		"extra":    maps.MapStr{"sex": "女"},
	})

	row := user.MustFind(id, QueryParam{})

	// 清空数据
	capsule.Query().Table(user.MetaData.Table.Name).Where("id", id).Delete()

	assert.Equal(t, row.Get("name"), "用户创建")
	assert.Equal(t, row.Dot().Get("extra.sex"), "女")

}

func TestModelMustSaveUpdate(t *testing.T) {
	user := Select("user")
	id := user.MustSave(maps.MapStr{
		"id":      1,
		"balance": 200,
	})

	row := user.MustFind(id, QueryParam{})

	// 恢复数据
	capsule.Query().Table(user.MetaData.Table.Name).Where("id", id).Update(maps.MapStr{"balance": 0})
	assert.Equal(t, any.Of(row.Get("balance")).CInt(), 200)
}

func TestModelMustUpdate(t *testing.T) {
	user := Select("user")
	id := user.MustUpdate(
		QueryParam{
			Wheres: []QueryWhere{
				{
					Column: "id",
					Value:  1,
				},
			},
		},
		maps.MapStr{
			"balance": 200,
		})

	row := user.MustFind(id, QueryParam{})

	// 恢复数据
	capsule.Query().Table(user.MetaData.Table.Name).Where("id", id).Update(maps.MapStr{"balance": 0})
	assert.Equal(t, any.Of(row.Get("balance")).CInt(), 200)
}

func TestModelMustDeleteSoft(t *testing.T) {
	user := Select("user")
	id := user.MustSave(maps.MapStr{
		"name":     "用户创建",
		"manu_id":  2,
		"type":     "user",
		"idcard":   "23082619820207006X",
		"mobile":   "13900004444",
		"password": "qV@uT1DI",
		"key":      "XZ12MiPp",
		"secret":   "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN",
		"status":   "enabled",
		"extra":    maps.MapStr{"sex": "女"},
	})
	user.MustDelete(id)
	row, _ := user.Find(id, QueryParam{})

	// 清空数据
	capsule.Query().Table(user.MetaData.Table.Name).Where("id", id).Delete()
	assert.Nil(t, row)
}

func TestModelMustDestory(t *testing.T) {
	user := Select("user")
	id := user.MustSave(maps.MapStr{
		"name":     "用户创建",
		"manu_id":  2,
		"type":     "user",
		"idcard":   "23082619820207006X",
		"mobile":   "13900004444",
		"password": "qV@uT1DI",
		"key":      "XZ12MiPp",
		"secret":   "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN",
		"status":   "enabled",
		"extra":    maps.MapStr{"sex": "女"},
	})
	user.MustDestroy(id)

	row, err := capsule.Query().Table(user.MetaData.Table.Name).Where("id", id).First()
	assert.True(t, row.IsEmpty())
	assert.Nil(t, err)
}

func TestModelMustInsert(t *testing.T) {
	columns := []string{"user_id", "province", "city", "location"}
	rows := [][]interface{}{
		{4, "北京市", "丰台区", "银海星月9号楼9单元9层1024室"},
		{4, "天津市", "塘沽区", "益海星云7号楼3单元1003室"},
	}
	address := Select("address")
	err := address.Insert(columns, rows)
	assert.Nil(t, err)
	capsule.Query().Table(address.MetaData.Table.Name).Where("user_id", 4).Delete()
}

func TestModelMustInsertError(t *testing.T) {
	columns := []string{"user_id", "province", "city", "location"}
	rows := [][]interface{}{
		{4, "北京市", "丰台区", "银海星月9号楼9单元9层1024室"},
		{4, "天津市", "塘沽区", "益海星云7号楼3单元1003室", 5028},
		{4, "天津市", "塘沽区", "益海星云7号楼3单元1002室"},
	}
	address := Select("address")
	assert.Panics(t, func() {
		address.Insert(columns, rows)
	})
}

func TestModelMustDeleteWhere(t *testing.T) {
	columns := []string{"name", "manu_id", "type", "idcard", "mobile", "password", "key", "secret", "status"}
	rows := [][]interface{}{
		{"用户创建1", 5, "user", "23082619820207006X", "13900004444", "qV@uT1DI", "XZ12MiP1", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
		{"用户创建2", 5, "user", "33082619820207006X", "13900005555", "qV@uT1DI", "XZ12MiP2", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
		{"用户创建3", 5, "user", "43082619820207006X", "13900006666", "qV@uT1DI", "XZ12MiP3", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
	}

	user := Select("user")
	user.Insert(columns, rows)
	param := QueryParam{Wheres: []QueryWhere{
		{
			Column: "manu_id",
			Value:  5,
		},
	}}
	effect := user.MustDeleteWhere(param)

	// 清理数据
	capsule.Query().Table(user.MetaData.Table.Name).Where("name", "like", "用户创建%").Delete()
	assert.Equal(t, effect, 3)
}

func TestModelMustDestoryWhere(t *testing.T) {
	columns := []string{"name", "manu_id", "type", "idcard", "mobile", "password", "key", "secret", "status"}
	rows := [][]interface{}{
		{"用户创建1", 5, "user", "23082619820207006X", "13900004444", "qV@uT1DI", "XZ12MiP1", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
		{"用户创建2", 5, "user", "33082619820207006X", "13900005555", "qV@uT1DI", "XZ12MiP2", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
		{"用户创建3", 5, "user", "43082619820207006X", "13900006666", "qV@uT1DI", "XZ12MiP3", "wBeYjL7FjbcvpAdBrxtDFfjydsoPKhRN", "enabled"},
	}

	user := Select("user")
	user.Insert(columns, rows)
	param := QueryParam{Wheres: []QueryWhere{
		{
			Column: "manu_id",
			Value:  5,
		},
	}}
	effect := user.MustDestoryWhere(param)

	// 清理数据
	assert.Equal(t, effect, 3)
}
