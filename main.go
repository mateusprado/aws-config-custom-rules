package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mateusprado/aws-config-custom-rules/evalutions"
)

func Handler(event events.ConfigEvent) error {

	result, err := evalutions.EvaluationRule(event)

	if err != nil {
		return err
	}

	return evalutions.RunEvaluation(result)

}

func main() {
	lambda.Start(Handler)
}
