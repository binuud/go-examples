package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "demo [arg]",
	Short:   "A brief description shown in command lists",
	Long:    `A detailed paragraph explaining what the command does when help is explicitly invoked.`,
	Example: `  demo my-argument --flag-value`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		// do actual work
	},
}

func NewGreetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "greet [name]",
		Short: "Print a friendly greeting",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := "world"
			if len(args) == 1 {
				name = args[0]
			}
			fmt.Printf("Hello, %s!\n", name)
			return nil
		},
	}
}

func NewServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Run the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")
			fmt.Printf("Serving on :%d\n", port)
			return nil
		},
	}
	cmd.Flags().Int("port", 8080, "port to listen on")
	return cmd
}

func NewDebugCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "debug",
		Short: "debug the HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			port, _ := cmd.Flags().GetInt("port")
			fmt.Printf("Serving on :%d\n", port)
			return nil
		},
	}
	cmd.Flags().Int("port", 8080, "port to listen on")
	return cmd
}

func NewShellCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "shell",
		Short: "Enter the interactive REPL shell",
		Run: func(cmd *cobra.Command, args []string) {
			// Pass rootCmd so the loop can feed strings into it
			startCobraREPL(rootCmd)
		},
	}
}

func NewAvatarCommand() *cobra.Command {
	avatarCommand := &cobra.Command{
		Use:   "avatar",
		Short: "*list*, *start*, *update*",
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Println("Available subcommands under avatar:")
			// Loop through all attached subcommands
			for _, subCmd := range cmd.Commands() {
				cmd.Printf("  - %s: %s\n", subCmd.Name(), subCmd.Short)
			}
		},
	}

	avatarCommand.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "lists an avatar",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Printf("Listing avatar\n")
		},
	})

	avatarCommand.AddCommand(&cobra.Command{
		Use:   "delete",
		Short: "deletes an avatar",
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Printf("Deleting avatar\n")

		},
	})

	avatarCommand.AddCommand(&cobra.Command{
		Use:   "start",
		Short: "starts an avatar",
		Run: func(cmd *cobra.Command, args []string) {

			cmd.Printf("Starting avatar\n")

		},
	})

	return avatarCommand
}

func startCobraREPL(rootCmd *cobra.Command) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Welcome to the Cobra-powered REPL!")
	fmt.Println("All standard CLI subcommands work here. Type 'exit' to quit.")
	fmt.Println("---------------------------------------------------------")

	for {
		fmt.Print("mycli (repl)> ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		if input == "exit" || input == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		// 4. Split the raw text input into an arguments slice (like os.Args)
		// Note: A simple strings.Fields works for standard words.
		// For strings with quoted spaces, use a library like 'google/shlex'.
		replArgs := strings.Fields(input)

		// 5. Intercept Cobra behavior for REPL context
		// Silence errors and usage printing so a typo doesn't kill or flood your loop
		rootCmd.SilenceErrors = true
		rootCmd.SilenceUsage = true

		// 6. Feed the arguments into Cobra
		rootCmd.SetArgs(replArgs)

		// 7. Execute! If it returns an error, catch it and print it locally
		if err := rootCmd.Execute(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func main() {
	rootCmd.AddCommand(NewServeCommand())
	rootCmd.AddCommand(NewDebugCommand())
	rootCmd.AddCommand(NewGreetCommand())
	rootCmd.AddCommand(NewShellCommand())
	rootCmd.AddCommand(NewAvatarCommand())
	fmt.Println("Started")
	Execute()
}
