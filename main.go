package main

import (
	"flag"
	"fmt"

	h "github.com/opencomputeproject/quantum-pci-ft/header"
	log "github.com/sirupsen/logrus"
)

// printHeader выводит информацию о заголовке прошивки
func printHeader(hdr *h.Header) {
	fmt.Println("Информация о заголовке:")
	fmt.Printf("PCI Vendor ID: 0x%04x\n", hdr.VendorId)
	fmt.Printf("PCI Device ID: 0x%04x\n", hdr.DeviceId)
	fmt.Printf("PCI HW Revision ID: 0x%04x\n", hdr.HardwareRevision)
	fmt.Printf("CRC16 образа: 0x%04x\n", hdr.CRC)
	fmt.Printf("Размер образа: %d байт\n", hdr.ImageSize)
}

func main() {
	fmt.Println("Quantum PCI FT - Инструмент для работы с прошивками PCI устройств")
	fmt.Println("===========================================================")
	
	c := &h.Config{}

	// Определение флагов командной строки
	flag.BoolVar(&c.Apply, "apply", false, "Создать новый файл прошивки с заголовком в начале")
	flag.StringVar(&c.InputPath, "input", "", "Путь к исходному файлу прошивки")
	flag.StringVar(&c.OutputPath, "output", "", "Путь к файлу прошивки с заголовком (будет перезаписан)")
	flag.IntVar(&c.VendorId, "vendor", 0, "PCI VEN_ID для добавления в заголовок")
	flag.IntVar(&c.DeviceId, "device", 0, "PCI DEV_ID для добавления в заголовок")
	flag.IntVar(&c.HardwareRevision, "hw", 0, "PCI REV_ID (ревизия оборудования) для добавления в заголовок")
	flag.Parse()

	// Проверка конфигурации
	if err := h.CheckConfig(c); err != nil {
		log.Fatal("Ошибка конфигурации: ", err)
	}

	// Открытие файлов
	if err := h.OpenFiles(c); err != nil {
		log.Fatal("Ошибка открытия файлов: ", err)
	}
	defer h.CloseFiles(c)

	// Проверка существующего заголовка
	oldHdr, err := h.ReadHeader(c)
	if err == nil {
		fmt.Println("\nВходной файл уже содержит заголовок:")
		printHeader(oldHdr)
		if c.Apply {
			fmt.Println("\n⚠️  Заголовок образа будет перезаписан новыми значениями")
		}
	}

	// Подготовка нового заголовка
	hdr, err := h.PrepareHeader(c)
	if err != nil {
		log.Fatal("Ошибка подготовки заголовка: ", err)
	}

	// Запись заголовка
	if err := h.WriteHeader(c, hdr); err != nil {
		log.Fatal("Ошибка записи заголовка: ", err)
	}

	// Вычисление CRC и копирование данных
	hdr.CRC, err = h.CalcCRC(c)
	if err != nil {
		log.Fatal("Ошибка вычисления CRC: ", err)
	}

	// Перезапись заголовка с вычисленным CRC
	if err := h.WriteHeader(c, hdr); err != nil {
		log.Fatal("Ошибка записи заголовка с CRC: ", err)
	}

	fmt.Println("\nНовый заголовок:")
	printHeader(hdr)
	
	if c.Apply {
		fmt.Printf("\n✅ Файл прошивки с заголовком успешно создан: %s\n", c.OutputPath)
	} else {
		fmt.Println("\n📋 Анализ завершен (используйте флаг -apply для создания файла)")
	}
}
