go run -x /usr/src/code/service/order/api/order.go -f /usr/src/code/service/order/api/etc/order.yaml &
go run -x /usr/src/code/service/order/rpc/order.go -f /usr/src/code/service/order/rpc/etc/order.yaml &

go run -x /usr/src/code/service/pay/api/pay.go -f /usr/src/code/service/pay/api/etc/pay.yaml &
go run -x /usr/src/code/service/pay/rpc/pay.go -f /usr/src/code/service/pay/rpc/etc/pay.yaml &

go run -x /usr/src/code/service/product/api/product.go -f /usr/src/code/service/product/api/etc/product.yaml &
go run -x /usr/src/code/service/product/rpc/product.go -f /usr/src/code/service/product/rpc/etc/product.yaml &

go run -x /usr/src/code/service/user/api/user.go -f /usr/src/code/service/user/api/etc/user.yaml &
go run -x /usr/src/code/service/user/rpc/user.go -f /usr/src/code/service/user/rpc/etc/user.yaml &

