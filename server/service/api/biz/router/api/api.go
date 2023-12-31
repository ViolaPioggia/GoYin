// Code generated by hertz generator. DO NOT EDIT.

package api

import (
	api "GoYin/server/service/api/biz/handler/api"
	"github.com/cloudwego/hertz/pkg/app/server"
)

/*
 This file will register all the routes of the services in the master idl.
 And it will update automatically when you use the "update" command for the idl.
 So don't modify the contents of the file, or your code will be deleted when it is updated.
*/

// Register register routes based on the IDL 'api.${HTTP Method}' annotation.
func Register(r *server.Hertz) {

	root := r.Group("/", rootMw()...)
	{
		_douyin := root.Group("/douyin", _douyinMw()...)
		{
			_comment := _douyin.Group("/comment", _commentMw()...)
			{
				_action := _comment.Group("/action", _actionMw()...)
				_action.POST("/", append(_comment0Mw(), api.Comment)...)
			}
			{
				_list := _comment.Group("/list", _listMw()...)
				_list.GET("/", append(_commentlistMw(), api.CommentList)...)
			}
		}
		{
			_favorite := _douyin.Group("/favorite", _favoriteMw()...)
			{
				_action0 := _favorite.Group("/action", _action0Mw()...)
				_action0.POST("/", append(_favorite0Mw(), api.Favorite)...)
			}
			{
				_list0 := _favorite.Group("/list", _list0Mw()...)
				_list0.GET("/", append(_favoritelistMw(), api.FavoriteList)...)
			}
		}
		{
			_feed := _douyin.Group("/feed", _feedMw()...)
			_feed.GET("/", append(_feed0Mw(), api.Feed)...)
		}
		{
			_message := _douyin.Group("/message", _messageMw()...)
			{
				_action1 := _message.Group("/action", _action1Mw()...)
				_action1.POST("/", append(_sentmessageMw(), api.SentMessage)...)
			}
			{
				_chat := _message.Group("/chat", _chatMw()...)
				_chat.GET("/", append(_chathistoryMw(), api.ChatHistory)...)
			}
		}
		{
			_publish := _douyin.Group("/publish", _publishMw()...)
			{
				_action2 := _publish.Group("/action", _action2Mw()...)
				_action2.POST("/", append(_publishvideoMw(), api.PublishVideo)...)
			}
			{
				_list1 := _publish.Group("/list", _list1Mw()...)
				_list1.GET("/", append(_videolistMw(), api.VideoList)...)
			}
		}
		{
			_relation := _douyin.Group("/relation", _relationMw()...)
			{
				_action3 := _relation.Group("/action", _action3Mw()...)
				_action3.POST("/", append(_action4Mw(), api.Action)...)
			}
			{
				_follow := _relation.Group("/follow", _followMw()...)
				{
					_list2 := _follow.Group("/list", _list2Mw()...)
					_list2.GET("/", append(_followinglistMw(), api.FollowingList)...)
				}
			}
			{
				_follower := _relation.Group("/follower", _followerMw()...)
				{
					_list3 := _follower.Group("/list", _list3Mw()...)
					_list3.GET("/", append(_followerlistMw(), api.FollowerList)...)
				}
			}
			{
				_friend := _relation.Group("/friend", _friendMw()...)
				{
					_list4 := _friend.Group("/list", _list4Mw()...)
					_list4.GET("/", append(_friendlistMw(), api.FriendList)...)
				}
			}
		}
		{
			_user := _douyin.Group("/user", _userMw()...)
			_user.GET("/", append(_getuserinfoMw(), api.GetUserInfo)...)
			{
				_login := _user.Group("/login", _loginMw()...)
				_login.POST("/", append(_login0Mw(), api.Login)...)
			}
			{
				_register := _user.Group("/register", _registerMw()...)
				_register.POST("/", append(_register0Mw(), api.Register)...)
			}
		}
	}
}
