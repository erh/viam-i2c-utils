package viami2cutils

import (
	"context"
	"fmt"
	"os/exec"

	toggleswitch "go.viam.com/rdk/components/switch"
	"go.viam.com/rdk/logging"
	"go.viam.com/rdk/resource"
)

var I2cOnOff = resource.NewModel("erh", "viam-i2c-utils", "i2c-on-off")

func init() {
	resource.RegisterComponent(toggleswitch.API, I2cOnOff,
		resource.Registration[toggleswitch.Switch, *Config]{
			Constructor: newViamI2cUtilsI2cOnOff,
		},
	)
}

type Config struct {
	Bus      int `json:"bus"`
	Address  int `json:"address"`
	Register int `json:"register"`
	OnValue  int `json:"on_value"`
	OffValue int `json:"off_value"`
}

// Validate ensures all parts of the config are valid and important fields exist.
func (cfg *Config) Validate(path string) ([]string, []string, error) {
	if cfg.Address == 0 {
		return nil, nil, fmt.Errorf("%s: address is required", path)
	}
	return nil, nil, nil
}

type viamI2cUtilsI2cOnOff struct {
	resource.AlwaysRebuild

	name   resource.Name
	logger logging.Logger
	cfg    *Config

	position uint32
}

func newViamI2cUtilsI2cOnOff(ctx context.Context, deps resource.Dependencies, rawConf resource.Config, logger logging.Logger) (toggleswitch.Switch, error) {
	conf, err := resource.NativeConfig[*Config](rawConf)
	if err != nil {
		return nil, err
	}
	return NewI2cOnOff(ctx, deps, rawConf.ResourceName(), conf, logger)
}

func NewI2cOnOff(ctx context.Context, deps resource.Dependencies, name resource.Name, conf *Config, logger logging.Logger) (toggleswitch.Switch, error) {
	s := &viamI2cUtilsI2cOnOff{
		name:     name,
		logger:   logger,
		cfg:      conf,
		position: 0,
	}
	return s, nil
}

func (s *viamI2cUtilsI2cOnOff) Name() resource.Name {
	return s.name
}

func (s *viamI2cUtilsI2cOnOff) i2cset(value int) error {
	cmd := exec.Command("i2cset", "-y",
		fmt.Sprintf("%d", s.cfg.Bus),
		fmt.Sprintf("0x%02x", s.cfg.Address),
		fmt.Sprintf("0x%02x", s.cfg.Register),
		fmt.Sprintf("0x%02x", value),
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("i2cset failed: %w: %s", err, string(output))
	}
	return nil
}

func (s *viamI2cUtilsI2cOnOff) SetPosition(ctx context.Context, position uint32, extra map[string]interface{}) error {
	var value int
	if position == 0 {
		value = s.cfg.OffValue
	} else {
		value = s.cfg.OnValue
	}
	if err := s.i2cset(value); err != nil {
		return err
	}
	s.position = position
	return nil
}

func (s *viamI2cUtilsI2cOnOff) GetPosition(ctx context.Context, extra map[string]interface{}) (uint32, error) {
	return s.position, nil
}

func (s *viamI2cUtilsI2cOnOff) GetNumberOfPositions(ctx context.Context, extra map[string]interface{}) (uint32, []string, error) {
	return 2, []string{"off", "on"}, nil
}

func (s *viamI2cUtilsI2cOnOff) DoCommand(ctx context.Context, cmd map[string]interface{}) (map[string]interface{}, error) {
	return nil, nil
}

func (s *viamI2cUtilsI2cOnOff) Close(context.Context) error {
	return nil
}
