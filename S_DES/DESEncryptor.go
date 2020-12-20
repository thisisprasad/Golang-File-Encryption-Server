package S_DES

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sync"
	"time"
)

const (
	bufferSize = 1000000
)

type DesEncryptor struct {
	filename                string
	encryptionFilename      string
	cipher                  DES_8encryption
	decryptionFileConnector *os.File
	decryptionFilename      string
	Fastmode                bool `json:"fastmode"`
	ConcurrentPaths         int  `json:"threadcount"`
	Chunksize               int  `json:"buffersize"`
}

func (encryptor *DesEncryptor) getBinaryByteArray(byteVal byte) []byte {
	var byteArray []byte
	for byteVal > 0 {
		// byteArray = append([]byte{(byte)(byteVal % 2)}, byteArray...)
		byteArray = append(byteArray, (byte)(byteVal%2))
		byteVal /= 2
	}
	for (int)(len(byteArray)) < 8 {
		// byteArray = append([]byte{0}, byteArray...)
		byteArray = append(byteArray, (byte)(0))
	}

	return byteArray
}

func (encryptor *DesEncryptor) convertBinaryByteArrayToByte(byteArray *[]byte) byte {
	var res byte = 0
	for i := len(*byteArray) - 1; i >= 0; i-- {
		if (*byteArray)[i] == (byte)(1) {
			res |= (1 << i)
		} else {
			mask := ^(1 << i)
			res &= (byte)(mask)
		}
	}

	return res
}

func (encryptor *DesEncryptor) encryptChunk(buffer *[]byte, bufferDataSize int, encryptionBuffer *[][]byte) {
	for i := 0; i < bufferDataSize; i++ {
		byteVal := (*buffer)[i]
		var binaryByteArray []byte = encryptor.getBinaryByteArray(byteVal)
		(*encryptionBuffer)[i] = encryptor.cipher.Encrypt(binaryByteArray)
	}
}

func (engine *DesEncryptor) decryptChunk(buffer *[]byte, bufferDataSize int, decryptionBuffer *[][]byte) {
	for i := 0; i < bufferDataSize; i++ {
		byteVal := (*buffer)[i]
		var binaryByteArray []byte = engine.getBinaryByteArray(byteVal)
		(*decryptionBuffer)[i] = engine.cipher.Decrypt(binaryByteArray)
	}
}

func (encryptor *DesEncryptor) writeEncryptionBufferToFile(encryptionBuffer *[][]byte,
	encryptionDataSize int,
	filename string) {
	// var byteVal byte
	// fmt.Println("enccryption buffer length:", encryptionDataSize)
	var byteArray []byte = make([]byte, encryptionDataSize)
	var cache []byte = make([]byte, 8)
	for i := 0; i < encryptionDataSize; i++ {
		copy(cache[:], (*encryptionBuffer)[i][:])
		byteVal := encryptor.convertBinaryByteArrayToByte(&cache)
		byteArray[i] = byteVal
	}

	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, ferr := os.Create(filename)
		if ferr != nil {
			log.Println("Problem encrypting file", ferr, "File not found in the system")
			return
		}
	}
	//	write byte array to buffer
	permissions := 0644
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, (os.FileMode)(permissions))
	if err != nil {
		log.Fatalln("Problem encrypting file", err)
	}
	_, err = file.WriteString((string)(byteArray))

	file.Close()
}

func (engine *DesEncryptor) writeDecryptionBufferToFile(decryptionBuffer *[][]byte,
	decryptionDataSize int,
	filename string) {

	// fmt.Println("Decryption buffer length:", decryptionDataSize)
	var byteArray []byte = make([]byte, decryptionDataSize)
	var cache []byte = make([]byte, 8)
	for i := 0; i < decryptionDataSize; i++ {
		copy(cache[:], (*decryptionBuffer)[i][:])
		byteVal := engine.convertBinaryByteArrayToByte(&cache)
		byteArray[i] = byteVal
	}

	_, err := engine.decryptionFileConnector.WriteString((string)(byteArray))
	if err != nil {
		log.Fatalln("Problem decrypting file", err)
	}
}

func (encryptor *DesEncryptor) concurrentChunkEncryption(waitGroup *sync.WaitGroup, buffer *[]byte, bufferDataSize int, encryptionBuffer *[][]byte) {
	defer waitGroup.Done()
	encryptor.encryptChunk(buffer, bufferDataSize, encryptionBuffer)
}

/**
Uses goroutines to encrypt the file in chunks
*/
func (encryptor *DesEncryptor) fastEncrypt(file *os.File) {
	//	Concurrent encryption of multiple chunks
	dataBuffer := make([][]byte, encryptor.ConcurrentPaths) //	holds the primary data to be encypted.
	for i := 0; i < len(dataBuffer); i++ {
		dataBuffer[i] = make([]byte, encryptor.Chunksize)
	}
	// var encryptionBuffer [][][]byte
	encryptionBuffer := make([][][]byte, encryptor.ConcurrentPaths)
	for i := 0; i < len(encryptionBuffer); i++ {
		encryptionBuffer[i] = make([][]byte, encryptor.Chunksize)
		for j := 0; j < len(encryptionBuffer[i][j]); j++ {
			encryptionBuffer[i][j] = make([]byte, 8)
		}
	}

	waitGroup := new(sync.WaitGroup)
	chunkMap := make([]int, len(dataBuffer)) //	stores the count of bytes w.r.t every chunk.
	var err error
	cnt := 0
	for err != io.EOF {
		err = nil
		bytesread := 0
		chunkCount := 0
		for i := 0; i < len(dataBuffer); i++ {
			bytesread, err = file.Read(dataBuffer[i])
			if err != nil {
				if err != io.EOF {
					log.Println("Problem reading from primary file", err)
				}
				break
			}
			waitGroup.Add(1)
			go encryptor.concurrentChunkEncryption(waitGroup, &dataBuffer[i], bytesread, &encryptionBuffer[i])
			chunkCount += 1
			chunkMap[i] = bytesread
		}
		waitGroup.Wait()

		for i := 0; i < chunkCount; i++ {
			encryptor.writeEncryptionBufferToFile(&encryptionBuffer[i], chunkMap[i], encryptor.encryptionFilename)
		}
		cnt += 1
	}
}

/**
Reads in a chunk of byte-data from file.
Encrypts it.
Writes the encrypted chunk of data in the encryption file
*/
func (encryptor *DesEncryptor) runEncryption(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		log.Println("Problem encrypting file", err)
		return
	}
	defer file.Close()

	//	Create encryption file
	_, err = os.Create(encryptor.encryptionFilename)
	if err != nil {
		log.Println("Problem encrypting file", err)
		return
	}

	var buffer []byte
	buffer = make([]byte, encryptor.Chunksize)
	var encryptionBuffer [][]byte
	encryptionBuffer = make([][]byte, encryptor.Chunksize)
	for i := 0; i < encryptor.Chunksize; i++ {
		encryptionBuffer[i] = make([]byte, 8)
	}

	if encryptor.Fastmode == true {
		log.Println("Fastmode enabled encryption...")
		concurrentStart := time.Now()
		encryptor.fastEncrypt(file)
		concurrentElapsed := time.Since(concurrentStart)
		log.Println("The total time for fast mode: ", concurrentElapsed)
	} else {
		log.Println("Using normal mode encryption...")
		normalStart := time.Now()
		for {
			bytesread, err := file.Read(buffer)
			if err != nil {
				if err != io.EOF {
					log.Fatalln("Problem while reading the contents of the file.")
				}
				break
			}
			encryptor.encryptChunk(&buffer, bytesread, &encryptionBuffer)
			//	Write encryptionBuffer into file.
			encryptor.writeEncryptionBufferToFile(&encryptionBuffer, bytesread, encryptor.encryptionFilename)
		}
		normalElapsed := time.Since(normalStart)
		log.Println("The total time for normal mode: ", normalElapsed)
	}
}

func (engine *DesEncryptor) concurrentChunkDecryption(waitGroup *sync.WaitGroup, buffer *[]byte, bufferDataSize int, decryptionBuffer *[][]byte) {
	defer waitGroup.Done()
	engine.decryptChunk(buffer, bufferDataSize, decryptionBuffer)
}

func (engine *DesEncryptor) fastDecrypt(file *os.File) {
	dataBuffer := make([][]byte, engine.ConcurrentPaths) //	holds the primary data to be decrypted.
	for i := 0; i < len(dataBuffer); i++ {
		dataBuffer[i] = make([]byte, engine.Chunksize)
	}
	decryptionBuffer := make([][][]byte, engine.ConcurrentPaths)
	for i := 0; i < len(decryptionBuffer); i++ {
		decryptionBuffer[i] = make([][]byte, engine.Chunksize)
		for j := 0; j < len(decryptionBuffer[i][j]); j++ {
			decryptionBuffer[i][j] = make([]byte, 8)
		}
	}
	waitGroup := new(sync.WaitGroup)
	chunkMap := make([]int, len(dataBuffer))
	var err error
	cnt := 0
	for err != io.EOF {
		err = nil
		bytesread := 0
		chunkCount := 0
		for i := 0; i < len(dataBuffer); i++ {
			bytesread, err = file.Read(dataBuffer[i])
			if err != nil {
				if err != io.EOF {
					log.Println("Problem reading from primary file", err)
				}
				break
			}
			waitGroup.Add(1)
			go engine.concurrentChunkDecryption(waitGroup, &dataBuffer[i], bytesread, &decryptionBuffer[i])
			chunkCount += 1
			chunkMap[i] = bytesread
		}
		waitGroup.Wait()

		for i := 0; i < chunkCount; i++ {
			//	write decryption buffer to file
			engine.writeDecryptionBufferToFile(&decryptionBuffer[i], chunkMap[i], engine.decryptionFilename)
		}
		cnt += 1
	}
}

/**
Decrypt the bytes of the encrypted file
*/
func (engine *DesEncryptor) runDecryption(filename string) {
	file, ferr := os.Open(engine.encryptionFilename)
	defer file.Close()
	if ferr != nil {
		log.Fatalln("Problem opening encrypted file...Filename:", engine.encryptionFilename, ferr)
		return
	}

	if engine.Fastmode == true {
		log.Println("Fastmode enabled decryption...")
		concurrentStart := time.Now()
		engine.fastDecrypt(file)
		concurrentElapsed := time.Since(concurrentStart)
		log.Println("The total decryption time for fast mode: ", concurrentElapsed)
	} else {
		log.Println("Decryption in normal mode...")
		var buffer []byte = make([]byte, engine.Chunksize)
		var decryptionBuffer [][]byte = make([][]byte, engine.Chunksize)
		for i := 0; i < engine.Chunksize; i++ {
			decryptionBuffer[i] = make([]byte, 8)
		}
		normalStart := time.Now()
		for {
			bytesread, err := file.Read(buffer)
			if err != nil {
				if err != io.EOF {
					log.Fatalln("Problem reading encrypted file", err)
				}
				break
			}
			engine.decryptChunk(&buffer, bytesread, &decryptionBuffer)
			//	Write decryptionBuffer into file
			engine.writeDecryptionBufferToFile(&decryptionBuffer, bytesread, engine.decryptionFilename)
		}
		normalElapsed := time.Since(normalStart)
		log.Println("Total decryption time in normal mode: ", normalElapsed)
	}

	// fmt.Println("Encryption filename:", engine.encryptionFilename)
}

/**
public File-encryption API.
@return: boolean value. True if encryption is successful and encrypted file
		is created successfully on the disk. Otherwise false is returned.
*/
func (encryptor *DesEncryptor) EncryptFile(filename string) bool {
	log.Println("file for encryption:", filename)
	log.Println("File-Encryption procedure started...")

	encryptor.filename = filename
	encryptor.encryptionFilename = filename + ".enc"
	encryptor.runEncryption(filename)

	//	Check whether a file of same size is created as the original
	//	file on the disk.
	primaryFileStat, err := os.Stat(filename)
	if err != nil {
		log.Println("Encryption failure for file:", filename)
		return false
	}
	primaryFileSize := primaryFileStat.Size()

	encryptedFileStat, err := os.Stat(encryptor.encryptionFilename)
	if err != nil {
		log.Println("Encryption failure for file:", encryptor.encryptionFilename)
		return false
	}
	encryptedFileSize := encryptedFileStat.Size()
	if encryptedFileSize != primaryFileSize {
		log.Println("Encryption Failure. Try again...")
		return false
	}

	log.Println("Deleting primary file:", filename)
	os.Remove(filename)

	log.Println("File-Encryption procedure complete...")
	return true
}

/**
Public Decryption API.
@return: boolean value. True if decryption is successful and decrypted file
		is created successfully on the disk. Otherwise false is returned.
*/
func (engine *DesEncryptor) DecryptFile(filename string) bool {
	log.Println("Decryption procedure started...")

	engine.decryptionFilename = filename
	engine.encryptionFilename = filename + ".enc"
	var err error
	engine.decryptionFileConnector, err = os.Create(engine.decryptionFilename)
	if err != nil {
		log.Fatalln("Problem decrypting the file", err)
	}
	permissions := 0644
	engine.decryptionFileConnector, err =
		os.OpenFile(engine.decryptionFilename, os.O_APPEND|os.O_WRONLY, (os.FileMode)(permissions))
	// file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, (os.FileMode)(permissions))
	// engine.cipher.Init("des_input.txt")
	engine.runDecryption(filename)
	engine.decryptionFileConnector.Close()

	//	Delete the encrypted file and keep the decrypted(which is the primary file) file
	encryptedFileStat, err := os.Stat(engine.encryptionFilename)
	if err != nil {
		log.Println("decryption failure for file:", engine.decryptionFilename)
		return false
	}
	encryptedFileSize := encryptedFileStat.Size()

	decryptedFileStat, err := os.Stat(engine.decryptionFilename)
	if err != nil {
		log.Println("decryption failure for file:", engine.decryptionFilename)
		return false
	}
	decryptedFileSize := decryptedFileStat.Size()
	if encryptedFileSize != decryptedFileSize {
		log.Println("Decryption failed! Try again")
		return false
	}

	log.Println("Removing encrypted file: ", engine.encryptionFilename)
	os.Remove(engine.encryptionFilename)

	log.Println("Decryption procedure complete...")
	return true
}

/**
Reads the DES configuration file. Initializes the struct members.
*/
func (engine *DesEncryptor) readConfigurationFile(configFilename string) {
	data, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatalln("Problem reading the DES configuration file...", err)
	}
	err = json.Unmarshal(data, &engine)
}

func (engine *DesEncryptor) Init(filename string) {
	log.Println("Initializing FES engine...")
	engine.readConfigurationFile("des_config.json")
	engine.cipher.Init(filename)
	log.Println("FES engine intialization complete...")
}
