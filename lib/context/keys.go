package context

type ContextKey struct {
	key string
}

var (
	BlockDatabase   = &ContextKey{key: "BlockDatabase"}
	ContentDatabase = &ContextKey{key: "ContentDatabase"}
	CacheDatabase   = &ContextKey{key: "CacheDatabase"}
	GrpcServer      = &ContextKey{key: "GrpcServer"}

	PrivateKey = &ContextKey{key: "PrivateKey"}
	PublicKey  = &ContextKey{key: "PublicKey"}
	Address    = &ContextKey{key: "Address"}
)
