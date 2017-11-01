package orawsgen

import "fmt"

type entryJava15 struct {
	key string
	val string
	hash int
	next *entryJava15
}

type hashMapJava15 struct {
}



func hashCode(s string) (int32) {
	if s==""{
		return 0
	}
	var hc int32=0
	for i:=0; i<len(s); i++ {
		hc = 31 * hc + int32(s[i])
	}
	return hc;
}

func hash(hc int32)(int32) {
	var h uint32=uint32(hc);
	//h += ~(h << 9)
	h+=^(h<<9)
	//h ^=  (h >>> 14)
	h=h^(h>>14)
	h+=(h << 4)
	//h ^=  (h >>> 10)
	h=h^(h>>10)
	return int32(h)
}


func main() {

	fmt.Println(hash(100))
}

/*


static int indexFor(int h, int length) {
return h & (length-1);
}

public void put(String key, String val) {
int hash = hash(key.hashCode());
int i = indexFor(hash, table.length);
//tohle hleda jestli uz to tam neni
for (Entry e = table[i]; e != null; e = e.next) {
Object k;
if (e.hash == hash && ((k = e.key) == key || key.equals(k))) {
return;
}
}
addEntry(hash, key, val, i);
}

void addEntry(int hash, String key, String val, int bucketIndex) {
Entry e = table[bucketIndex];
table[bucketIndex] = new Entry(hash, key, val, e);
if (size++ >= threshold) resize(2 * table.length);
}

void resize(int newCapacity) {
Entry[] newTable = new Entry[newCapacity];
transfer(newTable);
table = newTable;
threshold = (int)(newCapacity * 0.75f);
}

void transfer(Entry[] newTable) {
Entry[] src = table;
int newCapacity = newTable.length;
for (int j = 0; j < src.length; j++) {
Entry e = src[j];
if (e != null) {
src[j] = null;
do {
Entry next = e.next;
int i = indexFor(e.hash, newCapacity);
e.next = newTable[i];
newTable[i] = e;
e = next;
} while (e != null);
}
}
}

public String[] toKeyArray(){
String[] ret=new String[size];
int ii=0;
//zacina od konce pole
for (int j = table.length-1; j>=0; j--) {
Entry e = table[j];
if (e != null) {
do {
ret[ii]=e.key;
ii++;
Entry next = e.next;
e = next;
} while (e != null);
}
}
return ret;
}

public String[] toValArray(){
String[] ret=new String[size];
int ii=0;
//zacina od konce pole
for (int j = table.length-1; j>=0; j--) {
Entry e = table[j];
if (e != null) {
do {
ret[ii]=e.val;
ii++;
Entry next = e.next;
e = next;
} while (e != null);
}
}
return ret;
}


int size=0;
int threshold = (int)(16 * 0.75f);
Entry[] table=new Entry[16];


public static void main(String[] args) {
HashMapJava15 set;


set=new HashMapJava15();

set.put("poresultsOut","p1");
set.put("poerrnoOut","p2");
set.put("poerrtextOut","p3");

for(String s:set.toKeyArray()) {
System.out.println(s);
}
for(String s:set.toValArray()) {
System.out.println(s);
}
System.out.println("----------------");

/ *
    <ns0:poerrtextOut xsi:nil="1"/>
    <ns0:poerrnoOut>0</ns0:poerrnoOut>
    <ns0:poresultsOut><results/></ns0:poresultsOut>
 * /

set=new HashMapJava15();

set.put("posumannpremOut","p1");
set.put("posumfreqpremOut","p2");
set.put("popremsxmlOut","p3");
set.put("poerrnoOut","p4");
set.put("poerrtextOut","p5");
set.put("poerrxmlOut","p6");

for(String s:set.toKeyArray()) {
System.out.println(s);
}
for(String s:set.toValArray()) {
System.out.println(s);
}
System.out.println("----------------");

/ *
po_sum_ann_prem out number,--celkove rocni pojistne
po_sum_freq_prem out number,--celkove lhutni pojistne
po_prems_xml out xmltype, --seznam tarifu a jejich rocni pojistne
po_err_no out number,
po_err_text out varchar2,
po_err_xml out xmltype

<ns0:posumfreqpremOut xsi:nil="1"/>
<ns0:poerrtextOut>xml neobsahuje tag apng</ns0:poerrtextOut>
<ns0:popremsxmlOut><tariffs/></ns0:popremsxmlOut>
<ns0:poerrxmlOut><errors/></ns0:poerrxmlOut>
<ns0:poerrnoOut>-1</ns0:poerrnoOut>
<ns0:posumannpremOut xsi:nil="1"/>
 * /


}

}
*/