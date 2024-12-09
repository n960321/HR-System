/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "HRSystem",
	Short: "Run a HR system for managing employee attendance records",
	Long: `Run a simple HR backend system with the following features:
		1. User login
		2. Admin creates employee accounts 
		3. User changes password
		4. User clock in/out
		5. User retrieves attendance records
	`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Hello World")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
