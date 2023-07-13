package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

const (
	assumeRoleTimeout = time.Second * 5
)

func unsetEnvVars(env []string, unsetVarNames ...string) []string {
	// create a map of variables to unset (used as a set)
	shouldUnset := make(map[string]struct{})
	for _, varName := range unsetVarNames {
		shouldUnset[varName] = struct{}{}
	}
	// iterate over environment variables, unseting
	// any present in the shouldUnset map
	newEnv := []string{}
	for _, envVar := range env {
		envVarName := strings.Split(envVar, "=")[0]
		if _, unset := shouldUnset[envVarName]; unset {
			continue
		}
		newEnv = append(newEnv, envVar)
	}
	return newEnv
}

func toEnvVars(set map[string]string) []string {
	vars := []string{}
	for k, v := range set {
		vars = append(vars, fmt.Sprintf("%s=%s", k, v))
	}
	return vars
}

func errorOut(msg string) {
	fmt.Printf("\n\tUsage: %s <role-arn-to-assume> <command> [args...]\n\n", os.Args[0])
	fmt.Println(fmt.Sprintf("ERROR: %s", msg))
	os.Exit(1)
}

func main() {
	ctx := context.Background()

	if len(os.Args) < 3 {
		errorOut("not enough arguments")
	}

	roleArn := os.Args[1]
	command := os.Args[2]
	args := os.Args[3:]

	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		errorOut(fmt.Sprintf("unable to load AWS SDK config, %v", err))
	}

	assumeRoleContext, assumeRoleContextCancel := context.WithTimeout(ctx, assumeRoleTimeout)
	defer assumeRoleContextCancel()

	stsService := sts.NewFromConfig(cfg)
	assumeRoleOutput, err := stsService.AssumeRole(assumeRoleContext, &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(fmt.Sprintf("WAR-%d", time.Now().UnixNano())),
	})
	if err != nil {
		errorOut(fmt.Sprintf("failed to assume role \"%s\": %s", roleArn, err))
	}

	cmd := exec.Command(command, args...)
	cmd.Env = append(
		unsetEnvVars(os.Environ(), []string{
			"AWS_PROFILE",
			"AWS_ACCESS_KEY_ID",
			"AWS_SECRET_ACCESS_KEY",
			"AWS_SESSION_TOKEN",
		}...),
		toEnvVars(map[string]string{
			"AWS_ACCESS_KEY_ID":     aws.ToString(assumeRoleOutput.Credentials.AccessKeyId),
			"AWS_SECRET_ACCESS_KEY": aws.ToString(assumeRoleOutput.Credentials.SecretAccessKey),
			"AWS_SESSION_TOKEN":     aws.ToString(assumeRoleOutput.Credentials.SessionToken),
		})...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err = cmd.Start(); err != nil {
		errorOut(fmt.Sprintf("failed to start command: %v", err))
	}

	cmd.Wait()
}
