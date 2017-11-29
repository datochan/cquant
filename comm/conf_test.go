package comm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	assert.True(t, true)
	conf := new(Configure)
	conf.Parse("../configure.toml")
	assert.Equal(t, ":8000", conf.App.Addr)
	assert.Equal(t, "DEBUG", conf.App.Logger.Level)
	assert.Equal(t, "DatoQuant", conf.App.Logger.Name)
	assert.Equal(t, "mysql", conf.Db.Driver)
	assert.Equal(t, "root:123456@localhost/DatoQuant?charset=utf8&parseTime=True&loc=Local", conf.Db.Source)
}
