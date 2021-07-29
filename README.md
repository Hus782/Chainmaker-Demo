# Chainmaker-Demo
## A simple voting demo app using 长安链 (Chainmaker)
Smart contract is written in Solidity, client side is using the Golang SDK
### To run the project 
1. [Install]( https://docs.chainmaker.org.cn/tutorial/%E5%BF%AB%E9%80%9F%E5%85%A5%E9%97%A8.html#id5) Chainmaker and build a network
2. Clone the [Go SDK](https://docs.chainmaker.org.cn/tutorial/%E5%BF%AB%E9%80%9F%E5%85%A5%E9%97%A8.html#go-sdk)

3. Copy the repository to the directory in which you installed chainmaker-go and chainmaker-sdk-go
4. Copy the generated certificates from chainmaker-go/build/crypto-config/ to testdata/crypto-config
    ```
    rm -r testdata/crypto-config/*
    cp -r chainmaker-go/build/crypto-config/* testdata/crypto-config
    ```
5. Run the project
    ```
    go run main.go
    ```

