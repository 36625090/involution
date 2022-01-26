package framework

func (b *Backend) Validate(in interface{}) error {
	return b.validator.Struct(in)
}