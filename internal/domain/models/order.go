package models

type Order struct {
	OrderId    int32 `json:"orderId" protobuf:"int32 ,1,opt,name=orderId"`
	CustomerId int32 `json:"customerId" protobuf:"int32,2,opt,name=customerId"`
	ProductId  int32 `json:"productId" protobuf:"int3,3,opt,name=productId"`
	Quantity   int32 `json:"quantity" protobuf:"int32,4,opt,name=quantity"`
}

type Customer struct {
	CustomerId int32  `json:"customer_id"`
	Name       string `json:"name"`
	Surname    string `json:"surname"`
	Username   string `json:"username"`
}

type Product struct {
	ProductId int32  `json:"product_id"`
	Name      string `json:"name"`
}
