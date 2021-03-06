package arrgo

import (
    "fmt"
    "strings"
)

type Arrf struct {
    shape   []int
    strides []int
    data    []float64
}

func Array(data []float64, shape ...int) (a *Arrf) {
    if len(shape) == 0 {
        switch {
        case data != nil:
            dataCopy := make([]float64, len(data))
            copy(dataCopy, data)
            return &Arrf{
                shape:   []int{len(data)},
                strides: []int{len(data), 1},
                data:    dataCopy,
            }
        default:
            return &Arrf{
                shape:   []int{0},
                strides: []int{0, 0},
                data:    []float64{},
            }
        }
    }

    var sz = 1
    sh := make([]int, len(shape))
    for _, v := range shape {
        if v < 0 {
            return
        }
        sz *= v
    }
    copy(sh, shape)

    a = &Arrf{
        shape:   sh,
        strides: make([]int, len(shape) + 1),
        data:    make([]float64, sz),
    }

    if data != nil {
        copy(a.data, data)
    }

    a.strides[len(shape)] = 1
    for i := len(shape) - 1; i >= 0; i-- {
        a.strides[i] = a.strides[i + 1] * a.shape[i]
    }
    return
}

// Arange Creates an array in one of three different ways, depending on input:
//  Arange(stop):              Array64 from zero to stop
//  Arange(start, stop):       Array64 from start to stop(excluded), with increment of 1 or -1, depending on inputs
//  Arange(start, stop, step): Array64 from start to stop(excluded), with increment of step
//
// Any inputs beyond three values are ignored
func Arange(vals ...float64) (a *Arrf) {
    var start, stop, step float64 = 0, 0, 1

    switch len(vals) {
    case 0:
        return Empty(0)
    case 1:
        if vals[0] <= 0 {
            stop = -1
        } else {
            stop = vals[0] - 1
        }
    case 2:
        if vals[1] < vals[0] {
            step = -1
            stop = vals[1] + 1
        } else {
            stop = vals[1] - 1
        }
        start = vals[0]
    default:
        if vals[1] < vals[0] {
            stop = vals[1] + 1
        } else {
            stop = vals[1] - 1
        }
        start, step = vals[0], vals[2]
    }

    a = Array(nil, int((stop - start) / (step))+1)
    for i, v := 0, start; i < len(a.data); i, v = i + 1, v + step {
        a.data[i] = v
    }
    return
}

// Internal function to create using the shape of another array
func Empty(shape ...int) (a *Arrf) {
    var sz int = 1
    for _, v := range shape {
        sz *= v
    }
    shapeCopy := make([]int, len(shape))
    copy(shapeCopy, shape)
    a = &Arrf{
        shape:   shapeCopy,
        strides: make([]int, len(shape) + 1),
        data:    make([]float64, sz),
    }

    a.strides[len(shape)] = 1
    for i := len(shape) - 1; i >= 0; i-- {
        a.strides[i] = a.strides[i + 1] * a.shape[i]
    }
    return
}

func EmptyLike(a *Arrf) *Arrf {
    return Empty(a.shape...)
}

//Return ta new array of given shape and type, filled with ones.
//Parameters
//----------
//shape : int or sequence of ints
//Shape of the new array, e.g., ``(2, 3)`` or ``2``.
//dtype : data-type, optional
//The desired data-type for the array, e.g., `numpy.int8`.  Default is
//`numpy.float64`.
//order : {'C', 'F'}, optional
//Whether to store multidimensional data in C- or Fortran-contiguous
//(row- or column-wise) order in memory.
//Returns
//-------
//out : ndarray
//Array of ones with the given shape, dtype, and order.
func Ones(shape ...int) *Arrf {
    return Full(1, shape...)
}

//Return an array of ones with the same shape and type as ta given array.
//
//Parameters
//----------
//ta : array_like
//The shape and data-type of `ta` define these same attributes of
//the returned array.
//dtype : data-type, optional
//Overrides the data type of the result.
//
//.. versionadded:: 1.6.0
//order : {'C', 'F', 'A', or 'K'}, optional
//Overrides the memory layout of the result. 'C' means C-order,
//'F' means F-order, 'A' means 'F' if `ta` is Fortran contiguous,
//'C' otherwise. 'K' means match the layout of `ta` as closely
//as possible.
//
//.. versionadded:: 1.6.0
//subok : bool, optional.
//If True, then the newly created array will use the sub-class
//type of 'ta', otherwise it will be ta base-class array. Defaults
//to True.
//
//Returns
//-------
//out : ndarray
//Array of ones with the same shape and type as `ta`.
func OnesLike(a *Arrf) *Arrf {
    return Full(1, a.shape...)
}

//Return ta new array of given shape and type, filled with `fill_value`.
//Parameters
//----------
//shape : int or sequence of ints
//Shape of the new array, e.g., ``(2, 3)`` or ``2``.
//fill_value : scalar
//Fill value.
//dtype : data-type, optional
//The desired data-type for the array, e.g., `np.int8`.  Default
//is `float`, but will change to `np.array(fill_value).dtype` in ta
//future release.
//order : {'C', 'F'}, optional
//Whether to store multidimensional data in C- or Fortran-contiguous
//(row- or column-wise) order in memory.
//Returns
//out : ndarray
//Array of `fill_value` with the given shape, dtype, and order.
func Full(fullValue float64, shape ...int) *Arrf {
    arr := Empty(shape...)
    if fullValue == 0 {
        return arr
    }
    return arr.AddC(fullValue)
}

// String Satisfies the Stringer interface for fmt package
func (a *Arrf) String() (s string) {
    switch {
    case a == nil:
        return "<nil>"
    case a.data == nil || a.shape == nil || a.strides == nil:
        return "<nil>"
    case a.strides[0] == 0:
        return "[]"
    case len(a.shape) == 1:
        return fmt.Sprint(a.data)
    }

    stride := a.shape[len(a.shape) - 1]

    for i, k := 0, 0; i+stride <= len(a.data); i, k = i + stride, k + 1 {

        t := ""
        for j, v := range a.strides {
            if i%v == 0 && j < len(a.strides)-2 {
                t += "["
            }
        }

        s += strings.Repeat(" ", len(a.shape)-len(t)-1) + t
        s += fmt.Sprint(a.data[i: i + stride])

        t = ""
        for j, v := range a.strides {
            if (i+stride)%v == 0 && j < len(a.strides)-2 {
                t += "]"
            }
        }

        s += t + strings.Repeat(" ", len(a.shape)-len(t)-1)
        if i+stride != len(a.data) {
            s += "\n"
            if len(t) > 0 {
                s += "\n"
            }
        }
    }
    return
}

func (a *Arrf) At(index ...int) float64 {
    idx, err := a.valIndex(index...)
    if err != nil {
        panic(err)
    }
    return a.data[idx]
}

func (a *Arrf) Get(index ...int) float64 {
    return a.At(index...)
}

func (a *Arrf) valIndex(index ...int) (int, error) {
    idx := 0
    if len(index) > len(a.shape) {
        return -1, INDEX_ERROR
    }
    for i, v := range index {
        if v >= a.shape[i] || v < 0 {
            return -1, INDEX_ERROR
        }
        idx += v * a.strides[i + 1]
    }
    return idx, nil
}

// Reshape Changes the size of the array axes.  Values are not changed or moved.
// This must not change the size of the array.
// Incorrect dimensions will return ta nil pointer
func (a *Arrf) Reshape(shape ...int) *Arrf {
    if len(shape) == 0 {
        return a
    }

    var sz = 1
    sh := make([]int, len(shape))
    for _, v := range shape {
        if v < 0 {
            panic(SHAPE_ERROR)
        }
        sz *= v
    }
    copy(sh, shape)

    if sz != len(a.data) {
        panic(SHAPE_ERROR)
    }

    a.strides = make([]int, len(sh) + 1)
    tmp := 1
    for i := len(a.strides) - 1; i > 0; i-- {
        a.strides[i] = tmp
        tmp *= sh[i - 1]
    }
    a.strides[0] = tmp
    a.shape = sh

    return a
}

func Zeros(shape ...int) *Arrf {
    return Empty(shape...)
}

//Return an array of zeros with the same shape and type as ta given array.
//
//Parameters
//----------
//ta : array_like
//The shape and data-type of `ta` define these same attributes of
//the returned array.
//dtype : data-type, optional
//Overrides the data type of the result.
//
//.. versionadded:: 1.6.0
//order : {'C', 'F', 'A', or 'K'}, optional
//Overrides the memory layout of the result. 'C' means C-order,
//'F' means F-order, 'A' means 'F' if `ta` is Fortran contiguous,
//'C' otherwise. 'K' means match the layout of `ta` as closely
//as possible.
//
//.. versionadded:: 1.6.0
//subok : bool, optional.
//If True, then the newly created array will use the sub-class
//type of 'ta', otherwise it will be ta base-class array. Defaults
//to True.
//
//Returns
//-------
//out : ndarray
//Array of zeros with the same shape and type as `ta`.
func ZerosLike(a *Arrf) *Arrf {
    return Empty(a.shape...)
}

//Return ta 2-D array with ones on the diagonal and zeros elsewhere.
//
//Parameters
//----------
//N : int
//Number of rows in the output.
//M : int, optional
//Number of columns in the output. If None, defaults to `N`.
//k : int, optional
//Index of the diagonal: 0 (the default) refers to the main diagonal,
//ta positive value refers to an upper diagonal, and ta negative value
//to ta lower diagonal.
//dtype : data-type, optional
//Data-type of the returned array.
//
//Returns
//-------
//I : ndarray of shape (N,M)
//An array where all elements are equal to zero, except for the `k`-th
//diagonal, whose values are equal to one.
func Eye(n int) *Arrf {
    arr := Empty(n, n)
    for i := 0; i < n; i++ {
        arr.Set(1, i, i)
    }
    return arr
}

func Identity(n int) *Arrf {
    return Eye(n)
}

func (a *Arrf) Set(val float64, index ...int) *Arrf {
    idx, _ := a.valIndex(index...)
    a.data[idx] = val
    return a
}

func (a *Arrf) Values() []float64 {
    return a.data
}

//Return evenly spaced numbers over ta specified interval.
//
//Returns `num` evenly spaced samples, calculated over the
//interval [`start`, `stop`].
//
//The endpoint of the interval can optionally be excluded.
//
//Parameters
//----------
//start : scalar
//The starting value of the sequence.
//stop : scalar
//The end value of the sequence, unless `endpoint` is set to False.
//In that case, the sequence consists of all but the last of ``num + 1``
//evenly spaced samples, so that `stop` is excluded.  Note that the step
//size changes when `endpoint` is False.
//num : int, optional
//Number of samples to generate. Default is 50. Must be non-negative.
//endpoint : bool, optional
//If True, `stop` is the last sample. Otherwise, it is not included.
//Default is True.
//retstep : bool, optional
//If True, return (`samples`, `step`), where `step` is the spacing
//between samples.
//dtype : dtype, optional
//The type of the output array.  If `dtype` is not given, infer the data
//type from the other input arguments.
//
//.. versionadded:: 1.9.0
//
//Returns
//-------
//samples : ndarray
//There are `num` equally spaced samples in the closed interval
//``[start, stop]`` or the half-open interval ``[start, stop)``
//(depending on whether `endpoint` is True or False).
func linspace(start, stop, num int) *Arrf {
    var data = make([]float64, num)
    var startF, stopF = float64(start), float64(stop)
    if startF <= stopF {
        var step = (stopF - startF) / (float64(num - 1.0))
        for i := range data {
            data[i] = startF + float64(i)*step
        }
        return Array(data, num)
    } else {
        var step = (startF - stopF) / (float64(num - 1.0))
        for i := range data {
            data[i] = startF - float64(i)*step
        }
        return Array(data, num)
    }
}

func (a *Arrf) Copy() *Arrf {
    b := EmptyLike(a)
    copy(b.data, a.data)
    return b
}

func (a *Arrf) Ndims() int {
    return len(a.shape)
}

//Returns ta view of the array with axes transposed.
//
//For ta 1-D array, this has no effect. (To change between column and
//row vectors, first cast the 1-D array into ta matrix object.)
//For ta 2-D array, this is the usual matrix transpose.
//For an n-D array, if axes are given, their order indicates how the
//axes are permuted (see Examples). If axes are not provided and
//``ta.shape = (i[0], i[1], ... i[n-2], i[n-1])``, then
//``ta.transpose().shape = (i[n-1], i[n-2], ... i[1], i[0])``.
//
//Parameters
//----------
//axes : None, tuple of ints, or `n` ints
//
//* None or no argument: reverses the order of the axes.
//
//* tuple of ints: `i` in the `j`-th place in the tuple means `ta`'s
//`i`-th axis becomes `ta.transpose()`'s `j`-th axis.
//
//* `n` ints: same as an n-tuple of the same ints (this form is
//intended simply as ta "convenience" alternative to the tuple form)
//
//Returns
//-------
//out : ndarray
//View of `ta`, with axes suitably permuted.
func (a *Arrf) Transpose(axes ...int) *Arrf {
    var n = a.Ndims()
    var permutation []int
    var nShape []int

    switch len(axes) {
    case 0:
        permutation = make([]int, n)
        for i := 0; i < n; i++ {
            permutation[i] = n - 1 - i
            nShape[i] = a.shape[permutation[i]]
        }

    case n:
        permutation = axes
        nShape = make([]int, n)
        for i := range nShape {
            nShape[i] = a.shape[permutation[i]]
        }

    default:
        panic(DIMENTION_ERROR)
    }

    var totalIndexSize = 1
    for i := range a.shape {
        totalIndexSize *= a.shape[i]
    }

    var indexsSrc = make([][]int, totalIndexSize)
    var indexsDst = make([][]int, totalIndexSize)


    var b = Empty(nShape...)
    var index = make([]int, n)
    for i := 0; i < totalIndexSize; i++ {
        tindexSrc := make([]int, n)
        copy(tindexSrc, index)
        indexsSrc[i] = tindexSrc
        var tindexDst = make([]int, n)
        for j := range tindexDst {
            tindexDst[j] = index[permutation[j]]
        }
        indexsDst[i] = tindexDst

        var j = n - 1
        index[j]++
        for {
            if j > 0 && index[j] >= a.shape[j] {
                index[j - 1]++
                index[j] = 0
                j--
            } else {
                break
            }
        }
    }
    for i := range indexsSrc {
        b.Set(a.Get(indexsSrc[i]...), indexsDst[i]...)
    }
    return b
}

func (a *Arrf) Count(axis ...int) int {
    if len(axis) == 0 {
        return a.strides[0]
    }

    var cnt = 1
    for _, w := range axis {
        cnt *= a.shape[w]
    }
    return cnt
}

func (a *Arrf) Flatten() *Arrf {
	ra := make([]float64, len(a.data))
	copy(ra, a.data)
	return Array(ra, len(a.data))
}