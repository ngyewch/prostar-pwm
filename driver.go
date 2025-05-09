package prostar_pwm

import (
	"encoding/binary"
	"errors"
	"github.com/simonvetter/modbus"
	"github.com/x448/float16"
	"sync"
)

type Dev struct {
	mc     *modbus.ModbusClient
	unitId uint8
	mutex  *sync.Mutex
}

func New(mc *modbus.ModbusClient, unitId uint8, mutex *sync.Mutex) *Dev {
	return &Dev{
		mc:     mc,
		unitId: unitId,
		mutex:  mutex,
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

	r.SupplyVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0004)
	if err != nil {
		return RawADCData{}, err
	}
	r.GateDriveVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0005)
	if err != nil {
		return RawADCData{}, err
	}
	r.MeterBusSupplyVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0006)
	if err != nil {
		return RawADCData{}, err
	}
	r.InternalReferenceVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0007)
	if err != nil {
		return RawADCData{}, err
	}
	r.NegativeSupplyRailForCurrentMeasurement, err = dev.readInputRegisterFromFloat16ToFloat32(0x0008)
	if err != nil {
		return RawADCData{}, err
	}
	r.LoadFETGateVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0009)
	if err != nil {
		return RawADCData{}, err
	}
	r.ArrayFETGateVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x000a)
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

	r.ArrayCurrent, err = dev.readInputRegisterFromFloat16ToFloat32(0x0011)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.BatteryTerminalVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0012)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.ArrayVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0013)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.LoadVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0014)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.LoadCurrent, err = dev.readInputRegisterFromFloat16ToFloat32(0x0016)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.BatterySenseVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0017)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.BatteryVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0018)
	if err != nil {
		return FilteredADCData{}, err
	}
	r.BatteryCurrent, err = dev.readInputRegisterFromFloat16ToFloat32(0x0019)
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

	r.Heatsink, err = dev.readInputRegisterFromFloat16ToFloat32(0x001a)
	if err != nil {
		return TemperatureData{}, err
	}
	r.Battery, err = dev.readInputRegisterFromFloat16ToFloat32(0x001b)
	if err != nil {
		return TemperatureData{}, err
	}
	r.Ambient, err = dev.readInputRegisterFromFloat16ToFloat32(0x001c)
	if err != nil {
		return TemperatureData{}, err
	}
	r.Remote, err = dev.readInputRegisterFromFloat16ToFloat32(0x001d)
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
		v, err := dev.readInputRegister(0x0021)
		if err != nil {
			return ChargerStatus{}, err
		}
		if v != nil {
			v2 := ChargeState(*v)
			r.ChargeState = &v2
		}
	}
	{
		v, err := dev.readInputRegister(0x0022)
		if err != nil {
			return ChargerStatus{}, err
		}
		if v != nil {
			v2 := ArrayFault(*v)
			details := v2.Details()
			r.ArrayFault = &details
		}
	}
	r.BatteryVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0023)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.BatteryRegulatorReferenceVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0024)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.AhChargeResettable, err = dev.readInputRegisterFromUint32ToFloat32(0x0026, 10)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.AhChargeTotal, err = dev.readInputRegisterFromUint32ToFloat32(0x0028, 10)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.KWhChargeResettable, err = dev.readInputRegisterFromUint16ToFloat32(0x002a, 10)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.KWhChargeTotal, err = dev.readInputRegisterFromUint16ToFloat32(0x002b, 10)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.BatteryTemperatureFoldback100PercentOutputLimit, err = dev.readInputRegisterFromFloat16ToFloat32(0x002c)
	if err != nil {
		return ChargerStatus{}, err
	}
	r.BatteryTemperatureFoldback0PercentOutputLimit, err = dev.readInputRegisterFromFloat16ToFloat32(0x002d)
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
		v, err := dev.readInputRegister(0x002e)
		if err != nil {
			return LoadStatus{}, err
		}
		if v != nil {
			v2 := LoadState(*v)
			r.LoadState = &v2
		}
	}
	{
		v, err := dev.readInputRegister(0x002f)
		if err != nil {
			return LoadStatus{}, err
		}
		if v != nil {
			v2 := LoadFault(*v)
			details := v2.Details()
			r.LoadFault = &details
		}
	}
	r.LoadCurrentCompensatedLVDVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0030)
	if err != nil {
		return LoadStatus{}, err
	}
	r.LoadHVDVoltage, err = dev.readInputRegisterFromFloat16ToFloat32(0x0031)
	if err != nil {
		return LoadStatus{}, err
	}
	r.AhLoadResettable, err = dev.readInputRegisterFromUint32ToFloat32(0x0032, 10)
	if err != nil {
		return LoadStatus{}, err
	}
	r.AhLoadTotal, err = dev.readInputRegisterFromUint32ToFloat32(0x0034, 10)
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

	r.Hourmeter, err = dev.readInputRegisterFromUint32(0x0036)
	if err != nil {
		return MiscData{}, err
	}
	{
		v, err := dev.readInputRegisterFromUint32(0x0038)
		if err != nil {
			return MiscData{}, err
		}
		if v != nil {
			v2 := Alarm(*v)
			details := v2.Details()
			r.Alarm = &details
		}
	}
	r.DIPSwitch, err = dev.readInputRegister(0x003a)
	{
		v, err := dev.readInputRegister(0x003b)
		if err != nil {
			return MiscData{}, err
		}
		if v != nil {
			v2 := LEDState(*v)
			r.LEDState = &v2
		}
	}
	{
		v, err := dev.readInputRegister(0x004d)
		if err != nil {
			return MiscData{}, err
		}
		if v != nil {
			v2 := ChargeStatusLEDState(*v)
			r.ChargeStatusLEDState = &v2
		}
	}
	r.LightingShouldBeOn, err = dev.readInputRegister(0x004e)
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

	r.RegulationVoltageAt25C, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe000)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.FloatVoltageAt25C, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe001)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.TimeBeforeEnteringFloat, err = dev.readHoldingRegister(0xe002)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.TimeBeforeEnteringFloatDueToLowBattery, err = dev.readHoldingRegister(0xe003)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.VoltageTriggerForLowBatteryFloatTime, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe004)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.VoltageToCancelFloat, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe005)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.ExitFloatTime, err = dev.readHoldingRegister(0xe006)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.EqualizeVoltageAt25C, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe007)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.DaysBetweenEQCycles, err = dev.readHoldingRegister(0xe008)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.EqualizeTimeLimitAboveEVReg, err = dev.readHoldingRegister(0xe009)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.EqualizeTimeLimitAtEVEq, err = dev.readHoldingRegister(0xe00a)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.ReferenceChargeVoltageLimit, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe010)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.TemperatureCompensationCoefficient, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe01a)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.HighVoltageDisconnectAt25C, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe01b)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.HighVoltageReconnect, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe01c)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.MaximumChargeVoltageReference, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe01d)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.MaxBatteryTempCompensationLimit, err = dev.readHoldingRegisterFromUint16ToInt16(0xe01e)
	if err != nil {
		return ChargeSettings{}, err
	}
	r.MinBatteryTempCompensationLimit, err = dev.readHoldingRegisterFromUint16ToInt16(0xe01f)
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

	r.LowVoltageDisconnect, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe022)
	if err != nil {
		return LoadSettings{}, err
	}
	r.LoadHighVoltageReconnect, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe023)
	if err != nil {
		return LoadSettings{}, err
	}
	r.LoadHighVoltageDisconnect, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe024)
	if err != nil {
		return LoadSettings{}, err
	}
	r.LoadHighVoltageReconnect, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe025)
	if err != nil {
		return LoadSettings{}, err
	}
	r.LVDLoadCurrentCompensation, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe026)
	if err != nil {
		return LoadSettings{}, err
	}
	r.LVDWarningTimeout, err = dev.readHoldingRegister(0xe027)
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

	r.LEDGreenToGreenAndYellowLimit, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe030)
	if err != nil {
		return MiscSettings{}, err
	}
	r.LEDGreenAndYellowToYellowLimit, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe031)
	if err != nil {
		return MiscSettings{}, err
	}
	r.LEDYellowToYellowAndRedLimit, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe032)
	if err != nil {
		return MiscSettings{}, err
	}
	r.LEDYellowAndRedToRedFlashingLimit, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe033)
	if err != nil {
		return MiscSettings{}, err
	}
	r.ModbusID, err = dev.readHoldingRegister(0xe034)
	if err != nil {
		return MiscSettings{}, err
	}
	r.MeterbusID, err = dev.readHoldingRegister(0xe035)
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

	r.ChargeCurrentLimit, err = dev.readHoldingRegisterFromFloat16ToFloat32(0xe038)
	if err != nil {
		return PWMSettings{}, err
	}

	return r, nil
}

func (dev *Dev) readInputRegister(addr uint16) (*uint16, error) {
	v, err := dev.mc.ReadRegister(addr, modbus.INPUT_REGISTER)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &v, nil
}

func (dev *Dev) readHoldingRegister(addr uint16) (*uint16, error) {
	v, err := dev.mc.ReadRegister(addr, modbus.HOLDING_REGISTER)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &v, nil
}

func (dev *Dev) readHoldingRegisterFromUint16ToInt16(addr uint16) (*int16, error) {
	v, err := dev.mc.ReadRegister(addr, modbus.HOLDING_REGISTER)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	v2 := int16(v)

	return &v2, nil
}

func (dev *Dev) readInputRegisterFromFloat16ToFloat32(addr uint16) (*float32, error) {
	v, err := dev.mc.ReadRegister(addr, modbus.INPUT_REGISTER)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	f16 := float16.Frombits(v)
	f32 := f16.Float32()
	return &f32, nil
}

func (dev *Dev) readHoldingRegisterFromFloat16ToFloat32(addr uint16) (*float32, error) {
	v, err := dev.mc.ReadRegister(addr, modbus.HOLDING_REGISTER)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	f16 := float16.Frombits(v)
	f32 := f16.Float32()
	return &f32, nil
}

func (dev *Dev) readInputRegisterFromUint16ToFloat32(addr uint16, divisor float32) (*float32, error) {
	v, err := dev.mc.ReadRegister(addr, modbus.INPUT_REGISTER)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	f32 := float32(v) / divisor
	return &f32, nil
}

func (dev *Dev) readInputRegisterFromUint32(addr uint16) (*uint32, error) {
	v, err := dev.mc.ReadRegisters(addr, 2, modbus.INPUT_REGISTER)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	var b []byte
	b = binary.BigEndian.AppendUint16(b, v[0])
	b = binary.BigEndian.AppendUint16(b, v[1])
	u32 := binary.BigEndian.Uint32(b)
	return &u32, nil
}

func (dev *Dev) readInputRegisterFromUint32ToFloat32(addr uint16, divisor float32) (*float32, error) {
	v, err := dev.mc.ReadRegisters(addr, 2, modbus.INPUT_REGISTER)
	if err != nil {
		if errors.Is(err, modbus.ErrIllegalDataAddress) {
			return nil, nil
		} else {
			return nil, err
		}
	}
	var b []byte
	b = binary.BigEndian.AppendUint16(b, v[0])
	b = binary.BigEndian.AppendUint16(b, v[1])
	f32 := float32(int32(binary.BigEndian.Uint32(b))) / divisor
	return &f32, nil
}
