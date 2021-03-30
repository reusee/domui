package domui

import "strings"

func Grid(specs ...GridSpec) Spec {
	return grid([]any{
		"display", "grid",
	}, specs...)
}

func InlineGrid(specs ...GridSpec) Spec {
	return grid([]any{
		"display", "inline-grid",
	}, specs...)
}

func grid(styles []any, specs ...GridSpec) Spec {

	var apply func(GridSpec)
	apply = func(spec GridSpec) {
		switch spec := spec.(type) {

		case GridSpecs:
			for _, s := range spec {
				apply(s)
			}

		case RowsSpec:
			styles = append(styles, "grid-template-rows", spec)

		case ColumnsSpec:
			styles = append(styles, "grid-template-columns", spec)

		case AreasSpec:
			var areas []string
			for _, area := range spec {
				areas = append(areas, `"`+area+`"`)
			}
			styles = append(styles,
				"grid-template-areas",
				strings.Join(areas, "\n"),
			)

		case JustifyItemsSpec:
			styles = append(styles, "justify-items", spec(nil))

		case AlignItemsSpec:
			styles = append(styles, "align-items", spec(nil))

		case JustifyContentSpec:
			styles = append(styles, "justify-content", spec(nil))

		case AlignContentSpec:
			styles = append(styles, "align-content", spec(nil))

		case ColumnGapSpec:
			styles = append(styles, "column-gap", spec)
			styles = append(styles, "grid-column-gap", spec)

		case RowGapSpec:
			styles = append(styles, "row-gap", spec)
			styles = append(styles, "grid-row-gap", spec)

		case AutoColumnsSpec:
			styles = append(styles, "grid-auto-columns", spec)

		case AutoRowsSpec:
			styles = append(styles, "grid-auto-rows", spec)

		case AutoFlowSpec:
			styles = append(styles, "grid-auto-flow", spec)

		}
	}

	for _, spec := range specs {
		apply(spec)
	}

	return Styles(styles...)
}

type GridSpec interface {
	IsGridSpec()
}

type GridSpecs []GridSpec

func (_ GridSpecs) IsGridSpec() {}

// col

type ColumnsSpec string

func (_ ColumnsSpec) IsGridSpec() {}

func Cols(spec string) ColumnsSpec {
	return ColumnsSpec(spec)
}

// row

type RowsSpec string

func (_ RowsSpec) IsGridSpec() {}

func Rows(spec string) RowsSpec {
	return RowsSpec(spec)
}

// areas

type AreasSpec []string

func (_ AreasSpec) IsGridSpec() {}

func Areas(areas ...string) AreasSpec {
	return AreasSpec(areas)
}

// place

type Place func(Place) string

func (_ Place) Start() string {
	return "start"
}

func (_ Place) End() string {
	return "end"
}

func (_ Place) Center() string {
	return "center"
}

func (_ Place) Stretch() string {
	return "stretch"
}

func (_ Place) SpaceAround() string {
	return "space-around"
}

func (_ Place) SpaceBetween() string {
	return "space-between"
}

func (_ Place) SpaceEvenly() string {
	return "space-evenly"
}

// justify items

type JustifyItemsSpec Place

func (_ JustifyItemsSpec) IsGridSpec() {}

func JustifyItems(place Place) JustifyItemsSpec {
	return JustifyItemsSpec(place)
}

// align items

type AlignItemsSpec Place

func (_ AlignItemsSpec) IsGridSpec() {}

func AlignItems(place Place) AlignItemsSpec {
	return AlignItemsSpec(place)
}

// justify content

type JustifyContentSpec Place

func (_ JustifyContentSpec) IsGridSpec() {}

func JustifyContent(place Place) JustifyContentSpec {
	return JustifyContentSpec(place)
}

// align content

type AlignContentSpec Place

func (_ AlignContentSpec) IsGridSpec() {}

func AlignContent(place Place) AlignContentSpec {
	return AlignContentSpec(place)
}

// column gap

type ColumnGapSpec string

func (_ ColumnGapSpec) IsGridSpec() {}

func ColGap(spec string) ColumnGapSpec {
	return ColumnGapSpec(spec)
}

// row gap

type RowGapSpec string

func (_ RowGapSpec) IsGridSpec() {}

func RowGap(spec string) RowGapSpec {
	return RowGapSpec(spec)
}

// auto columns

type AutoColumnsSpec string

func (_ AutoColumnsSpec) IsGridSpec() {}

func AutoCols(spec string) AutoColumnsSpec {
	return AutoColumnsSpec(spec)
}

// auto rows

type AutoRowsSpec string

func (_ AutoRowsSpec) IsGridSpec() {}

func AutoRows(spec string) AutoRowsSpec {
	return AutoRowsSpec(spec)
}

// auto flow

type AutoFlowSpec string

func (_ AutoFlowSpec) IsGridSpec() {}

func AutoFlow(spec string) AutoFlowSpec {
	return AutoFlowSpec(spec)
}

// utils

func AutoDivide(max string) string {
	return "repeat(auto-fill, minmax(" + max + ", 1fr))"
}

var StretchAll = GridSpecs{
	JustifyItems(Place.Stretch),
	AlignItems(Place.Stretch),
	JustifyContent(Place.Stretch),
	AlignContent(Place.Stretch),
}

// testing

type TestGrid Spec

func (_ Def) TestGrid() TestGrid {

	return Div(

		Grid(
			Cols("1fr 2fr 3fr"),
			Rows("1fr 2fr"),
			Areas(
				"a a a",
				"b c c",
			),
			JustifyItems(Place.Stretch),
			AlignItems(Place.Stretch),
			JustifyContent(Place.Stretch),
			AlignContent(Place.Stretch),
			StretchAll,
			RowGap("1rem"),
			ColGap("2rem"),
		),
		Styles(
			"width", "100%",
			"height", "100%",
		),

		Div(
			S("foo"),
			GridItem(
				Col("3"),
				Row("1"),
			),
			Styles(
				"border", "1px solid red",
			),
		),
		Div(
			S("bar"),
			GridItem(
				Col("1"),
				Row("1"),
			),
			Styles(
				"border", "1px solid blue",
			),
		),
		Div(
			S("baz"),
			GridItem(
				Area("c"),
			),
			Styles(
				"border", "1px solid green",
			),
		),
		Div(
			S("qux"),
			GridItem(
				Area("b"),
				JustifySelf(Place.Center),
				AlignSelf(Place.End),
			),
			Styles(
				"border", "1px solid yellow",
			),
		),
	)
}
