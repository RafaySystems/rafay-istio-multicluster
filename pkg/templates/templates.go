package templates

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/RafaySystems/rafay-istio-multicluster/pkg/funcmap"
	"github.com/RafaySystems/rafay-istio-multicluster/pkg/log"
	"github.com/imdario/mergo"
	"github.com/pkg/errors"
)

type TmplOptions struct {
	// Teemplate files
	Templates []string

	// Value files
	Values []string

	// Value config merged
	Config map[string]interface{}

	// Add the environment map to the variables.
	Environment string

	// Time format
	TimeFormat string
}

// Templates
type Templates interface {
	Render(logger log.Logger) ([]byte, error)
}

// New
func New(toptions TmplOptions) Templates {
	return &toptions
}

// Render
func Render(logger log.Logger, tmplsOpts TmplOptions) ([]byte, error) {
	return tmplsOpts.Render(logger)
}

// Render
func (toptions *TmplOptions) Render(logger log.Logger) ([]byte, error) {
	var retData []byte
	globalContext := map[string]interface{}{}

	globalContext = toptions.Retyper(globalContext, retypeSingleElementSlice)

	files := toptions.Templates

	if toptions.Environment != "" || len(files) == 0 {
		v := make(map[string]string)
		for _, item := range os.Environ() {
			splits := strings.Split(item, "=")
			v[splits[0]] = strings.Join(splits[1:], "=")
		}
		globalContext[toptions.Environment] = v
	}

	mergo.Merge(&globalContext, toptions.Config)

	logger.Debugf("globalContext: %#v", globalContext)

	status := 0

	for _, file := range files {
		status |= func() (status int) {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("PANIC: %+v", r)
					status = 15
				}
			}()

			var err error
			var data []byte
			var r io.ReadCloser
			if file == "" {
				r = os.Stdin
			} else {
				r, err = os.Open(file)
				if err != nil {
					err = errors.Wrapf(err, " values '%v'", toptions.Values)
					logger.Warn(err)
					return 1
				}
				defer r.Close()
			}
			f, err := ioutil.ReadAll(r)
			if err != nil {
				err = errors.Wrapf(err, " values '%v'", toptions.Values)
				logger.Warn(err)
				return 2
			}
			data = f

			tmpl, err := template.New(file).
				Option(fmt.Sprintf("missingkey=%s", "error")).
				Funcs(funcmap.Map).
				Parse(string(data))
			if err != nil {
				err = errors.Wrapf(err, " values '%v'", toptions.Values)
				logger.Warn(err)
				return 4
			}

			var b bytes.Buffer
			err = tmpl.Execute(&b, globalContext)
			if err != nil {
				err = errors.Wrapf(err, " values '%v'", toptions.Values)
				logger.Warn(err)
				return 8
			}

			data = b.Bytes()
			retData = append(retData, data...)
			//fmt.Println(string(data))

			return 0
		}()
	}

	if status != 0 {
		return nil, fmt.Errorf("failed to render from templates")
	}

	fmt.Println("\n", string(retData))

	return retData, nil
}

func (toptions TmplOptions) typerep(d string) (result interface{}) {
	v := string(d)
	if parsedValue, err := strconv.ParseInt(v, 10, 64); err == nil {
		result = parsedValue
	} else if parsedValue, err := strconv.ParseFloat(v, 64); err == nil {
		result = parsedValue
	} else if parsedValue, err := strconv.ParseBool(v); err == nil {
		result = parsedValue
	} else if parsedValue, err := time.Parse(toptions.TimeFormat, v); err == nil {
		result = parsedValue
	} else {
		result = v
	}
	return
}

type retyperTmplOptions func(*retyperConfig)
type retyperConfig struct {
	retypeSingleElementSlice bool
}

func retypeSingleElementSlice(config *retyperConfig) { config.retypeSingleElementSlice = true }

func (toptions TmplOptions) retyping(source map[string]interface{}, config retyperConfig) map[string]interface{} {
	for k, v := range source {
		switch vt := v.(type) {
		case map[interface{}]interface{}:
			f := make(map[string]interface{}, len(vt))
			for e, s := range vt {
				f[fmt.Sprintf("%v", e)] = s
			}
			toptions.retyping(f, config)
			source[k] = f
		case map[string]string:
			f := make(map[string]interface{}, len(vt))
			for e, s := range vt {
				f[e] = s
			}
			toptions.retyping(f, config)
			source[k] = f
		case int:
			source[k] = int64(vt)
		case map[string]interface{}:
			toptions.retyping(vt, config)
		case []interface{}:
			if config.retypeSingleElementSlice && len(vt) == 1 {
				source[k] = vt[0]
				continue
			}
			kind := reflect.Invalid
			valid := func() bool {
				for i, value := range vt {
					t := reflect.TypeOf(value)
					if i == 0 {
						kind = t.Kind()
						continue
					}

					if kind == reflect.Invalid || t.Kind() != kind {
						return false
					}
				}
				return true
			}()
			if !valid {
				source[k] = vt
				continue
			}

			switch kind {
			case reflect.Bool:
				rt := make([]bool, len(vt))
				for i, value := range vt {
					rt[i] = value.(bool)
				}
				source[k] = rt
			case reflect.Int64:
				rt := make([]int64, len(vt))
				for i, value := range vt {
					rt[i] = value.(int64)
				}
				source[k] = rt
			case reflect.Float64:
				rt := make([]float64, len(vt))
				for i, value := range vt {
					rt[i] = value.(float64)
				}
				source[k] = rt
			case reflect.String:
				rt := make([]string, len(vt))
				for i, value := range vt {
					rt[i] = value.(string)
				}
				source[k] = rt
			default:
				source[k] = vt
			}
		case string:
			source[k] = toptions.typerep(vt)
		case bool, int64, float64:
			source[k] = vt
		default:
			source[k] = vt
			log.GetLogger().Warnf("WARNING: unexpected %[1]T %#[1]v", vt)
		}
	}
	return source
}

func (toptions TmplOptions) Retyper(source map[string]interface{}, options ...retyperTmplOptions) map[string]interface{} {
	config := retyperConfig{}
	for _, f := range options {
		f(&config)
	}
	d := toptions.retyping(source, config)
	return d
}
