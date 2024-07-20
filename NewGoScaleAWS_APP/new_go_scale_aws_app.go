package main

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigateway"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type NewGoScaleAwsAppStackProps struct {
	awscdk.StackProps
}

func NewNewGoScaleAwsAppStack(scope constructs.Construct, id string, props *NewGoScaleAwsAppStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// Create the DynamoDB table
	table := awsdynamodb.NewTable(stack, jsii.String("FEMUserTable"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("username"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		TableName: jsii.String("userTable"),
	})

	// Create the IAM role for the Lambda function
	lambdaRole := awsiam.NewRole(stack, jsii.String("myLambdaFuncServiceRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("lambda.amazonaws.com"), nil),
	})

	// Attach DynamoDB access policy to the role
	lambdaRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Actions:   jsii.Strings("dynamodb:GetItem", "dynamodb:PutItem"),
		Resources: jsii.Strings(*table.TableArn()),
	}))

	// Create the Lambda function
	myFunction := awslambda.NewFunction(stack, jsii.String("myLambdaFunc"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_PROVIDED_AL2023(),
		Code:    awslambda.Code_FromAsset(jsii.String("lambda/function.zip"), nil),
		Handler: jsii.String("main"),
		Role:    lambdaRole,
	})

	// Grant Lambda function read/write permissions to the DynamoDB table
	table.GrantReadWriteData(myFunction)

	// Create the API Gateway
	api := awsapigateway.NewRestApi(stack, jsii.String("myAPIGateway"), &awsapigateway.RestApiProps{
		DefaultCorsPreflightOptions: &awsapigateway.CorsOptions{
			AllowHeaders: jsii.Strings("Content-Type", "Authorization"),
			AllowMethods: jsii.Strings("GET", "POST", "PUT", "DELETE", "OPTIONS"),
			AllowOrigins: jsii.Strings("*"),
		},
		DeployOptions: &awsapigateway.StageOptions{
			LoggingLevel: awsapigateway.MethodLoggingLevel_INFO,
		},
	})

	// Create the Lambda integration
	integration := awsapigateway.NewLambdaIntegration(myFunction, nil)

	// Define the API Gateway resources and methods
	registerResource := api.Root().AddResource(jsii.String("register"), nil)
	registerResource.AddMethod(jsii.String("POST"), integration, nil)

	loginResource := api.Root().AddResource(jsii.String("login"), nil)
	loginResource.AddMethod(jsii.String("POST"), integration, nil)

	// Output the API Gateway endpoint
	awscdk.NewCfnOutput(stack, jsii.String("APIGatewayURL"), &awscdk.CfnOutputProps{
		Value: api.Url(),
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewNewGoScaleAwsAppStack(app, "NewGoScaleAwsAppStack", &NewGoScaleAwsAppStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to be deployed
func env() *awscdk.Environment {
	return nil
	// Uncomment if you know exactly what account and region you want to deploy the stack to
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }
	// Uncomment to specialize this stack for the AWS Account and Region that are implied by the current CLI configuration
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
