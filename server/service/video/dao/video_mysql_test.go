package dao

import (
	"GoYin/server/common/test"
	"GoYin/server/service/video/model"
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"testing"
)

func TestVideoLifecycleInMySQL(t *testing.T) {
	cleanUpFunc, db, err := test.RunMysqlInDocker(t)

	defer cleanUpFunc()

	if err != nil {
		t.Fatal(err)
	}

	dao := NewMysqlManager(db)

	ctx := context.Background()
	timeStamp := int64(1676323214)
	videoList := make([]*model.Video, 0)
	for i := int64(0); i < 5; i++ {
		v := &model.Video{
			ID:         100000 + i,
			AuthorId:   200000 + i%2,
			PlayUrl:    fmt.Sprintf("vidoe%d-fake-play-url", i),
			CoverUrl:   fmt.Sprintf("vidoe%d-fake-cover-url", i),
			Title:      fmt.Sprintf("video%d-tiltle", i),
			CreateTime: timeStamp + i,
		}
		videoList = append(videoList, v)
	}

	cases := []struct {
		name       string
		op         func() (string, error)
		wantErr    bool
		wantResult string
	}{
		{
			name: "publish video",
			op: func() (string, error) {
				for _, v := range videoList {
					if err = dao.PublishVideo(ctx, v); err != nil {
						return "", err
					}
				}
				return "", nil
			},
			wantErr:    false,
			wantResult: "",
		},
		{
			name: "duplicate publish video",
			op: func() (string, error) {
				err = dao.PublishVideo(ctx, videoList[0])
				return "", err
			},
			wantErr: true,
		},
		{
			name: "get video list by last time",
			op: func() (string, error) {
				videos, err := dao.GetBasicVideoListByLatestTime(ctx, videoList[0].AuthorId, timeStamp+2)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(videos)
				if err != nil {
					return "", err
				}
				return string(result), err
			},
			wantErr:    false,
			wantResult: `[{"ID":100002,"AuthorId":200000,"PlayUrl":"vidoe2-fake-play-url","CoverUrl":"vidoe2-fake-cover-url","Title":"video2-tiltle","CreateTime":1676323216},{"ID":100001,"AuthorId":200001,"PlayUrl":"vidoe1-fake-play-url","CoverUrl":"vidoe1-fake-cover-url","Title":"video1-tiltle","CreateTime":1676323215},{"ID":100000,"AuthorId":200000,"PlayUrl":"vidoe0-fake-play-url","CoverUrl":"vidoe0-fake-cover-url","Title":"video0-tiltle","CreateTime":1676323214}]`,
		},
		{
			name: "get video list by Author id",
			op: func() (string, error) {
				video, err := dao.GetPublishedVideoListByUserId(ctx, videoList[0].AuthorId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(video)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `[{"ID":100004,"AuthorId":200000,"PlayUrl":"vidoe4-fake-play-url","CoverUrl":"vidoe4-fake-cover-url","Title":"video4-tiltle","CreateTime":1676323218},{"ID":100002,"AuthorId":200000,"PlayUrl":"vidoe2-fake-play-url","CoverUrl":"vidoe2-fake-cover-url","Title":"video2-tiltle","CreateTime":1676323216},{"ID":100000,"AuthorId":200000,"PlayUrl":"vidoe0-fake-play-url","CoverUrl":"vidoe0-fake-cover-url","Title":"video0-tiltle","CreateTime":1676323214}]`,
		},
		{
			name: "get videoId list by Author id",
			op: func() (string, error) {
				video, err := dao.GetPublishedVideoIdListByUserId(ctx, videoList[0].AuthorId)
				if err != nil {
					return "", err
				}
				result, err := sonic.Marshal(video)
				if err != nil {
					return "", err
				}
				return string(result), nil
			},
			wantErr:    false,
			wantResult: `[100000,100002,100004]`,
		},
	}

	for _, cc := range cases {
		result, err := cc.op()
		if cc.wantErr {
			if err == nil {
				t.Errorf("%s:want error;got none", cc.name)
			} else {
				continue
			}
		}
		if err != nil {
			t.Errorf("%s:operation failed: %v", cc.name, err)
		}
		if result != cc.wantResult {
			t.Errorf("%s:result err: want %s,got %s", cc.name, cc.wantResult, result)
		}
	}
}
