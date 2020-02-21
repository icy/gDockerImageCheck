/*
  Purpose : Check if docker image does exist in Docker repository
  Author  : Ky-Anh Huynh
  Date    : 2020-02
  Support :

  - [x] ECR
  - [ ] as a library
*/

package main

import "fmt"
import "os"
import "regexp"
import "context"

import "github.com/aws/aws-sdk-go-v2/aws"
import "github.com/aws/aws-sdk-go-v2/aws/external"
import "github.com/aws/aws-sdk-go-v2/service/ecr"

func log2(msg string) {
	fmt.Fprintf(os.Stderr, msg)
}

func log2exit(retval int, msg string) {
	log2(msg)
	os.Exit(retval)
}

func GetEnvWithDefault(key string, default_value string) string {
	val := os.Getenv(key)
	if len(val) == 0 {
		val = default_value
	}
	return val
}

func ECRCheckImage(tag string, repo string, ecr_id string, region string, profile string) bool {
	cfg, err := external.LoadDefaultAWSConfig(
		external.WithSharedConfigProfile(profile),
		external.WithRegion(region),
	)
	if err != nil {
		log2(fmt.Sprintf(":: Error: '%s'.\n", err.Error()))
		return false
	}

	img := &ecr.ImageIdentifier{
		ImageTag: aws.String(tag),
	}

	input := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repo),
		RegistryId:     nil,
		ImageIds:       []ecr.ImageIdentifier{*img},
	}

	if len(ecr_id) > 0 {
		input.RegistryId = aws.String(ecr_id)
	}

	svc := ecr.New(cfg)
	req := svc.DescribeImagesRequest(input)
	if _, err := req.Send(context.TODO()); err == nil {
		log2(fmt.Sprintf(":: Info: Found ECR image %s:%s\n", repo, tag))
		return true
	} else {
		log2(fmt.Sprintf(":: Error: '%s'.\n", err.Error()))
		return false
	}
}

// <ECR_ID>.dkr.ecr.<AWS_REGION>.amazonaws.com/<REPO>:<TAG>
func ECRdetect(image string) (string, string, string, string) {
	re := regexp.MustCompile(`^([0-9]+)\.dkr\.ecr.([^.]+)\.amazonaws\.com/([^:]+):(.+)$`)
	result := re.FindStringSubmatch(image)
	if len(result) >= 5 {
		return result[1], result[2], result[3], result[4]
	}
	return "", "", "", ""
}

func main() {
	err := 0

	profile := GetEnvWithDefault("AWS_PROFILE", "default")

	for _, arg := range os.Args[1:] {
		ecr_id, region, repo, tag := ECRdetect(arg)
		if len(ecr_id) > 0 {
			if !ECRCheckImage(tag, repo, ecr_id, region, profile) {
				err += 1
			}
		}
	}

  /*
	region := GetEnvWithDefault("AWS_REGION", "eu-west-1")
	tag := GetEnvWithDefault("IMAGE_TAG", "")
	repo := GetEnvWithDefault("ECR_REPO", "")
	ecr_id := GetEnvWithDefault("ECR_ID", "")

	if len(repo) > 0 && len(tag) > 0 {
		if !ECRCheckImage(tag, repo, ecr_id, region, profile) {
			err += 1
		}
	}
	if err > 0 {
		err = 1
	}
  */

	os.Exit(err)
}
