/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    wuxing
 * @Date:    2021/12/14 3:49 下午
 * @package: astro
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package astro

/*
十天干：甲、乙、丙、丁、戊、己、庚、辛、壬、癸

阳天干：甲、丙、戊、庚、壬
阴天干：乙、丁、己、辛、癸

天干五行属性：
甲乙：木（东）
丙丁：火（南）
戊己：土（中）
庚辛：金（西）
壬癸：水（北）

天干相合：
甲己：土
乙庚：金
丙辛：水
丁壬：木
戊癸：火

-------------------------------------------------------------------------

十二地支：子、丑、寅、卯、辰、巳、午、未、申、酉、戌、亥

阳地支：子、寅、辰、午、申、戌
阴地支：丑、卯、巳、未、酉、亥

地支五行属性：
亥子：水（北）
寅卯：木（东）
巳午：火（南）
申酉：金（西）
辰戌丑未：土（中）

地支六合：
子丑：土
寅亥：木
戌卯：火
辰酉：金
巳申：水
午未：日月

===========================================================================
*/

/*
天干五合和地支六合是批八字里的一个重要内容。

天干五合是甲己相合为土，乙庚相合为金，丙辛相合为水，丁壬相合为木，戊癸相合为火。
天干五合可以结合年上起月法来记，一下子能记住两个知识要点。
年上起月法为：
甲己之年丙作首，乙庚之岁戊为头。
丙辛之岁寻庚上，丁壬壬寅顺水流。
若问戊癸何处起，甲寅之上好追求。
甲己合，丙做首丙火生土，所以甲己合化为土。
乙庚合，乙庚之岁戊为头，戊土生金，所以乙庚合化金。
丙辛合，丙辛之岁寻庚上，庚金生水，所以丙辛合化水。
丁壬合，丁壬壬寅顺水流，就是第一个月为壬月，水生木，所以丁壬合化木。
戊癸合，甲寅之上好追求，第一月是甲月，甲木生活，所以戊癸合化火。



地支六合为：子丑合化土，寅亥合化木，卯戌合化火，辰酉合化金，巳申合化水，午未合化土。
通过生肖记忆，鼠与牛为合，虎与猪为合，兔与狗为合， 龙与鸡为合， 蛇与猴为合， 马与羊为合。（婚姻六合）
黑鼠黄牛正相合，富贵荣华福禄多。
白虎黑猪上等婚，人口兴旺有精神。
白兔黄狗古来有，金玉满堂乐悠悠。
黄龙黑鸡更相投，福寿绵长永无休。
猴蛇相配满堂红，子孙三代不受穷。
红马黄羊在福中，百年长寿无大凶。
这套口诀正是民间断婚姻常用的。


六合说完，顺道说一下地支六害。
六害为：子未相害、丑午相害、寅巳相害，卯辰相害，申亥相害，酉戌相害。
这段话看似简单，背起来很难，我们在把他们转化为生肖，变成六害歌，同时也是民间断婚姻常说的。
自古白马怕青牛，羊鼠相逢一旦休，蛇遇猛虎如刀断，猪见猿猴不到头，龙逢兔儿云端去，金鸡见犬泪交流。

*/

func wuXingAttr(word string) string {
	switch word {
	case "壬", "癸", "亥", "子":
		return "水"
	case "甲", "乙", "寅", "卯":
		return "木"
	case "丙", "丁", "巳", "午":
		return "火"
	case "庚", "辛", "申", "酉":
		return "金"
	case "戊", "己", "辰", "戌", "丑", "未":
		return "土"
	default:
		return ""
	}
}

func wuXingAttrs(words []string) []string {
	var wuxings []string
	for _, w := range words {
		wuxings = append(wuxings, wuXingAttr(w))
	}
	return wuxings
}

func missWuXing(wuxings []string) []string {
	var all = map[string]bool{"金": true, "木": true, "水": true, "火": true, "土": true}
	var miss []string
	for _, wx := range wuxings {
		if all[wx] {
			delete(all, wx)
		}
	}
	for wx := range all {
		miss = append(miss, wx)
	}
	return miss
}

func direction(attr string) string {
	switch attr {
	case "金":
		return "西"
	case "木":
		return "东"
	case "水":
		return "北"
	case "火":
		return "南"
	case "土":
		return "中"
	default:
		return ""
	}
}

func Combine(words1, words2 []string) {

}
