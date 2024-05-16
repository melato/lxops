package cli

// ImageMetadataOps - edit image metadata
type ImageMetadataOps struct {
	Properties ImageMetadataOptions
	File       string `name:"f" usage:"metadata.yaml file"`
	OutputFile string `name:"o" usage:"output file, if other than input file"`
}

func (t *ImageMetadataOps) Init() error {
	return t.Properties.Init()
}

func (t *ImageMetadataOps) Configured() error {
	return t.Properties.Configured()
}

func (t *ImageMetadataOps) Update() error {
	var m *ImageMetadata
	var err error
	if t.File != "" {
		m, err = ReadImageMetadata(t.File)
	} else {
		m, err = NewImageMetadata()
	}
	if err != nil {
		return err
	}
	t.Properties.Apply(m)
	outputFile := t.OutputFile
	if outputFile == "" {
		outputFile = t.File
	}
	if outputFile == "" || outputFile == "." {
		return m.Print()
	} else {
		return m.WriteFile(t.File)
	}
}
