package file

import "os"

// PermissionSetter represent available file permissions to set
type PermissionSetter func(po *permissionSettings)

type permissionSettings struct {
	mode os.FileMode
}

// SetPermissions sets the permissions mode value
func SetPermissions(p os.FileMode) PermissionSetter {
	return func(opts *permissionSettings) {
		opts.mode = p
	}
}

func setDefaults(df os.FileMode) permissionSettings {
	return permissionSettings{
		mode: df,
	}
}
