package mega

import (
	"time"

	"github.com/OpenListTeam/OpenList/v4/pkg/utils"

	"github.com/OpenListTeam/OpenList/v4/internal/model"
	"github.com/t3rm1n4l/go-mega"
)

type MegaNode struct {
	n *mega.Node
}

func (m *MegaNode) GetSize() int64 {
	return m.n.GetSize()
}

func (m *MegaNode) GetName() string {
	return m.n.GetName()
}

func (m *MegaNode) CreateTime() time.Time {
	return m.n.GetTimeStamp()
}

func (m *MegaNode) GetHash() utils.HashInfo {
	//Meganz use md5, but can't get the original file hash, due to it's encrypted in the cloud
	return utils.HashInfo{}
}

func (m *MegaNode) ModTime() time.Time {
	return m.n.GetTimeStamp()
}

func (m *MegaNode) IsDir() bool {
	return m.n.GetType() == mega.FOLDER || m.n.GetType() == mega.ROOT
}

func (m *MegaNode) GetID() string {
	return m.n.GetHash()
}

func (m *MegaNode) GetPath() string {
	return ""
}

var _ model.Obj = (*MegaNode)(nil)
