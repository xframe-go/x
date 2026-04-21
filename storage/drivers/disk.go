package drivers

type LocalDisk struct {
	Config LocalConfig
}

func (d *LocalDisk) DriverName() string {
	return "local"
}

type S3Disk struct {
	Config S3Config
}

func (d *S3Disk) DriverName() string {
	return "s3"
}

type R2Disk struct {
	Config R2Config
}

func (d *R2Disk) DriverName() string {
	return "r2"
}

type NilDisk struct{}

func (d *NilDisk) DriverName() string {
	return "nil"
}

func NewLocalDisk(config LocalConfig) *LocalDisk {
	return &LocalDisk{Config: config}
}

func NewS3Disk(config S3Config) *S3Disk {
	return &S3Disk{Config: config}
}

func NewR2Disk(config R2Config) *R2Disk {
	return &R2Disk{Config: config}
}

func NewNilDisk() *NilDisk {
	return &NilDisk{}
}
