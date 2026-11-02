package main

func main() {
	app, err := Init()
	if err != nil {
		panic(err)
	}
	Entry(app)
}