package CORE

import(
	"encoding/base64"
	"testing"
	"fmt"
)

func TestSplitMac(t *testing.T){
	// if i == -1 { code path
	v1, v2 := SplitMac("t")
	if v1 != "t" || v2 != "" {
		fmt.Println("FAIL SplitMac 1")
		t.Fail()
	}
	// normal code path
	v1, v2  = SplitMac("t.t")
	if v1 != "t" || v2 != "t"{
		fmt.Println("FAIL SplitMac 2")
		t.Fail()
	}
	testing.Coverage()
}
func TestCreateHmac(t *testing.T){
	if string(CreateHmac("hello")) != string([]byte{0x2a,0x0a,0x47,0x5e,0x3b,0xd2,0x79,0x5e,0x9f,0x01,0xfa,0x4a,0xce,0x1c,0xab,0x0a,0x41,0xfb,0x09,0xf3,0x8c,0x57,0xed,0x1f,0x14,0xca,0xf8,0xc6,0xf5,0x16,0x1f,0xe0}) {
		fmt.Println("FAIL CreateHmac")
		t.Fail()	
	}
	testing.Coverage()
}
func TestCheckMac(t *testing.T){
	// normal code path
	if !CheckMac("hello",base64.RawURLEncoding.EncodeToString(CreateHmac("hello"))){
		fmt.Println("FAIL CheckMac 1")
		t.Fail()
	}
	// if err != nil code path
	if CheckMac("hel lo","hel lo"){
		fmt.Println("FAIL CheckMac 2")
		t.Fail()
	}
}
func TestValidLogin(t *testing.T){
	if ValidLogin("","a"){
		fmt.Println("FAIL ValidLogin 1")
		t.Fail()
	}
	if ValidLogin("a",""){
		fmt.Println("FAIL ValidLogin 2")
		t.Fail()
	}
	if !ValidLogin("a","a"){
		fmt.Println("FAIL ValidLogin 3")
		t.Fail()
	}
}

/*

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
*/