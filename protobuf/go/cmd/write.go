package cmd

import (
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
	"grpcdemo/pb"
	"log"
	"os"
)

var destinationFile string

func init() {
	writeCmd.Flags().StringVarP(&destinationFile, "file", "f", "", "The destination file where to save serialized protobuf message")
	rootCmd.AddCommand(writeCmd)
}

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write protobuf message to file",
	Long:  "Write protobuf message to file",
	Run: func(cmd *cobra.Command, args []string) {
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

		if err := os.WriteFile(destinationFile, bytes, 0644); err != nil {
			log.Fatalln(err)
		}
	},
}
