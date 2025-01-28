package go3s

import (
	"context"
	"fmt"
	"math"
	"math/big"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/dwdwow/golimiter"
)

const (
	PUBLIC_BASE_URL = "https://public-api.solscan.io"
	PRO_BASE_URL    = "https://pro-api.solscan.io/v2.0"
)

const (
	V2_MAX_REQUESTS_PER_MINUTE = 1000
	V3_MAX_REQUESTS_PER_MINUTE = 2000
)

var (
	V2Limiter = golimiter.NewReqLimiter(time.Minute, V2_MAX_REQUESTS_PER_MINUTE/2)
	V3Limiter = golimiter.NewReqLimiter(time.Minute, V3_MAX_REQUESTS_PER_MINUTE/2)
)

type RespData[D any] struct {
	Success bool `json:"success"`
	Data    D    `json:"data"`
}

type RespDataWithTotal[Item any] struct {
	Transactions []Transaction `json:"transactions"`
	Items        []Item        `json:"items"`
	Data         []Item        `json:"data"`
	Total        int64         `json:"total"`
}

// Basic error types
type Errors struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (e *Errors) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

type RespError struct {
	Success bool   `json:"success"`
	Errors  Errors `json:"errors"`
}

type ChainInfo struct {
	BlockHeight      int64 `json:"blockHeight"`
	CurrentEpoch     int64 `json:"currentEpoch"`
	AbsoluteSlot     int64 `json:"absoluteSlot"`
	TransactionCount int64 `json:"transactionCount"`
}

// Enums
type Flow string

const (
	FlowIn    Flow = "in"
	FlowOut   Flow = "out"
	FlowEmpty Flow = ""
)

type TinyPageSize int64

const (
	TinyPageSize12 TinyPageSize = 12
	TinyPageSize24 TinyPageSize = 24
	TinyPageSize36 TinyPageSize = 36
)

type SmallPageSize int64

const (
	SmallPageSize10 SmallPageSize = 10
	SmallPageSize20 SmallPageSize = 20
	SmallPageSize30 SmallPageSize = 30
	SmallPageSize40 SmallPageSize = 40
)

type LargePageSize int64

const (
	LargePageSize10  LargePageSize = 10
	LargePageSize20  LargePageSize = 20
	LargePageSize30  LargePageSize = 30
	LargePageSize40  LargePageSize = 40
	LargePageSize60  LargePageSize = 60
	LargePageSize100 LargePageSize = 100
)

type SortBy string

const (
	SortByBlockTime SortBy = "block_time"
)

type MarketSortBy string

const (
	MarketSortByVolume MarketSortBy = "volume"
	MarketSortByTrade  MarketSortBy = "trade"
)

type TokenSortBy string

const (
	TokenSortByPrice       TokenSortBy = "price"
	TokenSortByHolder      TokenSortBy = "holder"
	TokenSortByMarketCap   TokenSortBy = "market_cap"
	TokenSortByCreatedTime TokenSortBy = "created_time"
)

type Transfer struct {
	BlockID       int64        `json:"block_id" bson:"block_id"`
	TransID       string       `json:"trans_id" bson:"trans_id"`
	BlockTime     int64        `json:"block_time" bson:"block_time"`
	Time          string       `json:"time" bson:"time"`
	ActivityType  ActivityType `json:"activity_type" bson:"activity_type"`
	FromAddress   string       `json:"from_address" bson:"from_address"`
	ToAddress     string       `json:"to_address" bson:"to_address"`
	TokenAddress  string       `json:"token_address" bson:"token_address"`
	TokenDecimals int64        `json:"token_decimals" bson:"token_decimals"`
	Amount        big.Int      `json:"amount" bson:"amount"`
	Flow          Flow         `json:"flow" bson:"flow"`
}

type NFTCollectionSortBy string

const (
	NFTCollectionSortByItems      NFTCollectionSortBy = "items"
	NFTCollectionSortByFloorPrice NFTCollectionSortBy = "floor_price"
	NFTCollectionSortByVolumes    NFTCollectionSortBy = "volumes"
)

type NFTCollectionItemSortBy string

const (
	NFTCollectionItemSortByLastTrade    NFTCollectionItemSortBy = "last_trade"
	NFTCollectionItemSortByListingPrice NFTCollectionItemSortBy = "listing_price"
)

type SortOrder string

const (
	SortOrderAsc  SortOrder = "asc"
	SortOrderDesc SortOrder = "desc"
)

type AccountActivityType string

const (
	AccountActivityTypeTransfer      AccountActivityType = "ACTIVITY_SPL_TRANSFER"
	AccountActivityTypeBurn          AccountActivityType = "ACTIVITY_SPL_BURN"
	AccountActivityTypeMint          AccountActivityType = "ACTIVITY_SPL_MINT"
	AccountActivityTypeCreateAccount AccountActivityType = "ACTIVITY_SPL_CREATE_ACCOUNT"
)

type TokenType string

const (
	TokenTypeToken TokenType = "token"
	TokenTypeNFT   TokenType = "nft"
)

type ActivityType string

const (
	// Account Activity
	ActivityTypeSwap            ActivityType = "ACTIVITY_TOKEN_SWAP"
	ActivityTypeAggSwap         ActivityType = "ACTIVITY_AGG_TOKEN_SWAP"
	ActivityTypeAddLiquidity    ActivityType = "ACTIVITY_TOKEN_ADD_LIQ"
	ActivityTypeRemoveLiquidity ActivityType = "ACTIVITY_TOKEN_REMOVE_LIQ"
	ActivityTypeStake           ActivityType = "ACTIVITY_SPL_TOKEN_STAKE"
	ActivityTypeUnstake         ActivityType = "ACTIVITY_SPL_TOKEN_UNSTAKE"
	ActivityTypeWithdrawStake   ActivityType = "ACTIVITY_SPL_TOKEN_WITHDRAW_STAKE"
	ActivityTypeMint            ActivityType = "ACTIVITY_SPL_MINT"
	ActivityTypeInitMint        ActivityType = "ACTIVITY_SPL_INIT_MINT"

	// Token Activity
	ActivityTypeTransfer      ActivityType = "ACTIVITY_SPL_TRANSFER"
	ActivityTypeBurn          ActivityType = "ACTIVITY_SPL_BURN"
	ActivityTypeTokenMint     ActivityType = "ACTIVITY_SPL_MINT"
	ActivityTypeCreateAccount ActivityType = "ACTIVITY_SPL_CREATE_ACCOUNT"
)

type NFTActivityType string

const (
	NFTActivityTypeSold        NFTActivityType = "ACTIVITY_NFT_SOLD"
	NFTActivityTypeListing     NFTActivityType = "ACTIVITY_NFT_LISTING"
	NFTActivityTypeBidding     NFTActivityType = "ACTIVITY_NFT_BIDDING"
	NFTActivityTypeCancelBid   NFTActivityType = "ACTIVITY_NFT_CANCEL_BID"
	NFTActivityTypeCancelList  NFTActivityType = "ACTIVITY_NFT_CANCEL_LIST"
	NFTActivityTypeRejectBid   NFTActivityType = "ACTIVITY_NFT_REJECT_BID"
	NFTActivityTypeUpdatePrice NFTActivityType = "ACTIVITY_NFT_UPDATE_PRICE"
	NFTActivityTypeListAuction NFTActivityType = "ACTIVITY_NFT_LIST_AUCTION"
)

type BalanceChangeType string

const (
	BalanceChangeTypeInc BalanceChangeType = "inc"
	BalanceChangeTypeDec BalanceChangeType = "dec"
)

type TxStatus string

const (
	TxStatusSuccess TxStatus = "Success"
	TxStatusFail    TxStatus = "Fail"
)

type StakeRole string

const (
	StakeRoleStaker     StakeRole = "staker"
	StakeRoleWithdrawer StakeRole = "withdrawer"
)

type StakeAccountStatus string

const (
	StakeAccountStatusActive StakeAccountStatus = "active"
)

type StakeAccountType string

const (
	StakeAccountTypeActive StakeAccountType = "active"
)

type AccountType string

const (
	AccountTypeSystemAccount AccountType = "system_account"
)

type TxFilter string

const (
	TxFilterExceptVote TxFilter = "exceptVote"
	TxFilterAll        TxFilter = "all"
)

type TokenAccount struct {
	TokenAccount  string `json:"token_account" bson:"token_account"`
	TokenAddress  string `json:"token_address" bson:"token_address"`
	Amount        int64  `json:"amount" bson:"amount"`
	TokenDecimals int64  `json:"token_decimals" bson:"token_decimals"`
	Owner         string `json:"owner" bson:"owner"`
}

type ChildRouter struct {
	Token1         string `json:"token1" bson:"token1"`
	Token1Decimals int64  `json:"token1_decimals" bson:"token1_decimals"`
	Amount1        int64  `json:"amount1" bson:"amount1"`
	Token2         string `json:"token2" bson:"token2"`
	Token2Decimals int64  `json:"token2_decimals" bson:"token2_decimals"`
	Amount2        int64  `json:"amount2" bson:"amount2"`
}

type Router struct {
	Token1         string        `json:"token1" bson:"token1"`
	Token1Decimals int64         `json:"token1_decimals" bson:"token1_decimals"`
	Amount1        int64         `json:"amount1" bson:"amount1"`
	Token2         string        `json:"token2" bson:"token2"`
	Token2Decimals int64         `json:"token2_decimals" bson:"token2_decimals"`
	Amount2        int64         `json:"amount2" bson:"amount2"`
	ChildRouters   []ChildRouter `json:"child_routers" bson:"child_routers"`
}

type AmountInfo struct {
	Token1         string   `json:"token1" bson:"token1"`
	Token1Decimals int64    `json:"token1_decimals" bson:"token1_decimals"`
	Amount1        int64    `json:"amount1" bson:"amount1"`
	Token2         string   `json:"token2" bson:"token2"`
	Token2Decimals int64    `json:"token2_decimals" bson:"token2_decimals"`
	Amount2        int64    `json:"amount2" bson:"amount2"`
	Routers        []Router `json:"routers" bson:"routers"`
}

type DefiActivity struct {
	BlockID      int64        `json:"block_id" bson:"block_id"`
	TransID      string       `json:"trans_id" bson:"trans_id"`
	BlockTime    int64        `json:"block_time" bson:"block_time"`
	Time         string       `json:"time" bson:"time"`
	ActivityType ActivityType `json:"activity_type" bson:"activity_type"`
	FromAddress  string       `json:"from_address" bson:"from_address"`
	ToAddress    string       `json:"to_address" bson:"to_address"`
	Sources      []string     `json:"sources" bson:"sources"`
	Platform     []string     `json:"platform" bson:"platform"`
	Routers      Router       `json:"routers" bson:"routers"`
}

type AccountChangeActivity struct {
	BlockID       int64             `json:"block_id" bson:"block_id"`
	BlockTime     int64             `json:"block_time" bson:"block_time"`
	Time          string            `json:"time" bson:"time"`
	TransID       string            `json:"trans_id" bson:"trans_id"`
	Address       string            `json:"address" bson:"address"`
	TokenAddress  string            `json:"token_address" bson:"token_address"`
	TokenAccount  string            `json:"token_account" bson:"token_account"`
	TokenDecimals int64             `json:"token_decimals" bson:"token_decimals"`
	Amount        int64             `json:"amount" bson:"amount"`
	PreBalance    int64             `json:"pre_balance" bson:"pre_balance"`
	PostBalance   int64             `json:"post_balance" bson:"post_balance"`
	ChangeType    BalanceChangeType `json:"change_type" bson:"change_type"`
	Fee           int64             `json:"fee" bson:"fee"`
}

type ParsedCancelAllAndPlaceOrders struct {
	Type      string `json:"type" bson:"type"`
	Program   string `json:"program" bson:"program"`
	ProgramID string `json:"program_id" bson:"program_id"`
}

type Transaction struct {
	Slot               int64                           `json:"slot" bson:"slot"`
	Fee                int64                           `json:"fee" bson:"fee"`
	Status             TxStatus                        `json:"status" bson:"status"`
	Signer             []string                        `json:"signer" bson:"signer"`
	BlockTime          int64                           `json:"block_time" bson:"block_time"`
	TxHash             string                          `json:"tx_hash" bson:"tx_hash"`
	ParsedInstructions []ParsedCancelAllAndPlaceOrders `json:"parsed_instructions" bson:"parsed_instructions"`
	ProgramIDs         []string                        `json:"program_ids" bson:"program_ids"`
	Time               string                          `json:"time" bson:"time"`
}

type AccountStake struct {
	Amount               int64              `json:"amount" bson:"amount"`
	Role                 []StakeRole        `json:"role" bson:"role"`
	Status               StakeAccountStatus `json:"status" bson:"status"`
	Type                 StakeAccountType   `json:"type" bson:"type"`
	Voter                string             `json:"voter" bson:"voter"`
	ActiveStakeAmount    int64              `json:"active_stake_amount" bson:"active_stake_amount"`
	DelegatedStakeAmount int64              `json:"delegated_stake_amount" bson:"delegated_stake_amount"`
	SolBalance           int                `json:"sol_balance" bson:"sol_balance"`
	TotalReward          string             `json:"total_reward" bson:"total_reward"`
	StakeAccount         string             `json:"stake_account" bson:"stake_account"`
	ActivationEpoch      int64              `json:"activation_epoch" bson:"activation_epoch"`
	StakeType            int64              `json:"stake_type" bson:"stake_type"`
}

type AccountDetail struct {
	Account      string      `json:"account" bson:"account"`
	Lamports     int64       `json:"lamports" bson:"lamports"`
	Type         AccountType `json:"type" bson:"type"`
	Executable   bool        `json:"executable" bson:"executable"`
	OwnerProgram string      `json:"owner_program" bson:"owner_program"`
	// RentEpoch is too large, so we don't need it.
	// RentEpoch    int64       `json:"rent_epoch"`
	IsOncurve bool `json:"is_oncurve" bson:"is_oncurve"`
}

type Market struct {
	PoolID             string  `json:"pool_id" bson:"pool_id"`
	ProgramID          string  `json:"program_id" bson:"program_id"`
	Token1             string  `json:"token_1" bson:"token_1"`
	Token2             string  `json:"token_2" bson:"token_2"`
	TokenAccount1      string  `json:"token_account_1" bson:"token_account_1"`
	TokenAccount2      string  `json:"token_account_2" bson:"token_account_2"`
	TotalTrades24h     int64   `json:"total_trades_24h" bson:"total_trades_24h"`
	TotalTradesPrev24h int64   `json:"total_trades_prev_24h" bson:"total_trades_prev_24h"`
	TotalVolume24h     float64 `json:"total_volume_24h" bson:"total_volume_24h"`
	TotalVolumePrev24h float64 `json:"total_volume_prev_24h" bson:"total_volume_prev_24h"`
}

type Token struct {
	Address        string  `json:"address" bson:"address"`
	Decimals       int64   `json:"decimals" bson:"decimals"`
	Name           string  `json:"name" bson:"name"`
	Symbol         string  `json:"symbol" bson:"symbol"`
	MarketCap      float64 `json:"market_cap" bson:"market_cap"`
	Price          float64 `json:"price" bson:"price"`
	Price24hChange float64 `json:"price_24h_change" bson:"price_24h_change"`
	CreatedTime    int64   `json:"created_time" bson:"created_time"`
}

type TokenPrice struct {
	Date  int64   `json:"date" bson:"date"` // yyyymmdd format
	Price float64 `json:"price" bson:"price"`
}

type TokenHolder struct {
	Address  string  `json:"address" bson:"address"`
	Amount   big.Int `json:"amount" bson:"amount"`
	Decimals int64   `json:"decimals" bson:"decimals"`
	Owner    string  `json:"owner" bson:"owner"`
	Rank     int64   `json:"rank" bson:"rank"`
}

type TokenMeta struct {
	Supply         string  `json:"supply" bson:"supply"`
	Address        string  `json:"address" bson:"address"`
	Name           string  `json:"name" bson:"name"`
	Symbol         string  `json:"symbol" bson:"symbol"`
	Icon           string  `json:"icon" bson:"icon"`
	Decimals       int64   `json:"decimals" bson:"decimals"`
	Holder         int64   `json:"holder" bson:"holder"`
	Creator        string  `json:"creator" bson:"creator"`
	CreateTx       string  `json:"create_tx" bson:"create_tx"`
	CreatedTime    int64   `json:"created_time" bson:"created_time"`
	FirstMintTx    string  `json:"first_mint_tx" bson:"first_mint_tx"`
	FirstMintTime  int64   `json:"first_mint_time" bson:"first_mint_time"`
	Price          float64 `json:"price" bson:"price"`
	Volume24h      float64 `json:"volume_24h" bson:"volume_24h"`
	MarketCap      float64 `json:"market_cap" bson:"market_cap"`
	MarketCapRank  int64   `json:"market_cap_rank" bson:"market_cap_rank"`
	PriceChange24h float64 `json:"price_change_24h" bson:"price_change_24h"`
}

type TokenTop struct {
	Address        string  `json:"address" bson:"address"`
	Decimals       int64   `json:"decimals" bson:"decimals"`
	Name           string  `json:"name" bson:"name"`
	Symbol         string  `json:"symbol" bson:"symbol"`
	MarketCap      float64 `json:"market_cap" bson:"market_cap"`
	Price          float64 `json:"price" bson:"price"`
	Price24hChange float64 `json:"price_24h_change" bson:"price_24h_change"`
	CreatedTime    int64   `json:"created_time" bson:"created_time"`
}

type AccountKey struct {
	Pubkey   string `json:"pubkey" bson:"pubkey"`
	Signer   bool   `json:"signer" bson:"signer"`
	Source   string `json:"source" bson:"source"`
	Writable bool   `json:"writable" bson:"writable"`
}

type TransferInfo struct {
	SourceOwner      string                 `json:"source_owner" bson:"source_owner"`
	Source           string                 `json:"source" bson:"source"`
	Destination      string                 `json:"destination" bson:"destination"`
	DestinationOwner string                 `json:"destination_owner" bson:"destination_owner"`
	TransferType     string                 `json:"transfer_type" bson:"transfer_type"`
	TokenAddress     string                 `json:"token_address" bson:"token_address"`
	Decimals         int64                  `json:"decimals" bson:"decimals"`
	AmountStr        string                 `json:"amount_str" bson:"amount_str"`
	Amount           int64                  `json:"amount" bson:"amount"`
	ProgramID        string                 `json:"program_id" bson:"program_id"`
	OuterProgramID   string                 `json:"outer_program_id" bson:"outer_program_id"`
	InsIndex         int64                  `json:"ins_index" bson:"ins_index"`
	OuterInsIndex    int64                  `json:"outer_ins_index" bson:"outer_ins_index"`
	Event            string                 `json:"event" bson:"event"`
	Fee              map[string]interface{} `json:"fee" bson:"fee"`
}

type InstructionData struct {
	InsIndex           int64                    `json:"ins_index" bson:"ins_index"`
	ParsedType         string                   `json:"parsed_type" bson:"parsed_type"`
	Type               string                   `json:"type" bson:"type"`
	ProgramID          string                   `json:"program_id" bson:"program_id"`
	Program            string                   `json:"program" bson:"program"`
	OuterProgramID     *string                  `json:"outer_program_id,omitempty" bson:"outer_program_id,omitempty"`
	OuterInsIndex      int64                    `json:"outer_ins_index" bson:"outer_ins_index"`
	DataRaw            interface{}              `json:"data_raw" bson:"data_raw"` // can be string or map
	Accounts           []string                 `json:"accounts" bson:"accounts"`
	Activities         []map[string]interface{} `json:"activities" bson:"activities"`
	Transfers          []TransferInfo           `json:"transfers" bson:"transfers"`
	ProgramInvokeLevel int64                    `json:"program_invoke_level" bson:"program_invoke_level"`
}

type BalanceChange struct {
	Address      string `json:"address" bson:"address"`
	PreBalance   string `json:"pre_balance" bson:"pre_balance"`
	PostBalance  string `json:"post_balance" bson:"post_balance"`
	ChangeAmount string `json:"change_amount" bson:"change_amount"`
}

type TokenBalanceChange struct {
	Address      string `json:"address" bson:"address"`
	ChangeType   string `json:"change_type" bson:"change_type"`
	ChangeAmount string `json:"change_amount" bson:"change_amount"`
	Decimals     int64  `json:"decimals" bson:"decimals"`
	PostBalance  string `json:"post_balance" bson:"post_balance"`
	// if prebalance is 0, it is number, otherwise it is string
	PreBalance   any    `json:"pre_balance" bson:"pre_balance"`
	TokenAddress string `json:"token_address" bson:"token_address"`
	Owner        string `json:"owner" bson:"owner"`
	PostOwner    string `json:"post_owner" bson:"post_owner"`
	PreOwner     string `json:"pre_owner" bson:"pre_owner"`
}

type TransactionDetail struct {
	BlockID              int64                `json:"block_id" bson:"block_id"`
	Fee                  int64                `json:"fee" bson:"fee"`
	Reward               []interface{}        `json:"reward" bson:"reward"`
	SolBalChange         []BalanceChange      `json:"sol_bal_change" bson:"sol_bal_change"`
	TokenBalChange       []TokenBalanceChange `json:"token_bal_change" bson:"token_bal_change"`
	TokensInvolved       []string             `json:"tokens_involved" bson:"tokens_involved"`
	ParsedInstructions   []InstructionData    `json:"parsed_instructions" bson:"parsed_instructions"`
	ProgramsInvolved     []string             `json:"programs_involved" bson:"programs_involved"`
	Signer               []string             `json:"signer" bson:"signer"`
	Status               int64                `json:"status" bson:"status"`
	AccountKeys          []AccountKey         `json:"account_keys" bson:"account_keys"`
	ComputeUnitsConsumed int64                `json:"compute_units_consumed" bson:"compute_units_consumed"`
	Confirmations        *int64               `json:"confirmations,omitempty" bson:"confirmations,omitempty"`
	// if version is 0, it is number, otherwise it is string
	Version         any      `json:"version" bson:"version"`
	TxHash          string   `json:"tx_hash" bson:"tx_hash"`
	BlockTime       int64    `json:"block_time" bson:"block_time"`
	LogMessage      []string `json:"log_message" bson:"log_message"`
	RecentBlockHash string   `json:"recent_block_hash" bson:"recent_block_hash"`
	TxStatus        string   `json:"tx_status" bson:"tx_status"`
}

type TxActionData struct {
	AmmID          string  `json:"amm_id" bson:"amm_id"`
	AmmAuthority   *string `json:"amm_authority,omitempty" bson:"amm_authority,omitempty"`
	Account        string  `json:"account" bson:"account"`
	Token1         string  `json:"token_1" bson:"token_1"`
	Token2         string  `json:"token_2" bson:"token_2"`
	Amount1        any     `json:"amount_1" bson:"amount_1"`
	Amount1Str     string  `json:"amount_1_str" bson:"amount_1_str"`
	Amount2        any     `json:"amount_2" bson:"amount_2"`
	Amount2Str     string  `json:"amount_2_str" bson:"amount_2_str"`
	TokenDecimal1  int64   `json:"token_decimal_1" bson:"token_decimal_1"`
	TokenDecimal2  int64   `json:"token_decimal_2" bson:"token_decimal_2"`
	TokenAccount11 string  `json:"token_account_1_1" bson:"token_account_1_1"`
	TokenAccount12 string  `json:"token_account_1_2" bson:"token_account_1_2"`
	TokenAccount21 string  `json:"token_account_2_1" bson:"token_account_2_1"`
	TokenAccount22 string  `json:"token_account_2_2" bson:"token_account_2_2"`
	Owner1         string  `json:"owner_1" bson:"owner_1"`
	Owner2         string  `json:"owner_2" bson:"owner_2"`
}

type TxAction struct {
	Name           string       `json:"name" bson:"name"`
	ActivityType   string       `json:"activity_type" bson:"activity_type"`
	ProgramID      string       `json:"program_id" bson:"program_id"`
	Data           TxActionData `json:"data" bson:"data"`
	InsIndex       int64        `json:"ins_index" bson:"ins_index"`
	OuterInsIndex  int64        `json:"outer_ins_index" bson:"outer_ins_index"`
	OuterProgramID *string      `json:"outer_program_id,omitempty" bson:"outer_program_id,omitempty"`
}

type TxActionTransfer struct {
	SourceOwner      string `json:"source_owner" bson:"source_owner"`
	Source           string `json:"source" bson:"source"`
	Destination      string `json:"destination" bson:"destination"`
	DestinationOwner string `json:"destination_owner" bson:"destination_owner"`
	TransferType     string `json:"transfer_type" bson:"transfer_type"`
	TokenAddress     string `json:"token_address" bson:"token_address"`
	Decimals         int64  `json:"decimals" bson:"decimals"`
	AmountStr        string `json:"amount_str" bson:"amount_str"`
	Amount           int64  `json:"amount" bson:"amount"`
	ProgramID        string `json:"program_id" bson:"program_id"`
	OuterProgramID   string `json:"outer_program_id" bson:"outer_program_id"`
	InsIndex         int64  `json:"ins_index" bson:"ins_index"`
	OuterInsIndex    int64  `json:"outer_ins_index" bson:"outer_ins_index"`
}

type TransactionAction struct {
	TxHash     string             `json:"tx_hash" bson:"tx_hash"`
	BlockID    int64              `json:"block_id" bson:"block_id"`
	BlockTime  int64              `json:"block_time" bson:"block_time"`
	Time       string             `json:"time" bson:"time"`
	Fee        int64              `json:"fee" bson:"fee"`
	Transfers  []TxActionTransfer `json:"transfers" bson:"transfers"`
	Activities []TxAction         `json:"activities" bson:"activities"`
}

type BlockDetail struct {
	FeeRewards        int64  `json:"fee_rewards" bson:"fee_rewards"`
	TransactionsCount int64  `json:"transactions_count" bson:"transactions_count"`
	CurrentSlot       int64  `json:"current_slot" bson:"current_slot"`
	BlockHeight       int64  `json:"block_height" bson:"block_height"`
	BlockTime         int64  `json:"block_time" bson:"block_time"`
	Time              string `json:"time" bson:"time"`
	BlockHash         string `json:"block_hash" bson:"block_hash"`
	ParentSlot        int64  `json:"parent_slot" bson:"parent_slot"`
	PreviousBlockHash string `json:"previous_block_hash" bson:"previous_block_hash"`
}

type PoolMarket struct {
	PoolAddress    string `json:"pool_address" bson:"pool_address"`
	ProgramID      string `json:"program_id" bson:"program_id"`
	Token1         string `json:"token1" bson:"token1"`
	Token1Account  string `json:"token1_account" bson:"token1_account"`
	Token2         string `json:"token2" bson:"token2"`
	Token2Account  string `json:"token2_account" bson:"token2_account"`
	TotalVolume24h int64  `json:"total_volume_24h" bson:"total_volume_24h"`
	TotalTrade24h  int64  `json:"total_trade_24h" bson:"total_trade_24h"`
	CreatedTime    int64  `json:"created_time" bson:"created_time"`
}

type PoolMarketInfo struct {
	PoolAddress   string  `json:"pool_address" bson:"pool_address"`
	ProgramID     string  `json:"program_id" bson:"program_id"`
	Token1        string  `json:"token1" bson:"token1"`
	Token2        string  `json:"token2" bson:"token2"`
	Token1Account string  `json:"token1_account" bson:"token1_account"`
	Token2Account string  `json:"token2_account" bson:"token2_account"`
	Token1Amount  float64 `json:"token1_amount" bson:"token1_amount"`
	Token2Amount  float64 `json:"token2_amount" bson:"token2_amount"`
}

type PoolMarketDayVolume struct {
	Day    int64   `json:"day" bson:"day"` // yyyymmdd format
	Volume float64 `json:"volume" bson:"volume"`
}

type PoolMarketVolume struct {
	PoolAddress          string                `json:"pool_address" bson:"pool_address"`
	ProgramID            string                `json:"program_id" bson:"program_id"`
	TotalVolume24h       int64                 `json:"total_volume_24h" bson:"total_volume_24h"`
	TotalVolumeChange24h float64               `json:"total_volume_change_24h" bson:"total_volume_change_24h"`
	TotalTrades24h       int64                 `json:"total_trades_24h" bson:"total_trades_24h"`
	TotalTradesChange24h float64               `json:"total_trades_change_24h" bson:"total_trades_change_24h"`
	Days                 []PoolMarketDayVolume `json:"days" bson:"days"`
}

type APIUsage struct {
	RemainingCUs     int64   `json:"remaining_cus" bson:"remaining_cus"`
	UsageCUs         int64   `json:"usage_cus" bson:"usage_cus"`
	TotalRequests24h int64   `json:"total_requests_24h" bson:"total_requests_24h"`
	SuccessRate24h   float64 `json:"success_rate_24h" bson:"success_rate_24h"`
	TotalCU24h       int64   `json:"total_cu_24h" bson:"total_cu_24h"`
}

type NFTCreator struct {
	Address  string `json:"address" bson:"address"`
	Verified int64  `json:"verified" bson:"verified"`
	Share    int64  `json:"share" bson:"share"`
}

type NFTFile struct {
	URI  string `json:"uri" bson:"uri"`
	Type string `json:"type" bson:"type"`
}

type NFTProperties struct {
	Files    []NFTFile `json:"files" bson:"files"`
	Category string    `json:"category" bson:"category"`
}

type NFTAttribute struct {
	TraitType string `json:"trait_type" bson:"trait_type"`
	Value     string `json:"value" bson:"value"`
}

type NFTMetadata struct {
	Image                string         `json:"image" bson:"image"`
	TokenID              int64          `json:"tokenId" bson:"tokenId"`
	Name                 string         `json:"name" bson:"name"`
	Symbol               string         `json:"symbol" bson:"symbol"`
	Description          string         `json:"description" bson:"description"`
	SellerFeeBasisPoints int64          `json:"seller_fee_basis_points" bson:"seller_fee_basis_points"`
	Edition              int64          `json:"edition" bson:"edition"`
	Attributes           []NFTAttribute `json:"attributes" bson:"attributes"`
	Properties           NFTProperties  `json:"properties" bson:"properties"`
	Retried              int64          `json:"retried" bson:"retried"`
}

type NFTData struct {
	Name                 string       `json:"name" bson:"name"`
	Symbol               string       `json:"symbol" bson:"symbol"`
	URI                  string       `json:"uri" bson:"uri"`
	SellerFeeBasisPoints int64        `json:"sellerFeeBasisPoints" bson:"sellerFeeBasisPoints"`
	Creators             []NFTCreator `json:"creators" bson:"creators"`
	ID                   int64        `json:"id" bson:"id"`
}

type NFTInfo struct {
	Address       string      `json:"address" bson:"address"`
	Collection    string      `json:"collection" bson:"collection"`
	CollectionID  string      `json:"collectionId" bson:"collectionId"`
	CollectionKey string      `json:"collectionKey" bson:"collectionKey"`
	CreatedTime   int64       `json:"createdTime" bson:"createdTime"`
	Data          NFTData     `json:"data" bson:"data"`
	Meta          NFTMetadata `json:"meta" bson:"meta"`
	MintTx        string      `json:"mintTx" bson:"mintTx"`
}

type NFTActivity struct {
	BlockID            int64           `json:"block_id" bson:"block_id"`
	TransID            string          `json:"trans_id" bson:"trans_id"`
	BlockTime          int64           `json:"block_time" bson:"block_time"`
	Time               string          `json:"time" bson:"time"`
	ActivityType       NFTActivityType `json:"activity_type" bson:"activity_type"`
	FromAddress        string          `json:"from_address" bson:"from_address"`
	ToAddress          string          `json:"to_address" bson:"to_address"`
	TokenAddress       string          `json:"token_address" bson:"token_address"`
	MarketplaceAddress string          `json:"marketplace_address" bson:"marketplace_address"`
	CollectionAddress  string          `json:"collection_address" bson:"collection_address"`
	Amount             int64           `json:"amount" bson:"amount"`
	Price              int64           `json:"price" bson:"price"`
	CurrencyToken      string          `json:"currency_token" bson:"currency_token"`
	CurrencyDecimals   int64           `json:"currency_decimals" bson:"currency_decimals"`
}

type NFTCollection struct {
	CollectionID    string   `json:"collection_id" bson:"collection_id"`
	Name            string   `json:"name" bson:"name"`
	Symbol          string   `json:"symbol" bson:"symbol"`
	FloorPrice      float64  `json:"floor_price" bson:"floor_price"`
	Items           int64    `json:"items" bson:"items"`
	Marketplaces    []string `json:"marketplaces" bson:"marketplaces"`
	Volumes         float64  `json:"volumes" bson:"volumes"`
	TotalVolPrev24h float64  `json:"total_vol_prev_24h" bson:"total_vol_prev_24h"`
}

type NFTTradeInfo struct {
	TradeTime        int64  `json:"trade_time" bson:"trade_time"`
	Signature        string `json:"signature" bson:"signature"`
	MarketID         string `json:"market_id" bson:"market_id"`
	Type             string `json:"type" bson:"type"`
	Price            string `json:"price" bson:"price"`
	CurrencyToken    string `json:"currency_token" bson:"currency_token"`
	CurrencyDecimals int64  `json:"currency_decimals" bson:"currency_decimals"`
	Seller           string `json:"seller" bson:"seller"`
	Buyer            string `json:"buyer" bson:"buyer"`
}

type NFTCollectionMeta struct {
	Name   string `json:"name" bson:"name"`
	Family string `json:"family" bson:"family"`
}

type NFTMetaProperties struct {
	Files    []NFTFile    `json:"files" bson:"files"`
	Category string       `json:"category" bson:"category"`
	Creators []NFTCreator `json:"creators" bson:"creators"`
}

type NFTItemMetadata struct {
	Name                 string            `json:"name" bson:"name"`
	Symbol               string            `json:"symbol" bson:"symbol"`
	Description          string            `json:"description" bson:"description"`
	SellerFeeBasisPoints int64             `json:"seller_fee_basis_points" bson:"seller_fee_basis_points"`
	Image                string            `json:"image" bson:"image"`
	ExternalURL          string            `json:"external_url" bson:"external_url"`
	Collection           NFTCollectionMeta `json:"collection" bson:"collection"`
	Attributes           []NFTAttribute    `json:"attributes" bson:"attributes"`
	Properties           NFTMetaProperties `json:"properties" bson:"properties"`
}

type NFTItemData struct {
	Name                 string       `json:"name" bson:"name"`
	Symbol               string       `json:"symbol" bson:"symbol"`
	URI                  string       `json:"uri" bson:"uri"`
	SellerFeeBasisPoints int64        `json:"sellerFeeBasisPoints" bson:"sellerFeeBasisPoints"`
	Creators             []NFTCreator `json:"creators" bson:"creators"`
	ID                   int64        `json:"id" bson:"id"`
}

type NFTItemInfo struct {
	Address      string          `json:"address" bson:"address"`
	TokenName    string          `json:"token_name" bson:"token_name"`
	TokenSymbol  string          `json:"token_symbol" bson:"token_symbol"`
	CollectionID string          `json:"collection_id" bson:"collection_id"`
	Data         NFTItemData     `json:"data" bson:"data"`
	Meta         NFTItemMetadata `json:"meta" bson:"meta"`
	MintTx       string          `json:"mint_tx" bson:"mint_tx"`
	CreatedTime  int64           `json:"created_time" bson:"created_time"`
}

type NFTCollectionItem struct {
	TradeInfo NFTTradeInfo `json:"tradeInfo" bson:"tradeInfo"`
	Info      NFTItemInfo  `json:"info" bson:"info"`
}

type Client struct {
	Limiter *golimiter.ReqLimiter
	Headers map[string][]string
}

func NewClient(auth_token string, limiter *golimiter.ReqLimiter) *Client {
	if auth_token == "" {
		auth_token = os.Getenv("SOLSCAN_AUTH_TOKEN")
	}
	if limiter == nil {
		limiter = V2Limiter
	}
	return &Client{
		Limiter: limiter,
		Headers: map[string][]string{
			"content-type": {"application/json"},
			"token":        {auth_token},
		},
	}
}

func NewV2Client(auth_token string) *Client {
	return NewClient(auth_token, V2Limiter)
}

func NewV3Client(auth_token string) *Client {
	return NewClient(auth_token, V3Limiter)
}

func (c *Client) ChainInfo(ctx context.Context) (ChainInfo, error) {
	sg := SimpleGetter[ChainInfo]{
		BaseURL: PUBLIC_BASE_URL,
		Path:    "chaininfo",
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

type AccountTransfersParams struct {
	ActivityType      []AccountActivityType `json:"activity_type,omitempty"`
	TokenAccount      string                `json:"token_account,omitempty"`
	FromAddress       string                `json:"from,omitempty"`
	ToAddress         string                `json:"to,omitempty"`
	Token             string                `json:"token,omitempty"`
	AmountRange       []float64             `json:"amount,omitempty"`
	BlockTimeRange    []int64               `json:"block_time,omitempty"`
	ExcludeAmountZero bool                  `json:"exclude_amount_zero,omitempty"`
	Flow              Flow                  `json:"flow,omitempty"`
	SortBy            SortBy                `json:"sort_by" default:"block_time"`
	SortOrder         SortOrder             `json:"sort_order" default:"desc"`
	Page              int64                 `json:"page" default:"1"`
	PageSize          LargePageSize         `json:"page_size" default:"100"`
}

func (c *Client) AccountTransfers(ctx context.Context, address string, optParams *AccountTransfersParams) ([]Transfer, error) {
	params := CreateParams(optParams, "address", address)
	sg := SimpleGetter[[]Transfer]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/transfer",
		Params:  params,
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) AccountTransfersPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, address string, optParams *AccountTransfersParams) ([]Transfer, error) {
	params := CreateParams(optParams, "address", address)
	g := PagingGetter[[]Transfer]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/transfer",
		Params:  params,
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]Transfer]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[Transfer](int64(LargePageSize100)),
			ResultsHandler:    CreateSliceResultsHandler[Transfer](totalSize),
		},
	}
	return g.Do(ctx)
}

type AccountTokenAccountsParams struct {
	Type     TokenType     `json:"type" default:"token"`
	HideZero bool          `json:"hide_zero,omitempty"`
	Page     int64         `json:"page" default:"1"`
	PageSize SmallPageSize `json:"page_size" default:"40"`
}

func (c *Client) AccountTokenAccounts(ctx context.Context, address string, optParams *AccountTokenAccountsParams) ([]TokenAccount, error) {
	params := CreateParams(optParams, "address", address)
	sg := SimpleGetter[[]TokenAccount]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/token-accounts",
		Params:  params,
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) AccountTokenAccountsPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, address string, optParams *AccountTokenAccountsParams) ([]TokenAccount, error) {
	params := CreateParams(optParams, "address", address)
	g := PagingGetter[[]TokenAccount]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/token-accounts",
		Params:  params,
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]TokenAccount]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[TokenAccount](int64(SmallPageSize40)),
			ResultsHandler:    CreateSliceResultsHandler[TokenAccount](totalSize),
		},
	}
	return g.Do(ctx)
}

type AccountDefiActivitiesParams struct {
	ActivityType   []ActivityType `json:"activity_type,omitempty"`
	FromAddress    string         `json:"from_address,omitempty"`
	Platform       []string       `json:"platform,omitempty"`
	Source         []string       `json:"source,omitempty"`
	Token          string         `json:"token,omitempty"`
	BlockTimeRange []int64        `json:"block_time,omitempty"`
	Page           int64          `json:"page" default:"1"`
	PageSize       SmallPageSize  `json:"page_size" default:"40"`
	SortBy         SortBy         `json:"sort_by" default:"block_time"`
	SortOrder      SortOrder      `json:"sort_order" default:"desc"`
}

func (c *Client) AccountDefiActivities(ctx context.Context, address string, optParams *AccountDefiActivitiesParams) ([]DefiActivity, error) {
	params := CreateParams(optParams, "address", address)
	sg := SimpleGetter[[]DefiActivity]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/defi/activities",
		Params:  params,
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) AccountDefiActivitiesPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, address string, optParams *AccountDefiActivitiesParams) ([]DefiActivity, error) {
	params := CreateParams(optParams, "address", address)
	g := PagingGetter[[]DefiActivity]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/defi/activities",
		Params:  params,
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]DefiActivity]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[DefiActivity](int64(SmallPageSize40)),
			ResultsHandler:    CreateSliceResultsHandler[DefiActivity](totalSize),
		},
	}
	return g.Do(ctx)
}

type AccountBalanceChangesParams struct {
	Token          string        `json:"token,omitempty"`
	AmountRange    []float64     `json:"amount,omitempty"`
	BlockTimeRange []int64       `json:"block_time,omitempty"`
	Page           int64         `json:"page" default:"1"`
	PageSize       LargePageSize `json:"page_size" default:"100"`
	RemoveSpam     bool          `json:"remove_spam" default:"true"`
	Flow           Flow          `json:"flow,omitempty"`
	SortBy         SortBy        `json:"sort_by" default:"block_time"`
	SortOrder      SortOrder     `json:"sort_order" default:"desc"`
}

func (c *Client) AccountBalanceChanges(ctx context.Context, address string, optParams *AccountBalanceChangesParams) ([]AccountChangeActivity, error) {
	sg := SimpleGetter[[]AccountChangeActivity]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/balance_change",
		Params:  CreateParams(optParams, "address", address),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) AccountBalanceChangesPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, address string, optParams *AccountBalanceChangesParams) ([]AccountChangeActivity, error) {
	g := PagingGetter[[]AccountChangeActivity]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/balance_change",
		Params:  CreateParams(optParams, "address", address),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]AccountChangeActivity]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[AccountChangeActivity](int64(LargePageSize100)),
			ResultsHandler:    CreateSliceResultsHandler[AccountChangeActivity](totalSize),
		},
	}
	return g.Do(ctx)
}

type AccountTransactionsParams struct {
	Before string        `json:"before,omitempty"`
	Limit  SmallPageSize `json:"limit" default:"40"`
}

func (c *Client) AccountTransactions(ctx context.Context, address string, optParams *AccountTransactionsParams) ([]Transaction, error) {
	sg := SimpleGetter[[]Transaction]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/transactions",
		Params:  CreateParams(optParams, "address", address),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) AccountTransactionsPagingQuery(ctx context.Context, totalSize int64, address string, optParams *AccountTransactionsParams) ([]Transaction, error) {
	txs := []Transaction{}
	if optParams == nil {
		optParams = &AccountTransactionsParams{}
	}
	before := optParams.Before
	g := PagingGetter[[]Transaction]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/transactions",
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	pageNum := int(math.Ceil(float64(totalSize) / float64(SmallPageSize40)))
	for i := 0; i < pageNum; i++ {
		optParams.Before = before
		g.Params = CreateParams(optParams, "address", address)
		newTxs, err := g.Do(ctx)
		if err != nil {
			return nil, err
		}
		txs = append(txs, newTxs...)
		if len(newTxs) < int(SmallPageSize40) {
			break
		}
		before = newTxs[len(newTxs)-1].TxHash
	}
	return txs, nil
}

type AccountStakesParams struct {
	Page     int64         `json:"page" default:"1"`
	PageSize SmallPageSize `json:"page_size" default:"40"`
}

func (c *Client) AccountStakes(ctx context.Context, address string, optParams *AccountStakesParams) ([]AccountStake, error) {
	sg := SimpleGetter[[]AccountStake]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/stake",
		Params:  CreateParams(optParams, "address", address),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) AccountStakesPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, address string, optParams *AccountStakesParams) ([]AccountStake, error) {
	g := PagingGetter[[]AccountStake]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/stake",
		Headers: c.Headers,
		Limiter: c.Limiter,
		Params:  CreateParams(optParams, "address", address),
		PagingParams: &PagingParams[[]AccountStake]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[AccountStake](int64(SmallPageSize40)),
			ResultsHandler:    CreateSliceResultsHandler[AccountStake](totalSize),
		},
	}
	return g.Do(ctx)
}

func (c *Client) AccountDetail(ctx context.Context, address string) (AccountDetail, error) {
	sg := SimpleGetter[AccountDetail]{
		BaseURL: PRO_BASE_URL,
		Path:    "/account/detail",
		Params:  url.Values{"address": {address}},
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) AccountRewardsExport(ctx context.Context, address string, timeFrom, timeTo int64) ([]byte, error) {
	sg := SimpleGetter[[]byte]{
		BaseURL:           PRO_BASE_URL,
		Path:              "/account/reward/export",
		Params:            url.Values{"address": {address}, "time_from": {strconv.FormatInt(timeFrom, 10)}, "time_to": {strconv.FormatInt(timeTo, 10)}},
		Headers:           c.Headers,
		Limiter:           c.Limiter,
		RespBodyUnmarshal: ExportBodyUnmarshal,
	}
	return sg.Do(ctx)
}

type AccountTransfersExportParams struct {
	ActivityType      []AccountActivityType `json:"activity_type,omitempty"`
	TokenAccount      string                `json:"token_account,omitempty"`
	FromAddress       string                `json:"from,omitempty"`
	ToAddress         string                `json:"to,omitempty"`
	Token             string                `json:"token,omitempty"`
	AmountRange       []float64             `json:"amount,omitempty"`
	BlockTimeRange    []int64               `json:"block_time,omitempty"`
	ExcludeAmountZero bool                  `json:"exclude_amount_zero,omitempty"`
	Flow              Flow                  `json:"flow,omitempty"`
}

func (c *Client) AccountTransfersExport(ctx context.Context, address string, optParams *AccountTransfersExportParams) ([]byte, error) {
	sg := SimpleGetter[[]byte]{
		BaseURL:           PRO_BASE_URL,
		Path:              "/account/transfer/export",
		Params:            CreateParams(optParams, "address", address),
		Headers:           c.Headers,
		Limiter:           c.Limiter,
		RespBodyUnmarshal: ExportBodyUnmarshal,
	}
	return sg.Do(ctx)
}

type TokenTransfersParams struct {
	ActivityType      []ActivityType `json:"activity_type,omitempty"`
	FromAddress       string         `json:"from,omitempty"`
	ToAddress         string         `json:"to,omitempty"`
	AmountRange       []float64      `json:"amount,omitempty"`
	BlockTimeRange    []int64        `json:"block_time,omitempty"`
	ExcludeAmountZero bool           `json:"exclude_amount_zero,omitempty"`
	Page              int64          `json:"page" default:"1"`
	PageSize          LargePageSize  `json:"page_size" default:"100"`
	SortBy            SortBy         `json:"sort_by" default:"block_time"`
	SortOrder         SortOrder      `json:"sort_order" default:"desc"`
}

func (c *Client) TokenTransfers(ctx context.Context, address string, optParams *TokenTransfersParams) ([]Transfer, error) {
	sg := SimpleGetter[[]Transfer]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/transfer",
		Params:  CreateParams(optParams, "address", address),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) TokenTransfersPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, address string, optParams *TokenTransfersParams) ([]Transfer, error) {
	g := PagingGetter[[]Transfer]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/transfer",
		Params:  CreateParams(optParams, "address", address),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]Transfer]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[Transfer](int64(LargePageSize100)),
			ResultsHandler:    CreateSliceResultsHandler[Transfer](totalSize),
		},
	}
	return g.Do(ctx)
}

type TokenDefiActivitiesParams struct {
	FromAddress    string         `json:"from_address,omitempty"`
	Platform       []string       `json:"platform,omitempty"`
	Source         []string       `json:"source,omitempty"`
	ActivityType   []ActivityType `json:"activity_type,omitempty"`
	Token          string         `json:"token,omitempty"`
	BlockTimeRange []int64        `json:"block_time,omitempty"`
	Page           int64          `json:"page" default:"1"`
	PageSize       LargePageSize  `json:"page_size" default:"100"`
	SortBy         SortBy         `json:"sort_by" default:"block_time"`
	SortOrder      SortOrder      `json:"sort_order" default:"desc"`
}

func (c *Client) TokenDefiActivities(ctx context.Context, address string, optParams *TokenDefiActivitiesParams) ([]DefiActivity, error) {
	sg := SimpleGetter[[]DefiActivity]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/defi/activities",
		Params:  CreateParams(optParams, "address", address),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) TokenDefiActivitiesPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, address string, optParams *TokenDefiActivitiesParams) ([]DefiActivity, error) {
	g := PagingGetter[[]DefiActivity]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/defi/activities",
		Params:  CreateParams(optParams, "address", address),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]DefiActivity]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[DefiActivity](int64(LargePageSize100)),
			ResultsHandler:    CreateSliceResultsHandler[DefiActivity](totalSize),
		},
	}
	return g.Do(ctx)
}

type TokenMarketsParams struct {
	Program  string        `json:"program,omitempty"`
	Page     int64         `json:"page" default:"1"`
	PageSize LargePageSize `json:"page_size" default:"100"`
	SortBy   MarketSortBy  `json:"sort_by" default:"volume"`
}

func (c *Client) TokenMarkets(ctx context.Context, token_pair []string, optParams *TokenMarketsParams) ([]Market, error) {
	sg := SimpleGetter[[]Market]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/markets",
		Params:  CreateParams(optParams, "token[]", token_pair[0], "token[]", token_pair[1]),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) TokenMarketsPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, token_pair []string, optParams *TokenMarketsParams) ([]Market, error) {
	g := PagingGetter[[]Market]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/markets",
		Params:  CreateParams(optParams, "token[]", token_pair[0], "token[]", token_pair[1]),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]Market]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[Market](int64(LargePageSize100)),
			ResultsHandler:    CreateSliceResultsHandler[Market](totalSize),
		},
	}
	return g.Do(ctx)
}

type TokenListParams struct {
	SortBy    TokenSortBy   `json:"sort_by" default:"price"`
	SortOrder SortOrder     `json:"sort_order" default:"desc"`
	Page      int64         `json:"page" default:"1"`
	PageSize  LargePageSize `json:"page_size" default:"100"`
}

func (c *Client) TokenList(ctx context.Context, optParams *TokenListParams) ([]Token, error) {
	sg := SimpleGetter[[]Token]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/list",
		Params:  CreateParams(optParams),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) TokenListPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, optParams *TokenListParams) ([]Token, error) {
	g := PagingGetter[[]Token]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/list",
		Params:  CreateParams(optParams),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]Token]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[Token](int64(LargePageSize100)),
			ResultsHandler:    CreateSliceResultsHandler[Token](totalSize),
		},
	}
	return g.Do(ctx)
}

func (c *Client) TokenTrending(ctx context.Context, limit int64) ([]Token, error) {
	g := SimpleGetter[[]Token]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/trending",
		Params:  url.Values{"limit": {strconv.FormatInt(limit, 10)}},
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return g.Do(ctx)
}

// Time is YYYYMMDD
func (c *Client) TokenPrice(ctx context.Context, address, startTime, endTime string) ([]TokenPrice, error) {
	params := url.Values{"address": {address}}
	params.Add("time[]", startTime)
	params.Add("time[]", endTime)
	sg := SimpleGetter[[]TokenPrice]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/price",
		Params:  params,
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

type TokenHoldersParams struct {
	FromAmount string        `json:"from_amount,omitempty"`
	ToAmount   string        `json:"to_amount,omitempty"`
	Page       int64         `json:"page" default:"1"`
	PageSize   SmallPageSize `json:"page_size" default:"40"`
}

func (c *Client) TokenHolders(ctx context.Context, address string, optParams *TokenHoldersParams) (RespDataWithTotal[TokenHolder], error) {
	sg := SimpleGetter[RespDataWithTotal[TokenHolder]]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/holders",
		Params:  CreateParams(optParams, "address", address),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) TokenHoldersPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, address string, optParams *TokenHoldersParams) (RespDataWithTotal[TokenHolder], error) {
	g := PagingGetter[RespDataWithTotal[TokenHolder]]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/holders",
		Params:  CreateParams(optParams, "address", address),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[RespDataWithTotal[TokenHolder]]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateWithTotalItemsFinishChecker[TokenHolder](int64(SmallPageSize40)),
			ResultsHandler:    CreateWithTotalItemsResultsHandler[TokenHolder](totalSize),
		},
	}
	return g.Do(ctx)
}

func (c *Client) TokenMeta(ctx context.Context, address string) (TokenMeta, error) {
	sg := SimpleGetter[TokenMeta]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/meta",
		Params:  url.Values{"address": {address}},
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) TokenTop(ctx context.Context) ([]TokenTop, error) {
	sg := SimpleGetter[[]TokenTop]{
		BaseURL: PRO_BASE_URL,
		Path:    "/token/top",
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

type NFTNewsParams struct {
	Filter   string       `json:"filter" default:"created_time"`
	Page     int64        `json:"page" default:"1"`
	PageSize TinyPageSize `json:"page_size" default:"36"`
}

func (c *Client) NFTNews(ctx context.Context, optParams *NFTNewsParams) (RespDataWithTotal[NFTInfo], error) {
	sg := SimpleGetter[RespDataWithTotal[NFTInfo]]{
		BaseURL: PRO_BASE_URL,
		Path:    "/nft/news",
		Params:  CreateParams(optParams),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) NFTNewsPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, optParams *NFTNewsParams) (RespDataWithTotal[NFTInfo], error) {
	g := PagingGetter[RespDataWithTotal[NFTInfo]]{
		BaseURL: PRO_BASE_URL,
		Path:    "/nft/news",
		Params:  CreateParams(optParams),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[RespDataWithTotal[NFTInfo]]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateWithTotalDataFinishChecker[NFTInfo](int64(TinyPageSize36)),
			ResultsHandler:    CreateWithTotalDataResultsHandler[NFTInfo](totalSize),
		},
	}
	return g.Do(ctx)
}

type NFTActivitiesParams struct {
	FromAddress    string            `json:"from,omitempty"`
	ToAddress      string            `json:"to,omitempty"`
	Source         []string          `json:"source,omitempty"`
	ActivityType   []NFTActivityType `json:"activity_type,omitempty"`
	Token          string            `json:"token,omitempty"`
	Collection     string            `json:"collection,omitempty"`
	CurrencyToken  string            `json:"currency_token,omitempty"`
	PriceRange     []float64         `json:"price,omitempty"`
	BlockTimeRange []int64           `json:"block_time,omitempty"`
	Page           int64             `json:"page" default:"1"`
	PageSize       LargePageSize     `json:"page_size" default:"100"`
}

func (c *Client) NFTActivities(ctx context.Context, optParams *NFTActivitiesParams) ([]NFTActivity, error) {
	sg := SimpleGetter[[]NFTActivity]{
		BaseURL: PRO_BASE_URL,
		Path:    "/nft/activities",
		Params:  CreateParams(optParams),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) NFTActivitiesPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, optParams *NFTActivitiesParams) ([]NFTActivity, error) {
	g := PagingGetter[[]NFTActivity]{
		BaseURL: PRO_BASE_URL,
		Path:    "/nft/activities",
		Params:  CreateParams(optParams),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]NFTActivity]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[NFTActivity](int64(LargePageSize100)),
			ResultsHandler:    CreateSliceResultsHandler[NFTActivity](totalSize),
		},
	}
	return g.Do(ctx)
}

type NFTCollectionListParams struct {
	Collection string              `json:"collection,omitempty"`
	SortBy     NFTCollectionSortBy `json:"sort_by" default:"floor_price"`
	SortOrder  SortOrder           `json:"sort_order" default:"desc"`
	Page       int64               `json:"page" default:"1"`
	PageSize   SmallPageSize       `json:"page_size" default:"40"`
}

func (c *Client) NFTCollectionList(ctx context.Context, optParams *NFTCollectionListParams) ([]NFTCollection, error) {
	sg := SimpleGetter[[]NFTCollection]{
		BaseURL: PRO_BASE_URL,
		Path:    "/nft/collection/lists",
		Params:  CreateParams(optParams),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) NFTCollectionListPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, optParams *NFTCollectionListParams) ([]NFTCollection, error) {
	g := PagingGetter[[]NFTCollection]{
		BaseURL: PRO_BASE_URL,
		Path:    "/nft/collection/lists",
		Params:  CreateParams(optParams),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]NFTCollection]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[NFTCollection](int64(LargePageSize100)),
			ResultsHandler:    CreateSliceResultsHandler[NFTCollection](totalSize),
		},
	}
	return g.Do(ctx)
}

type NFTCollectionItemsParams struct {
	SortBy   NFTCollectionItemSortBy `json:"sort_by" default:"last_trade"`
	Page     int64                   `json:"page" default:"1"`
	PageSize TinyPageSize            `json:"page_size" default:"36"`
}

func (c *Client) NFTCollectionItems(ctx context.Context, collection string, optParams *NFTCollectionItemsParams) ([]NFTCollectionItem, error) {
	sg := SimpleGetter[[]NFTCollectionItem]{
		BaseURL: PRO_BASE_URL,
		Path:    "/nft/collection/items",
		Params:  CreateParams(optParams, "collection", collection),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) NFTCollectionItemsPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, collection string, optParams *NFTCollectionItemsParams) ([]NFTCollectionItem, error) {
	g := PagingGetter[[]NFTCollectionItem]{
		BaseURL: PRO_BASE_URL,
		Path:    "/nft/collection/items",
		Params:  CreateParams(optParams, "collection", collection),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]NFTCollectionItem]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[NFTCollectionItem](int64(TinyPageSize36)),
			ResultsHandler:    CreateSliceResultsHandler[NFTCollectionItem](totalSize),
		},
	}
	return g.Do(ctx)
}

type TxLastParams struct {
	Limit  LargePageSize `json:"limit" default:"100"`
	Filter TxFilter      `json:"filter" default:"all"`
}

func (c *Client) TxLast(ctx context.Context, optParams *TxLastParams) ([]Transaction, error) {
	sg := SimpleGetter[[]Transaction]{
		BaseURL: PRO_BASE_URL,
		Path:    "/transaction/last",
		Params:  CreateParams(optParams),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) TxDetail(ctx context.Context, tx string) (TransactionDetail, error) {
	sg := SimpleGetter[TransactionDetail]{
		BaseURL: PRO_BASE_URL,
		Path:    "/transaction/detail",
		Params:  url.Values{"tx": {tx}},
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) TxActions(ctx context.Context, tx string) (TransactionAction, error) {
	sg := SimpleGetter[TransactionAction]{
		BaseURL: PRO_BASE_URL,
		Path:    "/transaction/actions",
		Params:  url.Values{"tx": {tx}},
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) BlocksLast(ctx context.Context, limit LargePageSize) ([]BlockDetail, error) {
	sg := SimpleGetter[[]BlockDetail]{
		BaseURL: PRO_BASE_URL,
		Path:    "/block/last",
		Params:  url.Values{"limit": {strconv.FormatInt(int64(limit), 10)}},
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

type BlockTransactionsParams struct {
	Page     int64         `json:"page" default:"1"`
	PageSize LargePageSize `json:"page_size" default:"100"`
}

func (c *Client) BlockTransactions(ctx context.Context, block int64, optParams *BlockTransactionsParams) (RespDataWithTotal[Transaction], error) {
	sg := SimpleGetter[RespDataWithTotal[Transaction]]{
		BaseURL: PRO_BASE_URL,
		Path:    "/block/transactions",
		Params:  CreateParams(optParams, "block", strconv.FormatInt(block, 10)),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) BlockTransactionsPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency, block int64, optParams *BlockTransactionsParams) (RespDataWithTotal[Transaction], error) {
	g := PagingGetter[RespDataWithTotal[Transaction]]{
		BaseURL: PRO_BASE_URL,
		Path:    "/block/transactions",
		Params:  CreateParams(optParams, "block", strconv.FormatInt(block, 10)),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[RespDataWithTotal[Transaction]]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateWithTotalTransactionsDataFinishChecker(int64(LargePageSize100)),
			ResultsHandler:    CreateWithTotalTransactionsResultsHandler(totalSize),
		},
	}
	return g.Do(ctx)
}

func (c *Client) BlockDetail(ctx context.Context, block int64) (BlockDetail, error) {
	sg := SimpleGetter[BlockDetail]{
		BaseURL: PRO_BASE_URL,
		Path:    "/block/detail",
		Params:  url.Values{"block": {strconv.FormatInt(block, 10)}},
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

type PoolMarketListParams struct {
	Program   string        `json:"program,omitempty"`
	SortBy    string        `json:"sort_by" default:"created_time"`
	SortOrder SortOrder     `json:"sort_order" default:"desc"`
	Page      int64         `json:"page" default:"1"`
	PageSize  LargePageSize `json:"page_size" default:"100"`
}

func (c *Client) PoolMarketList(ctx context.Context, optParams *PoolMarketListParams) ([]PoolMarket, error) {
	sg := SimpleGetter[[]PoolMarket]{
		BaseURL: PRO_BASE_URL,
		Path:    "/market/list",
		Params:  CreateParams(optParams),
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) PoolMarketListPagingQuery(ctx context.Context, startPage, totalSize, maxConcurrency int64, optParams *PoolMarketListParams) ([]PoolMarket, error) {
	g := PagingGetter[[]PoolMarket]{
		BaseURL: PRO_BASE_URL,
		Path:    "/market/list",
		Params:  CreateParams(optParams),
		Headers: c.Headers,
		Limiter: c.Limiter,
		GetterOption: &GetterOption{
			RetryInterval: time.Second,
			MaxRetries:    100,
		},
		PagingParams: &PagingParams[[]PoolMarket]{
			StartPage:         startPage,
			TotalSize:         totalSize,
			MaxConcurrency:    maxConcurrency,
			DataFinishChecker: CreateSliceDataFinishChecker[PoolMarket](int64(LargePageSize100)),
			ResultsHandler:    CreateSliceResultsHandler[PoolMarket](totalSize),
		},
	}
	return g.Do(ctx)
}

func (c *Client) PoolMarketInfo(ctx context.Context, address string) (PoolMarketInfo, error) {
	sg := SimpleGetter[PoolMarketInfo]{
		BaseURL: PRO_BASE_URL,
		Path:    "/market/info",
		Params:  url.Values{"address": {address}},
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

// Time is YYYYMMDD
func (c *Client) PoolMarketVolume(ctx context.Context, address string, startTime, endTime string) (PoolMarketVolume, error) {
	params := url.Values{"address": {address}}
	params.Add("time[]", startTime)
	params.Add("time[]", endTime)
	sg := SimpleGetter[PoolMarketVolume]{
		BaseURL: PRO_BASE_URL,
		Path:    "/market/volume",
		Params:  params,
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}

func (c *Client) APIUsage(ctx context.Context) (APIUsage, error) {
	sg := SimpleGetter[APIUsage]{
		BaseURL: PRO_BASE_URL,
		Path:    "/monitor/usage",
		Headers: c.Headers,
		Limiter: c.Limiter,
	}
	return sg.Do(ctx)
}
