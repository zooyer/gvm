package utils

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
	"golang.org/x/sys/windows/registry"
)

func getAbsEnv(key string) (val string, err error) {
	reg, exists, err := registry.CreateKey(registry.CURRENT_USER, "Environment", registry.ALL_ACCESS)
	if err != nil {
		return
	}
	defer reg.Close()

	if exists {
		if val, _, err = reg.GetStringValue(key); err != nil {
			if err == registry.ErrNotExist {
				return "", nil
			}
			return
		}
	}

	return
}

func setAbsEnv(key, val string) (err error) {
	reg, exists, err := registry.CreateKey(registry.CURRENT_USER, "Environment", registry.ALL_ACCESS)
	if err != nil {
		return
	}
	defer reg.Close()

	if exists {
		if err = reg.SetStringValue(key, val); err == nil {
			text, err := syscall.UTF16PtrFromString("Environment")
			if err != nil {
				return err
			}

			win.SendMessage(win.HWND_BROADCAST, win.WM_SETTINGCHANGE, 0, uintptr(unsafe.Pointer(text)))
		}
	}

	return
}
