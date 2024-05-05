/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"net/http"

	"github.com/defenseunicorns/uds-generator/src/config"
	"github.com/defenseunicorns/uds-generator/src/config/lang"
	"github.com/defenseunicorns/uds-generator/src/pkg/api"
	"github.com/defenseunicorns/uds-generator/src/pkg/ui"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:
Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		fmt.Println("serve called")

		// Create a new ServeMux for WebSocket server
		mux := http.NewServeMux()
		mux.HandleFunc("/ws", handleWebSocket)

		// Start WebSocket server on a separate port
		go func() {
			log.Println("WebSocket server starting on :" + config.WebsocketPort + "...")
			if err := http.ListenAndServe(":"+config.WebsocketPort, mux); err != nil {
				log.Fatal("WebSocket server error:", err)
			}
		}()

		// Define API routes
		http.HandleFunc("/api/message", api.APIHandler)

		// SPA App defined on root
		http.HandleFunc("/", ui.StaticFileHandler)

		// Start server
		log.Println("Server starting on :" + config.ApiPort + "...")
		if err := http.ListenAndServe(":"+config.ApiPort, nil); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&config.ApiPort, "port", "p", "8080", lang.CmdServeFlagPort)
	serveCmd.Flags().StringVarP(&config.WebsocketPort, "websocket-port", "w", "8081", lang.CmdServeFlagwsPort)
}
