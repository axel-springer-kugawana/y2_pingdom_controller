/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"pingdom_controller/controller"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start Pingdom Controller",
	Long: "Start Pingdom Controller",
	Run: func(cmd *cobra.Command, args []string) {

		// Enable line numbers in logging
		log.SetFlags(log.LstdFlags | log.Lshortfile)

		npe := controller.NewPingdomEngine()
		go npe.Run()

		go controller.IngressInformerFactory(npe)

		r := mux.NewRouter()
		r.HandleFunc("/health", controller.Health).Methods("GET")

		log.Printf("\nListenAndServe...")
		err := http.ListenAndServe(":" + os.Getenv("PORT"), nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}
