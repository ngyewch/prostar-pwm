package prostar_pwm

import (
	"encoding/binary"
	"errors"
	"sync"

	"github.com/simonvetter/modbus"
	"github.com/x448/float16"
)

type Dev struct {
	mc               *modbus.ModbusClient
	unitId           uint8
	mutex            *sync.Mutex
	inputRegisters   *Registers
	holdingRegisters *Registers
}

func New(mc *modbus.ModbusClient, unitId uint8, mutex *sync.Mutex) *Dev {
	return &Dev{
		mc:               mc,
		unitId:           unitId,
		mutex:            mutex,
		inputRegisters:   NewRegisters(mc, modbus.INPUT_REGISTER),
		holdingRegisters: NewRegisters(mc, modbus.HOLDING_REGISTER),
	}
}

func (dev *Dev) requestSetup() error {
	err := dev.mc.SetUnitId(dev.unitId)
	if err != nil {
		return err
	}
	err = dev.mc.SetEncoding(modbus.BIG_ENDIAN, modbus.LOW_WORD_FIRST)
	if err != nil {
		return err
	}
	return nil
}

func (dev *Dev) ReadRawADCData() (RawADCData, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return RawADCData{}, err
	}

	var r RawADCData

	r.SupplyVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0004)
	if err != nil {
		return RawADCData{}, err
	}
	r.GateDriveVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0005)
	if err != nil {
		return RawADCData{}, err
	}
	r.MeterBusSupplyVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0006)
	if err != nil {
		return RawADCData{}, err
	}
	r.InternalReferenceVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0007)
	if err != nil {
		return RawADCData{}, err
	}
	r.NegativeSupplyRailForCurrentMeasurement, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0008)
	if err != nil {
		return RawADCData{}, err
	}
	r.LoadFETGateVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0009)
	if err != nil {
		return RawADCData{}, err
	}
	r.ArrayFETGateVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x000a)
	if err != nil {
		return RawADCData{}, err
	}

	return r, nil
}

func (dev *Dev) ReadFilteredADCData() (FilteredADCData, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return FilteredADCData{}, err
	}

	var r FilteredADCData

	r.ArrayCurrent, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0011)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.BatteryTerminalVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0012)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.ArrayVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0013)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.LoadVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0014)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.LoadCurrent, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0016)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.BatterySenseVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0017)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.BatteryVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0018)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.BatteryCurrent, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0019)
	if err != nil {
		return FilteredADCData{}, err
	}

	return r, nil
}

func (dev *Dev) ReadTemperatureData() (TemperatureData, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return TemperatureData{}, err
	}

	var r TemperatureData

	r.Heatsink, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x001a)
	if err != nil {
		return TemperatureData{}, err
	}
	r.Battery, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x001b)
	if err != nil {
		return TemperatureData{}, err
	}
	r.Ambient, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x001c)
	if err != nil {
		return TemperatureData{}, err
	}
	r.Remote, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x001d)
	if err != nil {
		return TemperatureData{}, err
	}

	return r, nil
}

func (dev *Dev) ReadChargerStatus() (ChargerStatus, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return ChargerStatus{}, err
	}

	var r ChargerStatus

	{
		v, err := dev.inputRegisters.ReadUint16Ptr(0x0021)
		if err != nil {
			return ChargerStatus{}, err
		}
		if v != nil {
			v2 := ChargeState(*v)
			r.ChargeState = &v2
		}
	}
	{
		v, err := dev.inputRegisters.ReadUint16Ptr(0x0022)
		if err != nil {
			return ChargerStatus{}, err
		}
		if v != nil {
			v2 := ArrayFault(*v)
			details := v2.Details()
			r.ArrayFault = &details
		}
	}
	r.BatteryVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0023)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.BatteryRegulatorReferenceVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0024)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.AhChargeResettable, err = dev.inputRegisters.ReadUint32AsFloat32Ptr(0x0026, WordOrderingHighFirst, 10)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.AhChargeTotal, err = dev.inputRegisters.ReadUint32AsFloat32Ptr(0x0028, WordOrderingHighFirst, 10)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.KWhChargeResettable, err = dev.inputRegisters.ReadUint16AsFloat32Ptr(0x002a, 10)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.KWhChargeTotal, err = dev.inputRegisters.ReadUint16AsFloat32Ptr(0x002b, 10)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.BatteryTemperatureFoldback100PercentOutputLimit, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x002c)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.BatteryTemperatureFoldback0PercentOutputLimit, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x002d)
	if err != nil {
		return ChargerStatus{}, err
	}

	return r, nil
}

func (dev *Dev) ReadLoadStatus() (LoadStatus, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return LoadStatus{}, err
	}

	var r LoadStatus

	{
		v, err := dev.inputRegisters.ReadUint16Ptr(0x002e)
		if err != nil {
			return LoadStatus{}, err
		}
		if v != nil {
			v2 := LoadState(*v)
			r.LoadState = &v2
		}
	}
	{
		v, err := dev.inputRegisters.ReadUint16Ptr(0x002f)
		if err != nil {
			return LoadStatus{}, err
		}
		if v != nil {
			v2 := LoadFault(*v)
			details := v2.Details()
			r.LoadFault = &details
		}
	}
	r.LoadCurrentCompensatedLVDVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0030)
	if err != nil {
		return LoadStatus{}, err
	}
	r.LoadHVDVoltage, err = dev.inputRegisters.ReadFloat16AsFloat32Ptr(0x0031)
	if err != nil {
		return LoadStatus{}, err
	}
	r.AhLoadResettable, err = dev.inputRegisters.ReadUint32AsFloat32Ptr(0x0032, WordOrderingHighFirst, 10)
	if err != nil {
		return LoadStatus{}, err
	}
	r.AhLoadTotal, err = dev.inputRegisters.ReadUint32AsFloat32Ptr(0x0034, WordOrderingHighFirst, 10)
	if err != nil {
		return LoadStatus{}, err
	}

	return r, nil
}

func (dev *Dev) ReadMiscData() (MiscData, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return MiscData{}, err
	}

	var r MiscData

	r.Hourmeter, err = dev.inputRegisters.ReadUint32Ptr(0x0036, WordOrderingHighFirst)
	if err != nil {
		return MiscData{}, err
	}
	{
		v, err := dev.inputRegisters.ReadUint32Ptr(0x0038, WordOrderingHighFirst)
		if err != nil {
			return MiscData{}, err
		}
		if v != nil {
			v2 := Alarm(*v)
			details := v2.Details()
			r.Alarm = &details
		}
	}
	r.DIPSwitch, err = dev.inputRegisters.ReadUint16Ptr(0x003a)
	{
		v, err := dev.inputRegisters.ReadUint16Ptr(0x003b)
		if err != nil {
			return MiscData{}, err
		}
		if v != nil {
			v2 := LEDState(*v)
			r.LEDState = &v2
		}
	}
	{
		v, err := dev.inputRegisters.ReadUint16Ptr(0x004d)
		if err != nil {
			return MiscData{}, err
		}
		if v != nil {
			v2 := ChargeStatusLEDState(*v)
			r.ChargeStatusLEDState = &v2
		}
	}
	r.LightingShouldBeOn, err = dev.inputRegisters.ReadUint16Ptr(0x004e)
	if err != nil {
		return MiscData{}, err
	}

	return r, nil
}

func (dev *Dev) ReadChargeSettings() (ChargeSettings, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return ChargeSettings{}, err
	}

	var r ChargeSettings

	r.RegulationVoltageAt25C, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe000)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.FloatVoltageAt25C, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe001)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.TimeBeforeEnteringFloat, err = dev.holdingRegisters.ReadUint16Ptr(0xe002)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.TimeBeforeEnteringFloatDueToLowBattery, err = dev.holdingRegisters.ReadUint16Ptr(0xe003)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.VoltageTriggerForLowBatteryFloatTime, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe004)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.VoltageToCancelFloat, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe005)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.ExitFloatTime, err = dev.holdingRegisters.ReadUint16Ptr(0xe006)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.EqualizeVoltageAt25C, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe007)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.DaysBetweenEQCycles, err = dev.holdingRegisters.ReadUint16Ptr(0xe008)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.EqualizeTimeLimitAboveEVReg, err = dev.holdingRegisters.ReadUint16Ptr(0xe009)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.EqualizeTimeLimitAtEVEq, err = dev.holdingRegisters.ReadUint16Ptr(0xe00a)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.ReferenceChargeVoltageLimit, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe010)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.TemperatureCompensationCoefficient, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe01a)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.HighVoltageDisconnectAt25C, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe01b)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.HighVoltageReconnect, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe01c)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.MaximumChargeVoltageReference, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe01d)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.MaxBatteryTempCompensationLimit, err = dev.holdingRegisters.ReadUint16AsInt16Ptr(0xe01e)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.MinBatteryTempCompensationLimit, err = dev.holdingRegisters.ReadUint16AsInt16Ptr(0xe01f)
	if err != nil {
		return ChargeSettings{}, err
	}

	return r, nil
}

func (dev *Dev) ReadLoadSettings() (LoadSettings, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return LoadSettings{}, err
	}

	var r LoadSettings

	r.LowVoltageDisconnect, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe022)
	if err != nil {
		return LoadSettings{}, err
	}
	r.LowVoltageReconnect, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe023)
	if err != nil {
		return LoadSettings{}, err
	}
	r.LoadHighVoltageDisconnect, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe024)
	if err != nil {
		return LoadSettings{}, err
	}
	r.LoadHighVoltageReconnect, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe025)
	if err != nil {
		return LoadSettings{}, err
	}
	r.LVDLoadCurrentCompensation, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe026)
	if err != nil {
		return LoadSettings{}, err
	}
	r.LVDWarningTimeout, err = dev.holdingRegisters.ReadUint16Ptr(0xe027)
	if err != nil {
		return LoadSettings{}, err
	}

	return r, nil
}

func (dev *Dev) ReadMiscSettings() (MiscSettings, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return MiscSettings{}, err
	}

	var r MiscSettings

	r.LEDGreenToGreenAndYellowLimit, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe030)
	if err != nil {
		return MiscSettings{}, err
	}
	r.LEDGreenAndYellowToYellowLimit, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe031)
	if err != nil {
		return MiscSettings{}, err
	}
	r.LEDYellowToYellowAndRedLimit, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe032)
	if err != nil {
		return MiscSettings{}, err
	}
	r.LEDYellowAndRedToRedFlashingLimit, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe033)
	if err != nil {
		return MiscSettings{}, err
	}
	r.ModbusID, err = dev.holdingRegisters.ReadUint16Ptr(0xe034)
	if err != nil {
		return MiscSettings{}, err
	}
	r.MeterbusID, err = dev.holdingRegisters.ReadUint16Ptr(0xe035)
	if err != nil {
		return MiscSettings{}, err
	}

	return r, nil
}

func (dev *Dev) ReadPWMSettings() (PWMSettings, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return PWMSettings{}, err
	}

	var r PWMSettings

	r.ChargeCurrentLimit, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe038)
	if err != nil {
		return PWMSettings{}, err
	}

	return r, nil
}

func (dev *Dev) ReadStatistics() (Statistics, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return Statistics{}, err
	}

	var r Statistics
	r.Hourmeter, err = dev.holdingRegisters.ReadUint32Ptr(0xe040, WordOrderingLowFirst)
	if err != nil {
		return Statistics{}, err
	}
	r.AhLoadResettable, err = dev.holdingRegisters.ReadUint32AsFloat32Ptr(0xe042, WordOrderingLowFirst, 10)
	if err != nil {
		return Statistics{}, err
	}
	r.AhLoadTotal, err = dev.holdingRegisters.ReadUint32AsFloat32Ptr(0xe044, WordOrderingLowFirst, 10)
	if err != nil {
		return Statistics{}, err
	}
	r.AhChargeResettable, err = dev.holdingRegisters.ReadUint32AsFloat32Ptr(0xe046, WordOrderingLowFirst, 10)
	if err != nil {
		return Statistics{}, err
	}
	r.AhChargeTotal, err = dev.holdingRegisters.ReadUint32AsFloat32Ptr(0xe048, WordOrderingLowFirst, 10)
	if err != nil {
		return Statistics{}, err
	}
	r.KWhcResettable, err = dev.holdingRegisters.ReadUint16AsFloat32Ptr(0xe04a, 10)
	if err != nil {
		return Statistics{}, err
	}
	r.KWhcTotal, err = dev.holdingRegisters.ReadUint16AsFloat32Ptr(0xe04b, 10)
	if err != nil {
		return Statistics{}, err
	}
	r.BatteryVoltageMinimum, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe04c)
	if err != nil {
		return Statistics{}, err
	}
	r.BatteryVoltageMaximum, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe04d)
	if err != nil {
		return Statistics{}, err
	}
	r.ArrayVoltageMaximum, err = dev.holdingRegisters.ReadFloat16AsFloat32Ptr(0xe04e)
	if err != nil {
		return Statistics{}, err
	}
	r.TimeSinceLastEqualize, err = dev.holdingRegisters.ReadUint16Ptr(0xe04f)
	if err != nil {
		return Statistics{}, err
	}

	return r, nil
}

func (dev *Dev) ReadLoggedData() ([]LoggedDataRecord, error) {
	dev.mutex.Lock()
	defer dev.mutex.Unlock()

	err := dev.requestSetup()
	if err != nil {
		return nil, err
	}

	var records []LoggedDataRecord

	for i := 0; i < 256; i++ {
		v, err := dev.mc.ReadRegisters(0x8000+uint16(i*16), 16, modbus.INPUT_REGISTER)
		if err != nil {
			if errors.Is(err, modbus.ErrIllegalDataAddress) {
				return nil, nil
			} else {
				return nil, err
			}
		}
		hourmeter := WordOrderingLowFirst.Uint32(v[0:2])
		if (hourmeter != 0x00000000) && (hourmeter != 0xffffffff) {
			records = append(records, LoggedDataRecord{
				Hourmeter:                  hourmeter,
				AlarmDaily:                 Alarm(WordOrderingLowFirst.Uint32(v[2:4])).Details(),
				LoadFaultDaily:             LoadFault(WordOrderingLowFirst.Uint32(v[4:6])).Details(),
				ArrayFaultDaily:            ArrayFault(WordOrderingLowFirst.Uint32(v[6:8])).Details(),
				BatteryVoltageMinimumDaily: float16.Frombits(v[8]).Float32(),
				BatteryVoltageMaximumDaily: float16.Frombits(v[9]).Float32(),
				AhChargeDaily:              float16.Frombits(v[10]).Float32(),
				AhLoadDaily:                float16.Frombits(v[11]).Float32(),
				ArrayVoltageMaximumDaily:   float16.Frombits(v[12]).Float32(),
				TimeInAbsorptionDaily:      v[13],
				TimeInEqualizeDaily:        v[14],
				TimeInFloatDaily:           v[15],
			})
		}
	}

	return records, nil
}

type Registers struct {
	mc      *modbus.ModbusClient
	regType modbus.RegType
}

func NewRegisters(mc *modbus.ModbusClient, regType modbus.RegType) *Registers {
	return &Registers{
		mc:      mc,
		regType: regType,
	}
}

func (r *Registers) ReadUint16(addr uint16) (uint16, error) {
	return r.mc.ReadRegister(addr, r.regType)
}

func (r *Registers) ReadUint16Ptr(addr uint16) (*uint16, error) {
	v, err := r.ReadUint16(addr)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &v, nil
}

func (r *Registers) ReadFloat16(addr uint16) (float16.Float16, error) {
	v, err := r.mc.ReadRegister(addr, r.regType)
	if err != nil {
		return 0, err
	}
	return float16.Frombits(v), nil
}

func (r *Registers) ReadFloat16AsFloat32(addr uint16) (float32, error) {
	v, err := r.ReadFloat16(addr)
	if err != nil {
		return 0, err
	}
	return v.Float32(), nil
}

func (r *Registers) ReadFloat16AsFloat32Ptr(addr uint16) (*float32, error) {
	v, err := r.ReadFloat16AsFloat32(addr)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &v, nil
}

func (r *Registers) ReadUint16AsFloat32(addr uint16, divisor float32) (float32, error) {
	v, err := r.ReadUint16(addr)
	if err != nil {
		return 0, err
	}
	return float32(v) / divisor, nil
}

func (r *Registers) ReadUint16AsFloat32Ptr(addr uint16, divisor float32) (*float32, error) {
	v, err := r.ReadUint16AsFloat32(addr, divisor)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &v, nil
}

func (r *Registers) ReadUint16AsInt16(addr uint16) (int16, error) {
	v, err := r.ReadUint16(addr)
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

func (r *Registers) ReadUint16AsInt16Ptr(addr uint16) (*int16, error) {
	v, err := r.ReadUint16AsInt16(addr)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &v, nil
}

func (r *Registers) ReadUint32(addr uint16, wordOrdering WordOrdering) (uint32, error) {
	b, err := r.mc.ReadRegisters(addr, 2, r.regType)
	if err != nil {
		return 0, err
	}
	return wordOrdering.Uint32(b), nil
}

func (r *Registers) ReadUint32Ptr(addr uint16, wordOrdering WordOrdering) (*uint32, error) {
	v, err := r.ReadUint32(addr, wordOrdering)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &v, nil
}

func (r *Registers) ReadUint32AsFloat32(addr uint16, wordOrdering WordOrdering, divisor float32) (float32, error) {
	v, err := r.ReadUint32(addr, wordOrdering)
	if err != nil {
		return 0, err
	}
	return float32(v) / divisor, nil
}

func (r *Registers) ReadUint32AsFloat32Ptr(addr uint16, wordOrdering WordOrdering, divisor float32) (*float32, error) {
	v, err := r.ReadUint32AsFloat32(addr, wordOrdering, divisor)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &v, nil
}

type WordOrdering int

const (
	WordOrderingHighFirst WordOrdering = iota
	WordOrderingLowFirst
)

func (wordOrdering WordOrdering) Uint32(v []uint16) uint32 {
	var b []byte
	switch wordOrdering {
	case WordOrderingHighFirst:
		b = binary.BigEndian.AppendUint16(b, v[0])
		b = binary.BigEndian.AppendUint16(b, v[1])
	case WordOrderingLowFirst:
		b = binary.BigEndian.AppendUint16(b, v[1])
		b = binary.BigEndian.AppendUint16(b, v[0])
	}
	return binary.BigEndian.Uint32(b)
}
