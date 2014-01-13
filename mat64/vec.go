// Copyright ©2013 The gonum Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mat64

import (
	"github.com/gonum/blas"
)

var (
	vector *Vec

	_ Matrix  = vector
	_ Mutable = vector

	// _ Cloner      = vector
	// _ Viewer      = vector
	// _ Subvectorer = vector

	// _ Adder     = vector
	// _ Suber     = vector
	_ Muler = vector
	// _ Dotter    = vector
	// _ ElemMuler = vector

	// _ Scaler  = vector
	// _ Applyer = vector

	// _ Normer = vector
	// _ Sumer  = vector

	// _ Stacker   = vector
	// _ Augmenter = vector

	// _ Equaler       = vector
	// _ ApproxEqualer = vector

	// _ RawMatrixLoader = vector
	// _ RawMatrixer     = vector
)

type Vec []float64

func (m Vec) At(r, c int) float64 {
	if c != 0 || r < 0 || r >= len(m) {
		panic(ErrIndexOutOfRange)
	}
	return m[r]
}

func (m Vec) Set(r, c int, v float64) {
	if c != 0 || r < 0 || r >= len(m) {
		panic(ErrIndexOutOfRange)
	}
	m[r] = v
}

func (m Vec) Dims() (r, c int) { return len(m), 1 }

func (m *Vec) Mul(a, b Matrix) {
	ar, ac := a.Dims()
	br, bc := b.Dims()

	if ac != br {
		panic(ErrShape)
	}

	var w Vec
	if m != a && m != b {
		w = *m
	}
	if len(w) == 0 {
		w = use(w, ar)
	} else if ar != len(w) || bc != 1 {
		panic(ErrShape)
	}

	bv := *b.(*Vec) // This is a temporary restriction.

	if a, ok := a.(RawMatrixer); ok {
		amat := a.RawMatrix()
		blasEngine.Dgemv(BlasOrder,
			blas.NoTrans,
			ar, ac,
			1.,
			amat.Data, amat.Stride,
			bv, 1,
			0.,
			w, 1)
		*m = w
		return
	}

	if a, ok := a.(Vectorer); ok {
		row := make([]float64, ac)
		for r := 0; r < ar; r++ {
			w[r] = blasEngine.Ddot(ac, a.Row(row, r), 1, bv, 1)
		}
		*m = w
		return
	}

	row := make([]float64, ac)
	for r := 0; r < ar; r++ {
		for i := range row {
			row[i] = a.At(r, i)
		}
		var v float64
		for i, e := range row {
			v += e * bv[i]
		}
		w[r] = v
	}
	*m = w
}

func (m *Vec) Scale(f float64, a Matrix) {
	ar, ac := a.Dims()

	w := *m

	if ac != 1 {
		panic(ErrShape)
	}

	if len(w) == 0 {
		w = make(Vec, ar)
	} else if ar != len(w) {
		panic(ErrShape)
	}

	switch a := a.(type) {
	case *Vec:
		if &w != a {
			copy(w, *a)
		}
		blasEngine.Dscal(ar, f, w, 1)
	case Vec:
		copy(w, a)
		blasEngine.Dscal(ar, f, w, 1)
	default:
		for r := 0; r < ar; r++ {
			w[r] = f * a.At(r, 1)
		}
	}

	*m = w
}

func (m Vec) Copy(a Matrix) (r, c int) {
	r, c = a.Dims()
	r = min(r, len(m))
	if c < 1 {
		return r, c
	}
	c = 1

	switch a := a.(type) {
	case *Vec:
		copy(m, *a)
	case Vec:
		copy(m, a)
	default:
		for i := 0; i < r; i++ {
			m[i] = a.At(i, 1)
		}
	}
	return r, c
}
