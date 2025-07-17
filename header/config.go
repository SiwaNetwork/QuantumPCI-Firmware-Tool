package header

import (
	"errors"
	"os"
)

// Config содержит конфигурацию программы
type Config struct {
	Apply            bool       // Флаг создания выходного файла
	InputPath        string     // Путь к входному файлу прошивки
	OutputPath       string     // Путь к выходному файлу с заголовком
	VendorId         int        // PCI Vendor ID для заголовка
	DeviceId         int        // PCI Device ID для заголовка
	HardwareRevision int        // PCI Hardware Revision ID
	InputFile        *os.File   // Дескриптор входного файла
	OutputFile       *os.File   // Дескриптор выходного файла
}

// CheckConfig проверяет структуру конфигурации на корректность значений
// PCI Vendor ID, Device ID должны помещаться в uint16 и быть положительными
// PCI Hardware Revision должен быть любым значением uint16
func CheckConfig(c *Config) error {
	if c.InputPath == "" {
		return errors.New("Не указан входной файл")
	}

	if c.Apply && c.OutputPath == "" {
		return errors.New("Не указан выходной файл")
	}

	if c.Apply && (c.VendorId <= 0 || c.VendorId > 65535) {
		return errors.New("Пустой или некорректный PCI Vendor ID")
	}

	if c.Apply && (c.DeviceId <= 0 || c.DeviceId > 65535) {
		return errors.New("Пустой или некорректный PCI Device ID")
	}

	if c.Apply && (c.HardwareRevision < 0 || c.HardwareRevision > 65535) {
		return errors.New("Некорректный PCI Device Revision ID")
	}

	return nil
}

// OpenFiles открывает файлы, указанные в конфигурации. Возвращает ошибку при неудаче
func OpenFiles(c *Config) error {
	var err error
	c.InputFile, err = os.Open(c.InputPath)
	if err != nil {
		return err
	}

	if c.Apply {
		c.OutputFile, err = os.Create(c.OutputPath)
		return err
	}

	return nil
}

// CloseFiles закрывает ранее открытые файлы
func CloseFiles(c *Config) {
	c.InputFile.Close()
	if c.Apply {
		c.OutputFile.Close()
	}
}
