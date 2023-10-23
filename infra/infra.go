package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsevents"
	"github.com/aws/aws-cdk-go/awscdk/v2/awseventstargets"
	"github.com/aws/aws-cdk-go/awscdklambdagoalpha/v2"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type InfraStackProps struct {
	awscdk.StackProps
}

type ApiConstructProps struct{}

func NewInfraStack(scope constructs.Construct, id string, props *InfraStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	table := awsdynamodb.NewTable(stack, jsii.String("Table"), &awsdynamodb.TableProps{
		PartitionKey: &awsdynamodb.Attribute{
			Name: jsii.String("id"),
			Type: awsdynamodb.AttributeType_STRING,
		},
		BillingMode: awsdynamodb.BillingMode_PAY_PER_REQUEST,
	})

	lambdaEnvironment := map[string]*string{
		"TABLE_NAME": table.TableName(),
	}

	lambda := awscdklambdagoalpha.NewGoFunction(stack, jsii.String("Lambda"), &awscdklambdagoalpha.GoFunctionProps{
		Entry:       jsii.String("../cmd/server"),
		Environment: &lambdaEnvironment,
		Description: jsii.String("Update a DynamoDB table with a lambda function."),
	})

	table.GrantReadWriteData(lambda)

	rule := awsevents.NewRule(stack, jsii.String("ScheduleRule"), &awsevents.RuleProps{
		Schedule: awsevents.Schedule_Rate(awscdk.Duration_Minutes(jsii.Number(1))),
	})

	rule.AddTarget(awseventstargets.NewLambdaFunction(lambda, nil))

	NewApi(stack, "Api", table)

	return stack
}

func NewApi(scope constructs.Construct, id string, table awsdynamodb.Table) {
	taskDef := awsecs.NewTaskDefinition(scope, jsii.String("TaskDefinition"), &awsecs.TaskDefinitionProps{
		Compatibility: awsecs.Compatibility_FARGATE,
		Cpu:           jsii.String("256"),
		MemoryMiB:     jsii.String("512"),
	})

	taskDef.AddContainer(jsii.String("Container"), &awsecs.ContainerDefinitionOptions{
		Image: awsecs.ContainerImage_FromAsset(jsii.String("../"), &awsecs.AssetImageProps{
			AssetName: jsii.String("api"),
		}),
		PortMappings: &[]*awsecs.PortMapping{
			{
				ContainerPort: jsii.Number(8000),
				HostPort:      jsii.Number(8000),
			},
		},
		Environment: &map[string]*string{
			"TABLE_NAME": table.TableName(),
		},
		Logging: awsecs.LogDriver_AwsLogs(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("api"),
		}),
	})

	cluster := awsecs.NewCluster(scope, jsii.String("Cluster"), &awsecs.ClusterProps{})

	service := awsecs.NewFargateService(scope, jsii.String("Service"), &awsecs.FargateServiceProps{
		Cluster:        cluster,
		TaskDefinition: taskDef,
		ServiceName:    jsii.String("api"),
	})

	table.GrantReadWriteData(service.TaskDefinition().TaskRole())

	alb := awselasticloadbalancingv2.NewApplicationLoadBalancer(scope, jsii.String("LoadBalancer"), &awselasticloadbalancingv2.ApplicationLoadBalancerProps{
		InternetFacing:   jsii.Bool(true),
		LoadBalancerName: jsii.String("api"),
		Vpc:              cluster.Vpc(),
	})

	listener := alb.AddListener(jsii.String("Listener"), &awselasticloadbalancingv2.BaseApplicationListenerProps{
		Port: jsii.Number(80),
		Open: jsii.Bool(true),
	})

	listener.AddTargets(jsii.String("Target"), &awselasticloadbalancingv2.AddApplicationTargetsProps{
		Port: jsii.Number(8000),
		Targets: &[]awselasticloadbalancingv2.IApplicationLoadBalancerTarget{
			service,
		},
		HealthCheck: &awselasticloadbalancingv2.HealthCheck{
			Interval: awscdk.Duration_Seconds(jsii.Number(60)),
			Path:     jsii.String("/ping"),
			Timeout:  awscdk.Duration_Seconds(jsii.Number(5)),
		},
	})

}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewInfraStack(app, "YearOfDecay2", &InfraStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	// return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String("533674317867"),
		Region:  jsii.String("us-east-1"),
	}

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
		Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	}
}
