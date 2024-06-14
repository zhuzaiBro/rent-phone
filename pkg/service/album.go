package service

import (
	. "rentServer/core/service"
	"rentServer/pkg/dao"
	"rentServer/pkg/model/album"
)

type AlbumService interface {
	Service[album.Album]
}

type albumServiceIMPL struct {
	ServiceIMPL[album.Album]
	AlbumDao dao.AlbumDao
}

func NewAlbumService() AlbumService {
	albumDao := dao.NewAlbumDao()
	return &albumServiceIMPL{
		ServiceIMPL: ServiceIMPL[album.Album]{
			BaseDao: albumDao,
		},
		AlbumDao: albumDao,
	}
}
