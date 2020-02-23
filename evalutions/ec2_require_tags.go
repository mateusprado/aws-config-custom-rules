package evalutions

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/configservice"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/mateusprado/aws-config-custom-rules/model"
)

type EvaluationResult struct {
	ResourceId     string
	ComplianceType string
	Token          string
	Time           string
}

func EvaluationRule(event events.ConfigEvent) (*EvaluationResult, error) {

	var invokingEvent model.InvokingEvent
	err := json.Unmarshal([]byte(event.InvokingEvent), &invokingEvent)
	if err != nil {
		return nil, err
	}

	// unmarshal the ConfigurationItem
	var instance ec2.Instance
	err = json.Unmarshal(invokingEvent.ConfigurationItem.Configuration, &instance)
	if err != nil {
		return nil, err
	}

	if event.EventLeftScope {
		return &EvaluationResult{
			ResourceId:     invokingEvent.ConfigurationItem.ResourceID,
			ComplianceType: configservice.ComplianceTypeNotApplicable,
			Time:           invokingEvent.NotificationCreationTime,
			Token:          event.ResultToken,
		}, nil
	}

	if resource, ok := invokingEvent.ConfigurationItem.Tags["Platform"]; ok && resource == "sre" {
		return &EvaluationResult{
			ResourceId:     invokingEvent.ConfigurationItem.ResourceID,
			ComplianceType: configservice.ComplianceTypeCompliant,
			Time:           invokingEvent.NotificationCreationTime,
			Token:          event.ResultToken,
		}, nil
	}

	return &EvaluationResult{
		ResourceId:     invokingEvent.ConfigurationItem.ResourceID,
		ComplianceType: configservice.ComplianceTypeNonCompliant,
		Time:           invokingEvent.NotificationCreationTime,
		Token:          event.ResultToken,
	}, nil

}

func RunEvaluation(evaluation *EvaluationResult) error {
	service := configservice.New(session.Must(session.NewSession()))

	orderingTimestamp, err := time.Parse(time.RFC3339, evaluation.Time)
	if err != nil {
		return err
	}

	_, err = service.PutEvaluations(
		&configservice.PutEvaluationsInput{
			Evaluations: []*configservice.Evaluation{
				{
					ComplianceResourceId:   &evaluation.ResourceId,
					ComplianceResourceType: aws.String("AWS::EC2::Instance"),
					ComplianceType:         aws.String(evaluation.ComplianceType),
					OrderingTimestamp:      &orderingTimestamp,
				},
			},
			ResultToken: aws.String(evaluation.Token),
		},
	)

	if err != nil {
		fmt.Println(err)
		return (err)
	}

	return nil
}
