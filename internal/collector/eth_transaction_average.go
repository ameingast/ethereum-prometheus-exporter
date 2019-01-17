package collector

import (
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/prometheus/client_golang/prometheus"
)

type EthTransactionAverage struct {
	rpc  *rpc.Client
	desc *prometheus.Desc
}

func NewEthTransactionAverage(rpc *rpc.Client) *EthTransactionAverage {
	return &EthTransactionAverage{
		rpc: rpc,
		desc: prometheus.NewDesc(
			"eth_average_transactions",
			"the average number of transactions per second over last 5 blocks",
			nil,
			nil,
		),
	}
}

func (collector *EthTransactionAverage) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.desc
}

func (collector *EthTransactionAverage) Collect(ch chan<- prometheus.Metric) {
	const amountOfBlocks = 5
	var blockNumberResult hexutil.Uint64 = 0
	var totalTransactions = uint64(0)

	if err := collector.rpc.Call(&blockNumberResult, "eth_blockNumber"); err != nil {
		ch <- prometheus.NewInvalidMetric(collector.desc, err)
		return
	}

	for i := uint64(0); i < amountOfBlocks; i++ {
		var txCountResult hexutil.Uint64 = 0
		block := uint64(blockNumberResult) - i
		if err := collector.rpc.Call(&txCountResult, "eth_getBlockTransactionCountByNumber", hexutil.EncodeUint64(block)); err != nil {
			ch <- prometheus.NewInvalidMetric(collector.desc, err)
			return
		}
		totalTransactions += uint64(txCountResult)
	}

	ch <- prometheus.MustNewConstMetric(collector.desc, prometheus.GaugeValue, float64(totalTransactions) / amountOfBlocks)
}
