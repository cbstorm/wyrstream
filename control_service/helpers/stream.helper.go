package helpers

import (
	"fmt"
	"os"
	"strings"

	"github.com/cbstorm/wyrstream/lib/entities"
	"github.com/cbstorm/wyrstream/lib/logger"
	"github.com/cbstorm/wyrstream/lib/minio_service"
	"github.com/cbstorm/wyrstream/lib/repositories"
	"github.com/cbstorm/wyrstream/lib/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StreamsHelper struct {
	streams   *[]*entities.StreamEntity
	list_dirs *map[string]*minio_service.BulkListDirResult
}

func NewStreamsHelper(streams *[]*entities.StreamEntity) *StreamsHelper {
	return &StreamsHelper{
		streams: streams,
	}
}

func (h *StreamsHelper) ResolveStreamLogs() error {
	if len(*h.streams) == 0 {
		return nil
	}
	stream_obj_ids := utils.Map(h.streams, func(e *entities.StreamEntity, i int) primitive.ObjectID {
		return e.Id
	})
	stream_logs := make([]*entities.StreamLogEntity, 0)
	if err := repositories.GetStreamLogRepository().Find(map[string]interface{}{
		"stream_obj_id": map[string]interface{}{
			"$in": stream_obj_ids,
		},
	}, &stream_logs); err != nil {
		return err
	}
	stream_logs_group_by_stream_id := utils.GroupBy(&stream_logs, func(e *entities.StreamLogEntity) string {
		return e.StreamObjId.Hex()
	})
	for _, v := range *h.streams {
		v.StreamLogs = (*stream_logs_group_by_stream_id)[v.Id.Hex()]
	}
	return nil
}

func (h *StreamsHelper) ResolvePublisher() error {
	if len(*h.streams) == 0 {
		return nil
	}
	publisher_ids := utils.Map(h.streams, func(e *entities.StreamEntity, i int) primitive.ObjectID {
		return e.PublisherId
	})
	publishers := make([]*entities.UserEntity, 0)
	if err := repositories.GetUserRepository().FindManyByIds(*publisher_ids, &publishers); err != nil {
		return err
	}
	publishers_key_by_id := utils.KeyBy(&publishers, func(a *entities.UserEntity) string {
		return a.Id.Hex()
	})
	for _, v := range *h.streams {
		v.Publisher = (*publishers_key_by_id)[v.PublisherId.Hex()]
	}
	return nil
}

func (h *StreamsHelper) ListStorageDirs() {
	dir_prefixs := utils.Map(h.streams, func(a *entities.StreamEntity, b int) string {
		return fmt.Sprintf("streams/%s/segments/", a.StreamId)
	})
	list_dirs := minio_service.GetMinioService().ListDirs(dir_prefixs)
	h.list_dirs = list_dirs
}

func (h *StreamsHelper) getListDirByStream(stream *entities.StreamEntity) (*[]string, error) {
	res := (*h.list_dirs)[fmt.Sprintf("streams/%s/segments/", stream.StreamId)]
	return res.Result, res.Error
}

func (h *StreamsHelper) GenerateHLSPlaylist(stream *entities.StreamEntity) (string, error) {
	segments, err := h.getListDirByStream(stream)
	if err != nil {
		return "", err
	}
	f_path := fmt.Sprintf("tmp/%s/playlist.m3u8", stream.StreamId)
	if err := utils.AssertDir(f_path); err != nil {
		return "", err
	}
	f, err := os.Create(f_path)
	defer func() {
		if err := f.Sync(); err != nil {
			logger.UnexpectedErrLogger.Error("Could not sync file %s due to an error: %v", f_path, err)
		}
		if err := f.Close(); err != nil {
			logger.UnexpectedErrLogger.Error("Could not close file %s due to an error: %v", f_path, err)
		}
	}()
	if err != nil {
		return "", err
	}
	if _, err := f.WriteString(fmt.Sprintf("#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-TARGETDURATION:5\n#EXT-X-MEDIA-SEQUENCE:%d\n", len(*segments))); err != nil {
		return "", err
	}
	t := "#EXTINF:5.000000,"
	for _, v := range *segments {
		s := strings.Replace(v, fmt.Sprintf("streams/%s", stream.StreamId), "..", 1)
		if _, err := f.WriteString(fmt.Sprintf("%s\n%s\n", t, s)); err != nil {
			return "", err
		}
	}
	if _, err := f.WriteString("#EXT-X-ENDLIST"); err != nil {
		return "", err
	}
	return f_path, nil
}
