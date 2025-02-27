package helper

import (
	"io"
	"log"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"

	"github.com/BeeTimeClock/BeeTimeClock-Server/core"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"golang.org/x/net/context"
)

func getMinioClient(env *core.Environment) (*minio.Client, error) {
	minioClient, err := minio.New(env.Storage.Endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(env.Storage.AccessKeyID, env.Storage.SecretAccessKey, ""),
	})

	if err != nil {
		return nil, err
	}

	return minioClient, nil
}

func prepareLocalFiles(env *core.Environment) error {
	return os.MkdirAll(env.UploadPath, os.ModePerm)
}

func SaveFile(env *core.Environment, fileName string, file multipart.File, fileMeta *multipart.FileHeader) error {
	if env.Storage.HasS3() {
		return s3UploadFile(env, fileName, file, fileMeta)
	} else {
		return localUploadFile(env, fileName, file, fileMeta)
	}
}

func localUploadFile(env *core.Environment, fileName string, file multipart.File, fileMeta *multipart.FileHeader) error {
	err := prepareLocalFiles(env)
	if err != nil {
		return err
	}

	out, err := os.Create(filepath.Join(env.UploadPath, fileName))
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func s3UploadFile(env *core.Environment, fileName string, file multipart.File, fileMeta *multipart.FileHeader) error {
	client, err := getMinioClient(env)
	if err != nil {
		return err
	}

	ctx := context.Background()
	_, err = client.PutObject(ctx, env.Storage.BucketName, fileName, file, fileMeta.Size, minio.PutObjectOptions{
		UserMetadata: map[string]string{
			"OriginalFilename": fileMeta.Filename,
		},
		ContentType: fileMeta.Header.Get("Content-Type"),
	})
	if err != nil {
		return err
	}

	return nil
}

func GetFile(env *core.Environment, fileName string) (io.Reader, *multipart.FileHeader, error) {
	if env.Storage.HasS3() {
		return s3GetFile(env, fileName)
	} else {
		return localGetFile(env, fileName)
	}
}

func localGetFile(env *core.Environment, fileName string) (io.Reader, *multipart.FileHeader, error) {
	fullPath := filepath.Join(env.UploadPath, fileName)
	stat, err := os.Stat(fullPath)
	if err != nil {
		return nil, nil, err
	}

	file, err := os.Open(fullPath)

	fileHeader := multipart.FileHeader{
		Filename: fileName,
		Header:   textproto.MIMEHeader{},
		Size:     stat.Size(),
	}

	return file, &fileHeader, err
}

func s3GetFile(env *core.Environment, fileName string) (*minio.Object, *multipart.FileHeader, error) {
	client, err := getMinioClient(env)
	if err != nil {
		return nil, nil, err
	}

	object, err := client.GetObject(context.Background(), env.Storage.BucketName, fileName, minio.GetObjectOptions{})
	if err != nil {
		return nil, nil, err
	}

	stat, err := object.Stat()
	if err != nil {
		return nil, nil, err
	}

	fileHeader := multipart.FileHeader{
		Filename: fileName,
		Header:   textproto.MIMEHeader{},
		Size:     stat.Size,
	}

	return object, &fileHeader, err
}

func ExistsFile(env *core.Environment, fileName string) (bool, error) {
	if env.Storage.HasS3() {
		return s3ExistsFile(env, fileName)
	} else {
		return localExistsFile(env, fileName)
	}
}

func localExistsFile(env *core.Environment, fileName string) (bool, error) {
	_, err := os.Stat(filepath.Join(env.UploadPath, fileName))
	return err == nil, nil
}

func s3ExistsFile(env *core.Environment, fileName string) (bool, error) {
	client, err := getMinioClient(env)
	if err != nil {
		return false, err
	}

	objectsCh := client.ListObjects(context.Background(), env.Storage.BucketName, minio.ListObjectsOptions{
		Prefix: fileName,
	})

	hasObject := false
	for range objectsCh {
		hasObject = true
		break
	}

	return hasObject, nil
}
