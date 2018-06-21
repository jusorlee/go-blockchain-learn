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
生成新的区块链
./blockpart4 createblockchain -address abcd
查询余额
./blockpart4 getbalance -address abcd
转账：
./blockpart4 send -from abcd -to leilei -amount 6
