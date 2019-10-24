package pwengine

import (
	"context"

	"pathwar.land/go/pkg/pwdb"
)

func (e *engine) ToolGenerateFakeData(context.Context, *Void) (*Void, error) {
	return &Void{}, pwdb.GenerateFakeData(e.db, e.snowflake, e.logger.Named("generate-fake-data"))
}