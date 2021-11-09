package kiface

type IDataPack interface {
	Pack(IMessage) ([]byte, error)
	UnPack([]byte) (IMessage, error)
}
