package bad

type Pw[E any, R any] struct {
	Entity E `db:""`
	RecCtx R `db:""`
}

// There should be a syntax error highlight below
var _ Pw[string]
