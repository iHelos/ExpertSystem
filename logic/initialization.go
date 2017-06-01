package logic

import "fmt"

func (es expertSystem) getAnswersInfluence(answerID int) ([]*AnswerInfluence, error){
	res, err := es.db.Query("select id, parameter_id, operation_id, value from answers_influence where answer_id = ?", answerID)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var result []*AnswerInfluence
	for res.Next() {
		var temp AnswerInfluence
		var paramID int
		err := res.Scan(&temp.ID, &paramID, &temp.Operation, &temp.Value)
		if err != nil {
			return nil, err
		}
		temp.Param = es.parameters[paramID]
		result = append(result, &temp)
	}
	return result, nil
}

func (es expertSystem) getAnswers(questionID int) ([]*Answer, error) {
	res, err := es.db.Query("select id, answer_text, default_value from answers where question_id = ?", questionID)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var result []*Answer
	for res.Next() {
		var tempAnswer Answer
		err := res.Scan(&tempAnswer.ID, &tempAnswer.AnswerText, &tempAnswer.Value)
		if err != nil {
			return nil, err
		}
		inf, err := es.getAnswersInfluence(tempAnswer.ID)
		if err != nil {
			return nil, err
		}
		tempAnswer.Influences = inf
		result = append(result, &tempAnswer)
	}
	return result, nil
}

func (es expertSystem) getQuestionConditions(questionID int) ([]*QuestionCondition, error) {
	res, err := es.db.Query("select id, parameter_id, condition_id, value from parameters_questions_conditions where question_id = ?", questionID)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var result []*QuestionCondition
	for res.Next() {
		var temp QuestionCondition
		var parameterID int
		err := res.Scan(&temp.ID, &parameterID, &temp.Condition, &temp.ConditionValue)
		if err != nil {
			return nil, err
		}
		temp.Parameter = es.parameters[parameterID]
		result = append(result, &temp)
	}
	return result, nil
}

func (es expertSystem) getPossibleQuestions(questionID int) ([]*Question, error) {
	res, err := es.db.Query("select to_question from next_question where from_question = ?", questionID)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var result []*Question
	for res.Next() {
		var questionID int
		err := res.Scan(&questionID)
		if err != nil {
			return nil, err
		}
		q, ok := es.allQuestions[questionID]
		if !ok {
			continue
		}
		result = append(result, q)
	}
	return result, nil
}

func (es *expertSystem) getAllQuestions() error {
	es.allQuestions = make(map[int]*Question)

	res, err := es.db.Query("select id, text from questions")
	if err != nil {
		return err
	}
	defer res.Close()

	for res.Next() {
		var temp Question
		err := res.Scan(&temp.ID, &temp.QuestionText)
		if err != nil {
			return err
		}
		es.allQuestions[temp.ID] = &temp
	}
	for k,v := range es.allQuestions {
		answers, err := es.getAnswers(k)
		if err != nil {
			return err
		}
		v.Answers = answers
		if len(answers) == 0 {
			delete(es.allQuestions, k)
			continue
		}
		conditions, err := es.getQuestionConditions(k)
		if err != nil {
			return err
		}
		v.Conditions = conditions
		nextQuestions, err := es.getPossibleQuestions(k)
		v.NextPossibleQuestion = nextQuestions
	}
	return nil
}

func (es *expertSystem) getStartQuestions() error {
	res, err := es.db.Query("select q.id from questions q LEFT JOIN next_question n on q.id = n.to_question where n.id is NULL")
	if err != nil {
		return err
	}
	defer res.Close()

	for res.Next() {
		var questionID int
		err := res.Scan(&questionID)
		if err != nil {
			return err
		}
		q, ok := es.allQuestions[questionID]
		if !ok {
			continue
		}
		es.startQuestions = append(es.startQuestions, q)
	}
	return nil
}

func (es *expertSystem) getParameters() error {
	es.parameters = make(map[int]*Parameter)

	res, err := es.db.Query("select id, name from parameters")
	if err != nil {
		return err
	}
	defer res.Close()

	for res.Next() {
		var temp Parameter
		err := res.Scan(&temp.ID, &temp.Name)
		if err != nil {
			return err
		}
		es.parameters[temp.ID] = &temp
	}
	return nil
}

func (es *expertSystem) getAttributeValueDependencies(attributeValueID int) ([]*AttributeCondition, error) {
	res, err := es.db.Query("select id, parameter_id, condition_id, condition_value, is_strict from parameters_attributes_conditions where attribute_value_id = ?", attributeValueID)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	var result []*AttributeCondition
	for res.Next() {
		var temp AttributeCondition
		var paramID int
		err := res.Scan(&temp.ID, &paramID, &temp.Condition, &temp.ConditionValue, &temp.isStrict)
		if err != nil {
			return nil, err
		}
		temp.Parameter = es.parameters[paramID]
		result = append(result, &temp)
	}
	return result, nil
}

func (es *expertSystem) getAttributeValues(attributeID int) (map[int]*AttributeValue, error) {
	res, err := es.db.Query("select id, value from attribute_values where attribute_id = ?", attributeID)
	if err != nil {
		return nil, err
	}
	defer res.Close()
	result := make(map[int]*AttributeValue)
	for res.Next() {
		var temp AttributeValue
		err := res.Scan(&temp.ID, &temp.Value)
		if err != nil {
			return nil, err
		}
		dependencies, err := es.getAttributeValueDependencies(temp.ID)
		if err != nil {
			return nil, err
		}
		temp.Dependencies = dependencies
		result[temp.ID] = &temp
	}
	return result, nil
}

func (es *expertSystem) getAttributes() error {
	es.attributes = make(map[int]*Attribute)

	res, err := es.db.Query("select id, name from attributes")
	if err != nil {
		return err
	}
	defer res.Close()

	for res.Next() {
		var temp Attribute
		err := res.Scan(&temp.ID, &temp.Name)
		if err != nil {
			return err
		}
		values, err := es.getAttributeValues(temp.ID)
		if err != nil {
			return err
		}
		temp.Values = values
		es.attributes[temp.ID] = &temp
	}
	return nil
}

func (es *expertSystem) getObjectsAttributeValues(objectID int) ([]*AttributeValue, error){
	res, err := es.db.Query("select oa.attribute_value_id, av.attribute_id from objects_attributes oa " +
		"left JOIN attribute_values av on oa.attribute_value_id = av.id where object_id = ?", objectID)
	if err != nil {
		return nil, err
	}
	defer res.Close()

	var result []*AttributeValue
	for res.Next() {
		var attrValueID int
		var attrID int
		err := res.Scan(&attrValueID, &attrID)
		if err != nil {
			return nil, err
		}
		attr, ok := es.attributes[attrID]
		if !ok {
			fmt.Println("no such attribute")
			continue
		}
		value, ok := attr.Values[attrValueID]
		if !ok {
			fmt.Printf("no such attribute value: v_ID: %d, a_ID: %d\n", attrValueID, attrID)
			continue
		}
		result = append(result, value)
	}
	return result, nil
}

func (es *expertSystem) getObjects() (error) {
	es.objects = make(map[int]*Object)

	res, err := es.db.Query("select id, name from objects")
	if err != nil {
		return err
	}
	defer res.Close()

	for res.Next() {
		var temp Object
		err := res.Scan(&temp.ID, &temp.Name)
		if err != nil {
			return err
		}
		values, err := es.getObjectsAttributeValues(temp.ID)
		if err != nil {
			return err
		}
		temp.Attributes = values
		es.objects[temp.ID] = &temp
	}
	return nil
}