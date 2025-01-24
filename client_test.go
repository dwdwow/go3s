package go3s

import (
	"context"
	"fmt"
	"testing"
)

func TestChainInfo(t *testing.T) {
	client := NewV2Client("")
	chainInfo, err := client.ChainInfo(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", chainInfo)
}

func TestAccountTransfers(t *testing.T) {
	client := NewV2Client("")
	transfers, err := client.AccountTransfers(
		context.Background(),
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Printf("%+v\n", activity)
	fmt.Println(len(transfers))
}

func TestAccountTransfersPagingQuery(t *testing.T) {
	client := NewV2Client("")
	transfers, err := client.AccountTransfersPagingQuery(
		context.Background(),
		1,
		588,
		1,
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(transfers))
}

func TestAccountTokens(t *testing.T) {
	client := NewV2Client("")
	tokens, err := client.AccountTokenAccounts(
		context.Background(),
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(tokens))
}

func TestAccountTokenAccountsPagingQuery(t *testing.T) {
	client := NewV2Client("")
	transfers, err := client.AccountTokenAccountsPagingQuery(
		context.Background(),
		1,
		388,
		10,
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(transfers))
}

func TestAccountDefiActivities(t *testing.T) {
	client := NewV2Client("")
	activities, err := client.AccountDefiActivities(
		context.Background(),
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(activities))
}

func TestAccountDefiActivitiesPagingQuery(t *testing.T) {
	client := NewV2Client("")
	activities, err := client.AccountDefiActivitiesPagingQuery(
		context.Background(),
		1,
		388,
		10,
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(activities))
}

func TestAccountBalanceChanges(t *testing.T) {
	client := NewV2Client("")
	activities, err := client.AccountBalanceChanges(
		context.Background(),
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(activities))
}

func TestAccountBalanceChangesPagingQuery(t *testing.T) {
	client := NewV2Client("")
	activities, err := client.AccountBalanceChangesPagingQuery(
		context.Background(),
		1,
		388,
		10,
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(activities))
}

func TestAccountTransactions(t *testing.T) {
	client := NewV2Client("")
	transactions, err := client.AccountTransactions(
		context.Background(),
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(transactions))
}

func TestAccountTransactionsPagingQuery(t *testing.T) {
	client := NewV2Client("")
	transactions, err := client.AccountTransactionsPagingQuery(
		context.Background(),
		1000,
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(transactions))
}

func TestAccountStakes(t *testing.T) {
	client := NewV2Client("")
	stakes, err := client.AccountStakes(
		context.Background(),
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(stakes))
}

func TestAccountStakesPagingQuery(t *testing.T) {
	client := NewV2Client("")
	stakes, err := client.AccountStakesPagingQuery(
		context.Background(),
		1,
		388,
		10,
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(stakes))
}

func TestAccountDetail(t *testing.T) {
	client := NewV2Client("")
	detail, err := client.AccountDetail(
		context.Background(),
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", detail)
}

func TestAccountRewardsExport(t *testing.T) {
	client := NewV2Client("")
	rewards, err := client.AccountRewardsExport(
		context.Background(),
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		1716672000,
		1716758400,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(rewards))
}

func TestAccountTransfersExport(t *testing.T) {
	client := NewV2Client("")
	transfers, err := client.AccountTransfersExport(
		context.Background(),
		"3zAQJcPLbfi2mwnPraQpfuNFh5h5PN7XLkNDJSZ5i7E5",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(transfers))
}

func TestTokenTransfers(t *testing.T) {
	client := NewV2Client("")
	transfers, err := client.TokenTransfers(
		context.Background(),
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(transfers))
}

func TestTokenTransfersPagingQuery(t *testing.T) {
	client := NewV2Client("")
	transfers, err := client.TokenTransfersPagingQuery(
		context.Background(),
		1,
		388,
		10,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(transfers))
}

func TestTokenDefiActivities(t *testing.T) {
	client := NewV2Client("")
	activities, err := client.TokenDefiActivities(
		context.Background(),
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(activities))
}

func TestTokenDefiActivitiesPagingQuery(t *testing.T) {
	client := NewV2Client("")
	activities, err := client.TokenDefiActivitiesPagingQuery(
		context.Background(),
		1,
		388,
		10,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(activities))
}

func TestTokenMarkets(t *testing.T) {
	client := NewV2Client("")
	markets, err := client.TokenMarkets(
		context.Background(),
		[]string{"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC", "So11111111111111111111111111111111111111112"},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(markets))
}

func TestTokenMarketsPagingQuery(t *testing.T) {
	client := NewV2Client("")
	markets, err := client.TokenMarketsPagingQuery(
		context.Background(),
		1,
		388,
		10,
		[]string{"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC", "So11111111111111111111111111111111111111112"},
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(markets))
}

func TestTokenList(t *testing.T) {
	client := NewV2Client("")
	tokens, err := client.TokenList(
		context.Background(),
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(tokens))
}

func TestTokenListPagingQuery(t *testing.T) {
	client := NewV2Client("")
	tokens, err := client.TokenListPagingQuery(
		context.Background(),
		1,
		388,
		10,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(tokens))
}

func TestTokenTrending(t *testing.T) {
	client := NewV2Client("")
	tokens, err := client.TokenTrending(
		context.Background(),
		100,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(tokens))
}

func TestTokenPrice(t *testing.T) {
	client := NewV2Client("")
	price, err := client.TokenPrice(
		context.Background(),
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		"20250101",
		"20250120",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(price)
}

func TestTokenHolders(t *testing.T) {
	client := NewV2Client("")
	holders, err := client.TokenHolders(
		context.Background(),
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(holders.Items), holders.Total)
}

func TestTokenHoldersPagingQuery(t *testing.T) {
	client := NewV2Client("")
	holders, err := client.TokenHoldersPagingQuery(
		context.Background(),
		1,
		388,
		10,
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(holders.Items), holders.Total)
}

func TestTokenMeta(t *testing.T) {
	client := NewV2Client("")
	meta, err := client.TokenMeta(
		context.Background(),
		"HeLp6NuQkmYB4pYWo2zYs22mESHXPQYzXbB8n4V98jwC",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(meta)
}

func TestTokenTop(t *testing.T) {
	client := NewV2Client("")
	tokens, err := client.TokenTop(
		context.Background(),
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(tokens))
}

func TestTxLast(t *testing.T) {
	client := NewV2Client("")
	txs, err := client.TxLast(
		context.Background(),
		&TxLastParams{Limit: LargePageSize100},
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(txs))
}

func TestTxDetail(t *testing.T) {
	client := NewV2Client("")
	tx, err := client.TxDetail(
		context.Background(),
		"3uf5w7XnMBd4xZTRTqSzi1L9S1QBFCMxAutTd4vCAdPWnjf5b815Ng2GfQwVRf4qHxHDDFWHT6tndHSD88HQkbSk",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", tx)
}

func TestTxActions(t *testing.T) {
	client := NewV2Client("")
	actions, err := client.TxActions(
		context.Background(),
		"3uf5w7XnMBd4xZTRTqSzi1L9S1QBFCMxAutTd4vCAdPWnjf5b815Ng2GfQwVRf4qHxHDDFWHT6tndHSD88HQkbSk",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", actions)
}

func TestBlockLast(t *testing.T) {
	client := NewV2Client("")
	blocks, err := client.BlocksLast(context.Background(), LargePageSize100)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", blocks)
}

func TestBlockTransactions(t *testing.T) {
	client := NewV2Client("")
	transactions, err := client.BlockTransactions(
		context.Background(),
		315989642,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(transactions.Transactions), transactions.Total)
}

func TestBlockTransactionsPagingQuery(t *testing.T) {
	client := NewV2Client("")
	transactions, err := client.BlockTransactionsPagingQuery(
		context.Background(),
		1,
		388,
		10,
		315989642,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(transactions.Transactions), transactions.Total)
}

func TestBlockDetail(t *testing.T) {
	client := NewV2Client("")
	block, err := client.BlockDetail(
		context.Background(),
		315989642,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", block)
}

func TestPoolMarketList(t *testing.T) {
	client := NewV2Client("")
	markets, err := client.PoolMarketList(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(markets))
}

func TestPoolMarketListPagingQuery(t *testing.T) {
	client := NewV2Client("")
	markets, err := client.PoolMarketListPagingQuery(
		context.Background(),
		1,
		388,
		10,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(markets))
}

func TestPoolMarketInfo(t *testing.T) {
	client := NewV2Client("")
	info, err := client.PoolMarketInfo(
		context.Background(),
		"44W73kGYQgXCTNkGxUmHv8DDBPCxojBcX49uuKmbFc9U",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", info)
}

func TestPoolMarketVolume(t *testing.T) {
	client := NewV2Client("")
	volume, err := client.PoolMarketVolume(
		context.Background(),
		"FDxGM9n4UQjUunjb43be1hs8oYFAPYziX1bWWc212dVU",
		"20250121",
		"20250123",
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", volume)
}

func TestAPIUsage(t *testing.T) {
	client := NewV2Client("")
	usage, err := client.APIUsage(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", usage)
}

func TestNFTNews(t *testing.T) {
	client := NewV2Client("")
	news, err := client.NFTNews(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(news.Data))
}

func TestNFTNewsPagingQuery(t *testing.T) {
	client := NewV2Client("")
	news, err := client.NFTNewsPagingQuery(
		context.Background(),
		1,
		388,
		10,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(news.Data))
}

func TestNFTActivities(t *testing.T) {
	client := NewV2Client("")
	activities, err := client.NFTActivities(
		context.Background(),
		&NFTActivitiesParams{
			Source: []string{"opensea"},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(activities))
}

func TestNFTActivitiesPagingQuery(t *testing.T) {
	client := NewV2Client("")
	activities, err := client.NFTActivitiesPagingQuery(
		context.Background(),
		1,
		388,
		10,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(activities))
}

func TestNFTCollectionList(t *testing.T) {
	client := NewV2Client("")
	collections, err := client.NFTCollectionList(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(collections))
}

func TestNFTCollectionListPagingQuery(t *testing.T) {
	client := NewV2Client("")
	collections, err := client.NFTCollectionListPagingQuery(
		context.Background(),
		1,
		388,
		10,
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(collections))
}

func TestNFTCollectionItems(t *testing.T) {
	client := NewV2Client("")
	items, err := client.NFTCollectionItems(
		context.Background(),
		"CY2E69dSG9vBsMoaXDvYmMDSMEP4SZtRY1rqVQ9tkNDu",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(items))
}

func TestNFTCollectionItemsPagingQuery(t *testing.T) {
	client := NewV2Client("")
	items, err := client.NFTCollectionItemsPagingQuery(
		context.Background(),
		1,
		388,
		10,
		"CY2E69dSG9vBsMoaXDvYmMDSMEP4SZtRY1rqVQ9tkNDu",
		nil,
	)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(len(items))
}
