```mermaid
architecture-beta
    group transaction_manager(cloud)[Transaction Manager]
    group api(cloud)[API]
    group core(cloud)[Core Engine]

    service user(server)[User] 
    service server(server)[Server] in api
    service transaction_server(server)[Transaction Manager] in transaction_manager
    service db(database)[Database] in transaction_manager
    service wal(database)[WAL] in transaction_manager
    service queue_manager(server)[Queue Manager] in core
    service queue_rebalancer(server)[Queue Rebalancer] in core
    service queue_allocator(server)[Queue Allocator] in core
    service consumer_manager(server)[Consumer Manager]
    service metrics(server)[Metrics Logging]
    service message_store(database)[Message Store] in transaction_manager

    user:B --> T:server
    server:R --> L:transaction_server
    transaction_server:B --> T:db
    transaction_server:T --> B:wal
    transaction_server:R --> L:queue_manager
    queue_manager:R --> L:queue_rebalancer
    queue_manager:B --> T:queue_allocator
    queue_manager:R --> L:consumer_manager
    queue_manager:T --> B:message_store
    queue_manager:L --> L:metrics
 
```