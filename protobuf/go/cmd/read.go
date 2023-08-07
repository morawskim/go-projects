package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
	"grpcdemo/pb"
	"log"
	"os"
)

var fileToRead string

func init() {
	readCmd.Flags().StringVarP(&fileToRead, "file", "f", "", "Source file to read from")
	rootCmd.AddCommand(readCmd)
}

var readCmd = &cobra.Command{
	Use:   "read",
	Short: "Read protobuf message from file",
	Long:  `Read protobuf message from file`,
	Run: func(cmd *cobra.Command, args []string) {
		var loadedOrder pb.Order
		file, err := os.ReadFile(fileToRead)
		if err != nil {
			log.Fatalln(err)
		}

		if err := proto.Unmarshal(file, &loadedOrder); err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("%v\n", &loadedOrder)
	},
}
