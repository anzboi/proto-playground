package examples

import "github.com/spf13/cobra"

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: "example",
		Run: func(cmd *cobra.Command, args []string) {
			Run()
		},
	}
	return cmd
}

func Run() {
	ExampleMarshal()
}
