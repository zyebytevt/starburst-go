package vseeface

import (
	"os/exec"
	"time"
)

func sendKeyShortcut(key string) error {
	exec.Command("/usr/bin/xdotool", "keydown", "--window", vsfWindowId, "CTRL").Run()
	time.Sleep(time.Millisecond * 100)
	exec.Command("/usr/bin/xdotool", "keydown", "--window", vsfWindowId, "SHIFT").Run()
	time.Sleep(time.Millisecond * 100)
	exec.Command("/usr/bin/xdotool", "keydown", "--window", vsfWindowId, key).Run()
	time.Sleep(time.Millisecond * 100)
	exec.Command("/usr/bin/xdotool", "keyup", "--window", vsfWindowId, key).Run()
	time.Sleep(time.Millisecond * 100)
	exec.Command("/usr/bin/xdotool", "keyup", "--window", vsfWindowId, "SHIFT").Run()
	time.Sleep(time.Millisecond * 100)
	exec.Command("/usr/bin/xdotool", "keyup", "--window", vsfWindowId, "CTRL").Run()

	return nil
}
