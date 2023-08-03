package main

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"grpcdemo/pb"
	"log"
	"os"
)

const filename = "protobufer-message.bin"

func main() {
	order := pb.Order{
		Number: "LI/FOO/123456789",
		Status: pb.OrderStatus_ORDER_STATUS_NEW,
		Items: []*pb.OrderItem{
			{Sku: "456789", Quantity: 2, UnitPrice: 499},
			{Sku: "AZ456", Quantity: 1, UnitPrice: 999},
		},
		CreatedAt: timestamppb.Now(),
	}

	bytes, err := proto.Marshal(&order)
	if err != nil {
		log.Fatalln(err)
	}

	if err := os.WriteFile(filename, bytes, 0644); err != nil {
		log.Fatalln(err)
	}

	var loadedOrder pb.Order
	file, err := os.ReadFile(filename)
	if err != nil {
		log.Fatalln(err)
	}

	if err := proto.Unmarshal(file, &loadedOrder); err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%v\n", &loadedOrder)
}
