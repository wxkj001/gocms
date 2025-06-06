package extend

import (
	"gocms/extend/internal/oss"

	"go.uber.org/fx"
	"xorm.io/xorm"
)

// extend模块
var ExtendModule = fx.Module("extendModule", fx.Provide(NewExtend))

type ExtendResult struct {
	fx.Out
	Extend *Extend
}
type ExtendParams struct {
	fx.In
	*Extend
}
type Extend struct {
	S3 *oss.S3
}

func NewExtend(db *xorm.Engine) (ExtendResult, error) {
	return ExtendResult{Extend: &Extend{
		S3: oss.NewS3(),
	}}, nil
}
