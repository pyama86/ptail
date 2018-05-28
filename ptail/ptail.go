package ptail

import (
	"bytes"
	"os"

	"golang.org/x/sync/errgroup"
)

const BUFSIZE = 4096

type middlewareFunc func([]byte) error
type Ptail struct {
	fileName string
	readLine int
	funcs    []middlewareFunc
}

func NewPtail(f string, r int) Ptail {
	return Ptail{
		fileName: f,
		readLine: r,
		funcs:    []middlewareFunc{},
	}
}

func (p *Ptail) Use(f middlewareFunc) {
	p.funcs = append(p.funcs, f)
}

func (p *Ptail) Execute() error {
	read := 0
	fp, err := os.Open(p.fileName)
	if err != nil {
		return err
	}
	defer fp.Close()

	s, err := fp.Stat()
	if err != nil {
		return err
	}

	start := int(s.Size() - BUFSIZE)
	sep := []byte("\n")
	buf := make([]byte, BUFSIZE)

	if start < 0 {
		start = 0
		buf = make([]byte, s.Size())
	}

	var pos int
	var eg errgroup.Group
loop:
	for {
		_, err := fp.ReadAt(buf, int64(start))
		if err != nil {
			return err
		}

		// 最初に改行が見つかる位置以降の行を処理する
		firstSep := bytes.Index(buf, sep)
		if start == 0 {
			pos = 0
		} else {
			pos = firstSep + 1
		}

		for _, l := range bytes.Split(buf[pos:], sep) {
			eg.Go(func() error {
				for _, f := range p.funcs {
					if err := f(l); err != nil {
						return err
					}
				}
				return nil
			})

			read++
			if p.readLine <= read {
				break loop
			}
		}

		// ファイルの先頭
		if start == 0 {
			break
		}

		// 最初の改行までの位置を処理していないので、次回処理する
		start -= (BUFSIZE - firstSep)
		buf = make([]byte, BUFSIZE)
		if start < 0 {
			buf = make([]byte, BUFSIZE+start)
			start = 0
		}
	}
	if err := eg.Wait(); err != nil {
		return err
	}

	return nil
}
