package main

import (
	"flag"
	"fmt"
	"github.com/iHelos/ExpertSystem/logic"
	"github.com/iHelos/ExpertSystem/utils"
)

func main(){

	//Подключение к базе
	dbUser := flag.String("dbuser", "root", "mysql login")
	dbPassword := flag.String("dbpwd", "12345", "mysql password")
	dbHost := flag.String("dbhost", "127.0.0.1", "mysql host")
	dbPort := flag.String("dbport", "3306", "mysql port")
	dbName := flag.String("dbName", "expert_system", "logic name")

	flag.Parse()

	es, err := logic.InitExpertSystem(
		*dbHost,
		*dbPort,
		*dbUser,
		*dbPassword,
		*dbName,
	)
	if err != nil {
		fmt.Printf("Initialization problem: %s\n", err)
		return
	}
	utils.CheckInterrupt(
		func() {
			es.PrintResult()
		},
	)
	utils.CLS()
	for es.NextQuestion() {
		utils.CLS()
	}
}