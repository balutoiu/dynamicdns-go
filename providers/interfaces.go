package providers

type DNSClient interface {
	UpdateIP() error
}
