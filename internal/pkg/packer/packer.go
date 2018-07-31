package packer

func Pack(inputFolder, outputFile string) error {
	_, err := openOutput(outputFile)
	defer closeOutput()
	if err != nil {
		return err
	}
	_, _, err = readInputFolder(inputFolder)

	return err
}
