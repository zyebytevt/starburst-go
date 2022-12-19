package dbus

import (
	"github.com/godbus/dbus"
	"github.com/zyebytevt/starburst-go/lib"
)

func callCallback(button *lib.Button) error {
	dest := button.Config.Parameters["destination"].(string)
	path := button.Config.Parameters["path"].(string)
	method := button.Config.Parameters["method"].(string)

	var params []any

	if check := button.Config.Parameters["params"]; check != nil {
		params = check.([]any)
	} else {
		params = nil
	}

	return connection.Object(dest, dbus.ObjectPath(path)).Call(method, 0, params).Err
}
