package shapefile

import (
	"github.com/mercatormaps/go-shapefile/dbf"
	"github.com/mercatormaps/go-shapefile/shp"
	"golang.org/x/text/encoding"
)

// Option funcs can be passed to NewScanner().
type Option func(*options)

// PointPrecision sets shp.PointPrecision.
func PointPrecision(p uint) Option {
	return func(o *options) {
		o.shp = append(o.shp, shp.PointPrecision(p))
	}
}

// CharacterEncoding sets dbf.CharacterEncoding.
func CharacterEncoding(enc encoding.Encoding) Option {
	return func(o *options) {
		o.dbf = append(o.dbf, dbf.CharacterEncoding(enc))
	}
}

// FilterFields sets dbf.FilterFields.
func FilterFields(names ...string) Option {
	return func(o *options) {
		o.dbf = append(o.dbf, dbf.FilterFields(names...))
	}
}

// Options for shp and dbf parsing.
type options struct {
	shp []shp.Option
	dbf []dbf.Option
}
