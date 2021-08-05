package remail

func buildRecepient(recepients []Recepient) (ret []string) {
	ret = make([]string, len(recepients))

	for row, cc := range recepients {
		if cc.Name != "" {
			ret[row] = cc.Name + "<" + cc.Address + ">"
		} else {
			ret[row] = cc.Address
		}
	}

	return
}

func mustBuildBody(body MessageBody) string {
	return string(body.Body)
}
