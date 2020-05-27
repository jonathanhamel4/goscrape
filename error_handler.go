package main

func VerifyError(err error) {
	if err != nil {
		panic(err)
	}
}