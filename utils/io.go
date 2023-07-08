package utils

import "io"

// 写入全部
func WriteAll(w io.Writer, buf []byte) error {
	for len(buf) > 0 {
		n, err := w.Write(buf)
		if nil != err {
			return err
		}

		buf = buf[n:]
	}

	return nil
}

func IoSwap(src io.ReadWriteCloser, dst io.ReadWriteCloser) {
	go func() {
		defer dst.Close()

		io.Copy(dst, src)
	}()

	io.Copy(src, dst)
}