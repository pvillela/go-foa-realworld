package bad

type Pw[E any, R any] struct {
	Entity E `db:""`
	RecCtx R `db:""`
}

// There should be a syntax error highlight below saying:
// "got 1 arguments but 2 type parameters [compiler(WrongTypeArgCount)]"
var _ Pw[string]
