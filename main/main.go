package main

import (
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"
	lab "graphics-lab-1"
	"log"
	"unsafe"
)

const (
	ImpreciseRGBWarn = "Warning: The RGB color can be imprecise due to rounding"
)

var (
	rgbColor  lab.RGB
	cmykColor *lab.CMYK
	hlsColor  *lab.HLS

	rgbChangedProg  bool
	cmykChangedProg bool
	hlsChangedProg  bool

	sample                             *walk.ImageView
	rgbRTE, rgbGTE, rgbBTE             *walk.NumberEdit
	cmykCTE, cmykMTE, cmykYTE, cmykKTE *walk.NumberEdit
	hlsHTE, hlsLTE, hlsSTE             *walk.NumberEdit
	rgbRS, rgbGS, rgbBS                *walk.Slider
	cmykCS, cmykMS, cmykYS, cmykKS     *walk.Slider
	hlsHS, hlsLS, hlsSS                *walk.Slider
	mainSB                             *walk.StatusBarItem

	chooseColor = win.CHOOSECOLOR{
		LStructSize:  uint32(unsafe.Sizeof(win.CHOOSECOLOR{})),
		LpCustColors: &[16]win.COLORREF{},
		Flags:        win.CC_ANYCOLOR | win.CC_FULLOPEN | win.CC_RGBINIT,
	}
)

func updateColors(color interface{}) {
	rgbColor = lab.ToRGB(color)
	cmykColor = lab.ToCMYK(color)
	hlsColor = lab.ToHLS(color)
}

func updateRGBText() {
	_ = rgbRTE.SetValue(float64(rgbColor.R()))
	_ = rgbGTE.SetValue(float64(rgbColor.G()))
	_ = rgbBTE.SetValue(float64(rgbColor.B()))
}

func updateCMYKText() {
	_ = cmykCTE.SetValue(cmykColor.C() * 100)
	_ = cmykMTE.SetValue(cmykColor.M() * 100)
	_ = cmykYTE.SetValue(cmykColor.Y() * 100)
	_ = cmykKTE.SetValue(cmykColor.K() * 100)
}

func updateHLSText() {
	_ = hlsHTE.SetValue(hlsColor.H())
	_ = hlsLTE.SetValue(hlsColor.L() * 100)
	_ = hlsSTE.SetValue(hlsColor.S() * 100)
}

func updateSampleSliders() {
	brush, _ := walk.NewSolidColorBrush(walk.Color(rgbColor))
	sample.SetBackground(brush)

	rgbRS.SetValue(int(rgbColor.R()))
	rgbGS.SetValue(int(rgbColor.G()))
	rgbBS.SetValue(int(rgbColor.B()))

	cmykCS.SetValue(int(cmykColor.C() * 100))
	cmykMS.SetValue(int(cmykColor.M() * 100))
	cmykYS.SetValue(int(cmykColor.Y() * 100))
	cmykKS.SetValue(int(cmykColor.K() * 100))

	hlsHS.SetValue(int(hlsColor.H()))
	hlsLS.SetValue(int(hlsColor.L() * 100))
	hlsSS.SetValue(int(hlsColor.S() * 100))
}

func updateAll() {
	updateRGBText()
	updateCMYKText()
	updateHLSText()
	updateSampleSliders()
}

func progChangeGuard(handler func()) {
	rgbChangedProg = true
	cmykChangedProg = true
	hlsChangedProg = true

	handler()

	rgbChangedProg = false
	cmykChangedProg = false
	hlsChangedProg = false
}

func fetchRGBFromInputs() {
	if rgbChangedProg {
		return
	}
	progChangeGuard(func() {
		newRGB := lab.NewRGB(byte(rgbRTE.Value()), byte(rgbGTE.Value()), byte(rgbBTE.Value()))
		updateColors(newRGB)
		updateCMYKText()
		updateHLSText()
		updateSampleSliders()
		_ = mainSB.SetText("")
	})
}

func fetchRGBFromSlider() {
	if rgbChangedProg {
		return
	}
	progChangeGuard(func() {
		newRGB := lab.NewRGB(byte(rgbRS.Value()), byte(rgbGS.Value()), byte(rgbBS.Value()))
		updateColors(newRGB)
		updateAll()
		_ = mainSB.SetText("")
	})
}

func fetchCMYKFromInputs() {
	if cmykChangedProg {
		return
	}
	progChangeGuard(func() {
		newCMYK := lab.NewCMYK(cmykCTE.Value()/100, cmykMTE.Value()/100, cmykYTE.Value()/100, cmykKTE.Value()/100)
		updateColors(newCMYK)
		updateRGBText()
		updateHLSText()
		updateSampleSliders()
		_ = mainSB.SetText(ImpreciseRGBWarn)
	})
}

func fetchCMYKFromSlider() {
	if cmykChangedProg {
		return
	}
	progChangeGuard(func() {
		newCMYK := lab.NewCMYK(float64(cmykCS.Value())/100,
			float64(cmykMS.Value())/100,
			float64(cmykYS.Value())/100,
			float64(cmykKS.Value())/100)
		updateColors(newCMYK)
		updateAll()
		_ = mainSB.SetText(ImpreciseRGBWarn)
	})
}

func fetchHLSFromInputs() {
	if hlsChangedProg {
		return
	}
	progChangeGuard(func() {
		newHLS := lab.NewHLS(hlsHTE.Value(), hlsLTE.Value()/100, hlsSTE.Value()/100)
		updateColors(newHLS)
		updateRGBText()
		updateCMYKText()
		updateSampleSliders()
		_ = mainSB.SetText(ImpreciseRGBWarn)
	})
}

func fetchHLSFromSlider() {
	if hlsChangedProg {
		return
	}
	progChangeGuard(func() {
		newHLS := lab.NewHLS(float64(hlsHS.Value()), float64(hlsLS.Value())/100, float64(hlsSS.Value())/100)
		updateColors(newHLS)
		updateAll()
		_ = mainSB.SetText(ImpreciseRGBWarn)
	})
}

func main() {
	progChangeGuard(func() {
		updateColors(rgbColor)
	})

	mw := MainWindow{
		Title:  "Computer Graphics Lab 1",
		Size:   Size{Width: 600, Height: 400},
		Layout: VBox{},

		Children: []Widget{
			VSplitter{
				Children: []Widget{
					GroupBox{
						Title:  "Chosen color and color picking (via palette)",
						Layout: HBox{},
						Children: []Widget{
							Label{Text: "Color sample (click it to open the palette):"},
							ImageView{
								AssignTo:   &sample,
								Background: SolidColorBrush{Color: walk.RGB(0, 0, 0)},
								MinSize:    Size{Width: 100, Height: 100},
								MaxSize:    Size{Width: 100, Height: 100},
								OnMouseDown: func(x, y int, button walk.MouseButton) {
									if button == walk.LeftButton {
										chooseColor.RgbResult = win.COLORREF(rgbColor)
										if win.ChooseColor(&chooseColor) {
											progChangeGuard(func() {
												updateColors(lab.RGB(chooseColor.RgbResult))
												updateAll()
											})
										}
									}
								},
							},
						},
					},
					GroupBox{
						Title:  "Color picking using different models",
						Layout: HBox{},
						Children: []Widget{
							GroupBox{
								Title:  "RGB",
								Layout: Grid{Columns: 2},
								Children: []Widget{
									Label{Text: "Red"},
									NumberEdit{
										AssignTo:       &rgbRTE,
										MinValue:       0,
										MaxValue:       255,
										Value:          float64(rgbColor.R()),
										Decimals:       0,
										OnValueChanged: fetchRGBFromInputs,
									},
									Slider{
										AssignTo:       &rgbRS,
										ColumnSpan:     2,
										MinValue:       0,
										MaxValue:       255,
										Value:          rgbColor.R(),
										OnValueChanged: fetchRGBFromSlider,
									},

									Label{Text: "Green"},
									NumberEdit{
										AssignTo:       &rgbGTE,
										MinValue:       0,
										MaxValue:       255,
										Value:          float64(rgbColor.G()),
										Decimals:       0,
										OnValueChanged: fetchRGBFromInputs,
									},
									Slider{
										AssignTo:       &rgbGS,
										ColumnSpan:     2,
										MinValue:       0,
										MaxValue:       255,
										Value:          rgbColor.G(),
										OnValueChanged: fetchRGBFromSlider,
									},

									Label{Text: "Blue"},
									NumberEdit{
										AssignTo:       &rgbBTE,
										MinValue:       0,
										MaxValue:       255,
										Value:          float64(rgbColor.B()),
										Decimals:       0,
										OnValueChanged: fetchRGBFromInputs,
									},
									Slider{
										AssignTo:       &rgbBS,
										ColumnSpan:     2,
										MinValue:       0,
										MaxValue:       255,
										Value:          rgbColor.B(),
										OnValueChanged: fetchRGBFromSlider,
									},
								},
							},
							GroupBox{
								Title:  "CMYK",
								Layout: Grid{Columns: 2},
								Children: []Widget{
									Label{Text: "Cyan (%)"},
									NumberEdit{
										AssignTo:       &cmykCTE,
										MinValue:       0,
										MaxValue:       100,
										Value:          cmykColor.C() * 100,
										Decimals:       5,
										OnValueChanged: fetchCMYKFromInputs,
									},
									Slider{
										AssignTo:       &cmykCS,
										ColumnSpan:     2,
										MinValue:       0,
										MaxValue:       100,
										Value:          int(cmykColor.C() * 100),
										OnValueChanged: fetchCMYKFromSlider,
									},

									Label{Text: "Magenta (%)"},
									NumberEdit{
										AssignTo:       &cmykMTE,
										MinValue:       0,
										MaxValue:       100,
										Value:          cmykColor.M() * 100,
										Decimals:       5,
										OnValueChanged: fetchCMYKFromInputs,
									},
									Slider{
										AssignTo:       &cmykMS,
										ColumnSpan:     2,
										MinValue:       0,
										MaxValue:       100,
										Value:          int(cmykColor.M() * 100),
										OnValueChanged: fetchCMYKFromSlider,
									},

									Label{Text: "Yellow (%)"},
									NumberEdit{
										AssignTo:       &cmykYTE,
										MinValue:       0,
										MaxValue:       100,
										Value:          cmykColor.Y() * 100,
										Decimals:       5,
										OnValueChanged: fetchCMYKFromInputs,
									},
									Slider{
										AssignTo:       &cmykYS,
										ColumnSpan:     2,
										MinValue:       0,
										MaxValue:       100,
										Value:          int(cmykColor.Y() * 100),
										OnValueChanged: fetchCMYKFromSlider,
									},

									Label{Text: "Key (%)"},
									NumberEdit{
										AssignTo:       &cmykKTE,
										MinValue:       0,
										MaxValue:       100,
										Value:          cmykColor.K() * 100,
										Decimals:       5,
										OnValueChanged: fetchCMYKFromInputs,
									},
									Slider{
										AssignTo:       &cmykKS,
										ColumnSpan:     2,
										MinValue:       0,
										MaxValue:       100,
										Value:          int(cmykColor.K() * 100),
										OnValueChanged: fetchCMYKFromSlider,
									},
								},
							},
							GroupBox{
								Title:  "HLS",
								Layout: Grid{Columns: 2},
								Children: []Widget{
									Label{Text: "Hue (˚)"},
									NumberEdit{
										AssignTo:       &hlsHTE,
										MinValue:       0,
										MaxValue:       360,
										Value:          hlsColor.H(),
										Decimals:       5,
										OnValueChanged: fetchHLSFromInputs,
									},
									Slider{
										AssignTo:       &hlsHS,
										ColumnSpan:     2,
										MinValue:       0,
										MaxValue:       360,
										Value:          int(hlsColor.H()),
										OnValueChanged: fetchHLSFromSlider,
									},

									Label{Text: "Lightness (%)"},
									NumberEdit{
										AssignTo:       &hlsLTE,
										MinValue:       0,
										MaxValue:       100,
										Decimals:       5,
										Value:          hlsColor.L() * 100,
										OnValueChanged: fetchHLSFromInputs,
									},
									Slider{
										AssignTo:       &hlsLS,
										ColumnSpan:     2,
										MinValue:       0,
										MaxValue:       100,
										Value:          int(hlsColor.L() * 100),
										OnValueChanged: fetchHLSFromSlider,
									},

									Label{Text: "Saturation (%)"},
									NumberEdit{
										AssignTo:       &hlsSTE,
										MinValue:       0,
										MaxValue:       100,
										Value:          hlsColor.S() * 100,
										Decimals:       5,
										OnValueChanged: fetchHLSFromInputs,
									},
									Slider{
										AssignTo:       &hlsSS,
										ColumnSpan:     2,
										MinValue:       0,
										MaxValue:       100,
										Value:          int(hlsColor.S() * 100),
										OnValueChanged: fetchHLSFromSlider,
									},
								},
							},
						},
					},
				},
			},
		},

		StatusBarItems: []StatusBarItem{
			{
				AssignTo: &mainSB,
			},
		},
	}

	if _, err := mw.Run(); err != nil {
		log.Fatal(err)
	}
}
