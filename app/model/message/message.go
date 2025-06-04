package message

import (
	"douyin-backend/app/global/variable"
	"douyin-backend/app/model"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type MsgModel struct {
	*gorm.DB   `gorm:"-" json:"-"`
	ID         int64     `json:"id" db:"bigint"`
	TxUID      int64     `json:"tx_uid" db:"bigint"`
	RxUID      int64     `json:"rx_uid" db:"bigint"`
	MsgType    int       `json:"msg_type" db:"int"`
	MsgData    string    `json:"msg_data" db:"text"`
	ReadState  int       `json:"read_state" db:"int"`
	CreateTime time.Time `json:"create_time" db:"int"`
	DeleteTime time.Time `json:"delete_time" db:"int"`
}

func CreateMsgFactory(sqlType string) *MsgModel {
	return &MsgModel{
		DB: model.UseDbConn(sqlType),
	}
}

func (m *MsgModel) SendMsg(txUid, rxUid int64, msgType int, msgData string, readState int, createTime int) bool {
	sql := `INSERT INTO tb_messages (tx_uid, rx_uid, msg_type, msg_data, read_state, create_time) VALUES (?, ?, ?, ?, ?, ?);`
	result := m.Exec(sql, txUid, rxUid, msgType, msgData, readState, createTime)
	if result.RowsAffected > 0 {
		return true
	} else {
		return false
	}
}
func (m *MsgModel) GetAllMsg(uid int64) (allMsg map[string][]Message, ok bool) {
	allMsg = make(map[string][]Message)
	sql := `SELECT * FROM tb_messages WHERE tx_uid=? OR rx_uid=?;`
	var res []Message
	result := m.Raw(sql, uid, uid).Find(&res)
	if result.Error != nil {
		variable.ZapLog.Error("GetAllMsg 查询出错", zap.Error(result.Error))
		ok = false
		return
	}
	// 如果查询结果不为空
	if len(res) > 0 {
		for _, msg := range res {
			uidStr := strconv.FormatInt(uid, 10)
			if msg.TxUID == uidStr {
				allMsg[msg.RxUID] = append(allMsg[msg.RxUID], msg)
			} else {
				allMsg[msg.TxUID] = append(allMsg[msg.TxUID], msg)
			}
		}
	}
	ok = true
	return
}
