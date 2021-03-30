package domui

import (
	"testing"
)

func TestPatch(t *testing.T) {

	t.Run("different kind", func(t *testing.T) {
		WithTestApp(
			t,
			func(app *App) {
				html := app.HTML()
				if html != "<div></div>" {
					t.Fatal()
				}
				app.Update(func() string {
					return "bar"
				})
				app.Render()
				html = app.HTML()
				if html != "foo" {
					t.Fatalf("got %s", html)
				}
			},
			func() string {
				return "foo"
			},
			func(s string) Root {
				if s == "foo" {
					return Div()
				}
				return Text("foo")
			},
		)
	})

	t.Run("different tag", func(t *testing.T) {
		WithTestApp(
			t,
			func(app *App) {
				html := app.HTML()
				if html != "<div></div>" {
					t.Fatal()
				}
				app.Update(func() string {
					return "bar"
				})
				app.Render()
				html = app.HTML()
				if html != "<p></p>" {
					t.Fatalf("got %s", html)
				}
			},
			func() string {
				return "foo"
			},
			func(s string) Root {
				if s == "foo" {
					return Div()
				}
				return P()
			},
		)
	})

	t.Run("patch", func(t *testing.T) {
		WithTestApp(
			t,
			func(app *App) {
				html := app.HTML()
				if html != `<div id="foo" class="foo" foo="foo" style="font-size: 42px; font-weight: 42;"></div>` {
					t.Fatalf("got %s", html)
				}

				app.Update(func() (string, int) {
					return "bar", 2
				})
				app.Render()
				html = app.HTML()
				if html != `<div id="bar" class="bar qux" style="font-size: 2px; display: block; font-weight: 2;" attr2="bar" bar="bar"></div>` {
					t.Fatalf("got %s", html)
				}

				app.Update(func() int {
					return 3
				})
				app.Render()
				html = app.HTML()
				if html != `<div></div>` {
					t.Fatalf("got %s\n", html)
				}
			},
			func() (string, int) {
				return "foo", 42
			},
			func(s string, i int) Root {
				if i == 3 {
					return Div()
				}

				return Div(
					ID(s),
					Style(sp(`font-size: %dpx`, i)),
					Styles(
						"font-weight", i,
					),
					Class(s),
					Attr(s, s),

					If(
						i == 2,
						Styles("display", "block"),
						Class("qux"),
						Attr("attr2", s),
					),
				)
			},
		)
	})

	t.Run("test patch event", func(t *testing.T) {
		m := make(map[int]int)
		WithTestApp(
			t,
			func(app *App) {
				app.element.Call("click")
				if m[1] != 1 {
					t.Fatal()
				}

				app.Update(func() int {
					return 2
				})
				app.Render()
				app.element.Call("click")
				if m[1] != 1 {
					t.Fatal()
				}
				if m[2] != 1 {
					t.Fatal()
				}

				app.Update(func() int {
					return 3
				})
				app.Render()
				app.element.Call("click")
				if m[1] != 1 {
					t.Fatal()
				}
				if m[2] != 1 {
					t.Fatal()
				}
				if m[3] != 1 {
					t.Fatal()
				}

				app.Update(func() int {
					return 4
				})
				app.Render()
				app.element.Call("click")
				if m[1] != 1 {
					t.Fatal()
				}
				if m[2] != 1 {
					t.Fatal()
				}
				if m[3] != 1 {
					t.Fatal()
				}
				if m[4] != 0 {
					t.Fatal()
				}
			},
			func(i int) Root {
				if i == 4 {
					return Div()
				}
				return Div(
					On("click", func() {
						m[i]++
					}),
				)
			},
			func() int {
				return 1
			},
		)
	})

	t.Run("patch children", func(t *testing.T) {
		WithTestApp(
			t,
			func(app *App) {
				html := app.HTML()
				if html != `<div><p>0</p></div>` {
					t.Fatalf("got %s", html)
				}

				app.Update(func() int {
					return 2
				})
				app.Render()
				html = app.HTML()
				if html != `<div><p>0</p><p>1</p></div>` {
					t.Fatalf("got %s", html)
				}

				app.Update(func() int {
					return 1
				})
				app.Render()
				html = app.HTML()
				if html != `<div><p>0</p></div>` {
					t.Fatalf("got %s", html)
				}

				app.Update(func() int {
					return 3
				})
				app.Render()
				html = app.HTML()
				if html != `<div><p>0</p><p>1</p><p>2</p></div>` {
					t.Fatalf("got %s", html)
				}

				app.Update(func() int {
					return 0
				})
				app.Render()
				html = app.HTML()
				if html != `<div></div>` {
					t.Fatalf("got %s", html)
				}

				app.Update(func() int {
					return 3
				})
				app.Render()
				html = app.HTML()
				if html != `<div><p>0</p><p>1</p><p>2</p></div>` {
					t.Fatalf("got %s", html)
				}

			},
			func(n int) Root {
				var children Specs
				for i := 0; i < n; i++ {
					children = append(children, P(Text(sp("%d", i))))
				}
				return Div(
					children...,
				)
			},
			func() int {
				return 1
			},
		)
	})

}
