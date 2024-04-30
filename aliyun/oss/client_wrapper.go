package aliyunOss

import (
	"errors"
	"fmt"
	aliyunOss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/oneliang/util-golang/file"
	"log"
	"path/filepath"
)

type ClientWrapper struct {
	endpoint        string
	accessKeyId     string
	accessKeySecret string
	bucket          string
	ossClient       *aliyunOss.Client
}

func NewClientWrapper(endpoint string, accessKeyId string, accessKeySecret string, bucket string) (*ClientWrapper, error) {
	ossClient, err := aliyunOss.New(endpoint, accessKeyId, accessKeySecret)
	if err != nil {
		log.Fatalf("connect to oss error, endpoint:%s, error:%v", endpoint, err)
		return nil, err
	}
	clientWrapper := &ClientWrapper{
		endpoint:        endpoint,
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		bucket:          bucket,
		ossClient:       ossClient,
	}
	return clientWrapper, nil
}
func (this *ClientWrapper) uploadFileToOss(objectKey string, filePath string) error {
	log.Printf("upload file to oss, object key:%s, file path:%s", objectKey, filePath)
	clientBucket, err := this.ossClient.Bucket(this.bucket)
	if err != nil {
		log.Printf("get bucket error, bucket:%s, error%v", this.bucket, err)
		return err
	}
	objectAcl := aliyunOss.ObjectACL(aliyunOss.ACLPublicRead)
	if err = clientBucket.PutObjectFromFile(objectKey, filePath, objectAcl); err != nil {
		log.Printf("put object from file error:%v", err)
		return err
	}
	return nil
}

func (this *ClientWrapper) UploadFile(objectKey string, filePath string) (string, error) {
	if !file.Exists(filePath) {
		errorMessage := fmt.Sprintf("upload file check error, file has not exists:%s", filePath)
		log.Printf(errorMessage)
		return "", errors.New(errorMessage)
	}
	baseFile := filepath.Base(filePath)
	if err := this.uploadFileToOss(objectKey, filePath); err != nil {
		log.Printf("upload file to oss error:%v", err)
		return "", err
	}
	fileUrl := fmt.Sprintf("https://%s.%s/%s",
		this.bucket,
		this.endpoint,
		objectKey,
		baseFile)
	return fileUrl, nil
}

func (this *ClientWrapper) downloadFileFromOss(objectKey string, filePath string) error {
	log.Printf("download file from oss, object key:%s, file path:%s", objectKey, filePath)
	clientBucket, err := this.ossClient.Bucket(this.bucket)
	if err != nil {
		log.Printf("get bucket error, bucket:%s, error%v", this.bucket, err)
		return err
	}
	objectAcl := aliyunOss.ObjectACL(aliyunOss.ACLPublicRead)
	if err = clientBucket.GetObjectToFile(objectKey, filePath, objectAcl); err != nil {
		log.Printf("get object to file error:%v", err)
		return err
	}
	return nil
}

func (this *ClientWrapper) DownloadFile(objectKey string, filePath string) error {
	if err := this.downloadFileFromOss(objectKey, filePath); err != nil {
		log.Printf("download file from oss error:%v", err)
		return err
	}
	return nil
}
