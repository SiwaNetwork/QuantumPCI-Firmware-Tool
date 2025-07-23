package header

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	crc "github.com/sigurn/crc16"
)

// hdrMagic - массив из 4 константных байтов-идентификаторов
var hdrMagic = [4]byte{'S', 'H', 'I', 'W'}

// ErrNoHeader возвращается, когда в файле нет заголовка
var ErrNoHeader = errors.New("Заголовок не найден")

const hdrSize = 16

// Header представляет структуру заголовка прошивки
type Header struct {
	Magic            [4]byte // Магические байты для идентификации
	VendorId         uint16  // ID производителя PCI устройства
	DeviceId         uint16  // ID PCI устройства
	ImageSize        uint32  // Размер образа прошивки (без заголовка)
	HardwareRevision uint16  // Ревизия оборудования
	CRC              uint16  // Контрольная сумма CRC16
}

// firmwareImageSize вычисляет размер образа прошивки
func firmwareImageSize(c *Config) (uint32, error) {
	stat, err := c.InputFile.Stat()
	if err != nil {
		return 0, err
	}

	pos, err := c.InputFile.Seek(0, 1)
	if err != nil {
		return 0, err
	}

	return uint32(stat.Size() - pos), nil
}

// ReadHeader пытается прочитать заголовок из входного файла
func ReadHeader(c *Config) (*Header, error) {
	buf := make([]byte, hdrSize)
	n, err := c.InputFile.Read(buf)
	if err != nil {
		return nil, err
	}
	if n == hdrSize {
		hdr := &Header{}
		binary.Read(bytes.NewReader(buf), binary.BigEndian, hdr)
		if hdr.Magic == hdrMagic {
			return hdr, nil
		}
	}
	// Если заголовок не найден, возвращаем указатель файла в начало
	c.InputFile.Seek(0, 0)
	return nil, ErrNoHeader
}

// PrepareHeader создает структуру заголовка из значений конфигурации
func PrepareHeader(c *Config) (*Header, error) {
	imageSize, err := firmwareImageSize(c)
	if err != nil {
		return nil, err
	}

	hdr := &Header{
		Magic:            [4]byte{'S', 'H', 'I', 'W'},
		VendorId:         uint16(c.VendorId),
		DeviceId:         uint16(c.DeviceId),
		ImageSize:        imageSize,
		HardwareRevision: uint16(c.HardwareRevision),
	}

	return hdr, nil
}

// WriteHeader записывает структуру заголовка в начало выходного файла
func WriteHeader(c *Config, hdr *Header) error {
	if !c.Apply {
		return nil
	}
	h := new(bytes.Buffer)
	binary.Write(h, binary.BigEndian, hdr)
	_, err := c.OutputFile.WriteAt(h.Bytes(), 0)

	if err != nil {
		return err
	}

	// Перемещаем указатель в конец файла
	_, err = c.OutputFile.Seek(0, 2)
	return err
}

// CalcCRC вычисляет CRC16 входного файла и копирует данные в выходной файл при необходимости
func CalcCRC(c *Config) (uint16, error) {
	crcTable := crc.MakeTable(crc.CRC16_ARC)
	buf := make([]byte, 16384) // Буфер 16KB для эффективного чтения
	crc16 := uint16(0xFFFF)

	n, err := c.InputFile.Read(buf)
	for ; n > 0 && (err == nil || err == io.EOF); n, err = c.InputFile.Read(buf) {
		crc16 = crc.Update(crc16, buf[:n], crcTable)
		if c.Apply {
			_, err = c.OutputFile.Write(buf[:n])
			if err != nil {
				break
			}
		}
	}

	if err != nil && err != io.EOF {
		return 0, err
	}

	crc16 = crc.Complete(crc16, crcTable)

	return crc16, nil
}
