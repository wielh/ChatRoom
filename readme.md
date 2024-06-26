# Chatroom API 介紹

## 1. 這是一個怎麼樣的程式

  使用 go 與 grpc 做的範例 API, 有以下的功能 :
  
 * 支援 google 登入
 * 使用 JWT 保存登入訊息 
 * 使用微服務將不同功能分離, 彼此之間用 grpc 通信
 * 使用 elastic 紀錄 90 日以內的 log

## 2. 使用到的相關技術
  go, gin, grpc, websocket, postgresql, go-pg, elastic 與 kibana
  
## 3. 部署注意事項
 * postgreSQL 是在 windows 上的本地執行, 範例數據可從 setting.sql 匯入，而
 * 要啟用 googleLogin 功能，記得要去更改 googleVerifyID 與 googleVerifyPassword，
 關於申請 googleVerifyID 方法，參見 [GoogleAPI申請](https://blog.hungwin.com.tw/aspnet-google-login/)
 * 要使用 grpc 在微服務之間通信，需要安裝  protobuf-compiler, protoc-gen-ts，安裝完成後使用以下指令編譯 proto 檔案

 ```
   protoc --go_out=. --go-grpc_out=. .\proto\*.proto
```

 * elastic 是紀錄與搜尋 log 的引擎程式，而 kibana 是 elastic 的可視化程式。搜尋 log 的功能主要是在左上角 "discover" 欄位。
   詳細請參見[kibana教學](https://medium.com/%E7%A8%8B%E5%BC%8F%E4%B9%BE%E8%B2%A8/elk-%E6%95%99%E5%AD%B8%E8%88%87%E4%BB%8B%E7%B4%B9-c54af6f06e61)