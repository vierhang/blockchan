package BLC

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CLI struct {
}

// PrintUsage 用法展示
func PrintUsage() {
	fmt.Println("Usage:")
	//初始化区块链
	fmt.Printf("\tcreageblockchain -address address -- 创建区块链 \n")
	// 添加区块
	fmt.Printf("\taddblock -data DATA --添加区块\n")
	// 打印完整区块信息
	fmt.Printf("\tprintchain --输出区块链信息\n")
	// 通过命令行转账
	fmt.Printf("\t-from FROM -to TO -amount AMOUNT --发起转账\n")
	fmt.Printf("\t转账参数说明\n")
	fmt.Printf("\t\t-from FROM -- 转账原地址\n")
	fmt.Printf("\t\t-to TO -- 转账目标地址\n")
	fmt.Printf("\t\t-amount AMOUNT -- 转账金额\n")

	//查询余额
	fmt.Printf("\t getbalance -address FROM -- 查询指定地址余额\n")
	fmt.Printf("\t 查询余额地址参数说明\n")
	fmt.Printf("\t\t  -address FROM -- 指定地址余额\n")

}

func (c *CLI) getBalance(from string) {
	// 查找该地址UTXO
	// 获取区块链对象
	blockChanObj := NewBlockChain()
	defer blockChanObj.DB.Close()
	amount := blockChanObj.getBalance(from)
	fmt.Printf("\t地址【%s】的余额：【%d】\n", from, amount)
}

// 发起交易
func (c *CLI) send(from, to, amount []string) {
	if !dbExists() {
		fmt.Println("数据库不存在")
		os.Exit(1)
	}
	// 获取区块链对象
	blockchainObj := NewBlockChain()
	defer blockchainObj.DB.Close()
	blockchainObj.MineNewBlock(from, to, amount)
}

// CreateBlockChain 初始化区块链
func (c *CLI) CreateBlockChain(address string) {
	CreateBlockChainWithGenesisBLock(address)
}

//添加区块
func (c *CLI) AddBlock(txs []*Transaction) {
	if !dbExists() {
		fmt.Println("dbExists")
		os.Exit(1)
	}
	NewBlockChain().AddBlock(txs)
}

// PrintChain 打印完整区块链信息
func (c *CLI) PrintChain() {
	if !dbExists() {
		fmt.Println("dbExists")
		os.Exit(1)
	}
	NewBlockChain().PrintChan()
}

// IsValidArgs 参数数量检测函数
func IsValidArgs() {
	if len(os.Args) < 2 {
		PrintUsage()
		fmt.Println("os.Args < 2")
		os.Exit(1)
	}
}

// Run 命令行运行函数
func (c *CLI) Run() {
	IsValidArgs()
	// 新建相关命令
	// 添加区块
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)
	createBLCWithGenesisBlockCmd := flag.NewFlagSet("cerateblockchain", flag.ExitOnError)
	// 查询余额命令
	getBalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)

	//发起交易
	// bcli send -from "[\"test\"]" -to "[\"b\"]" -amount "[\"20\"]"
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	// 数据参数处理
	flagAddBlockArg := addBlockCmd.String("data", "sent 100 btc to player", "添加区块数据")
	// 创建区块时指定的矿工地址
	flagCreageBlockchainArg := createBLCWithGenesisBlockCmd.String("address", "weihang", "指定接受系统奖励的矿工地址")

	//发起交易参数
	flagSendFromArg := sendCmd.String("from", "", "转账原地址")
	flagSendToArg := sendCmd.String("to", "", "转账目标地址")
	flagSendAmountArg := sendCmd.String("amount", "", "转账金额")
	//查询余额命令行参数
	flagGetBalanceArg := getBalanceCmd.String("address", "", "要查询余额的地址")
	// 判断命令
	switch os.Args[1] {
	case "getbalance":
		if err := getBalanceCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse getbalance failed ! %v\n", err)
		}
	case "send":
		if err := sendCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse sendCmd failed ! %v\n", err)
		}

	case "addblock":
		if err := addBlockCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse addBlockCmd err %v \n", err)
		}
	case "printchain":
		if err := printChainCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse printChainCmd err %v \n", err)
		}
	case "createblockchain":
		if err := createBLCWithGenesisBlockCmd.Parse(os.Args[2:]); err != nil {
			log.Panicf("parse createBLCWithGenesisBlockCmd err %v \n", err)
		}
	default:
		PrintUsage()
		fmt.Println(os.Args[1])
		fmt.Println("os.Args[1] switch error")
		os.Exit(1)
	}
	// 查询余额
	if getBalanceCmd.Parsed() {
		if *flagGetBalanceArg == "" {
			fmt.Println("请输入查询地址")
			os.Exit(1)
		}
		c.getBalance(*flagGetBalanceArg)
	}

	// 发起转账
	if sendCmd.Parsed() {
		if *flagSendFromArg == "" {
			fmt.Println("原地址不能为空")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendToArg == "" {
			fmt.Println("目标地址不能为空")
			PrintUsage()
			os.Exit(1)
		}
		if *flagSendAmountArg == "" {
			fmt.Println("转账金额不能为空")
			PrintUsage()
			os.Exit(1)
		}
		fmt.Printf("\tFROM:[%s]\n", JSONToSlice(*flagSendFromArg))
		fmt.Printf("\tTO:[%s]\n", JSONToSlice(*flagSendToArg))
		fmt.Printf("\tAMOUNT:[%s]\n", JSONToSlice(*flagSendAmountArg))
		c.send(JSONToSlice(*flagSendFromArg), JSONToSlice(*flagSendToArg), JSONToSlice(*flagSendAmountArg))
	}
	// 添加区块命令
	if addBlockCmd.Parsed() {
		if *flagAddBlockArg == "" {
			PrintUsage()
			os.Exit(1)
		}
		c.AddBlock([]*Transaction{})
	}
	// 输出区块信息
	if printChainCmd.Parsed() {
		c.PrintChain()
	}
	// 创建区块链命令
	if createBLCWithGenesisBlockCmd.Parsed() {
		if *flagCreageBlockchainArg == "" {
			fmt.Println("flagCreageBlockchainArg == ''")
			PrintUsage()
			os.Exit(1)
		}
		c.CreateBlockChain(*flagCreageBlockchainArg)
	}

}
