package logic

import "fmt"

type AttributeValue struct {
	ID int
	Value string
	Dependencies []*AttributeCondition
	isStrict bool
}

func (av *AttributeValue) Valid() (result bool){
	result = !av.isStrict
	for _, v := range av.Dependencies {
		tempResult, isStrict := v.check()
		result = result && tempResult && isStrict
		if result == false {
			if isStrict == true {
				v.isStrict = true
				return
			}
		}
	}
	return
}

type AttributeCondition struct {
	ID int
	Parameter *Parameter
	ConditionValue int
	Condition string
	isStrict bool
}

func (q *AttributeCondition) check() (result, isStrict bool){
	switch q.Condition {
	case "EQ":
		result = q.Parameter.Value == q.ConditionValue
	case "GT":
		result = q.Parameter.Value > q.ConditionValue
	case "LT":
		result = q.Parameter.Value < q.ConditionValue
	case "LTE":
		result = q.Parameter.Value <= q.ConditionValue
	case "GTE":
		result = q.Parameter.Value >= q.ConditionValue
	case "NEQ":
		result = q.Parameter.Value != q.ConditionValue
	default:
		fmt.Println("No such condition")
		result = false
	}
	if result == false {
		isStrict = q.isStrict
	}
	return
}

type Attribute struct {
	ID int
	Name string
	Values	map[int]*AttributeValue
}

type Parameter struct {
	ID int
	Name string
	Value int
}

type AnswerInfluence struct {
	ID int
	Param *Parameter
	Operation string
	Value int
}

type Answer struct {
	ID int
	AnswerText string
	Influences []*AnswerInfluence
	Value int
}

func (i *AnswerInfluence) Do(answerValue int) {
	switch i.Operation {
	case "+":
		i.Param.Value += i.Value
	case "-":
		i.Param.Value -= i.Value
	case "=":
		i.Param.Value = i.Value
	case "+?":
		i.Param.Value += answerValue
	case "-?":
		i.Param.Value -= answerValue
	case "=?":
		i.Param.Value = answerValue
	default:
		return
	}
}

func (a *Answer) DoInfluence(){
	for _, v := range a.Influences {
		v.Do(a.Value)
	}
}

type QuestionCondition struct {
	ID int
	Parameter *Parameter
	ConditionValue int
	Condition string
}

func (q *QuestionCondition) check() (result bool){
	switch q.Condition {
	case "EQ":
		result = q.Parameter.Value == q.ConditionValue
	case "GT":
		result = q.Parameter.Value > q.ConditionValue
	case "LT":
		result = q.Parameter.Value < q.ConditionValue
	case "LTE":
		result = q.Parameter.Value <= q.ConditionValue
	case "GTE":
		result = q.Parameter.Value >= q.ConditionValue
	case "NEQ":
		result = q.Parameter.Value != q.ConditionValue
	default:
		fmt.Println("No such condition")
		result = false
	}
	return
}

type Question struct {
	ID int
	QuestionText string
	Answers []*Answer
	isAsked bool
	Conditions []*QuestionCondition
	NextPossibleQuestion []*Question
}

func (q *Question) validQuestion() (result bool){
	result = !q.isAsked
	for _, v := range q.Conditions {
		result = result && v.check()
	}
	return
}

type Object struct {
	ID int
	Name string
	Attributes []*AttributeValue
}

func(o *Object) getResult() float64 {
	var result float64
	for _,v := range o.Attributes {
		if v.Valid() {
			result++
			continue
		}
		if v.isStrict {
			return 0
		}
	}
	return result/float64(len(o.Attributes))
}
