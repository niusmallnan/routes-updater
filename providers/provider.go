package providers

const (
	MetadataURL            = "http://%s/2016-07-29"
	DefaultMetadataAddress = "169.254.169.250"
)

type Provider interface {
	NewInst(string) (*Provider, error)
	Start()
	Reload() error
}
