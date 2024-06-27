package models

func CreateFile(file File) (File, error) {
	if err := DB.Create(&file).Error; err != nil {
		return File{}, err
	}

	return file, nil
}

func UpdateFile(id string, file File) (File, error) {
	var updatedFile File
	if err := DB.Where("id = ?", id).First(&updatedFile).Error; err != nil {
		return File{}, err
	}

	if err := DB.Model(&updatedFile).Updates(file).Error; err != nil {
		return File{}, err
	}

	return updatedFile, nil
}

func GetFile(id string) (File, error) {
	var file File
	if err := DB.Where("id = ?", id).First(&file).Error; err != nil {
		return File{}, err
	}

	return file, nil
}

func DeleteFile(id string) error {
	var file File
	if err := DB.Where("id = ?", id).Delete(&file).Error; err != nil {
		return err
	}

	return nil
}
