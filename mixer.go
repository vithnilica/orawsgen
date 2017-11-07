package main

import (
	"sort"
	"fmt"
)

type entryJava15 struct {
	key  string
	val  string
	hash int32
	next *entryJava15
}

type hashMapJava15 struct {
	table     []*entryJava15
	size      int
	threshold int
}

func hashCodeString(s string) (int32) {
	if s == "" {
		return 0
	}
	var hc int32 = 0
	for i := 0; i < len(s); i++ {
		hc = 31*hc + int32(s[i])
	}
	return hc;
}

func hash15(hc int32) (int32) {
	var h uint32 = uint32(hc);
	//h += ~(h << 9)
	h += ^(h << 9)
	//h ^=  (h >>> 14)
	h = h ^ (h >> 14)
	h += (h << 4)
	//h ^=  (h >>> 10)
	h = h ^ (h >> 10)
	return int32(h)
}

func indexFor(h int32, length int) (int) {
	return int(h) & (length - 1);
}

func (hm *hashMapJava15) Put(key string, val string) {
	var h int32 = hash15(hashCodeString(key));
	var i int = indexFor(h, len(hm.table))
	//tohle hleda jestli uz to tam neni
	for e := hm.table[i]; e != nil; e = e.next {
		if (e.hash == h) && (e.key == key) {
			e.val = val
			return;
		}
	}
	hm.addEntry(h, key, val, i);
}

func (hm *hashMapJava15) addEntry(hash int32, key string, val string, bucketIndex int) {
	e := hm.table[bucketIndex];
	ne := &entryJava15{key: key, val: val, hash: hash, next: e}
	hm.table[bucketIndex] = ne
	if (hm.size >= hm.threshold) {
		hm.size++
		hm.resize(2 * len(hm.table));
	} else {
		hm.size++
	}
}

func (hm *hashMapJava15) resize(newCapacity int) {
	var newTable []*entryJava15 = make([]*entryJava15, newCapacity, newCapacity)

	//presune data ze stare do nove tabulky
	for j := 0; j < len(hm.table); j++ {
		e := hm.table[j];
		if (e != nil) {
			hm.table[j] = nil;
		}
		for ; e != nil; {
			next := e.next;
			i := indexFor(e.hash, newCapacity);
			e.next = newTable[i];
			newTable[i] = e;
			e = next;
		}
	}

	hm.table = newTable;
	hm.threshold = (int)(float64(newCapacity) * 0.75);
}

func (hm *hashMapJava15) toValArray() ([]string) {
	ret := make([]string, hm.size, hm.size)
	var ii int = 0;
	//zacina od konce pole
	for j := len(hm.table) - 1; j >= 0; j-- {
		e := hm.table[j];
		for ; e != nil; {
			ret[ii] = e.val
			ii++;
			next := e.next;
			e = next;
		}
	}
	return ret;
}

func newHashMapJava15() (*hashMapJava15) {
	var hm *hashMapJava15
	hm = new(hashMapJava15)
	hm.size = 0
	hm.threshold = (int)(float64(16) * 0.75)
	hm.table = make([]*entryJava15, 16, 16)
	return hm
}

type stringArr2Slice [][2]string

func (p stringArr2Slice) Len() int           { return len(p) }
func (p stringArr2Slice) Less(i, j int) bool { return p[i][0] < p[j][0] }
func (p stringArr2Slice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p stringArr2Slice) Sort() { sort.Sort(p) }

func MixerJava15(parr15 [][2]string) ([]string) {
	//ohjebak pro outparametry "setrideny" nejdrive podle abecedy a pak pomoci hashmapy v jave
	m15 := newHashMapJava15()

	//setrideni podle abecedy (v jave to delaj pomoci TreeMap)
	var slice stringArr2Slice = stringArr2Slice(parr15)

	fmt.Printf("faze 1 %v \n", slice)
	sort.Sort(&slice)

	fmt.Printf("faze 2 %v \n", slice)

	//naplneni hash mapy
	for i := range parr15 {
		m15.Put(parr15[i][0], parr15[i][1])
	}
	//vyplivnuti hodnot serazenych tak jak to dela hasmapa v jave 1.5
	return m15.toValArray()
}
