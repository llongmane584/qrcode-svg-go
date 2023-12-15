package main

import (
	"fmt"
	"image/png"
	"os"
	"strings"

	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"github.com/makiuchi-d/gozxing/qrcode/decoder"
)

func main() {
	// QRコードに含める内容
	content := "祇園精舍の鐘の声、諸行無常の響きあり。娑羅双樹の花の色、盛者必衰の理をあらはす。驕れる人も久しからず、ただ春の夜の夢のごとし。猛き者もつひにはほろびぬ、ひとへに風の前の塵に同じ。遠く異朝をとぶらへば、秦の趙高、漢の王莽、梁の朱忌、唐の祿山、これらは皆舊主先皇の政にもしたがはず、樂しみをきはめ、諌めをも思ひ入れず、天下の亂れん事を悟らずして、民間の愁ふるところを知らざりしかば、久しからずして、亡じにし者どもなり。近く本朝をうかがふに、承平の將門、天慶の純友、康和の義親、平治の信賴、これらはおごれる心もたけき事も、皆とりどりにこそありしかども、まぢかくは六波羅の入道、前太政大臣平朝臣清盛公と申しし人のありさま、傳へ承るこそ心もことばも及ばれね。その先祖を尋ぬれば桓武天皇第五の皇子、一品式部卿葛原親王九代の後胤、讃岐守正盛が孫、刑部卿忠盛朝臣の嫡男なり。かの親王の御子、高見王、無官無位にして失せ給ひぬ。その御子、高望王の時、初めて平の姓を賜はつて、上総介に成り給ひしより、たちまちに王氏を出でて人臣に列なる、その子鎮守府将軍良望、後には國香と改む。國香より正盛に至る六代は、諸国の受領たりしかども、殿上の仙籍をば未だ赦されず。"

	hints := make(map[gozxing.EncodeHintType]interface{})
	hints[gozxing.EncodeHintType_ERROR_CORRECTION] = decoder.ErrorCorrectionLevel_M
	hints[gozxing.EncodeHintType_QR_VERSION] = 40
	// Shift_JISで作るとLINEで読めるがOSベンダーが提供するカメラでは読めないあるいは文字化けする
	//hints[gozxing.EncodeHintType(gozxing.EncodeHintType_CHARACTER_SET)] = "Shift_JIS"
	hints[gozxing.EncodeHintType(gozxing.EncodeHintType_CHARACTER_SET)] = "UTF-8"

	// QRコードのビットマップを生成
	qrWriter := qrcode.NewQRCodeWriter()
	bitMatrix, err := qrWriter.Encode(content, gozxing.BarcodeFormat_QR_CODE, 600, 600, hints)
	if err != nil {
		panic(err)
	}

	// 生成したQRコードをファイルに保存
	file, err := os.Create("qr.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = png.Encode(file, bitMatrix)
	if err != nil {
		panic(err)
	}

	// SVGフォーマットに変換
	svg := toSVGString(bitMatrix)
	err = os.WriteFile("qr.svg", []byte(svg), 0644)
	if err != nil {
		panic(err)
	}
}

// toSVGString はQRコードのビットマトリックスをSVGフォーマットに変換します。
func toSVGString(matrix *gozxing.BitMatrix) string {
	var svg strings.Builder
	svg.WriteString(`<svg xmlns="http://www.w3.org/2000/svg" version="1.1" width="200" height="200">`)
	for y := 0; y < matrix.GetHeight(); y++ {
		for x := 0; x < matrix.GetWidth(); x++ {
			if matrix.Get(x, y) {
				svg.WriteString(fmt.Sprintf(`<rect x="%d" y="%d" width="1" height="1" fill="#000000"/>`, x, y))
			}
		}
	}
	svg.WriteString(`</svg>`)
	return svg.String()
}
