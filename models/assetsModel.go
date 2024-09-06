package models

import (
	"bytes"
	"fmt"
	"go-nuxt-blogs/pkg/errs"
	"io"
	"net/http"
	"strings"

	"github.com/fzzp/gotk"
)

// FileType 多媒体文件类型定义
type FileType int

const (
	_ FileType = iota
	IMAGE
	VIDEO
	AUDIO
	DOCUMENT
)

func (ft FileType) String() string {
	if ft >= IMAGE || ft <= DOCUMENT {
		// NOTE: 第一个是空
		return [...]string{"", "IMAGE", "VIDEO", "AUDIO", "DOCUMENT"}[ft]
	} else {
		return ""
	}
}

func StringAsFileType(s string) FileType {
	dataMap := map[string]FileType{
		"IMAGE":    IMAGE,
		"VIDEO":    VIDEO,
		"AUDIO":    AUDIO,
		"DOCUMENT": DOCUMENT,
	}
	ft, exists := dataMap[s]
	if exists {
		return ft
	}
	return 0
}

type Assets struct {
	ID          int64    `json:"id"`
	UserId      int64    `json:"userId"`
	Data        []byte   `json:"-"`
	FileType    FileType `json:"fileType"`
	Filename    string   `json:"filename"`
	Size        int64    `json:"size"`
	Description string   `json:"description"`
	CreatedAt   string   `json:"-"`
	UpdatedAt   string   `json:"-"`
	DeletedAt   *string  `json:"-"`
}

// BlobUploader 管理二进制文件上传
type BlobUploader struct {
	MaxFileSize      int64
	AllowedFileTypes []string
}

func NewBlobUploader(maxFileSize int64, allowTypes ...string) BlobUploader {
	if maxFileSize <= 0 {
		maxFileSize = 10 << 20 // 10mb
	}
	return BlobUploader{
		MaxFileSize:      maxFileSize,
		AllowedFileTypes: allowTypes,
	}
}

func (u *BlobUploader) GetBlob(r *http.Request) ([]*Assets, *gotk.ApiError) {
	err := r.ParseMultipartForm(int64(u.MaxFileSize))
	if err != nil {
		return nil, errs.ErrBadRequest.AsException(err, "解析表单数据失败")
	}

	var assets []*Assets

	for _, fHeaders := range r.MultipartForm.File {
		for _, hdr := range fHeaders {
			assets, err = func(assets []*Assets) ([]*Assets, error) {
				var a Assets
				infile, err := hdr.Open()
				if err != nil {
					return nil, err
				}
				defer infile.Close()

				if hdr.Size > int64(u.MaxFileSize) {
					return nil, errs.ErrBadRequest.AsException(err, fmt.Sprintf("上传文件必须小于 %d", u.MaxFileSize))
				}

				buff := make([]byte, 512)
				_, err = infile.Read(buff) // 读取 512 字节，判断文件类型
				if err != nil {
					return nil, errs.ErrBadRequest.AsException(err, "获取文件类型错误")
				}

				allowed := false
				mimeType := http.DetectContentType(buff)
				if len(u.AllowedFileTypes) > 0 {
					for _, v := range u.AllowedFileTypes {
						if strings.EqualFold(mimeType, v) {
							allowed = true
						}
					}
				} else {
					allowed = true
				}

				if !allowed {
					return nil, errs.ErrBadRequest.AsException(err, "不支持文件类型："+mimeType)
				}

				_, err = infile.Seek(0, 0)
				if err != nil {
					return nil, errs.ErrServerError.AsException(err, "设置文件偏移量错误")
				}

				dataBuff := new(bytes.Buffer)
				n, err := io.Copy(dataBuff, infile)
				if err != nil {
					return nil, errs.ErrServerError.AsException(err, "读取文件错误")
				}

				if n == 0 {
					return nil, errs.ErrServerError.AsException(err, "没有读取到文件")

				}
				a.Data = dataBuff.Bytes()

				// 设置其他信息
				a.Size = hdr.Size
				a.Filename = hdr.Filename

				assets = append(assets, &a)

				return assets, nil
			}(assets)

			if err != nil {
				return assets, errs.ErrServerError.AsException(err, "读取文件失败")
			}
		}
	}
	return assets, nil
}
