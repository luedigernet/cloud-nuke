package aws

import (
	"time"

	"github.com/gruntwork-io/cloud-nuke/config"

	"github.com/aws/aws-sdk-go/aws"
	awsgo "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/gruntwork-io/cloud-nuke/logging"
	"github.com/gruntwork-io/go-commons/errors"
)

func getAllLambdaFunctions(session *session.Session, excludeAfter time.Time, configObj config.Config) (names []*string, err error) {
	svc := lambda.New(session)

	result, err := svc.ListFunctions(nil)

	if err != nil {
		return nil, errors.WithStackTrace(err)
	}

	for _, lambda := range result.Functions {
		include, err := shouldIncludeLambdaFunction(lambda, excludeAfter, configObj)
		if err != nil {
			panic(err)
		}
		if include {
			names = append(names, lambda.FunctionName)
		}
	}

	return names, nil
}

func nukeAllLambdaFunctions(session *session.Session, names []*string) error {
	svc := lambda.New(session)

	if len(names) == 0 {
		logging.Logger.Infof("No Lambda Functions to nuke in region %s", *session.Config.Region)
		return nil
	}

	logging.Logger.Infof("Deleting all Lambda Functions in region %s", *session.Config.Region)
	deletedNames := []*string{}

	for _, name := range names {
		params := &lambda.DeleteFunctionInput{
			FunctionName: name,
		}

		_, err := svc.DeleteFunction(params)

		if err != nil {
			logging.Logger.Errorf("[Failed] %s: %s", *name, err)
		} else {
			deletedNames = append(deletedNames, name)
			logging.Logger.Infof("Deleted Lambda Function: %s", awsgo.StringValue(name))
		}
	}

	logging.Logger.Infof("[OK] %d Lambda Function(s) deleted in %s", len(deletedNames), *session.Config.Region)
	return nil
}

func shouldIncludeLambdaFunction(lambdaFunction *lambda.FunctionConfiguration, excludeAfter time.Time, configObj config.Config) (result bool, err error) {

	if lambdaFunction == nil {
		return false, nil
	}

	layout := "2006-01-02T15:04:05.000+0000"
	lm, err := time.Parse(layout, *lambdaFunction.LastModified)
	if err != nil {
		return false, err
	}

	if excludeAfter.Before(aws.TimeValue(&lm)) {
		return false, nil
	}

	return config.ShouldInclude(
		aws.StringValue(lambdaFunction.FunctionName),
		configObj.Lambda.IncludeRule.NamesRegExp,
		configObj.Lambda.ExcludeRule.NamesRegExp,
	), nil
}
