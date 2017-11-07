package main

import (
	"testing"
)

func testHash(hc int32, result int32, t *testing.T) {
	x := hash15(hc)
	if (x != result) {
		t.Error("testHash", hc, result, x)
	}
}

func TestSomeHash(t *testing.T) {
	testHash(1, -8139033, t)
	testHash(-1, 8662, t)
	testHash(123, -803083, t)
	testHash(-2147483648, 2143418240, t)
	testHash(2147483647, -2147341994, t)

	testHash(378740719, -177873679, t)
	testHash(-824446717, -2034951654, t)
	testHash(-245393993, 1435927871, t)
	testHash(1402489406, 1394845835, t)
	testHash(3079315, -979621435, t)
	testHash(-982148929, 2124480805, t)
	testHash(208254106, -919277600, t)
	testHash(1246571290, -1357702441, t)
	testHash(772425373, -1317033193, t)
	testHash(1245631111, -1771047320, t)
	testHash(1763629341, -501754279, t)
	testHash(-141154685, 2146767548, t)

}

func testHashCode(s string, result int32, t *testing.T) {
	hc := hashCodeString(s)
	if (hc != result) {
		t.Error("testHashCode", s, result, hc)
	}
}

func TestHashCodes(t *testing.T) {
	testHashCode("", 0, t)
	testHashCode("a", 97, t)
	testHashCode("Ahoj", 2039906, t)
	testHashCode("Lhfsdilvbdsbds;dolhqMX", -1418708814, t)
	testHashCode("ccxcz", 94496399, t)
	testHashCode("ZZZZ", 2770560, t)
	testHashCode("zxxcvvb_fdf", 1419738242, t)

	testHashCode("accountBalance", 378740719, t)
	testHashCode("bookingFrequency", -824446717, t)
	testHashCode("contractEndDate", -245393993, t)
	testHashCode("contractStartDate", 1402489406, t)
	testHashCode("debt", 3079315, t)
	testHashCode("freqPremium", -982148929, t)
	testHashCode("freqPremiumEmployer", 208254106, t)
	testHashCode("lastAnniversary", 1246571290, t)
	testHashCode("nextAnniversary", 772425373, t)
	testHashCode("paymentMethod", 1245631111, t)
	testHashCode("premiumsBilledTo", 1763629341, t)
	testHashCode("premiumsPaidTo", -141154685, t)
}

func TestHashMapMixer1(t *testing.T) {
	orig := [...]string{"posumannpremOut", "posumfreqpremOut", "popremsxmlOut", "poerrnoOut", "poerrtextOut", "poerrxmlOut"}
	mixed := [...]string{"posumfreqpremOut", "poerrtextOut", "poerrxmlOut", "popremsxmlOut", "poerrnoOut", "posumannpremOut"}

	m := newHashMapJava15()
	for i := range orig {
		m.Put(orig[i], orig[i])
	}
	hashmapmixed := m.toValArray()
	for i := range mixed {
		if mixed[i] != hashmapmixed[i] {
			t.Error("TestHashMapMixer1", i, mixed[i], hashmapmixed[i])
		}
	}
}

func TestHashMapMixer2(t *testing.T) {
	orig := [...]string{"poresultsOut", "poerrnoOut", "poerrtextOut"}
	mixed := [...]string{"poerrtextOut", "poerrnoOut", "poresultsOut"}

	m := newHashMapJava15()
	for i := range orig {
		m.Put(orig[i], orig[i])
	}
	hashmapmixed := m.toValArray()
	for i := range mixed {
		if mixed[i] != hashmapmixed[i] {
			t.Error("TestHashMapMixer2", i, mixed[i], hashmapmixed[i])
		}
	}
}

func TestHashMapMixer3(t *testing.T) {
	orig := [...]string{"accountBalance", "bookingFrequency", "contractEndDate", "contractStartDate", "debt", "freqPremium", "freqPremiumEmployer", "lastAnniversary", "nextAnniversary", "paymentMethod", "premiumsBilledTo", "premiumsPaidTo"}
	mixed := [...]string{"contractEndDate", "premiumsPaidTo", "contractStartDate", "bookingFrequency", "premiumsBilledTo", "paymentMethod", "nextAnniversary", "lastAnniversary", "freqPremium", "debt", "accountBalance", "freqPremiumEmployer"}

	m := newHashMapJava15()
	for i := range orig {
		m.Put(orig[i], orig[i])
	}
	hashmapmixed := m.toValArray()
	for i := range mixed {
		if mixed[i] != hashmapmixed[i] {
			t.Error("TestHashMapMixer3", i, mixed[i], hashmapmixed[i])
		}
	}
}

func TestMixerTypeItems1(t *testing.T) {
	orig := [...]string{"freqPremium", "freqPremiumEmployer", "bookingFrequency", "paymentMethod", "contractStartDate", "contractEndDate", "premiumsPaidTo", "premiumsBilledTo", "lastAnniversary", "nextAnniversary", "debt", "accountBalance"}
	mixed := [...]string{"contractEndDate", "premiumsPaidTo", "contractStartDate", "bookingFrequency", "premiumsBilledTo", "paymentMethod", "nextAnniversary", "lastAnniversary", "freqPremium", "debt", "accountBalance", "freqPremiumEmployer"}

	var parr15 [][2]string
	for i := range orig {
		parr15 = append(parr15, [2]string{orig[i], orig[i]})
	}

	hashmapmixed := MixerJava15(parr15)

	for i := range mixed {
		if mixed[i] != hashmapmixed[i] {
			t.Error("TestMixerTypeItems1", i, mixed[i], hashmapmixed[i])
		}
	}
}

func TestMixerOutParams1(t *testing.T) {
	//poradi podle baliku
	orig := [...]string{"posumannpremOut", "posumfreqpremOut", "popremsxmlOut", "poerrnoOut", "poerrtextOut", "poerrxmlOut"}
	//poradi podle wsdl
	mixed := [...]string{"posumfreqpremOut", "poerrtextOut", "popremsxmlOut", "poerrxmlOut", "poerrnoOut", "posumannpremOut"}

	var parr15 [][2]string
	for i := range orig {
		parr15 = append(parr15, [2]string{orig[i], orig[i]})
	}

	hashmapmixed := MixerJava15(parr15)

	for i := range mixed {
		if mixed[i] != hashmapmixed[i] {
			t.Error("TestMixerOutParams1", i, mixed[i], hashmapmixed[i])
		}
	}
}
