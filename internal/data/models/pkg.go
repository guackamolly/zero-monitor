package models

type Package struct {
	Name        string
	Description string
	Version     string
}

func NewPackage(
	name string,
	description string,
	version string,
) Package {
	return Package{
		Name:        name,
		Description: description,
		Version:     version,
	}
}
