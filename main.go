package main

func main() {
	a := App{}
	a.Initialize(getEnv())
	a.Run(":8080")
}
