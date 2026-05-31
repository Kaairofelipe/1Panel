package files

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"time"

	"github.com/1Panel-dev/1Panel/agent/global"
	"github.com/1Panel-dev/1Panel/agent/utils/cmd"
	"github.com/1Panel-dev/1Panel/agent/utils/common"
)

type TarGzArchiver struct {
}

func NewTarGzArchiver() ShellArchiver {
	return &TarGzArchiver{}
}

func (t TarGzArchiver) Extract(ctx context.Context, filePath, dstDir string, secret string) error {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return fmt.Errorf("failed to create destination dir: %w", err)
	}
	var err error
	if len(secret) != 0 {
		script := "openssl enc -d -aes-256-cbc -k \"$1\" -in \"$2\" | tar -zxvf - -C \"$3\" > /dev/null 2>&1"
		global.LOG.Debug(fmt.Sprintf("openssl enc -d -aes-256-cbc -k '******' -in '%s' | tar -zxvf - -C '%s' > /dev/null 2>&1", filePath, dstDir))
		err = cmd.NewCommandMgr(cmd.WithContext(ctx)).RunBashCWithArgs(script, secret, filePath, dstDir)
	} else {
		global.LOG.Debug(fmt.Sprintf("tar -zxvf '%s' -C '%s' > /dev/null 2>&1", filePath, dstDir))
		err = cmd.NewCommandMgr(cmd.WithContext(ctx)).Run("tar", "-zxvf", filePath, "-C", dstDir)
	}
	if err != nil {
		return err
	}
	return nil
}

func (t TarGzArchiver) Compress(ctx context.Context, sourcePaths []string, dstFile string, secret string) error {
	tmpFile := path.Join(global.Dir.TmpDir, fmt.Sprintf("%s%s.tar.gz", common.RandStr(50), time.Now().Format("20060102150405")))
	op := NewFileOp()
	var err error
	defer func() {
		_ = op.DeleteFile(tmpFile)
		if err != nil {
			_ = op.DeleteFile(dstFile)
		}
	}()

	aheadDir := filepath.Dir(sourcePaths[0])
	if len(aheadDir) == 0 {
		aheadDir = "/"
	}

	relativePaths := make([]string, len(sourcePaths))
	for i, sp := range sourcePaths {
		relativePaths[i] = filepath.Base(sp)
	}

	if len(secret) != 0 {
		script := "tar -zcf - -C \"$1\" \"${@:4}\" | openssl enc -aes-256-cbc -salt -k \"$2\" -out \"$3\""
		global.LOG.Debug(fmt.Sprintf("tar -zcf - -C \"%s\" ... | openssl enc -aes-256-cbc -salt -k ****** -out \"%s\"", aheadDir, tmpFile))
		args := append([]string{aheadDir, secret, tmpFile}, relativePaths...)
		argsWithScript := append([]string{script}, args...)
		err = cmd.NewCommandMgr(cmd.WithContext(ctx)).RunBashCWithArgs(argsWithScript...)
	} else {
		global.LOG.Debug(fmt.Sprintf("tar -zcf \"%s\" -C \"%s\" ...", tmpFile, aheadDir))
		args := append([]string{"-zcf", tmpFile, "-C", aheadDir}, relativePaths...)
		err = cmd.NewCommandMgr(cmd.WithContext(ctx)).Run("tar", args...)
	}
	if err != nil {
		return err
	}
	if err = op.Mv(tmpFile, dstFile); err != nil {
		return err
	}
	return nil
}
