package CORE

import(
	"github.com/Esseh/retrievable"
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
	testing.Coverage()
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
	testing.Coverage()
}

func TestGetLocationName(t *testing.T){
	// if err != nil codepath
	if _, err := GetLocationName("nonexistant-country","nonexistant-region"); err == nil {
		fmt.Println("FAIL GetLocationName 1")
		t.Fail()
	}
	// r.Code == region codepath
	if s, _ := GetLocationName("US","CA"); s != "California, United States"{
		fmt.Println("FAIL GetLocationName 2")
		t.Fail()
	}
	// normal codepath
	if s, _ := GetLocationName("US",""); s != "United States" {
		fmt.Println("FAIL GetLocationName 3")
		t.Fail()
	}
	testing.Coverage()
}

func TestEncrypt(t *testing.T){
	// err != nil codepath
	if _, err := Encrypt([]byte("hello"), []byte{1,2,3}); err == nil {
		fmt.Println("FAIL Encrypt 1")
		t.Fail()
	}
	// for len(data) < b.Blocksize codepath
	s1, _ := Encrypt([]byte("hoi"), EncryptKey)
	s2, _ := Decrypt(s1,EncryptKey)
	
	if "hoi" != string(s2) {
		fmt.Println("FAIL Encrypt 2")
		t.Fail()
	}
	// normal code path
	s3, _ := Encrypt([]byte("hoiiiiiiiiiiiiii"), EncryptKey)
	s4, _ := Decrypt(s3,EncryptKey)
	if "hoiiiiiiiiiiiiii" != string(s4) {
		fmt.Println(string(s4))
		fmt.Println("FAIL Encrypt 3")
		t.Fail()
	}
	testing.Coverage()
}

func TestDecrypt(t *testing.T){
	// err != nil	1
	if _, err := Decrypt("",[]byte{1,2,3}); err == nil {
		fmt.Println("FAIL Decrypt 1")
		t.Fail()
	}
	// err != nil	2
	if _, err := Decrypt("122312[p[[p",EncryptKey); err == nil {
		fmt.Println("FAIL Decrypt 2")
		t.Fail()
	}
	// normal path
	s1, _ := Encrypt([]byte("hoi"), EncryptKey)
	s2, _ := Decrypt(s1,EncryptKey)
	if "hoi" != string(s2) {
		fmt.Println("FAIL Decrypt 3")
		t.Fail()
	}
	testing.Coverage()
}

func TestGetAvatarPath(t *testing.T){
	if GetAvatarPath(int64(0)) != "users/0/avatar" {
		fmt.Println("FAIL GetAvatarPath")
		t.Fail()
	}
	testing.Coverage()
}

func TestEscapeString(t *testing.T){
	if EscapeString("<ScRiPt>") != "<&#115;&#99;&#114;&#105;&#112;&#116;>" {
		fmt.Println("FAIL EscapeString 1")
		t.Fail()
	}
	if EscapeString("<iFrAmE>") != "<&#105;&#102;&#114;&#97;&#109;&#101;>" {
		fmt.Println("FAIL EscapeString 2")	
		t.Fail()
	}
	testing.Coverage()
}

func TestInc(t *testing.T){
	if Inc("1") != "2" {
		fmt.Println("FAIL Inc")
		t.Fail()
	}
	testing.Coverage()
}

func TestToInt(t *testing.T){
	if ToInt(retrievable.IntID(int64(0))) != int64(0) {
		fmt.Println("Fail ToInt")
		t.Fail()
	}
	testing.Coverage()
}

/*

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

//Exactly What it Says on the Tin
func ServeTemplateWithParams(res http.ResponseWriter, templateName string, params interface{}) error {
	return TPL.ExecuteTemplate(res, templateName, &params)
}
*/