package aozoraConvert

func latexCmd(cmd string) string {

	return `\` + cmd

}

func latexOpt(opt string) string {

	return `[` + opt + `]`

}

func latexArg(arg string) string {

	return `{` + arg + `}`

}
