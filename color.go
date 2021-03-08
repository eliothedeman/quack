package quack

// EnableColor will enable terminal color output via escape codes.
func EnableColor(o *options) {
	o.colorEnabled = true
}

// DisableColor will disable terminal color output.
func DisableColor(o *options) {
	o.colorEnabled = false
}
