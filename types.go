package prostar_pwm

import "fmt"

type RawADCData struct {
	SupplyVoltage                           *float32 // V, vdd_actual   3.3V Supply Voltage
	GateDriveVoltage                        *float32 // V, adc_fgdrive  Gate Drive Voltage
	MeterBusSupplyVoltage                   *float32 // V, adc_pmeter   MeterBus Supply Voltage
	InternalReferenceVoltage                *float32 // V, adc_vrefint  Internal Reference Voltage
	NegativeSupplyRailForCurrentMeasurement *float32 // V, adc_FN3      Negative Supply rail for current measurement
	LoadFETGateVoltage                      *float32 // V, adc_gload    Load FET gate voltage
	ArrayFETGateVoltage                     *float32 // V, adc_gatepv   Array FET gate voltage
}

type FilteredADCData struct {
	ArrayCurrent           *float32 // A, adc_ia       Array Current
	BatteryTerminalVoltage *float32 // V, adc_vbterm   Battery Terminal Voltage
	ArrayVoltage           *float32 // V, adc_va       Array Voltage
	LoadVoltage            *float32 // V, adc_vl       Load Voltage
	LoadCurrent            *float32 // A, adc_il       Load Current
	BatterySenseVoltage    *float32 // V, adc_vbsense  Battery Sense Voltage
	BatteryVoltage         *float32 // V, adc_vb_f_1m  Battery Voltage, slow filter (60s)
	BatteryCurrent         *float32 // A, adc_ib_f_1m  Battery Current (net), slow filter (60s)
}

type TemperatureData struct {
	Heatsink *float32 // ºC, T_hs    Heatsink Temperature
	Battery  *float32 // ºC, T_batt  Battery Temperature (Either Ambient or RTS is connected)
	Ambient  *float32 // ºC, T_amb   Ambient (local) Temperature
	Remote   *float32 // ºC, T_rts   Remote Temperature Sensor Temperature
}

type ChargerStatus struct {
	ChargeState                                     *ChargeState       // charge_state         Charge State
	ArrayFault                                      *ArrayFaultDetails // array_fault          Array Fault Bitfield
	BatteryVoltage                                  *float32           // V,  vb_f             Battery Voltage, slow filter (25s)
	BatteryRegulatorReferenceVoltage                *float32           // V,  vb_ref           Battery Regulator Reference Voltage
	AhChargeResettable                              *float32           // Ah, Ahc_r            Ah Charge Resettable
	AhChargeTotal                                   *float32           // Ah, Ahc_t            Ah Charge Total
	KWhChargeResettable                             *float32           // kWh, kWhc_r          kWh Charge Resettable
	KWhChargeTotal                                  *float32           // kWh, kWhc_t          kWh Charge Total
	BatteryTemperatureFoldback100PercentOutputLimit *float32           // ºC, Tb_lo_limit_100  Battery Temp Foldback 100% Output Limit
	BatteryTemperatureFoldback0PercentOutputLimit   *float32           // ºC, Tb_lo_limit_0    Battery Temp Foldback 0% Output Limit
}

type ChargeState uint16

const (
	ChargeStateStart ChargeState = iota
	ChargeStateNightCheck
	ChargeStateDisconnect
	ChargeStateNight
	ChargeStateFault
	ChargeStateBulk
	ChargeStateAbsorption
	ChargeStateFloat
	ChargeStateEqualize
)

func (v ChargeState) String() string {
	switch v {
	case ChargeStateStart:
		return "Start"
	case ChargeStateNightCheck:
		return "Night Check"
	case ChargeStateDisconnect:
		return "Disconnect"
	case ChargeStateNight:
		return "Night"
	case ChargeStateFault:
		return "Fault"
	case ChargeStateBulk:
		return "Bulk"
	case ChargeStateAbsorption:
		return "Absorption"
	case ChargeStateFloat:
		return "Float"
	case ChargeStateEqualize:
		return "Equalize"
	default:
		return fmt.Sprintf("0x%04x", uint16(v))
	}
}

type ArrayFault uint16

func (v ArrayFault) Details() ArrayFaultDetails {
	return ArrayFaultDetails{
		Raw:                          uint16(v),
		OvercurrentPhase1:            checkBit(uint16(v), 0),
		FETsShorted:                  checkBit(uint16(v), 1),
		SoftwareBug:                  checkBit(uint16(v), 2),
		BatteryHighVoltageDisconnect: checkBit(uint16(v), 3),
		ArrayHighVoltageDisconnect:   checkBit(uint16(v), 4),
		EEPROMSettingEdit:            checkBit(uint16(v), 5),
		RTSShorted:                   checkBit(uint16(v), 6),
		RTSWasValidNowDisconnected:   checkBit(uint16(v), 7),
		LocalTemperatureSensorFailed: checkBit(uint16(v), 8),
		BatteryLowVoltageDisconnect:  checkBit(uint16(v), 9),
		DIPSwitchChanged:             checkBit(uint16(v), 10),
		ProcessorSupplyFault:         checkBit(uint16(v), 11),
	}
}

type ArrayFaultDetails struct {
	Raw                          uint16
	OvercurrentPhase1            bool
	FETsShorted                  bool
	SoftwareBug                  bool
	BatteryHighVoltageDisconnect bool
	ArrayHighVoltageDisconnect   bool
	EEPROMSettingEdit            bool
	RTSShorted                   bool
	RTSWasValidNowDisconnected   bool
	LocalTemperatureSensorFailed bool
	BatteryLowVoltageDisconnect  bool
	DIPSwitchChanged             bool
	ProcessorSupplyFault         bool
}

type LoadStatus struct {
	LoadState                        *LoadState        // load_state
	LoadFault                        *LoadFaultDetails // load_fault
	LoadCurrentCompensatedLVDVoltage *float32          // V_lvd
	LoadHVDVoltage                   *float32          // V_lhvd
	AhLoadResettable                 *float32          // Ahl_r
	AhLoadTotal                      *float32          // Ahl_t
}

type LoadState uint16

const (
	LoadStateStart LoadState = iota
	LoadStateLoadOn
	LoadStateLVDWarning
	LoadStateLVD
	LoadStateFault
	LoadStateDisconnect
	LoadStateLoadOff
	LoadStateOverride
)

func (v LoadState) String() string {
	switch v {
	case LoadStateStart:
		return "Start"
	case LoadStateLoadOn:
		return "Load On"
	case LoadStateLVDWarning:
		return "LVD Warning"
	case LoadStateLVD:
		return "LVD"
	case LoadStateFault:
		return "Fault"
	case LoadStateDisconnect:
		return "Disconnect"
	case LoadStateLoadOff:
		return "Load Off"
	case LoadStateOverride:
		return "Override"
	default:
		return fmt.Sprintf("0x%04x", uint16(v))
	}
}

type LoadFault uint16

func (v LoadFault) Details() LoadFaultDetails {
	return LoadFaultDetails{
		Raw:                     uint16(v),
		ExternalShortCircuit:    checkBit(uint16(v), 0),
		Overcurrent:             checkBit(uint16(v), 1),
		FETsShorted:             checkBit(uint16(v), 2),
		SoftwareBug:             checkBit(uint16(v), 3),
		HighVoltageDisconnect:   checkBit(uint16(v), 4),
		HeatsinkOverTemperature: checkBit(uint16(v), 5),
		DIPSwitchChanged:        checkBit(uint16(v), 6),
		EEPROMSettingEdit:       checkBit(uint16(v), 7),
		FP10Fault:               checkBit(uint16(v), 8),
		ProcessorSupplyFault:    checkBit(uint16(v), 9),
	}
}

type LoadFaultDetails struct {
	Raw                     uint16
	ExternalShortCircuit    bool
	Overcurrent             bool
	FETsShorted             bool
	SoftwareBug             bool
	HighVoltageDisconnect   bool
	HeatsinkOverTemperature bool
	DIPSwitchChanged        bool
	EEPROMSettingEdit       bool
	FP10Fault               bool
	ProcessorSupplyFault    bool
}

type MiscData struct {
	Hourmeter            *uint32               // hours, hourmeter
	Alarm                *AlarmDetails         // alarm
	DIPSwitch            *uint16               // dip_switch
	LEDState             *LEDState             // led_state
	ChargeStatusLEDState *ChargeStatusLEDState // charge_led_state
	LightingShouldBeOn   *uint16               // lighting_should_be_on
}

type Alarm uint32

func (v Alarm) Details() AlarmDetails {
	return AlarmDetails{
		Raw:                              uint32(v),
		RTSOpen:                          checkBit(uint32(v), 0),
		RTSShort:                         checkBit(uint32(v), 1),
		RTSDisconnected:                  checkBit(uint32(v), 2),
		HeatsinkTemperatureSensorOpen:    checkBit(uint32(v), 3),
		HeatsinkTemperatureSensorShorted: checkBit(uint32(v), 4),
		HeatsinkHot:                      checkBit(uint32(v), 5),
		CurrentLimit:                     checkBit(uint32(v), 6),
		IOffset:                          checkBit(uint32(v), 7),
		BatterySenseOutOfRange:           checkBit(uint32(v), 8),
		BatterySenseDisconnected:         checkBit(uint32(v), 9),
		Uncalibrated:                     checkBit(uint32(v), 10),
		BatteryTemperatureOutOfRange:     checkBit(uint32(v), 11),
		FP10SupplyOutOfRange:             checkBit(uint32(v), 12),
		FETOpen:                          checkBit(uint32(v), 13),
		IAOffset:                         checkBit(uint32(v), 14),
		ILOffset:                         checkBit(uint32(v), 15),
		SupplyOutOfRange:                 checkBit(uint32(v), 16),
		Reset:                            checkBit(uint32(v), 19),
		LVD:                              checkBit(uint32(v), 20),
		LogTimeout:                       checkBit(uint32(v), 21),
		EEPROMAccessFailure:              checkBit(uint32(v), 22),
	}
}

type AlarmDetails struct {
	Raw                              uint32
	RTSOpen                          bool
	RTSShort                         bool
	RTSDisconnected                  bool
	HeatsinkTemperatureSensorOpen    bool
	HeatsinkTemperatureSensorShorted bool
	HeatsinkHot                      bool
	CurrentLimit                     bool
	IOffset                          bool
	BatterySenseOutOfRange           bool
	BatterySenseDisconnected         bool
	Uncalibrated                     bool
	BatteryTemperatureOutOfRange     bool
	FP10SupplyOutOfRange             bool
	FETOpen                          bool
	IAOffset                         bool
	ILOffset                         bool
	SupplyOutOfRange                 bool
	Reset                            bool
	LVD                              bool
	LogTimeout                       bool
	EEPROMAccessFailure              bool
}

type LEDState uint16

// TODO

type ChargeStatusLEDState uint16

// TODO
