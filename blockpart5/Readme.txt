目录：
block.go 		//区块相关的信息
blockchain.go 		//列表相关信息
proofofwork.go    	//POW，工作量证明
utils.go        	//工具
blockchain.go       	//把区块链信息写入到数据库
cli.go              	//命令接口
transaction.go		//交易


运行：
go build
./执行文件名
创建钱包
./blockpart5 createwallet
生成新的区块链
./blockpart5 createblockchain -address 15wEYRaK3fpgkjtBXoTNAgzZvHvfGKSLRo
查询余额
./blockpart5 getbalance -address 15wEYRaK3fpgkjtBXoTNAgzZvHvfGKSLRo
转账：
./blockpart5 send -from 15wEYRaK3fpgkjtBXoTNAgzZvHvfGKSLRo -to 1ADdcLTbpd8X24yyJan9E739DRiX2hVCtE -amount 6
查询地址列表
./blockpart5 listaddresses
