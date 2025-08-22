package service

import (
	"errors"
	"fmt"
	"io"
	"kgm2flac-backend/internal/utils"
	"os"
	"path/filepath"
	common "unlock-music.dev/cli/algo/common"
	"unlock-music.dev/cli/algo/kgm"
)

type DecryptService struct{}

func NewDecryptService() *DecryptService {
	return &DecryptService{}
}

func (s *DecryptService) DecryptKgmFile(inPath string) (outPath string, cleanup func(), err error) {
	// 从原decryptKgmPureGo函数迁移
	in, err := os.Open(inPath)
	if err != nil {
		return "", func() {}, err
	}
	defer in.Close()

	dec := kgm.NewDecoder(&common.DecoderParams{Reader: in})
	if err := dec.Validate(); err != nil {
		return "", func() {}, fmt.Errorf("不是有效的 KGM/KGMA/VPR 文件: %w", err)
	}

	outPath = filepath.Join(os.TempDir(), fmt.Sprintf("kgm_dec_%s.bin", utils.RandHex(8)))
	out, e := os.Create(outPath)
	if e != nil {
		return "", func() {}, e
	}
	defer out.Close()

	buf := make([]byte, 64*1024)
	for {
		n, e := dec.Read(buf)
		if n > 0 {
			if _, werr := out.Write(buf[:n]); werr != nil {
				return "", func() {}, werr
			}
		}
		if errors.Is(e, io.EOF) {
			break
		}
		if e != nil {
			return "", func() {}, e
		}
	}
	return outPath, func() { _ = os.Remove(outPath) }, nil
}
