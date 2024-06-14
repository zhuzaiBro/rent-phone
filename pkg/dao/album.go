package dao

import (
	. "rentServer/core/dao"
	"rentServer/initilization/db"
	"rentServer/pkg/model/album"
)

var _ AlbumDao = (*albumDao)(nil)

type AlbumDao interface {
	Dao[album.Album]
}

type albumDao struct {
	DaoIMPL[album.Album]
}

func NewAlbumDao() AlbumDao {
	return &albumDao{
		DaoIMPL: DaoIMPL[album.Album]{
			Db:    db.DB,
			Model: &album.Album{},
		},
	}
}
