package lightning

type PaymentState string
type LightningNodeType string

const (
	PaymentSucceeded PaymentState = "succeeded"
	PaymentFailed    PaymentState = "failed"
	PaymentPending   PaymentState = "pending"

	NodeTypeLnd LightningNodeType = "LND"
	NodeTypeCln LightningNodeType = "CLN"
)

type PaymentStatus struct {
	State         PaymentState
	FailureReason *string

	Hash     string
	Preimage string

	FeeMsat uint64
}

type PaymentUpdate struct {
	IsLastUpdate bool
	Update       PaymentStatus
}

type LightningInfo struct {
	Pubkey      string
	BlockHeight uint32
	Version     string
	Network     string
	Synced      bool
}

type AddInvoiceResponse struct {
	PaymentRequest string
	PaymentHash    []byte
}

type LightningNode interface {
	Connect() error
	//Name() string
	//NodeType() LightningNodeType
	//PaymentStatus(preimageHash string) (*PaymentStatus, error)

	//SendPayment(invoice string, feeLimit uint64, timeout int32) (<-chan *PaymentUpdate, error)
	//PayInvoice(invoice string, maxParts uint32, timeoutSeconds int32) (int64, error)
	CreateInvoice(value int64, preimage []byte, expiry int64, memo string) (*AddInvoiceResponse, error)

	NewAddress() (string, error)

	GetInfo() (*LightningInfo, error)
	//ListChannels() (*lnrpc.ListChannelsResponse, error)
	//GetChannelInfo(chanId uint64) (*lnrpc.ChannelEdge, error)
}
