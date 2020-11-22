package register

type Options struct {
	Address        []string
	NameSpace      string
	ServerInstance string
}

type Option func(*Options)
