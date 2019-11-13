package core

func Init() {
	InitializeConfig(DefaultConfigLocation)
	initializeHasher()
}
