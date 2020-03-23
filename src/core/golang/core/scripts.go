package core

type ScriptRunSettings struct {
	// VM settings
	CpuAllocation      float64
	MemBytesAllocation int64
	DiskSizeBytes      int64

	// Environment settings
	KotlinContainerVersion string
	GrchiveCoreVersion     string

	// Compile settings
	CompileOnly    bool
	ScriptChecksum string
}
