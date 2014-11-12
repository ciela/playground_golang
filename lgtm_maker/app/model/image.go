package model

import (
	"github.com/naoina/genmai"
)

type Image struct {
	Id        int64     `db:"pk" json:"id"`
	Key       string    `json:"key"`
	Format    string    `json:"format"`
	CreatedAt time.Time `json:"created_at"`

	genmai.TimeStamp
}

func (m *Image) BeforeInsert() error {
	// FIXME: This method is auto-generated by Kocha.
	//        You can remove this method if unneeded.
	return m.TimeStamp.BeforeInsert()
}

func (m *Image) AfterInsert() error {
	// FIXME: This method is auto-generated by Kocha.
	//        You can remove this method if unneeded.
	return nil
}

func (m *Image) BeforeUpdate() error {
	// FIXME: This method is auto-generated by Kocha.
	//        You can remove this method if unneeded.
	return m.TimeStamp.BeforeUpdate()
}

func (m *Image) AfterUpdate() error {
	// FIXME: This method is auto-generated by Kocha.
	//        You can remove this method if unneeded.
	return nil
}

func (m *Image) BeforeDelete() error {
	// FIXME: This method is auto-generated by Kocha.
	//        You can remove this method if unneeded.
	return nil
}

func (m *Image) AfterDelete() error {
	// FIXME: This method is auto-generated by Kocha.
	//        You can remove this method if unneeded.
	return nil
}
