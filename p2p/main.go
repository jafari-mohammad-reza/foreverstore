package p2p

type Transport interface {
	ListenAndAccept() error
}

type HandShaker interface {
	HandShake() error
}
