# TRON protobuf protocol

## TRON使用Google protobuf协议，协议内容涉及到账户，区块，传输多个层面。

+	账户有基本账户、资产发布账户和合约账户三种类型。一个账户包含：账户名称，账户类型，地址，余额，投票，其他资产6种属性。
+	更进一步的，基本账户可以申请成为验证节点，验证节点具有额外的属性，投票统计数目，公钥，URL，以及历史表现等参数。

   3种`Account`类型：`Normal`，`AssetIssue`，`Contract`。

    enum AccountType {   
      Normal = 0;   
      AssetIssue = 1;   
      Contract = 2; 
     }

   一个`Account`包含7种参数：  
   `account_name`：该账户的名称——比如： ”_SicCongsAccount_”。  
   `type`:该账户的类型——比如：  _0_ 代表的账户类型是`Normal`。  
   `balance`:该账户的TRX余额——比如：_4213312_。  
   `votes`:账户所得投票数——比如：_{(“0x1b7w…9xj3”,323),(“0x8djq…j12m”,88),…,(“0x82nd…mx6i”,10001)}_。  
   `asset`：除TRX以外账户上的其他资产——比如：_{<”WishToken”,66666>,<”Dogie”,233>}_。
   `latest_operation_time`: 该账户的最新活跃时间。
   
    // Account 
    message Account {   
      message Vote {     
        bytes vote_address = 1;     
        int64 vote_count = 2;   
       }   
       bytes accout_name = 1;   
       AccountType type = 2;   
       bytes address = 3;   
       int64 balance = 4;   
       repeated Vote votes = 5;   
       map<string, int64> asset = 6; 
       int64 latest_operation_time = 10;
     }

   一个`Witness`包含8种参数：  
   `address`：验证节点的地址——比如：_“0xu82h…7237”_。   
   `voteCount`：验证节点所得投票数——比如：_234234_。  
   `pubKey`：验证节点的公钥——比如：_“0xu82h…7237”_。  
   `url`：验证节点的url链接。  
   `totalProduce`：验证节点产生的区块数——比如：_2434_。  
   `totalMissed`：验证节点丢失的区块数——比如：_7_。  
   `latestBlockNum`：最新的区块高度——比如：_4522_。
   `isJobs`：布尔表类型标志位。

    // Witness 
    message Witness {   
      bytes address = 1;   
      int64 voteCount = 2;   
      bytes pubKey = 3;   
      string url = 4;   
      int64 totalProduced = 5;   
      int64 totalMissed = 6;   
      int64 latestBlockNum = 7;
      bool isJobs = 9;
      }

+	一个区块由区块头和多笔交易构成。区块头包含时间戳，交易字典树的根，父哈希，签名等区块基本信息。

   一个`block`包含`transactions`和`block_header`。  
   `transactions`：区块里的交易信息。  
   `block_header`：区块的组成部分之一。

    // block 
    message Block {   
      repeated Transaction transactions = 1;   
      BlockHeader block_header = 2;
      }

   `BlockHeader` 包括`raw_data`和`witness_signature`。  
   `raw_data`：`raw`信息。  
   `witness_signature`：区块头到验证节点的签名。

   message `raw`包含6种参数：  
   `timestamp`：该消息体的时间戳——比如：_14356325_。  
   `txTrieRoot`：Merkle Tree的根——比如：_“7dacsa…3ed”_。  
   `parentHash`：上一个区块的哈希值——比如：_“7dacsa…3ed”_。  
   `number`：区块高度——比如：_13534657_。  
   `witness_id`：验证节点的id——比如：_“0xu82h…7237”_。  
   `witness_address`：验证节点的地址——比如：_“0xu82h…7237”_。

    message BlockHeader {   
      message raw {     
        int64 timestamp = 1;     
        bytes txTrieRoot = 2;     
        bytes parentHash = 3;    
        //bytes nonce = 5;    
        //bytes difficulty = 6;     
        uint64 number = 7;     
        uint64 witness_id = 8;     
        bytes witness_address = 9;   
       }   
       raw raw_data = 1;   
       bytes witness_signature = 2;
      }

   消息体 `ChainInventory` 包括 `BlockId` 和 `remain_num`。  
   `BlockId`: block的身份信息。  
   `remain_num`：在同步过程中，剩余的区块数量。
   
   A `BlockId` contains 2 parameters:  
   `hash`: 该区块的哈希值。    
   `number`: 哈希值和高度即为当前区块块号。
   
    message ChainInventory {
      message BlockId {
        bytes hash = 1;
        int64 number = 2;
       }
     repeated BlockId ids = 1;
     int64 remain_num = 2;
     }
           
+	交易合约有多种类型，包括账户创建合约、账户更新合约、转账合约、转账断言合约、资产投票合约、见证节点投票合约、见证节点创建合约、见证节点更新合约、资产发布合约、参与资产发布和与部署合约11种类型。

   `AccountCreatContract`包含3种参数：  
   `type`：账户类型——比如：_0_ 代表的账户类型是`Normal`。  
   `account_name`： 账户名称——比如： _"SiCongsaccount”_。  
   `owner_address`：合约持有人地址——比如： _“0xu82h…7237”_。

    message AccountCreateContract {   
      AccountType type = 1;   
      bytes account_name = 2;   
      bytes owner_address = 3;  
     }
   `AccountUpdateContract`包含2种参数：  
   `account_name`： 账户名称——比如： _"SiCongsaccount”_。  
   `owner_address`：合约持有人地址——比如： _“0xu82h…7237”_。
   
    message AccountUpdateContract {
      bytes account_name = 1;
      bytes owner_address = 2;
     }
     
   `TransferContract`包含3种参数：  
   `amount`：TRX数量——比如：_12534_。  
   `to_address`： 接收方地址——比如：_“0xu82h…7237”_。  
   `owner_address`：合约持有人地址——比如：_“0xu82h…7237”_。

    message TransferContract {   
      bytes owner_address = 1;   
      bytes to_address = 2;   
      int64 amount = 3;
      }

   `TransferAssetContract`包含4种参数：  
   `asset_name`：资产名称——比如：_”SiCongsaccount”_。  
   `to_address`：接收方地址——比如：_“0xu82h…7237”_。  
   `owner_address`：合约持有人地址——比如：_“0xu82h…7237”_。  
   `amount`：目标资产数量——比如：_12353_。

    message TransferAssetContract {   
      bytes asset_name = 1;   
      bytes owner_address = 2;   
      bytes to_address = 3;   
      int64 amount = 4; 
     }

   `VoteAssetContract`包含4种参数：  
   `vote_address`：投票人地址——比如：_“0xu82h…7237”_。  
   `support`：投票赞成与否——比如：_true_。  
   `owner_address`：合约持有人地址——比如：_“0xu82h…7237”_。  
   `count`：投票数目——比如：_2324234_。

    message VoteAssetContract {   
      bytes owner_address = 1;   
      repeated bytes vote_address = 2;   
      bool support = 3;   
      int32 count = 5;
     }

   `VoteWitnessContract`包含4种参数：  
   `vote_address`：投票人地址——比如：_“0xu82h…7237”_。  
   `support`：投票赞成与否——比如：_true_。   
   `owner_address`：合约持有人地址——比如：_“0xu82h…7237”_。  
   `count`：投票数目——比如：_32632_。 

    message VoteWitnessContract {   
      bytes owner_address = 1;   
      repeated bytes vote_address = 2;   
      bool support = 3;   
      int32 count = 5; 
     }

   `WitnessCreateContract`包含3种参数：  
   `private_key`：合约的私钥——比如：_“0xu82h…7237”_。  
   `owner_address`：合约持有人地址——比如：_“0xu82h…7237”_。  
   `url`：合约的url链接。
   
    message WitnessCreateContract {   
      bytes owner_address = 1;   
      bytes private_key = 2;   
      bytes url = 12;
      }

   `WitnessUpdateContract`包含2种参数：  
   `owner_address`：合约持有人地址——比如：_“0xu82h…7237”_。  
   `update_url`：合约的url链接。
   
    message WitnessUpdateContract {
      bytes owner_address = 1;
      bytes update_url = 12;
     }
   
   `AssetIssueContract`包含11种参数：  
   `name`：合约名称——比如：_“SiCongcontract”_。  
   `total_supply`：合约的赞成总票数——比如：_100000000_。  
   `owner_address`：合约持有人地址——比如：_“0xu82h…7237”_。  
   `trx_num`：对应TRX数量——比如：_232241_。  
   `num`： 对应的自定义资产数目。  
   `start_time`：开始时间——比如：_20170312_。  
   `end_time`：结束时间——比如：_20170512_。  
   `vote_score`：合约的评分——比如：_12343_。  
   `description`：合约的描述——比如：_”trondada”_。  
   `url`：合约的url地址链接。

    message AssetIssueContract {   
      bytes owner_address = 1;   
      bytes name = 2;   
      int64 total_supply = 4;   
      int32 trx_num = 6;   
      int32 num = 8;   
      int64 start_time = 9;   
      int64 end_time = 10; 
      int32 vote_score = 16;  
      bytes description = 20;   
      bytes url = 21; 
     }
     
   `ParticipateAssetIssueContract`包含4种参数：  
   `owner_address`：合约持有人地址——比如：_“0xu82h…7237”_。  
   `to_address`：接收方地址——比如：_“0xu82h…7237”_。  
   `asset_name`: 目标资产的名称。  
   `amount`： 小部分数量。
   
   `DeployContract`包含2种参数：  
   `script`：脚本。  
   `owner_address`：合约持有人地址——比如：_“0xu82h…7237”_。

    message DeployContract {   
      bytes owner_address = 1;   
      bytes script = 2; 
     }

   消息体 `Result` 包含 `fee` and `ret`2个参数.   
   `ret`: 交易结果。  
   `fee`: 交易扣除的费用。
   
   `code`是`ret`的类型定义，有`SUCCESS`和`FAILED`两种类型。
    
    message Result {
      enum code {
        SUCESS = 0;
        FAILED = 1;
       }
       int64 fee = 1;
       code ret = 2;
     }
      
+	每一个交易还包含多个输入与多个输出，以及其他一些相关属性。其中交易内的输入，交易本身，区块头均需签名。

   消息体 `Transaction`包括`raw_data`和`signature`。  
   `raw_data`: 消息体`raw`。  
   `signature`: 所有输入节点的签名。

   `raw_data`包含8种参数：  
   `type`：消息体raw的交易类型。  
   `vin`： 输入值。  
   `vout`： 输出值。  
   `expiration`：过期时间——比如：_20170312_。  
   `data`： 数据。  
   `contract`： 该交易内的合约。  
   `script`： 脚本。
   `timestamp`：该消息体的时间戳。

   消息体 `Contract`包含`type`和`parameter`。  
   `type`：合约的类型。  
   `parameter`：任意参数。

   有八种账户类型合约：`AccountCreateContract`，`TransferContract`，`TransferAssetContract`，`VoteAssetContract`，`VoteWitnessContract`，`WitnessCreateContract`，`AssetIssueContract` 和`DeployContract`。

   `TransactionType`包括`UtxoType`和`ContractType`。

    message Transaction {   
      enum TranscationType {     
        UtxoType = 0;     
        ContractType = 1;   
       }   
       message Contract {    
         enum ContractType {       
           AccountCreateContract = 0;       
           TransferContract = 1;       
           TransferAssetContract = 2;       
           VoteAssetContract = 3;       
           VoteWitnessContract = 4;      
           WitnessCreateContract = 5;       
           AssetIssueContract = 6;       
           DeployContract = 7;     
          }     
          ContractType type = 1;     
          google.protobuf.Any parameter = 2;  
        }   
        message raw {     
          TranscationType type = 2;     
          repeated TXInput vin = 5;     
          repeated TXOutput vout = 7;     
          int64 expiration = 8;     
          bytes data = 10;     
          repeated Contract contract = 11;     
          bytes scripts = 16; 
          in64 timestamp = 17; 
         }   
         raw raw_data = 1;   
         repeated bytes signature = 5; 
     }

   消息体 `TXOutputs`由`outputs`构成。  
   `outputs`: 元素为`TXOutput`的数组。

    message TXOutputs {   
      repeated TXOutput outputs = 1;
      }

   消息体 `TXOutput`包括`value`和`pubKeyHash`。  
   `value`：输出值。  
   `pubKeyhash`：公钥的哈希。

    message TXOutput {   
      int64 value = 1;   
      bytes pubKeyHash = 2;
      }

   消息体 `TXIutput`包括`raw_data`和`signature`。  
   `raw_data`：消息体`raw`。  
   `signature`：`TXInput`的签名。

   消息体 `raw`包含`txID`，`vout`和 `pubKey`。  
   `txID`：交易ID。  
   `Vout`：上一个输出的值。  
   `pubkey`:公钥。

    message TXInput {   
      message raw {     
        bytes txID = 1;     
        int64 vout = 2;     
        bytes pubKey = 3;   
       }   
       raw raw_data = 1;  
       bytes signature = 4; }

+	传输涉及的协议Inventory主要用于传输中告知接收方传输数据的清单。

   `Inventory`包括`type`和`ids`。  
   `type`：清单类型——比如：_0_ 代表`TRX`。  
   `ids`：清单中的物品ID。

   `InventoryType`包含`TRX`和 `BLOCK`。  
   `TRX`：交易。  
   `BLOCK`：区块。

    // Inventory 
    message Inventory {   
      enum InventoryType {     
        TRX = 0;     
        BLOCK = 1;   
       }   
       InventoryType type = 1;   
       repeated bytes ids = 2;
      }

   消息体 `Items`包含4种参数：  
   `type`：物品类型——比如：_1_ 代表 `TRX`。  
   `blocks`：物品中区块。  
   `blockheaders`：区块头。  
   `transactions`：交易。

   `Items`有四种类型，分别是 `ERR`， `TRX`，`BLOCK` 和`BLOCKHEADER`。  
   `ERR`：错误。  
   `TRX`：交易。  
   `BLOCK`：区块。  
   `BLOCKHEADER`：区块头。

    message Items {   
      enum ItemType {     
        ERR = 0;     
        TRX = 1;     
        BLOCK = 2;     
        BLOCKHEADER = 3;   
       }   
       ItemType type = 1;   
       repeated Block blocks = 2;   
       repeated BlockHeader block_headers = 3;   
       repeated Transaction transactions = 4; 
     }

   `Inventory`包含`type`和`items`。  
   `type`：物品种类。  
   `items`：物品清单。

    message InventoryItems {   
      int32 type = 1;   
      repeated bytes items = 2; 
     }
     
   消息体 `BlockInventory` 包含 `type`。  
   `type`: 清单种类.  
   
   有三种类型:`SYNC`, `ADVTISE`, `FETCH`。
     
    // Inventory
     message BlockInventory {
       enum Type {
         SYNC = 0;
         ADVTISE = 1;
         FETCH = 2;
       }
       
   消息体 `BlockId` 包括 `ids` and `type`。  
   `ids`: 区块身份信息。  
   `type`: 区块类型。 
      
   `ids` 包含2种参数：  
   `hash`: 区块的哈希值。  
   `number`: 哈希值和区块高度即为当前区块号。
       
     message BlockId {
        bytes hash = 1;
        int64 number = 2;
      }
      repeated BlockId ids = 1;
      Type type = 2;
     }
   
   `ReasonCode` 有15种可能断开的原因:  
   `REQUESTED`  
   `TCP_ERROR`  
   `BAD_PROTOCOL`   
   `USELESS_PEER`  
   `TOO_MANY_PEERS`  
   `DUPLICATE_PEER`  
   `INCOMPATIBLE_PROTOCOL`  
   `NULL_IDENTITY`  
   `PEER_QUITING`  
   `UNEXPECTED_IDENTITY`  
   `LOCAL_IDENTITY`  
   `PING_TIMEOU`  
   `USER_REASON`  
   `RESET`  
   `UNKNOWN` 
   
      enum ReasonCode {
        REQUESTED = 0;
        TCP_ERROR = 1;
        BAD_PROTOCOL = 2;
        USELESS_PEER = 3;
        TOO_MANY_PEERS = 4;
        DUPLICATE_PEER = 5;
        INCOMPATIBLE_PROTOCOL = 6;
        NULL_IDENTITY = 7;
        PEER_QUITING = 8;
        UNEXPECTED_IDENTITY = 9;
        LOCAL_IDENTITY = 10;
        PING_TIMEOUT = 11;
        USER_REASON = 12;
        RESET = 16;
        UNKNOWN = 255;
      }
       
   消息体`DisconnectMessage`包含`reason`。  
   `DisconnectMessage`：断开连接是的消息。  
   `reason`：断开连接时的原因。
   
   消息体`HelloMessage`包含2个参数：  
   `from`请：求建立连接的节点。  
   `version`：建立连接的节点。
         
+	钱包服务RPC和区块链浏览器。

   `Wallet`钱包服务包含多个RPC。  
   __`Getbalance`__：获取`Account`的余额。  
   __`CreatTransaction`__：通过`TransferContract`创建交易。  
   __`BroadcastTransaction`__：广播`Transaction`。  
   __`CreateAccount`__：通过`AccountCreateContract`创建账户。  
   __`CreatAssetIssue`__：通过`AssetIssueContract`发布一个资产。  
   __`ListAccounts`__：通过`ListAccounts`查看账户列表。  
   __`UpdateAccount`__：通过`UpdateAccountContract`发布一个资产。  
   __`VoteWitnessAccount`__：通过`VoteWitnessContract`发布一个资产。  
   __`WitnessList`__：通过`WitnessList`查看见证节点列表。  
   __`UpdateWitness`__：通过`WitnessUpdateContract`发布一个资产。  
   __`CreateWitness`__：通过`WitnessCreateContract`发布一个资产。  
   __`TransferAsset`__：通过`TransferAssetContract`发布一个资产。  
   __`ParticipateAssetIssue`__：通过`ParticipateAssetIssueContract`发布一个资产。  
   __`ListNodes`__：通过`ListNodes`查看节点列表。  
   __`GetAssetIssueList`__：通过`GetAssetIssueList`查看资产发布节点列表。  
   __`GetAssetIssueByAccount`__：通过`Account`获取发行资产。  
   __`GetAssetIssueByName`__：通过`Name`获取发行资产。  
   __`GetNowBlock`__：获取区块。  
   __`GetBlockByNum`__：根据块号获取区块。  
   __`TotalTransaction`__：查看总交易量。  
   
       service Wallet {
       
         rpc GetAccount (Account) returns (Account) {
       
         };
       
         rpc CreateTransaction (TransferContract) returns (Transaction) {
       
         };
       
         rpc BroadcastTransaction (Transaction) returns (Return) {
       
         };
       
         rpc ListAccounts (EmptyMessage) returns (AccountList) {
       
         };
       
         rpc UpdateAccount (AccountUpdateContract) returns (Transaction) {
              
         };
       
         rpc CreateAccount (AccountCreateContract) returns (Transaction) {
       
         };
       
         rpc VoteWitnessAccount (VoteWitnessContract) returns (Transaction) {
       
         };
       
         rpc CreateAssetIssue (AssetIssueContract) returns (Transaction) {
       
         };
       
         rpc WitnessList (EmptyMessage) returns (WitnessList) {
       
         };
       
         rpc UpdateWitness (WitnessUpdateContract) returns (Transaction) {
       
         };
       
         rpc CreateWitness (WitnessCreateContract) returns (Transaction) {
       
         };
       
         rpc TransferAsset (TransferAssetContract) returns (Transaction) {
       
         }
       
         rpc ParticipateAssetIssue (ParticipateAssetIssueContract) returns (Transaction) {
       
         }
       
         rpc ListNodes (EmptyMessage) returns (NodeList) {
       
         }
         rpc GetAssetIssueList (EmptyMessage) returns (AssetIssueList) {
       
         }
         rpc GetAssetIssueByAccount (Account) returns (AssetIssueList) {
       
         }
         rpc GetAssetIssueByName (BytesMessage) returns (AssetIssueContract) {
       
         }
         rpc GetNowBlock (EmptyMessage) returns (Block) {
       
         }
         rpc GetBlockByNum (NumberMessage) returns (Block) {
       
         }
         rpc TotalTransaction (EmptyMessage) returns (NumberMessage) {
       
         }
       };

   `AccountList`： 区块链浏览器中的账户列表。  
   消息体 `AccountList` 包含1个参数：  
   `account`：
   
        message AccountList {
          repeated Account accounts = 1;
        }  
   
   `WitnessList`：区块链浏览器中的见证节点列表。  
   消息体 `WitnessList` 包含1个参数：  
   `witnesses`：
      
           message WitnessList {
             repeated Witness witnesses = 1;
           }
           
   `AssetIssueList`：区块链浏览器中的发布资产列表。  
   消息体 `AssetIssueList` 包含1个参数:  
   `assetIssue`：
      
           message AssetIssueList {
             repeated AssetIssueContract assetIssue = 1;
           }
      
     
   `NodeList`： 分布节点图中的节点列表。  
   消息体 `NodeList` 包含1个参数：  
   `nodes`：
      
            message NodeList {
              repeated Node nodes = 1;
            }
      
   `Address`: 节点地址。  
   消息体`Address` 包含2个参数：  
   `host`：节点所有者。  
   `port`：节点的端口号。
      
             message Address {
               bytes host = 1;
               int32 port = 2;
             }
                  
    
   消息体`Return`只含有一个参数：  
   `result`: 布尔表类型标志位。  

              message `Return` {   
                 bool result = 1; 
              }
              
+  网络UDP消息结构。

  `Endpoint`：网络中节点信息存储结构.  
  消息体`Endpoint` 包含3个参数：  
   `address`：节点地址。  
   `port`：端口号。  
   `nodeId`： 节点ID信息。

    message Endpoint {
      bytes address = 1;
      int32 port = 2;
      bytes nodeId = 3;
     }

   `PingMessage`：节点建立连接时所发送的消息。  
   消息体`PingMessage` 包含4个参数：  
   `from`：消息来自的节点。  
   `to`： 消息发送的节点。  
   `version`： 网络版本。  
   `timestamp`：消息创建时的时间戳。

    message PingMessage {
       Endpoint from = 1;
       Endpoint to = 2;
       int32 version = 3;
       int64 timestamp = 4;
     }
    
   `PongMessage`：连接建立成功时的回复消息。  
   消息体`PongMessage` 包含3个参数：  
   `from`：消息来自的节点。  
   `echo`：  
   `timestamp`：消息创建时的时间戳。
   
    message PongMessage {
      Endpoint from = 1;
      int32 echo = 2;
      int64 timestamp = 3;
     }
   
   `FindNeighbours`：节点查询相邻节点时所发送的消息。  
   消息体`FindNeighbours` 包含3个参数：  
   `from`: 消息来自的节点。  
   `targetId`: 目标节点的信息。  
   `timestamp`: 消息创建时的时间戳。

    message FindNeighbours {
      Endpoint from = 1;
      bytes targetId = 2;
      int64 timestamp = 3;
     }

   `Neighbour`：相邻接点回复消息。  
   消息体`Neighbours` 包含3个参数：  
   `from`: 消息来自的节点。  
   `neighbours`: 相邻节点。  
   `timestamp`: 消息创建时的时间戳。
   
    message Neighbours {
      Endpoint from = 1;
      repeated Endpoint neighbours = 2;
      int64 timestamp = 3;
     }
     
# 详细的协议见附属文件。详细协议随着程序的迭代随时都可能发生变化，请以最新的版本为准。