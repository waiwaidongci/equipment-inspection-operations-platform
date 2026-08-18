package application

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/example/inspection-platform/internal/domain"
)

type Service struct {
	devices domain.DeviceRepository
	records domain.RecordRepository
	repairs domain.RepairRepository
}

func NewService(devices domain.DeviceRepository, records domain.RecordRepository, repairs domain.RepairRepository) *Service {
	return &Service{devices: devices, records: records, repairs: repairs}
}

type DeviceInput struct {
	Code         string `json:"code"`
	Name         string `json:"name"`
	Category     string `json:"category"`
	Model        string `json:"model"`
	SerialNumber string `json:"serialNumber"`
	Manufacturer string `json:"manufacturer"`
	Location     string `json:"location"`
	InstallDate  string `json:"installDate"`
	WarrantyEnd  string `json:"warrantyEnd"`
	Owner        string `json:"owner"`
	Status       string `json:"status"`
	Remark       string `json:"remark"`
}

func (s *Service) Create(ctx context.Context, input DeviceInput) (*domain.Device, error) {
	if err := validateDeviceInput(&input); err != nil {
		return nil, err
	}
	_, err := s.devices.FindByCode(ctx, input.Code)
	if err == nil {
		return nil, fmt.Errorf("%w: device code %q already exists", domain.ErrConflict, input.Code)
	}
	if !errors.Is(err, domain.ErrNotFound) {
		return nil, err
	}
	now := time.Now().UTC()
	device := &domain.Device{
		Code: input.Code, Name: input.Name, Category: input.Category, Model: input.Model,
		SerialNumber: input.SerialNumber, Manufacturer: input.Manufacturer, Location: input.Location,
		InstallDate: input.InstallDate, WarrantyEnd: input.WarrantyEnd, Owner: input.Owner,
		Status: input.Status, Remark: input.Remark, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.devices.Create(ctx, device); err != nil {
		return nil, err
	}
	return device, nil
}

func (s *Service) Update(ctx context.Context, id int64, input DeviceInput) (*domain.Device, error) {
	device, err := s.devices.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := validateDeviceInput(&input); err != nil {
		return nil, err
	}
	if device.Code != input.Code {
		_, err := s.devices.FindByCode(ctx, input.Code)
		if err == nil {
			return nil, fmt.Errorf("%w: device code %q already exists", domain.ErrConflict, input.Code)
		}
		if !errors.Is(err, domain.ErrNotFound) {
			return nil, err
		}
	}
	device.Code, device.Name, device.Category, device.Model = input.Code, input.Name, input.Category, input.Model
	device.SerialNumber, device.Manufacturer, device.Location = input.SerialNumber, input.Manufacturer, input.Location
	device.InstallDate, device.WarrantyEnd, device.Owner = input.InstallDate, input.WarrantyEnd, input.Owner
	device.Status, device.Remark = input.Status, input.Remark
	device.UpdatedAt = time.Now().UTC()
	if err := s.devices.Update(ctx, device); err != nil {
		return nil, err
	}
	return device, nil
}

func (s *Service) Get(ctx context.Context, id int64) (*domain.Device, error) {
	return s.devices.FindByID(ctx, id)
}

func (s *Service) List(ctx context.Context, filter domain.DeviceFilter) (domain.Page, error) {
	items, total, err := s.devices.List(ctx, filter)
	if err != nil {
		return domain.Page{}, err
	}
	return newPage(items, total, filter.Page, filter.PageSize), nil
}

func (s *Service) History(ctx context.Context, id int64) (domain.DeviceHistory, error) {
	if _, err := s.devices.FindByID(ctx, id); err != nil {
		return domain.DeviceHistory{}, err
	}
	records, err := s.records.ListByDevice(ctx, id, 100, 0)
	if err != nil {
		return domain.DeviceHistory{}, err
	}
	repairs, err := s.repairs.ListByDevice(ctx, id)
	if err != nil {
		return domain.DeviceHistory{}, err
	}
	return domain.DeviceHistory{Records: records, Repairs: repairs}, nil
}

func validateDeviceInput(input *DeviceInput) error {
	input.Code = strings.TrimSpace(input.Code)
	input.Name = strings.TrimSpace(input.Name)
	input.Category = strings.TrimSpace(input.Category)
	input.Location = strings.TrimSpace(input.Location)
	if input.Code == "" || input.Name == "" || input.Category == "" || input.Location == "" {
		return fmt.Errorf("%w: code, name, category and location are required", domain.ErrInvalid)
	}
	if input.Status == "" {
		input.Status = "active"
	}
	if input.Status != "active" && input.Status != "inactive" {
		return fmt.Errorf("%w: status must be active or inactive", domain.ErrInvalid)
	}
	return nil
}

func newPage(items any, total int64, page, pageSize int) domain.Page {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	totalPages := (total + int64(pageSize) - 1) / int64(pageSize)
	return domain.Page{Items: items, Total: total, Page: page, PageSize: pageSize, TotalPages: totalPages}
}
