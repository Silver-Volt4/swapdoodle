package globals

import (
	"context"
	"fmt"
	"strings"

	"github.com/PretendoNetwork/nex-go/v2/types"
	"github.com/minio/minio-go/v7"
)

func S3StatObject(bucket, key string) (minio.ObjectInfo, error) {
	return MinIOClient.StatObject(context.TODO(), bucket, key, minio.StatObjectOptions{})
}

func S3SetFileContent(bucket, key string, content string) (minio.ObjectInfo, error) {
	MinIOClient.PutObject(context.TODO(), bucket, key, strings.NewReader(content), int64(len(content)), minio.PutObjectOptions{})
	return MinIOClient.StatObject(context.TODO(), bucket, key, minio.StatObjectOptions{})
}

func S3GetNotificationKey(pid types.PID) string {
	return fmt.Sprintf("notifications/%d", pid)
}

func S3GetLetterKey(dataID types.UInt32) string {
	return fmt.Sprintf("letters/%d.bin", dataID)
}

func S3ObjectSize(bucket, key string) (uint64, error) {
	info, err := S3StatObject(bucket, key)
	if err != nil {
		return 0, err
	}

	return uint64(info.Size), nil
}
