package handledb

type Config struct{
	RawDbDir   string     `json:"rawDbDir""`
	ToDbDir    string     `json:"todbDir"`
	WalletDir  []string   `json:"walletDir"`
}