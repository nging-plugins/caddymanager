package handler

import (
	"context"
	"errors"
	"time"

	"github.com/admpub/nging/v5/application/library/charset"
	"github.com/admpub/nging/v5/application/library/cron"
	"github.com/admpub/nging/v5/application/library/notice"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"

	conf "github.com/nging-plugins/servermanager/application/library/config"
)

// CommandJob 计划任务调用方式
func CommandJob(id string) cron.Runner {
	return func(timeout time.Duration) (out string, runingErr string, onErr error, isTimeout bool) {
		idN := param.AsUint(id)
		if idN < 1 {
			onErr = errors.New(`Invalid ID: ` + id)
			return
		}
		m, result, err := ExecCommand(idN)
		if err != nil {
			onErr = err
			runingErr = remoteCmdResultPrefix() + err.Error()
			return
		}
		if len(result) > 0 {
			out += remoteCmdResultPrefix() + result + com.StrLF + com.StrLF
		}
		if m.Remote == `Y` || m.Id == 0 {
			return
		}

		wOut := cron.NewOutputWriter()
		wErr := cron.NewOutputWriter()
		noticeSender := notice.CustomOutputNoticer(wOut, wErr)
		env := conf.ParseEnvSlice(m.Env)
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		params := cron.CmdParams(m.Command)
		cmd := com.CreateCmdWithContext(ctx, params, func(b []byte) (e error) {
			if com.IsWindows {
				b, e = charset.Convert(`gbk`, `utf-8`, b)
				if e != nil {
					return e
				}
			}
			noticeSender(string(b), 1)
			return nil
		})
		if len(m.WorkDirectory) > 0 {
			cmd.Dir = m.WorkDirectory
		}
		if len(env) > 0 {
			cmd.Env = env
		}
		cmd.Stderr = com.CmdResultCapturer{Do: func(b []byte) (e error) {
			if com.IsWindows {
				b, e = charset.Convert(`gbk`, `utf-8`, b)
				if e != nil {
					return e
				}
			}
			noticeSender(string(b), 0)
			return nil
		}}
		if e := cmd.Run(); e != nil {
			isTimeout = errors.Is(e, context.DeadlineExceeded)
			noticeSender(e.Error(), 0)
		}
		out += localCmdResultPrefix() + wOut.String()
		runingErr += wErr.String()
		return
	}
}
