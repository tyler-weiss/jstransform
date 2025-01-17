package transform

import (
	"errors"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"time"

	jsonpath "github.com/GannettDigital/PaesslerAG_jsonpath"
	"github.com/antchfx/xmlquery"
	"github.com/microcosm-cc/bluemonday"
)

var durationRe = regexp.MustCompile(`^([\d]*?):?([\d]*):([\d]*)$`)

// duration is a transformOperation which changes from a string duration like
// "MM:SS" to a number of seconds as an integer.
type duration struct {
	re *regexp.Regexp
}

func (c *duration) init(args map[string]string) error {
	c.re = durationRe
	return nil
}

func (c *duration) transform(raw interface{}) (interface{}, error) {
	if array, ok := raw.([]interface{}); ok && len(array) == 1 {
		raw = array[0]
	}

	if raw == nil {
		return 0, nil
	}
	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("duration only supports strings")
	}

	matches := c.re.FindStringSubmatch(in)
	if matches == nil || len(matches) != 4 {
		return nil, errors.New("duration transform input did not match 'MM:SS' or 'HH:MM:SS'")
	}
	hours, err := strconv.Atoi(matches[1])
	if err != nil && matches[1] != "" {
		return nil, err
	}
	minutes, err := strconv.Atoi(matches[2])
	if err != nil && matches[2] != "" {
		return nil, err
	}
	seconds, err := strconv.Atoi(matches[3])
	if err != nil && matches[3] != "" {
		return nil, err
	}

	minutes += 60 * hours
	seconds += 60 * minutes

	return seconds, nil
}

// changeCase is a transformOperation which changes the case of strings.
type changeCase struct {
	args map[string]string
}

func (c *changeCase) init(args map[string]string) error {
	if err := requiredArgs([]string{"to"}, args); err != nil {
		return err
	}
	value := args["to"]
	if value != "lower" && value != "upper" {
		return errors.New("the argument 'to' is required and must be either 'lower' or 'upper'")
	}

	c.args = args
	return nil
}

func (c *changeCase) transform(raw interface{}) (interface{}, error) {
	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("changeCase only supports strings")
	}

	switch c.args["to"] {
	case "lower":
		return strings.ToLower(in), nil
	case "upper":
		return strings.ToUpper(in), nil
	}
	return nil, errors.New("unknown error in changeCase")
}

// notEmpty is a transformOperation which returns a boolean depending on if the passed in value is considered "empty", as the definition changes on a per-type basis.
type valueExists struct {
	args map[string]string
}

func (n *valueExists) init(args map[string]string) error {
	return nil
}

func (c *valueExists) transform(raw interface{}) (interface{}, error) {
	switch v := raw.(type) {
	case string:
		if len(v) > 0 {
			return true, nil
		}
	case []*xmlquery.Node:
		if len(v) > 0 {
			return true, nil
		}
	default:
		return nil, fmt.Errorf("received unsupported type: %T", raw)
	}

	return false, nil
}

// inverse is a transformOperation which flips the value of a boolean.
type inverse struct {
	args map[string]string
}

func (i *inverse) init(args map[string]string) error {
	return nil
}

func (i *inverse) transform(raw interface{}) (interface{}, error) {
	in, ok := raw.(bool)
	if !ok {
		return nil, errors.New("inverse only supports booleans")
	}
	return !in, nil
}

// max is a transformOperation which retrieves a field from the maximum item in
// an array. The maxiumum item is determined by comparing values in a defined
// number field on the array items.
type max struct {
	args map[string]string
}

func (m *max) init(args map[string]string) error {
	if err := requiredArgs([]string{"by", "return"}, args); err != nil {
		return err
	}
	m.args = args
	return nil
}

func (m *max) transform(in interface{}) (interface{}, error) {
	inArray, ok := in.([]interface{})
	if !ok {
		return nil, errors.New("input must be an array")
	}
	byArg := strings.Replace(m.args["by"], "@", "$", 1)
	returnArg := strings.Replace(m.args["return"], "@", "$", 1)

	var largest float64
	var largestIndex int
	for i, item := range inArray {
		byRaw, err := jsonpath.Get(byArg, item)
		if err != nil {
			return nil, fmt.Errorf("failed extracting 'by' field: %v", err)
		}
		by, ok := byRaw.(float64)
		if !ok {
			byInt, ok := byRaw.(int)
			if !ok {
				return nil, errors.New("by field is not a number")
			}
			by = float64(byInt)
		}
		if by > largest {
			largest = by
			largestIndex = i
		}
	}

	rawReturn, err := jsonpath.Get(returnArg, inArray[largestIndex])
	if err != nil {
		return nil, fmt.Errorf("failed extracting 'return' field: %v", err)
	}

	return rawReturn, nil
}

// replace is a transformOperation which performs a regex based find/replace on
// a string value.
type replace struct {
	args  map[string]string
	regex *regexp.Regexp
}

func (r *replace) init(args map[string]string) error {
	if err := requiredArgs([]string{"regex", "new"}, args); err != nil {
		return err
	}
	re, err := regexp.Compile(args["regex"])
	if err != nil {
		return fmt.Errorf("failed to parse regex %q: %v", args["regex"], err)
	}

	r.regex = re
	r.args = args
	return nil
}

func (r *replace) transform(raw interface{}) (interface{}, error) {
	if r.regex == nil {
		return nil, errors.New("init was not run")
	}
	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("replace only supports strings")
	}

	return r.regex.ReplaceAllString(in, r.args["new"]), nil
}

// split is a transformOperation which splits a string based on a given split string.
type split struct {
	args map[string]string
}

func (s *split) init(args map[string]string) error {
	if err := requiredArgs([]string{"on"}, args); err != nil {
		return err
	}

	s.args = args
	return nil
}

func (s *split) transform(raw interface{}) (interface{}, error) {
	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("split only supports strings")
	}

	// An empty string input should result in an empty array.
	// strings.Split() will return []string{""} instead of an empty array.
	// https://play.golang.org/p/8ySv_t37haN
	if in == "" {
		return []interface{}{}, nil
	}

	splits := strings.Split(in, s.args["on"])

	// Return []interface{} to avoid messing up type casts later in the process
	var interfaceSplits []interface{}
	for _, s := range splits {
		interfaceSplits = append(interfaceSplits, s)
	}
	return interfaceSplits, nil
}

// timeParse is a transformOperation which formats a date string into the layout.
type timeParse struct {
	args map[string]string
}

func (t *timeParse) init(args map[string]string) error {
	if err := requiredArgs([]string{"format", "layout"}, args); err != nil {
		return err
	}

	t.args = args
	return nil
}

func (t *timeParse) transform(raw interface{}) (interface{}, error) {
	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("timeParse only supports strings")
	}
	parsedTime, err := time.Parse(t.args["format"], in)
	if err != nil {
		return nil, fmt.Errorf("time could not be parsed using supplied format")
	}
	return parsedTime.Format(t.args["layout"]), nil
}

// currentTime is a transformOperation which returns the current time in a
// specified format.
type currentTime struct {
	args map[string]string
}

func (c *currentTime) init(args map[string]string) error {
	if err := requiredArgs([]string{"format"}, args); err != nil {
		return err
	}
	c.args = args
	return nil
}

func (c *currentTime) transform(_ interface{}) (interface{}, error) {
	timeFmt := c.args["format"]
	switch c.args["format"] {
	case "RFC3339":
		timeFmt = time.RFC3339
	}
	return time.Now().Format(timeFmt), nil
}

// toCamelCase is a transformOperation which converts strings with dashes to camelCase.
type toCamelCase struct {
	args map[string]string
}

func (c *toCamelCase) init(args map[string]string) error {
	if err := requiredArgs([]string{"delimiter"}, args); err != nil {
		return err
	}
	c.args = args
	return nil
}

func (c *toCamelCase) transform(raw interface{}) (interface{}, error) {
	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("toCamelCase only supports input of type string")
	}

	arr := strings.Split(in, c.args["delimiter"])

	for i, cap := range arr {
		if i == 0 {
			arr[0] = strings.ToLower(cap)
			continue
		}
		arr[i] = strings.Title(cap)
	}

	return strings.Join(arr, ""), nil
}

// removeHTML is a transformOperation which removes all html from a string.
type removeHTML struct {
	args map[string]string
}

func (c *removeHTML) init(args map[string]string) error {
	return nil
}

func (c *removeHTML) transform(raw interface{}) (interface{}, error) {
	in, ok := raw.(string)
	if !ok {
		return nil, errors.New("removeHTML only supports input of type string")
	}

	p := bluemonday.NewPolicy().AddSpaceWhenStrippingTag(true)
	sanitized := p.Sanitize(in)

	s := strings.ReplaceAll(strings.TrimSpace(sanitized), "  ", " ")
	return html.UnescapeString(s), nil
}

// convertToFloat64 is a transformOperation which converts various types to float64.
type convertToFloat64 struct {
	args map[string]string
}

func (c *convertToFloat64) init(args map[string]string) error {
	return nil
}

func (c *convertToFloat64) transform(raw interface{}) (interface{}, error) {
	switch in := raw.(type) {
	case string:
		return strconv.ParseFloat(in, 64)
	case int:
		return float64(in), nil
	case float64:
		return in, nil
	default:
		return nil, fmt.Errorf("convertToFloat64 only supports strings, int, and float64, raw type: %T", raw)
	}
}

// convertToInt64 is a transformOperation which converts various types to int64.
type convertToInt64 struct {
	args map[string]string
}

func (c *convertToInt64) init(args map[string]string) error {
	return nil
}

func (c *convertToInt64) transform(raw interface{}) (interface{}, error) {
	switch in := raw.(type) {
	case string:
		return strconv.ParseInt(in, 10, 64)
	case float32:
		return int64(in), nil
	case float64:
		return int64(in), nil
	case int:
		return int64(in), nil
	case int8:
		return int64(in), nil
	case int16:
		return int64(in), nil
	case int32:
		return int64(in), nil
	case int64:
		return int64(in), nil
	default:
		return nil, fmt.Errorf("convertToInt64 only supports strings, int, and float64, raw type: %T", raw)
	}
}

// requiredArgs checks the given args map to make sure it contains the required args
// and only the required args.
func requiredArgs(required []string, args map[string]string) error {
	if len(args) != len(required) {
		return fmt.Errorf("expected args %v but got %d args", required, len(args))
	}
	for _, arg := range required {
		if _, ok := args[arg]; !ok {
			return fmt.Errorf("argument %q is required", arg)
		}
	}
	return nil
}
