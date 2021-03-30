package domui

func GridItem(specs ...GridItemSpec) Spec {
	var styles []any

	var apply func(spec GridItemSpec)
	apply = func(spec GridItemSpec) {
	}

	for _, spec := range specs {
		switch spec := spec.(type) {

		case GridItemSpecs:
			for _, s := range spec {
				apply(s)
			}

		case ColumnSpec:
			styles = append(styles,
				"grid-column", spec,
			)

		case RowSpec:
			styles = append(styles,
				"grid-row", spec,
			)

		case AreaSpec:
			styles = append(styles,
				"grid-area", spec,
			)

		case JustifySelfSpec:
			styles = append(styles,
				"justify-self", spec(nil),
			)

		case AlignSelfSpec:
			styles = append(styles,
				"align-self", spec(nil),
			)

		}
	}

	return Styles(styles...)
}

type GridItemSpec interface {
	IsGridItemSpec()
}

type GridItemSpecs []GridItemSpec

func (_ GridItemSpecs) IsGridItemSpec() {}

// column

type ColumnSpec string

func (_ ColumnSpec) IsGridItemSpec() {}

func Col(spec string) ColumnSpec {
	return ColumnSpec(spec)
}

// row

type RowSpec string

func (_ RowSpec) IsGridItemSpec() {}

func Row(spec string) RowSpec {
	return RowSpec(spec)
}

// area

type AreaSpec string

func (_ AreaSpec) IsGridItemSpec() {}

func Area(spec string) AreaSpec {
	return AreaSpec(spec)
}

// justify self

type JustifySelfSpec Place

func (_ JustifySelfSpec) IsGridItemSpec() {}

func JustifySelf(place Place) JustifySelfSpec {
	return JustifySelfSpec(place)
}

// align self

type AlignSelfSpec Place

func (_ AlignSelfSpec) IsGridItemSpec() {}

func AlignSelf(place Place) AlignSelfSpec {
	return AlignSelfSpec(place)
}
