package main

func main() {
	router := buildRouter()
	_ = router.Run(":9876")
}
