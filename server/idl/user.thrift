namespace go user

include "base.thrift"

struct douyin_user_register_request {
    1: string username, // Username, up to 32 characters
    2: string password, // Password, up to 32 characters
}

struct douyin_user_register_response {
    1: base.douyin_base_response base_resp,
    2: i64 user_id,                         // User id
    3: string token,                        // User authentication token
}

struct douyin_user_login_request {
    1: string username, // Username, up to 32 characters
    2: string password, // Password, up to 32 characters
}

struct douyin_user_login_response {
    1: base.douyin_base_response base_resp,
    2: i64 user_id,                         // User id
    3: string token,                        // User authentication token
}

struct douyin_get_user_request {
    1: i64 viewer_id, // User id of viewer,set to zero when unclear
    2: i64 owner_id,  // User id of owner.
}

struct douyin_get_user_response {
    1: base.douyin_base_response base_resp,
    2: base.User user,                      // User Information
}

struct douyin_batch_get_user_request {
    1: i64 viewer_id,       // User id of viewer,set to zero when unclear
    2: list<i64> owner_id_list, // User id list of info owners.
}

struct douyin_batch_get_user_resonse {
    1: base.douyin_base_response base_resp,
    2: list<base.User> user_list,
}

struct douyin_get_relation_follow_list_request {
    1: i64 viewer_id, // User id of viewer,set to zero when unclear
    2: i64 owner_id,  // User id of owner.
}

struct douyin_get_relation_follow_list_response {
    1: base.douyin_base_response base_resp,
    2: list<base.User> user_list,           // List of user information
}

struct douyin_get_relation_follower_list_request {
    1: i64 viewer_id, // User id of viewer,set to zero when unclear
    2: i64 owner_id,  // User id of owner.
}

struct douyin_get_relation_follower_list_response {
    1: base.douyin_base_response base_resp,
    2: list<base.User> user_list,           // List of user information
}

struct douyin_get_relation_friend_list_request {
    1: i64 viewer_id, // User id of viewer,set to zero when unclear
    2: i64 owner_id,  // User id of owner.
}

struct douyin_get_relation_friend_list_response {
    1: base.douyin_base_response base_resp,
    2: list<base.FriendUser> user_list,     // List of user information
}

service UserService {
    douyin_user_register_response Register(1: douyin_user_register_request req),
    douyin_user_login_response Login(1: douyin_user_login_request req),
    douyin_get_user_response GetUserInfo(1: douyin_get_user_request req),
    douyin_batch_get_user_resonse BatchGetUserInfo(1: douyin_batch_get_user_request req),
    douyin_get_relation_follow_list_response GetFollowList(1: douyin_get_relation_follow_list_request req),
    douyin_get_relation_follower_list_response GetFollowerList(1: douyin_get_relation_follower_list_request req),
    douyin_get_relation_friend_list_response GetFriendList(1: douyin_get_relation_friend_list_request req),
}