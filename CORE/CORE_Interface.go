package CORE
import(
	"net/http"
	"html/template"
	"github.com/Esseh/retrievable"
	"golang.org/x/net/context"
	"github.com/pariz/gountries"
	"crypto/aes"
	"encoding/base64"
	"time"
	"strings"
	"strconv"
	"bytes"
	"io"
	"crypto/hmac"
	"crypto/sha256"
	"regexp"
)

const (
	// The bucket location for cloud storage.
	GCSBucket  = "csci150project.appspot.com"
	// The domain of our website
	DomainPath = "http://localhost:8080/"
	MaxAvatarSize = 500
	// Time allowed for cookie session
	SessionTime = 7 * 24 * 60 * 60
	// Key for HMAC encoding
	HMAC_Key = "csci150project2016"
)
var (
	// Global Template file.
	TPL *template.Template
	// This key needs to be exactly 32 bytes long
	// TODO This should not be in our git repo
	EncryptKey = []byte{33, 44, 160, 6, 124, 138, 93, 47, 177, 135, 163, 154, 42, 14, 58, 17, 85, 133, 174, 207, 255, 52, 3, 26, 145, 21, 169, 65, 106, 108, 0, 66}
)

func SplitMac(value string) (string, string) {
	i := strings.LastIndex(value, ".")
	if i == -1 {
		return value, ""
	}
	return value[:i], value[i+1:]
}

func CheckMac(value, mac string) bool {
	derivedMac, err := CreateHmac(value)
	if err != nil {
		return false
	}
	macData, err := base64.RawURLEncoding.DecodeString(mac)
	if err != nil {
		return false
	}
	return hmac.Equal(derivedMac, macData)
}

// Encodes HMAC value
func CreateHmac(value string) ([]byte, error) {
	mac := hmac.New(sha256.New, []byte(HMAC_Key))
	_, err := io.WriteString(mac, value)
	if err != nil {
		return []byte{}, err
	}
	return mac.Sum(nil), nil
}

// Checks if a password username combination is valid. It does not ensure that it is correct or that it exists.
func ValidLogin(username,password string) bool {
	return password != "" && username != ""
}

// Takes in country and region headers in order to generate a human readable name.
func GetLocationName(country, region string) (string, error) {
	c, err := gountries.New().FindCountryByAlpha(country)
	if err != nil { return "", err }
	for _, r := range c.SubDivisions() {
		if r.Code == region {
			return r.Name + ", " + c.Name.BaseLang.Common, nil
		}
	}
	return c.Name.BaseLang.Common, nil
}

// Encrypts data based on a key
func Encrypt(data []byte, key []byte) (string, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	for len(data) < b.BlockSize() {
		data = append(data, '=')
	}
	res := make([]byte, len(data))
	b.Encrypt(res, data)
	finalValue := base64.StdEncoding.EncodeToString(res)
	return finalValue, nil
}
// Decrypts data based on a key
func Decrypt(data string, key []byte) ([]byte, error) {
	b, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	strData, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	res := make([]byte, len(strData))
	b.Decrypt(res, strData)
	return bytes.TrimRight(res, "="), nil
}
// Generates the Avatar Path for a URL
func GetAvatarPath(userID int64) string {
	return "users/" + strconv.FormatInt(userID, 10) + "/avatar"
}
/// Parses markdown to produce HTML.
func EscapeString(inp string) string {
	data := []byte(inp)                                    // Convert to Byte
	regex, _ := regexp.Compile("[sS][cC][rR][iI][pP][tT]") // Escape Script Tag
	data = regex.ReplaceAll(data, []byte("&#115;&#99;&#114;&#105;&#112;&#116;"))
	regex2, _ := regexp.Compile("[iI][fF][rR][aA][mM][eE]") // Escape Iframe Tag
	data = regex2.ReplaceAll(data, []byte("&#105;&#102;&#114;&#97;&#109;&#101;"))
	return string(data)
}
// Increments a string integer
func Inc(inp string) string {
	i, _ := strconv.ParseInt(inp, 10, 64)
	return strconv.FormatInt(i+1, 10)
}
//Finds corresponding SVG template
func FindSVG(name string) (ret template.HTML, err error) {
	buf := bytes.NewBuffer([]byte{})
	err = TPL.ExecuteTemplate(buf, ("svg-" + name), nil)
	ret = template.HTML(buf.String())
	return
}
//Finds corresponding template
func FindTemplate(name string) (ret template.HTML, err error) {
	buf := bytes.NewBuffer([]byte{})
	err = TPL.ExecuteTemplate(buf, (name), nil)
	ret = template.HTML(buf.String())
	return
}
// Converts retrievable integer type to int64
func ToInt(id retrievable.IntID) int64 {
	return int64(id)
}
// Data Outputed by AddCtx?
type ContextData struct {
	Ctx  context.Context
	Data interface{}
}
// Adds data to context?
func AddCtx(ctx context.Context, data interface{}) *ContextData {
	return &ContextData{
		Ctx:  ctx,
		Data: data,
	}
}
// Gets the date from a time object.
func GetDate(t time.Time) string {
	return t.Format("2006-01-02")
}
// Gets the Avatar URL
func GetAvatarURL(userID retrievable.IntID) string {
	return "https://storage.googleapis.com/" + GCSBucket + "/" + GetAvatarPath(int64(userID))
}
//gets the Year from a submitted time.Time
func YearFromTime(t time.Time) int {
	return t.Year()
}
//gets the Month from a submitted time.Time
func MonthFromTime(t time.Time) time.Month {
	return t.Month()
}
//gets the Day from a submitted time.Time
func DayFromTime(t time.Time) int {
	return t.Day()
}
//Exactly What it Says on the Tin
func ServeTemplateWithParams(res http.ResponseWriter, templateName string, params interface{}) error {
	return TPL.ExecuteTemplate(res, templateName, &params)
}