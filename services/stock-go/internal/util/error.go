package util

func ok(err error) {
	if err != nil {
		panic(err)
	}
}
