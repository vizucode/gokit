package convert

func ToString(i interface{}) string {
	return string(interfaceToBytes(i))
}
