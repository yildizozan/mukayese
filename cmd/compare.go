/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yildizozan/mukayese/internal"
	"log"
)

// compareCmd represents the compare command
var compareCmd = &cobra.Command{
	Use:   "compare",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(2)(cmd, args); err != nil {
			return err
		}
		if err := cobra.MaximumNArgs(2)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		var hashSummaryCurrent map[string]string
		var hashSummaryPrevious map[string]string

		if !internal.IsDirectory(args[0]) {
			log.Fatal("current path is not directory")
		}

		if !internal.IsDirectory(args[1]) {
			log.Fatal("previous path is not directory")
		}

		hashSummaryCurrent = make(map[string]string)
		hashSummaryPrevious = make(map[string]string)

		internal.ListFilesChecksums(hashSummaryCurrent, args[0])
		internal.ListFilesChecksums(hashSummaryPrevious, args[1])

		fmt.Printf("Current: \n")
		for key, val := range hashSummaryCurrent {
			fmt.Printf("%s@sha256:%s\n", key, val)
		}

		fmt.Printf("Previous: \n")
		for key, val := range hashSummaryPrevious {
			fmt.Printf("%s@sha256:%s\n", key, val)
		}

		added := make(map[string]string)
		changed := make(map[string]string)
		deleted := make(map[string]string)

		for currKey, currVal := range hashSummaryCurrent {
			exist := true
			for prevKey, prevVal := range hashSummaryPrevious {
				if currKey == prevKey {
					exist = false
					if currVal != prevVal {
						changed[currKey] = currVal
					}
				}
			}
			if exist {
				added[currKey] = currVal
			}
		}

		// Determine deteled files
		for prevKey, prevVal := range hashSummaryPrevious {
			exist := true
			for currKey := range hashSummaryCurrent {
				if prevKey == currKey {
					exist = false
					break
				}
			}
			if exist {
				deleted[prevKey] = prevVal
			}
		}

		// Added
		fmt.Printf("Added: \n")
		for key, val := range added {
			fmt.Printf("%s@sha256:%s\n", key, val)
		}
		fmt.Println()

		// Changed
		fmt.Printf("Changed: \n")
		for key, val := range changed {
			fmt.Printf("%s@sha256:%s\n", key, val)
		}
		fmt.Println()

		// Deleted
		fmt.Printf("Deleted: \n")
		for key, val := range deleted {
			fmt.Printf("%s@sha256:%s\n", key, val)
		}
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// compareCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	compareCmd.Flags().BoolP("files", "f", false, "Compare files")
}
