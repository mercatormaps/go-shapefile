package dbf

import (
	"fmt"
	"io"
	"sync"

	"github.com/mercatormaps/go-shapefile/dbf/dbase5"
	"github.com/pkg/errors"
)

type Version uint

const (
	DBaseLevel5 Version = 3
	DBaseLevel7 Version = 4
)

type Scanner struct {
	in io.Reader

	versionOnce sync.Once
	version     Version

	headerOnce sync.Once
	header     Header

	scanOnce  sync.Once
	recordsCh chan *Record
	num       uint32

	errOnce sync.Once
	err     error
}

type Header interface {
	RecordLen() uint16
	NumRecords() uint32
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{
		in:        r,
		recordsCh: make(chan *Record),
	}
}

func (s *Scanner) Version() (Version, error) {
	var err error
	s.versionOnce.Do(func() {
		buf := make([]byte, 1)
		var n int
		if n, err = s.in.Read(buf); err != nil {
			return
		} else if n != len(buf) {
			err = fmt.Errorf("read %d bytes but expecting %d", n, len(buf))
			return
		}

		// dBase version number is first 3 bits
		s.version = Version(((buf[0]>>0)&1)<<0 | ((buf[0]>>1)&1)<<1 | ((buf[0]>>2)&1)<<2)
	})
	return s.version, err
}

func (s *Scanner) Header() (Header, error) {
	var err error
	if _, err = s.Version(); err != nil {
		return nil, errors.Wrap(err, "failed to parse version number")
	}

	s.headerOnce.Do(func() {
		switch s.version {
		case DBaseLevel5:
			s.header, err = dbase5.DecodeHeader(s.in)
		case DBaseLevel7:
			err = fmt.Errorf("dBase Level 7 is not supported")
		default:
			err = fmt.Errorf("unsupported version")
		}
	})
	return s.header, err
}

func (s *Scanner) Scan() error {
	if _, err := s.Header(); err != nil {
		return errors.Wrap(err, "failed to parse header")
	}

	s.scanOnce.Do(func() {
		go func() {
			defer close(s.recordsCh)

			for s.num < s.header.NumRecords() {
				rec, err := s.record()
				if err == io.EOF {
					s.setErr(fmt.Errorf("unexpected end of file"))
					return
				} else if err != nil {
					s.setErr(err)
					return
				}
				s.decodeRecord(rec)
			}

			buf := make([]byte, 1)
			if n, err := s.in.Read(buf); err != nil {
				s.setErr(fmt.Errorf("unexpected end of file"))
				return
			} else if n != len(buf) {
				s.setErr(fmt.Errorf("read %d bytes but expecting %d", n, len(buf)))
				return
			}

			if buf[0] != 0x1A {
				s.setErr(fmt.Errorf("missing file terminator"))
			}
		}()
	})
	return nil
}

func (s *Scanner) Record() *Record {
	rec, ok := <-s.recordsCh
	if !ok {
		return nil
	}
	return rec
}

func (s *Scanner) Err() error {
	return s.err
}

func (s *Scanner) decodeRecord(buf []byte) {
	s.recordsCh <- &Record{}
	s.num++
}

func (s *Scanner) record() ([]byte, error) {
	buf := make([]byte, s.header.RecordLen())
	n, err := s.in.Read(buf)
	if err == io.EOF && s.num != (s.header.NumRecords()-1) {
		return nil, NewError(
			fmt.Errorf("unexpected end of file: read %d records but expecting %d", s.num, s.header.NumRecords()),
			s.num)
	} else if err != nil {
		return nil, NewError(err, s.num)
	} else if n != len(buf) {
		return nil, NewError(fmt.Errorf("read %d bytes but expecting %d", n, len(buf)), s.num)
	}
	return buf, nil
}

func (s *Scanner) setErr(err error) {
	s.errOnce.Do(func() {
		s.err = err
	})
}
