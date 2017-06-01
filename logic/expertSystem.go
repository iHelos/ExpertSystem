package logic

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"github.com/iHelos/ExpertSystem/utils"
	"os"
	"strconv"
	"strings"
)

type expertSystem struct {
	db *sql.DB
	allQuestions map[int]*Question
	parameters map[int]*Parameter
	attributes map[int]*Attribute
	objects map[int]*Object
	startQuestions []*Question
	curQuestion *Question

}

func InitExpertSystem(
	host,
	port,
	login,
	password,
	dbName string,
) (*expertSystem, error) {
	db, err := initConnection(
		host,
		port,
		login,
		password,
		dbName,
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't connect to logic: %s", err))
	}

	es := &expertSystem{
		db:db,
	}
	err = es.getParameters()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't get params: %s", err))
	}
	err = es.getAllQuestions()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't get questions: %s", err))
	}
	err = es.getStartQuestions()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't get start questions: %s", err))
	}
	err = es.getAttributes()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't get attributes: %s", err))
	}
	err = es.getObjects()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't get objects: %s", err))
	}
	return es, nil
}

func (es *expertSystem) getNextQuestion() (*Question, bool) {
	if es.curQuestion != nil {
		for _, v := range es.curQuestion.NextPossibleQuestion {
			if v.validQuestion() {
				v.isAsked = true
				es.curQuestion = v
				return v, false
			}
		}
	}
	for _, v := range es.startQuestions {
		if v.validQuestion() {
			v.isAsked = true
			es.curQuestion = v
			return v, false
		}
	}
	return nil, true
}

func (es *expertSystem) PrintParams(){
	for _, v := range es.parameters {
		fmt.Printf("%s: %d\n", v.Name, v.Value)
	}
}

func (es *expertSystem) PrintResult(){
	for _, v := range es.objects {
		fmt.Printf("%s: %.2f\n", v.Name, v.getResult())
	}
}

func (es *expertSystem) PrintAttributes(){
	for _, a := range es.attributes {
		fmt.Printf("%s\n", a.Name)
		for _, v := range a.Values {
			fmt.Printf("%s: %t\n", v.Value, v.Valid())
		}
		fmt.Println("\n")
	}
}

func (es *expertSystem) PrintObjects(){
	for _, o := range es.objects {
		fmt.Printf("%s\n", o.Name)
		for _, v := range o.Attributes {
			fmt.Printf("%s: %t\n", v.Value, v.Valid())
		}
		fmt.Println("\n")
	}
}

func (es *expertSystem) NextQuestion() bool{
	question, isEnded := es.getNextQuestion()
	if isEnded {
		return false
	}
	reader := bufio.NewReader(os.Stdin)
	fmt.Println(question.QuestionText)
	for i,v := range question.Answers {
		fmt.Printf("%d) %s\n", i + 1, v.AnswerText)
	}
	for {
		text, _ := reader.ReadString('\n')
		text = strings.Split(text, "\n")[0]
		canExit := false
		switch text {
		case "0":
			fmt.Println("Что 0?")
		case "clear":
			utils.CLS()
			fmt.Println(question.QuestionText)
			for i,v := range question.Answers {
				fmt.Printf("%d) %s\n", i + 1, v.AnswerText)
			}
		case "objects":
			es.PrintObjects()
		case "result":
			es.PrintObjects()
		case "params":
			es.PrintParams()
		case "attrs":
			es.PrintAttributes()
		default:
			ans, err := strconv.Atoi(text)
			if err != nil {
				fmt.Println("Ждем число, а не текст!")
				break
			}
			answersLen := len(question.Answers)
			if answersLen == 1 {
				question.Answers[0].Value = ans
				question.Answers[0].DoInfluence()
				canExit = true
				break
			} else {
				if ans - 1 < answersLen {
					question.Answers[ans - 1].DoInfluence()
					canExit = true
					break
				}
				fmt.Println("Выберите один из существующих ответов!")
			}
		}
		if canExit {
			break
		}
	}
	return true
}


