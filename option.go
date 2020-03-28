package coven

import "strings"

type StructOption struct {
	BannedFields []string
	AliasFields  map[string]string
}

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

type structOption struct {
	BannedFields map[string]struct{}
	AliasFields  map[string]string
	NestedOption map[string]*structOption
}

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
