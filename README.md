# Projects Chat

Projects With An ID with list of users that can access to it.

## ENDPOINTS 
 - projects/:project_id/chat
    WebSocket->
        each message should be recorded with user, date and message
        
 - projects/:project_id/messages
        returns all chat messages of a project in json format
        