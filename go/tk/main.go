package main

import (
	tk "modernc.org/tk9.0"
)

func main() {
	tk.Pack(
		tk.TExit(),
		tk.Padx("1m"), tk.Pady("2m"), tk.Ipadx("1m"), tk.Ipady("1m"))
	app := tk.App.Center()
	app.WmTitle("demo")
	app.Wait()
}
