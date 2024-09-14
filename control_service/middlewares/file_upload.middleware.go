package middlewares

import (
	"bufio"

	"github.com/cbstorm/wyrstream/control_service/common"
	"github.com/cbstorm/wyrstream/lib/dtos"
)

func FileUploadMiddleware(file_key string) HttpMiddleware {
	return func(c common.IHttpContext) error {
		form, err := c.MultipartForm()
		if err != nil {
			return common.ResponseError(c, err)
		}
		feat := c.Queries()["feat"]
		files := form.File[file_key]
		f_inp := make([]*dtos.FileInput, 0)
		for _, f := range files {
			fr, err := f.Open()
			bs := make([]byte, f.Size)
			if err != nil {
				return common.ResponseError(c, err)
			}
			if _, err := bufio.NewReader(fr).Read(bs); err != nil {
				return common.ResponseError(c, err)
			}
			f_inp = append(f_inp, &dtos.FileInput{Size: f.Size, Bytes: bs, Feat: feat, MimeType: f.Header.Get("Content-Type"), Name: f.Filename})
			if err := fr.Close(); err != nil {
				return common.ResponseError(c, err)
			}
		}
		c.Locals(common.ReqFileKey{}, &f_inp)
		return c.Next()
	}
}
