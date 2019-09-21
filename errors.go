// +build go1.13
// Copyright (c) 2019 hankei6km
// Licensed under the MIT License. See LICENSE in the project root.

package ac

// errors.go と errors_compat.go は以下のソースコードの記事を参考に作成.
// [サポートページ：WEB&#43;DB PRESS Vol.112：｜gihyo.jp … 技術評論社](https://gihyo.jp/magazine/wdpress/archive/2019/vol112/support)
//  -「Goに入りては…… ── When In Go...」で使用されたソースコード

// 1.13以前でも実行時エラーにならないようになっているが、
// go version -m が実行できないので、実際には使うことはできない.
import "fmt"

func wrapf(err error, format string, a ...interface{}) error {
	// return fmt.Errorf(format, a...)
	return fmt.Errorf("%s: %w", fmt.Sprintf(format, a...), err)
}
