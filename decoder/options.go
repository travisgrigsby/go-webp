// The MIT License (MIT)
//
// Copyright (c) 2019 Amangeldy Kadyl
//
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the "Software"), to deal in
// the Software without restriction, including without limitation the rights to
// use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
// the Software, and to permit persons to whom the Software is furnished to do so,
// subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
// FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
// COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
// IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
// CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package decoder

/*
#cgo LDFLAGS: -lwebp
#include <webp/decode.h>
*/
import "C"
import (
	"errors"
	"image"
)

// Options specifies webp decoding parameters
type Options struct {
	BypassFiltering        bool
	NoFancyUpsampling      bool
	Crop                   image.Rectangle
	Scale                  image.Rectangle
	UseThreads             bool
	Flip                   bool
	DitheringStrength      int
	AlphaDitheringStrength int

	// These two are optimizations that require a little extra work on the caller side.

	// if nil, DefaultImageFactory will be used.  If non-nil, decode will return an image that must be put back into the pool
	// when you're done with it
	ImageFactory ImageFactory
	// if nil, a default buffer will be used.  If non-nil, decode will use this buffer to store data from the reader.
	// The idea is that this buffer be reused, so either pass this back in next time you call decode, or put it back into
	// a pool when you're done with it.
	Buffer []byte
}

// GetConfig build WebPDecoderConfig for libwebp
func (o *Options) GetConfig() (*C.WebPDecoderConfig, error) {
	config := C.WebPDecoderConfig{}

	if C.WebPInitDecoderConfig(&config) == 0 {
		return nil, errors.New("cannot init decoder config")
	}

	if o.BypassFiltering {
		config.options.bypass_filtering = 1
	}

	if o.NoFancyUpsampling {
		config.options.no_fancy_upsampling = 1
	}

	// проверяем надо ли кропнуть
	if o.Crop.Max.X > 0 && o.Crop.Max.Y > 0 {
		config.options.use_cropping = 1
		config.options.crop_left = C.int(o.Crop.Min.X)
		config.options.crop_top = C.int(o.Crop.Min.Y)
		config.options.crop_width = C.int(o.Crop.Max.X)
		config.options.crop_height = C.int(o.Crop.Max.Y)
	}

	// проверяем надо ли заскейлить
	if o.Scale.Max.X > 0 && o.Scale.Max.Y > 0 {
		config.options.use_scaling = 1
		config.options.scaled_width = C.int(o.Scale.Max.X)
		config.options.scaled_height = C.int(o.Scale.Max.Y)
	}

	if o.UseThreads {
		config.options.use_threads = 1
	}

	config.options.dithering_strength = C.int(o.DitheringStrength)

	if o.Flip {
		config.options.flip = 1
	}

	config.options.alpha_dithering_strength = C.int(o.AlphaDitheringStrength)

	return &config, nil
}

type ImageFactory interface {
	Get(width, height int) *image.NRGBA
}

type DefaultImageFactory struct{}

func (d *DefaultImageFactory) Get(width, height int) *image.NRGBA {
	return image.NewNRGBA(image.Rect(0, 0, width, height))
}
