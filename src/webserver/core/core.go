package core

func Init() {
	InitializeConfig(DefaultConfigLocation)
	InitializeHasher()
}
