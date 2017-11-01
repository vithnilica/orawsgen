package orawsgen

import "testing"

func testHash(hc int32, result int32,t *testing.T){
	x:=hash(hc)
	if(x!=result){
		t.Error("testHash", hc, result, x)
	}
}

func TestSomeHash(t *testing.T) {
	testHash(1,-8139033,t)
	testHash(-1,8662,t)
	testHash(123,-803083,t)
	testHash(-2147483648,2143418240,t)
	testHash(2147483647,-2147341994,t)
}


func testHashCode(s string, result int32,t *testing.T){
	hc:=hashCode(s)
	if(hc!=result){
		t.Error("testHashCode", s, result, hc)
	}
}

func TestHashCodes(t *testing.T){
	testHashCode("",0, t)
	testHashCode("a",97, t)
	testHashCode("Ahoj",2039906, t)
	testHashCode("Lhfsdilvbdsbds;dolhqMX",-1418708814, t)
	testHashCode("ccxcz",94496399, t)
	testHashCode("ZZZZ",2770560, t)
	testHashCode("zxxcvvb_fdf",1419738242, t)
}






