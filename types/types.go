package types

// minerPledgeInput
type MinerPledgeRequestInput struct {
	Account string `gorm:"column:Account"`
}

// minerPledgeOutput
type MinerPledgeRequestOutput struct {
	MinerPledge string `gorm:"column:MinerPledge"`
}

// AccountTotalRewardRequestInput
type AccountTotalRewardRequestInput struct {
	Account   string `gorm:"column:Account"`
	StartTime string `gorm:"column:StartTime"`
	EndTime   string `gorm:"column:EndTime"`
}

// AccountDetailRequestOutput
type AccountTotalRewardRequestOutput struct {
	TodayBlockRewards      string `gorm:"column:TodayBlockRewards"`     // 当天的奖励
	Today25PercentRewards  string `gorm:"column:Today25PercentRewards"` // 当天奖励25%的释放
	Today180PercentRewards string `gorm:"column:Today25PercentRewards"` // 当天累计1/180释放
	TotalTodayRewards      string `gorm:"column:Total25PercentRewards"` // 累计当天总释放
	PunishFee              string `gorm:"column:PunishFee"`             // 惩罚
	MinerPower             string `gorm:"column:MinerPower"`            // 矿工算力
}

// AccountRequestInput
type AccountDetailRequestInput struct {
	Account     string `gorm:"column:Account"`
	StartTime   string `gorm:"column:StartTime"`
	EndTime     string `gorm:"column:EndTime"`
	PageSize    string `gorm:"column:PageSize"`
	CurrentPage string `gorm:"column:CurrentPage"`
}

// normal account
type AccountInfoOutput struct {
	Id          string `gorm:"column:Id"` // normalAddress workerId(workerAddress) minerId(minerAddress)
	Balance     string `gorm:"column:Balance"`
	BlockHeight int64  `gorm:"column:BlockHeight"`
	Fee         string `gorm:"column:Fee"` // 费用
	MinerTip    string `gorm:"column:MinerTip"`
	SendIn      string `gorm:"column:Send"` // 转账（入)
	SendOut     string `gorm:"column:Send"` // 转账（出)
	Send        string `gorm:"column:Send"` // 转账
}

// worker account
type WorkerInfoOutput struct {
	AccountInfoOutput
	PreCommitSectors   string `gorm:"column:preCommitSectors"`
	ProveCommitSectors string `gorm:"column:proveCommitSectors"`
}

// miner account
type MinerInfoOutput struct {
	AccountInfoOutput
	PunishFee             string `gorm:"column:punishFee"` // 惩罚
	PreCommitDeposits     string `gorm:"column:preCommitDeposits"`
	PreCommitSectors      string `gorm:"column:preCommitSectors"`
	ProveCommitSectors    string `gorm:"column:proveCommitSectors"`
	BlockReward           string `gorm:"column:blockReward"`
	TAG                   string `gorm:"column:tag"`
	MinerAvailableBalance string `gorm:"column:minerAvailableBalance"`
	LockedFunds           string `gorm:"column:lockedFunds"`
	InitialPledge         string `gorm:"column:initialPledge"`
	SubLockFunds          string `gorm:"column:SubLockFunds"`
	WithdrawBalance       string `gorm:"column:WithdrawBalance"`
	FlagBlockIsNull       bool   `gorm:"column:FlagBlockIsNull"`
}

type ServiceRegisterOutput struct {
	IP   string `json:"IP"`
	Port string `json:"Port"`
}
