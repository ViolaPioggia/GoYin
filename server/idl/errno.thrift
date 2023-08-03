namespace go errno

enum Err {
    Success              = 0,
    ParamsErr            = 1,
    ServiceErr           = 2,
    RPCInteractionErr    = 10000,
    InteractionServerErr = 10001,
    RPCSocialityErr      = 20000,
    SocialityServerErr   = 20001,
    RPCUserErr           = 30000,
    UserServerErr        = 30001,
    UserAlreadyExistErr  = 30002,
    UserNotFoundErr      = 30003,
    AuthorizeFailErr     = 30004,
    RPCVideoErr          = 40000,
    VideoServerErr       = 40001,
    RPCChatErr           = 50000,
    ChatServerErr        = 50001,
}