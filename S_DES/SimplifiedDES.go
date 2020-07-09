package S_DES

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type DES_8encryption struct {
	p4                   []int //	4-length permutation
	p8                   []int //	8-length permutation
	p10                  []int //	10-length permuation
	initialPermutation   []int
	inversePermutation   []int
	expansionPermutation []int
	s0                   [][]int //	s0 matrix
	s1                   [][]int //	s1-matrix
	key                  []byte
	key1                 []byte
	key2                 []byte
}

/**
Sequence of reading input is strict
e.g.:
key:0010010111
P10:3 5 2 7 4 10 1 9 8 6
P8:6 3 7 4 8 5 10 9
P4:2 4 3 1
IP:2 6 3 1 4 8 5 7
IP-1:4 1 3 5 7 2 8 6
E/P:4 1 2 3 2 3 4 1
S0:1 0 3 2 3 2 1 0 0 2 1 3 3 1 3 2
S1:0 1 2 3 2 0 1 3 3 0 1 0 2 1 0 3
*/
func (cipher *DES_8encryption) readFile(configFile string) {
	var s string
	file, err := os.Open(configFile)
	if err != nil {
		log.Fatalln("Error opening config file", configFile, "Error: ", err)
	}
	defer file.Close()

	//	Key
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	s = scanner.Text()
	cipher.key = []byte(strings.Split(s, ":")[1])
	//	Converting char value to integer value
	for i := 0; i < len(cipher.key); i++ {
		cipher.key[i] = cipher.key[i] - '0'
	}

	//	p10 permutation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.p10 = StringToIntArray(s)

	//	p8 permutation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.p8 = StringToIntArray(s)

	//	p4 permuation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.p4 = StringToIntArray(s)

	//	initial permutation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.initialPermutation = StringToIntArray(s)

	//	inverse permutation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.inversePermutation = StringToIntArray(s)

	//	expansion permuation
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.expansionPermutation = StringToIntArray(s)

	//	s0 matrix
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.s0 = StringTo2DIntArray(s, 4, 4)

	//	s1 matrix
	scanner.Scan()
	s = strings.Split(scanner.Text(), ":")[1]
	cipher.s1 = StringTo2DIntArray(s, 4, 4)
}

/**
Generates intermediate keys of DES algorithm
*/
func (cipher *DES_8encryption) generateIntermediateKeys() {
	var p10key = cipher.applyPermutation(cipher.key, cipher.p10)
	var leftHalf []byte = p10key[0 : len(p10key)/2]
	var rightHalf []byte = p10key[len(p10key)/2:]
	cipher.circularLeftShift(&leftHalf, 1)
	cipher.circularLeftShift(&rightHalf, 1)

	var combinedKey []byte = make([]byte, len(p10key))
	copy(combinedKey[:len(p10key)/2], leftHalf[:])
	copy(combinedKey[len(p10key)/2:], rightHalf[:])
	cipher.key1 = cipher.applyPermutation(combinedKey, cipher.p8)

	leftHalf = combinedKey[0 : len(combinedKey)/2]
	rightHalf = combinedKey[len(combinedKey)/2:]
	cipher.circularLeftShift(&leftHalf, 2)
	cipher.circularLeftShift(&rightHalf, 2)
	cipher.key2 = cipher.applyPermutation(combinedKey, cipher.p8)
}

func (cipher *DES_8encryption) Init(configFile string) {
	log.Println("Initializing cipher...")

	cipher.readFile(configFile)
	cipher.generateIntermediateKeys()

	log.Println("cipher intialization complete...")
}

func (cipher *DES_8encryption) XOR(op1 []byte, op2 []byte) []byte {
	if len(op1) != len(op2) {
		var cache []byte
		if len(op1) < len(op2) {
			cache = op1[0:]
		} else {
			cache = op2[0:]
		}
		for len(cache) != MaxInt(len(op1), len(op2)) {
			cache = append([]byte{0}, cache...)
		}
	}

	var res []byte = make([]byte, len(op1))
	for i := 0; i < len(op1); i++ {
		if op1[i] == op2[i] {
			res[i] = 0
		} else {
			res[i] = 1
		}
	}

	return res
}

func (cipher *DES_8encryption) circularLeftShift(data *[]byte, shiftBy int) {
	shiftBy %= len(*data)
	var cache []byte = make([]byte, shiftBy)
	for i := 0; i < shiftBy; i++ {
		cache[i] = (*data)[i]
	}
	for i := shiftBy; i < len(*data); i++ {
		(*data)[i-shiftBy] = (*data)[i]
	}
	pos := 0
	for i := len(*data) - shiftBy; i < len(*data); i++ {
		(*data)[i] = cache[pos]
		pos += 1
	}
}

func (cipher *DES_8encryption) applyPermutation(data []byte, permutation []int) []byte {
	var res []byte = make([]byte, len(permutation))
	pos := 0
	for i := 0; i < len(permutation); i++ {
		ch := data[permutation[i]-1] //	0-indexed
		res[pos] = ch
		pos += 1
	}

	return res
}

func (cipher *DES_8encryption) fk(leftHalf []byte, rightHalf []byte, key []byte) []byte {
	var res []byte

	var epBits []byte = cipher.applyPermutation(rightHalf, cipher.expansionPermutation)
	var XORBits []byte = cipher.XOR(epBits, key)
	var leftNibble []byte = XORBits[0 : len(XORBits)/2]
	var rightNibble []byte = XORBits[len(XORBits)/2:]
	var row0 int = (int)(leftNibble[3] + leftNibble[0]*2)
	var col0 int = (int)(leftNibble[2] + leftNibble[1]*2)
	var row1 int = (int)(rightNibble[3] + rightNibble[0]*2)
	var col1 int = (int)(rightNibble[2] + rightNibble[1]*2)
	var val1 int = cipher.s0[row0][col0]
	var val2 int = cipher.s1[row1][col1]

	var i1, i2 []byte
	for len(i1) != 2 {
		i1 = append([]byte{(byte)(val1 % 2)}, i1...)
		val1 /= 2
	}
	for len(i2) != 2 {
		i2 = append([]byte{(byte)(val2 % 2)}, i2...)
		val2 /= 2
	}

	var cache []byte
	cache = append(cache, i1...)
	cache = append(cache, i2...)
	cache = cipher.applyPermutation(cache, cipher.p4)
	XORBits = cipher.XOR(leftHalf, cache)

	res = append(res, XORBits...)
	res = append(res, rightHalf...)

	return res
}

/**
encryption API exposed to client
*/
func (cipher *DES_8encryption) Encrypt(plainText []byte) []byte {
	// log.Println("Input byte is:", plainText)

	var ipBits []byte = cipher.applyPermutation(plainText, cipher.initialPermutation)
	var leftHalf []byte = ipBits[0 : len(ipBits)/2]
	var rightHalf []byte = ipBits[len(ipBits)/2:]
	var fkBits []byte = cipher.fk(leftHalf, rightHalf, cipher.key1)
	// fmt.Println("fkbits:", fkBits)

	rightHalf = fkBits[0 : len(fkBits)/2]
	leftHalf = fkBits[len(fkBits)/2:]
	fkBits = cipher.fk(leftHalf, rightHalf, cipher.key2)

	var cipherData []byte = cipher.applyPermutation(fkBits, cipher.inversePermutation)
	// fmt.Println("cipherData:", cipherData)

	return cipherData
}

func (cipher *DES_8encryption) Decrypt(encryptedData []byte) []byte {
	// log.Println("Input encryption data:", encryptedData)
	var ipInverseBits []byte = cipher.applyPermutation(encryptedData, cipher.initialPermutation)
	var leftHalf []byte = ipInverseBits[0 : len(ipInverseBits)/2]
	var rightHalf []byte = ipInverseBits[len(ipInverseBits)/2:]

	var fkBits []byte = cipher.fk(leftHalf, rightHalf, cipher.key2)

	rightHalf = fkBits[0 : len(fkBits)/2]
	leftHalf = fkBits[len(fkBits)/2:]
	fkBits = cipher.fk(leftHalf, rightHalf, cipher.key1)

	var decryptedData []byte = cipher.applyPermutation(fkBits, cipher.inversePermutation)
	// fmt.Println("decrypted data:", decryptedData)

	return decryptedData
}
