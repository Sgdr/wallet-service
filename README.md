# wallet-service
Simple wallets service which is implemented on go.
Wallet service provide some simple operations with accounts: view accounts, view payments for account, create payment.
Api is available in internal/api/swagger.yml and by <service_url>/doc when service is running.

## How to build and run

for development purposes

 go run main.go [parameters...] 
    
    Parameters:
       -config-path     path to config, by default "./config/config.yml"

Additionally to all, service can read parameters from environment variables
       
Service needs PostgreSQl DB. You can run it in docker:
      
      docker run -p 5432:5432 \
                -e POSTGRES_USER=wallet_user \
                -e POSTGRES_PASSWORD=wallet_user_pass \
                -e POSTGRES_DB=wallet \
                -v "/Users/Shared/postgres/:/var/lib/postgresql/data/:Z" \
                --name psgr -ti postgres:11.2-alpine
 
It creates db user "wallet_user" with password "wallet_user_pass" and database "wallet". This credentials are set in
default config.yml

After starting service automatically creates db structure. There is not any data there. It doesn't not corrupt exist data. 

The endpoints of service will be available by localhost:<http_port>
http_port - parameter from config

Available endpoints is described localhost:<http_port/doc

## How to run test

There is few test and one of them (for payment's creation method) requires DB payment_test_db with user "wallet_user",
password "wallet_user_pass" on localhost:5432.
Run it with command in directory wallet-service

     go test ./...

In future it should be more suitable: creation docker test containers will avoid starting db before test.