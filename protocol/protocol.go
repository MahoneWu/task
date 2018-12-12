package protocol

import (
    "bytes"
    "encoding/binary"
)

const (
    ConstHeader = "Headers"
    ConstHeaderLength = 7
    ConstMLength = 4
)

//package data
func Enpack(message []byte) []byte {
    return append(append([]byte(ConstHeader), IntToBytes(len(message))...), message...)
}

//unpacking data
func Depack(buffer []byte) ([]byte, []byte, bool){
    //get buffer data length
    length := len(buffer)
    //define flag to mark whether the tag still needs to be parase，if true then need to be parase
    var uncompletedFlag bool
    var i int
    data := make([]byte, 32)

    for i = 0; i < length; i++ {
        if length < i + ConstHeaderLength + ConstMLength{
            break
        }
        //if equal to protocal header ,then resolved
        if string(buffer[i:i+ConstHeaderLength]) == ConstHeader {
            //get ConstMLength value
            messageLength := ByteToInt(buffer[i+ConstHeaderLength : i+ConstHeaderLength+ConstMLength])
            //if whole buffer length less than ,represent the mesage passed in is incomplete
            if length < i+ConstHeaderLength+ConstMLength+messageLength {
                break
            }
            //get the value of message
            data = buffer[i+ConstHeaderLength+ConstMLength : i+ConstHeaderLength+ConstMLength+messageLength]
            if length > i + ConstHeaderLength + ConstMLength + messageLength{
                uncompletedFlag = true
            }
            //get unprocessed data byte
            uncompletedBuffer := buffer[i+ConstHeaderLength+ConstMLength+messageLength:]
            return data, uncompletedBuffer, uncompletedFlag
        }
    }
    return make([]byte, 0), buffer[:], false
}

//byte convert to Integer
func ByteToInt(n []byte) int {
    bytesbuffer := bytes.NewBuffer(n)
    var x int32
    binary.Read(bytesbuffer, binary.BigEndian, &x)

    return int(x)
}

//Integer convert to byte
func IntToBytes(n int) []byte {
    x := int32(n)
    bytesBuffer := bytes.NewBuffer([]byte{})
    binary.Write(bytesBuffer, binary.BigEndian, x)
    return bytesBuffer.Bytes()
}
