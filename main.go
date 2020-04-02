package main

import "github.com/elahe-dastan/applifier/cmd"

func main() {
	cmd.Execute()
}

//package main
//
//import (
//	"fmt"
//
//	"github.com/c-bata/go-prompt"
//	"github.com/elahe-dastan/applifier/config"
//	"github.com/elahe-dastan/applifier/internal/client"
//	"github.com/elahe-dastan/applifier/internal/server"
//	log "github.com/sirupsen/logrus"
//)
//
//func completer(d prompt.Document) []prompt.Suggest {
//	s := []prompt.Suggest{
//		{Text: "client", Description: "Store the username and age"},
//		{Text: "server", Description: "Store the article text posted by user"},
//	}
//	return prompt.FilterHasPrefix(s, d.GetWordBeforeCursor(), true)
//}
//
//func main() {
//	p := prompt.New(exec, completer)
//	//in := prompt.Input(">>> ", completer,
//	//	prompt.OptionTitle("sql-prompt"),
//	//	prompt.
//	//		prompt.OptionHistory([]string{"SELECT * FROM users;"}),
//	//	prompt.OptionPrefixTextColor(prompt.Yellow),
//	//	prompt.OptionPreviewSuggestionTextColor(prompt.Blue),
//	//	prompt.OptionSelectedSuggestionBGColor(prompt.LightGray),
//	//	prompt.OptionSuggestionBGColor(prompt.DarkGray))
//	//fmt.Println("Your input: " + in)
//	p.Run()
//}
//
//func exec(s string) {
//	if s == "client" {
//		fmt.Println("Hello from client")
//		cli := client.New()
//		//c := config.ReadClient()
//
//		if err := cli.Connect("127.0.0.1:8080"); err != nil {
//			log.Fatal(err)
//		}
//
//		//if _, err := cli.WhoAmI(); err != nil {
//		//	log.Println(err)
//		//}
//		//
//		//if _, err := cli.ListClientIDs(); err != nil {
//		//	log.Println(err)
//		//}
//	} else if s == "server" {
//		fmt.Println("Hello from server")
//		s := server.New()
//		c := config.ReadServer()
//
//		if err := s.Start(c); err != nil {
//			log.Println(err)
//		}
//	}
//}
