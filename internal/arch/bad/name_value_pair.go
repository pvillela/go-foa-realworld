package bad

// There should be a syntax error highlight below saying:
// "generic type cannot be alias [compiler(BadDecl)]"
type NameValuePair[N any, V any] = struct {
	Name  N
	Value V
}
