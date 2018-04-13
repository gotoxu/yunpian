package yunpian

type clientRequest interface {
	Verify() error
}
