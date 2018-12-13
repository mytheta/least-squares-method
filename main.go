package main

import (
	"image/color"
	"math/rand"
	"time"

	"gonum.org/v1/plot/vg"

	"github.com/skelterjohn/go.matrix"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func main() {

	var class []float64

	Dot := matrix.Product
	Inv := matrix.Inverse

	// 図の生成
	p, err := plot.New()
	if err != nil {
		panic(err)
	}

	//任意の点
	dots := make(plotter.XYs, 2)

	//クラス1
	x1, y1 := 8.0, 2.0
	dots[0].X = x1
	dots[0].Y = y1

	//クラス2
	x2, y2 := 3.0, 6.0
	dots[1].X = x2
	dots[1].Y = y2

	//各クラスのサンプル
	n := 100
	class1, plotdata1 := randomPoints(n, x1, y1)
	class2, plotdata2 := randomPoints(n, x2, y2)
	class = append(class1, class2...)
	mClass := matrix.MakeDenseMatrix(class, n*2, 3) // 配列から行列に変換

	//教師データ作成
	b := append(train(1, n), train(-1, n)...)
	mb := matrix.MakeDenseMatrix(b, n*2, 1)

	tClass := mClass.Transpose() // 転置行列

	w := Dot(tClass, mClass)
	w = Inv(w).DenseMatrix()
	w = Dot(w, tClass)
	w = Dot(w, mb)

	a1 := w.Get(0, 0)
	b1 := w.Get(1, 0)
	c1 := w.Get(2, 0)

	//最終境界線のplot--------------------------------------------------------
	lastBorder := plotter.NewFunction(func(x float64) float64 {
		//x2 = -(w1 / w2)*x1 - w0 / w2
		return -(b1/c1)*x - (a1 / c1)
	})
	lastBorder.Color = color.RGBA{B: 255, A: 255}
	//----------------------------------------------------------------------

	//label
	p.Title.Text = "Points Example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// Draw a grid behind the data
	p.Add(plotter.NewGrid())

	// Make a scatter plotter and set its style.
	s, err := plotter.NewScatter(plotdata1)
	if err != nil {
		panic(err)
	}

	y, err := plotter.NewScatter(plotdata2)
	if err != nil {
		panic(err)
	}

	r, err := plotter.NewScatter(dots)
	if err != nil {
		panic(err)
	}

	s.GlyphStyle.Color = color.RGBA{R: 255, B: 128, A: 55}
	y.GlyphStyle.Color = color.RGBA{R: 155, B: 128, A: 255}
	r.GlyphStyle.Color = color.RGBA{R: 128, B: 0, A: 0}
	p.Add(s)
	p.Add(y)
	p.Add(r)
	p.Add(lastBorder)
	p.Legend.Add("class1", s)
	p.Legend.Add("class2", y)

	// Axis ranges
	p.X.Min = 0
	p.X.Max = 10
	p.Y.Min = 0
	p.Y.Max = 10

	// Save the plot to a PNG file.
	if err := p.Save(6*vg.Inch, 6*vg.Inch, "report.png"); err != nil {
		panic(err)
	}
}

//ガウス分布
func random(axis float64) float64 {
	//分散
	dispersion := 1.0
	rand.Seed(time.Now().UnixNano())
	return rand.NormFloat64()*dispersion + axis
}

//学習データの生成
func randomPoints(n int, x, y float64) ([]float64, plotter.XYs) {
	matrix := make([][]float64, n)
	pts := make(plotter.XYs, n)
	var gyo []float64

	for i := range matrix {
		l := random(x)
		m := random(y)
		gyo = append(gyo, 1.0)
		gyo = append(gyo, l)
		gyo = append(gyo, m)
		pts[i].X = l
		pts[i].Y = m
	}
	return gyo, pts
}

//学習
func train(index, n int) []float64 {
	var array []float64

	for i := 0; i < n; i++ {
		array = append(array, float64(index))
	}

	return array
}
