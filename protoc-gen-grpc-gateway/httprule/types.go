package httprule

type template struct {
	segments []segment
	verb     string
}

type segment interface {
	// Stringer
	//	compile() (ops []op)
}

type wildcard struct{}

type deepWildcard struct{}

type literal string

type variable struct {
	path     string
	segments []segment
}
