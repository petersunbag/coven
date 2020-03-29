package coven

import "strings"

// StructOption describe options in converting struct.
// BannedFields lists fields that are not allowed to be converted in dst struct.
// AliasFields lists fields that are supposed to use alias field-name in dst struct when converting.
// Both BannedFields and AliasFields support nested fields.
// For example, `A.B` in BannedFields or key of AliasFields means `A` field in dst struct is a struct, and the option is applied to `B` field in `A`.
type StructOption struct {
	BannedFields []string
	AliasFields  map[string]string
}

// convert StructOption to structOption, and call structOption.parse()
func (o *StructOption) convert() *structOption {
	var oo structOption
	if o.BannedFields != nil {
		oo.BannedFields = make(map[string]struct{})
		for _, f := range o.BannedFields {
			oo.BannedFields[f] = struct{}{}
		}
	}
	if o.AliasFields != nil {
		oo.AliasFields = make(map[string]string)
	}
	for f, a := range o.AliasFields {
		oo.AliasFields[f] = a
	}
	oo.parse()
	return &oo
}

// structOption is an inner use version of StructOption, which use NestedOption to represent nested options.
type structOption struct {
	BannedFields map[string]struct{}
	AliasFields  map[string]string
	NestedOption map[string]*structOption
}

// parse can translate nested option represented by string such as `A.B` into NestedOption.
func (o *structOption) parse() {
	o.NestedOption = make(map[string]*structOption)

	for f := range o.BannedFields {
		fields := strings.Split(f, ".")
		if len(fields) > 1 {
			o.NestedOption[fields[0]] = &structOption{
				BannedFields: map[string]struct{}{
					strings.Join(fields[1:], "."): {},
				}}
			delete(o.BannedFields, f)
		}
	}

	for f, a := range o.AliasFields {
		fields := strings.Split(f, ".")
		if len(fields) > 1 {
			af := map[string]string{
				strings.Join(fields[1:], "."): a,
			}
			if nest, ok := o.NestedOption[fields[0]]; ok {
				nest.AliasFields = af
			} else {
				o.NestedOption[fields[0]] = &structOption{
					AliasFields: af,
				}
			}
			delete(o.AliasFields, f)
		}
	}

	for _, nest := range o.NestedOption {
		nest.parse()
	}
}
