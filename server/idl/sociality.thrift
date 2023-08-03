namespace go sociality

include "base.thrift"

struct douyin_relation_action_request {
    1: i64 user_id,    // User Id
    2: i64 to_user_id, // The other party's user id
    3: i8 action_type, // 1-Follow, 2-Unfollow
}

struct douyin_relation_action_response {
    1: base.douyin_base_response base_resp,
}

struct douyin_get_relation_id_list_request {
    1: i64 viewer_id, // User id of viewer,set to zero when unclear
    2: i64 owner_id,  // User id of owner.
    3: i8  option, //FollowingList=0 FollowerList=1 FriendsList=2
}

struct douyin_get_relation_id_list_response {
    1: base.douyin_base_response base_resp,
    2: list<i64> user_id_list,           // List of user information
}

struct douyin_get_social_info_request{
    1: i64 viewer_id
    2: i64 owner_id
}

struct douyin_get_social_info_response{
    1: base.douyin_base_response base_resp
    2: base.SocialInfo social_info
}

struct douyin_batch_get_social_info_request{
    1: i64 viewer_id
    2: list<i64> owner_id_list
}

struct douyin_batch_get_social_info_response{
    1: base.douyin_base_response base_resp
    2: list<base.SocialInfo> social_info_list
}


service SocialityService {
    douyin_relation_action_response Action(1: douyin_relation_action_request req),
    douyin_get_relation_id_list_response GetRelationIdList(1: douyin_get_relation_id_list_request req),
    douyin_get_social_info_response GetSocialInfo(1:douyin_get_social_info_request req),
    douyin_batch_get_social_info_response BatchGetSocialInfo(1:douyin_batch_get_social_info_request req),
}