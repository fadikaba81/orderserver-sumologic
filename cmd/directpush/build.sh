go build -o order-api
sudo setcap 'cap_net_bind_service=+ep' ./order-api
./order-api