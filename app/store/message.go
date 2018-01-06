package store

import "github.com/nanu-c/textsecure-qml/app/helpers"

type Message struct {
	ID         int64
	SID        int64
	Source     string
	Message    string
	Outgoing   bool
	SentAt     uint64
	ReceivedAt uint64
	HTime      string
	CType      int
	Attachment string
	IsSent     bool
	IsRead     bool
	Flags      int
}

func SaveMessage(m *Message) error {
	res, err := db.NamedExec(messagesInsert, m)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	m.ID = id
	return nil
}

func UpdateMessageSent(m *Message) error {
	_, err := db.NamedExec("UPDATE messages SET issent = :issent, sentat = :sentat WHERE id = :id", m)
	if err != nil {
		return err
	}
	return err
}

func UpdateMessageRead(m *Message) error {
	_, err := db.NamedExec("UPDATE messages SET isread = :isread WHERE id = :id", m)
	if err != nil {
		return err
	}
	return err
}
func LoadMessagesFromDB() error {
	err := db.Select(&AllGroups, groupsSelect)
	if err != nil {
		return err
	}
	for _, g := range AllGroups {
		Groups[g.GroupID] = g
	}

	err = db.Select(&AllSessions, sessionsSelect)
	if err != nil {
		return err
	}
	for _, s := range AllSessions {
		s.When = helpers.HumanizeTimestamp(s.Timestamp)
		s.Active = !s.IsGroup || (Groups[s.Tel] != nil && Groups[s.Tel].Active)
		SessionsModel.Sess = append(SessionsModel.Sess, s)
		SessionsModel.Len++
		err = db.Select(&s.Messages, messagesSelectWhere, s.ID)
		s.Len = len(s.Messages)
		if err != nil {
			return err
		}
		for _, m := range s.Messages {
			m.HTime = helpers.HumanizeTimestamp(m.SentAt)
		}
	}
	return nil
}

func DeleteMessage(id int64) error {
	_, err := db.Exec("DELETE FROM messages WHERE id = ?", id)
	return err
}
func (s *Session) GetMessages(i int) *Message {
	//FIXME when is index -1 ?
	if i == -1 || i >= len(s.Messages) {
		return &Message{}
	}
	return s.Messages[i]
}
func (m *Message) GetName() string {
	return TelToName(m.Source)
}
