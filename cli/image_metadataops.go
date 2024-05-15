package cli

// ImageMetadataOps - edit image metadata
type ImageMetadataOps struct {
	Properties ImageMetadataOptions
	File       string `name:"f" usage:"metadata.yaml file`
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
	f := m.GetFields()
	t.Properties.Override(f)
	t.Properties.SetImageDescription(f, t.Properties.Variant)
	m.SetFields(f)
	m.UpdateDates(t.Properties.Date, t.Properties.ExpiryDays)
	if t.File != "" {
		return m.WriteFile(t.File)
	} else {
		return m.Print()
	}
}
