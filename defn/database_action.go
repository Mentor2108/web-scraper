package defn

type DatabaseAction int

const (
	DatabaseActionCreate DatabaseAction = iota
	DatabaseActionRead
	DatabaseActionUpdate
	DatabaseActionDelete
)
