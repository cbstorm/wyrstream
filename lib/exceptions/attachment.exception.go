package exceptions

var (
	IMAGE_SCALE_EXCEPTION = func() *Exception { return NewException("IMAGE_SCALE_EXCEPTION").SetStatus(400) }
)
