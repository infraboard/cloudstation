package cli

import (
	"fmt"
	"net"
	"os"
	"path"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/infraboard/cloudstation/pkg/oss"
	"github.com/infraboard/cloudstation/pkg/oss/provider/ali"
)

const (
	// BuckName todo
	defaultBuckName = "devcloud-station"
	defaultEndpoint = ""
	defaultALIAK    = ""
	defaultALISK    = ""
)

var (
	buckName       string
	uploadFilePath string
	bucketEndpoint string
)

// uploadCmd represents the start command
var uploadCmd = &cobra.Command{
	Use:   "upload",
	Short: "上传文件到中转站",
	Long:  `上传文件到中转站`,
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := getProvider()
		if err != nil {
			return err
		}
		if uploadFilePath == "" {
			return fmt.Errorf("upload file path is missing")
		}

		// 为了防止文件都堆在一个文件夹里面 无法查看
		// 我们采用日期进行编码
		day := time.Now().Format("20060102")

		// 为了防止不同用户同一时间上传相同的文件
		// 我们采用用户的主机名作为前置
		hn, err := os.Hostname()
		if err != nil {
			ipAddr := getOutBindIp()
			if ipAddr == "" {
				hn = "unknown"
			} else {
				hn = ipAddr
			}
		}

		fn := path.Base(uploadFilePath)
		ok := fmt.Sprintf("%s/%s/%s", day, hn, fn)
		err = p.UploadFile(buckName, ok, uploadFilePath)
		if err != nil {
			return err
		}
		return nil
	},
}

func getOutBindIp() string {
	conn, err := net.Dial("udp", "baidu.com:80")
	if err != nil {
		return ""
	}
	defer conn.Close()

	addr := strings.Split(conn.LocalAddr().String(), ":")
	if len(addr) == 0 {
		return ""
	}

	return addr[0]
}

func getProvider() (p oss.Provider, err error) {
	switch ossProvider {
	case "ali":
		fmt.Printf("上传云商: 阿里云[%s]\n", defaultEndpoint)
		if aliAccessID == "" {
			aliAccessID = defaultALIAK
		}
		if aliAccessKey == "" {
			aliAccessKey = defaultALISK
		}
		fmt.Printf("上传用户: %s\n", aliAccessID)
		p, err = ali.NewProvider(bucketEndpoint, aliAccessID, aliAccessKey)
		return
	case "tx":
		return nil, fmt.Errorf("not impl")
	default:
		return nil, fmt.Errorf("unknown oss privier options [ali/tx]")
	}
}

func init() {
	uploadCmd.PersistentFlags().StringVarP(&uploadFilePath, "file_path", "f", "", "upload file path")
	uploadCmd.PersistentFlags().StringVarP(&buckName, "bucket_name", "b", defaultBuckName, "upload oss bucket name")
	RootCmd.AddCommand(uploadCmd)
}
