/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    word_test
 * @Date:    2021/12/17 6:14 δΈε
 * @package: word
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package word

import (
	"fmt"
	"testing"
)

func Test_CheckWord(t *testing.T) {
	ds := CheckWord("ζ")
	fmt.Println(ds[0].GetRadicals())
}
