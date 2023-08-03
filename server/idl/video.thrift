namespace go video

include "base.thrift"

struct douyin_feed_request {
    1: i64 latest_time, // Optional parameter, limit the latest submission timestamp of the returned video, accurate to seconds, and leave it blank to indicate the current time
    2: i64 viewer_id,   // Optional parameter, user id of viewer,set to zero when unclear
}

struct douyin_feed_response {
    1: base.douyin_base_response base_resp,
    2: list<base.Video> video_list,         // Video list
    3: i64 next_time,                       // In the video returned this time, publish the earliest time as the latest_time in the next request
}

struct douyin_publish_action_request {
    1: i64 user_id,      // User Id
    2: string play_url,  // Video play url
    3: string cover_url, // Video cover url
    4: string title,     // Video title
}

struct douyin_publish_action_response {
    1: base.douyin_base_response base_resp,
}

struct douyin_get_published_list_request {
    1: i64 viewer_id, // User id of viewer,set to zero when unclear
    2: i64 owner_id,  // User id of owner
}

struct douyin_get_published_list_response {
    1: base.douyin_base_response base_resp,
    2: list<base.Video> video_list,         // List of videos posted by the user
}

struct douyin_get_favorite_list_request {
    1: i64 viewer_id, // User id of viewer,set to zero when unclear
    2: i64 owner_id,  // User id of owner.
}

struct douyin_get_favorite_list_response {
    1: base.douyin_base_response base_resp,
    2: list<base.Video> video_list,         // List of videos posted by the user
}

struct douyin_get_published_video_id_list_request{
    1: i64 user_id,
}

struct douyin_get_published_video_id_list_response{
    1: base.douyin_base_response base_resp,
    2: list<i64> video_id_list,
}

service VideoService {
    douyin_feed_response Feed(1: douyin_feed_request req),
    douyin_publish_action_response PublishVideo(1: douyin_publish_action_request req),
    douyin_get_published_list_response GetPublishedVideoList(1: douyin_get_published_list_request req),
    douyin_get_favorite_list_response GetFavoriteVideoList(1: douyin_get_favorite_list_request req),
    douyin_get_published_video_id_list_response GetPublishedVideoIdList(1: douyin_get_published_video_id_list_request req),
}