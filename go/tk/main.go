package main

import (
	_ "embed"

	tk "modernc.org/tk9.0"
	_ "modernc.org/tk9.0/themes/azure"
)

func main() {
	tk.ActivateTheme("azure light")
	tk.Pack(
		tk.TExit(),
		tk.Padx("1m"), tk.Pady("2m"), tk.Ipadx("1m"), tk.Ipady("1m"))
	app := tk.App.Center()
	app.WmTitle("demo")
	app.Wait()
}
