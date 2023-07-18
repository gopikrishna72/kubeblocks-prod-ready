#
# Licensed to the Apache Software Foundation (ASF) under one
# or more contributor license agreements.  See the NOTICE file
# distributed with this work for additional information
# regarding copyright ownership.  The ASF licenses this file
# to you under the Apache License, Version 2.0 (the
# "License"); you may not use this file except in compliance
# with the License.  You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing,
# software distributed under the License is distributed on an
# "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
# KIND, either express or implied.  See the License for the
# specific language governing permissions and limitations
# under the License.
#

### --- General broker settings --- ###

# The metadata store URL
# Examples:
# * zk:my-zk-1:2181,my-zk-2:2181,my-zk-3:2181
# * my-zk-1:2181,my-zk-2:2181,my-zk-3:2181 (will default to ZooKeeper when the schema is not specified)
# * zk:my-zk-1:2181,my-zk-2:2181,my-zk-3:2181/my-chroot-path (to add a ZK chroot path)
metadataStoreUrl=

# Configuration file path for metadata store. It's supported by RocksdbMetadataStore and EtcdMetadataStore for now
metadataStoreConfigPath=

# Event topic to sync metadata between separate pulsar clusters on different cloud platforms.
metadataSyncEventTopic=

# Event topic to sync configuration-metadata between separate pulsar clusters on different cloud platforms.
configurationMetadataSyncEventTopic=

# The metadata store URL for the configuration data. If empty, we fall back to use metadataStoreUrl
configurationMetadataStoreUrl=

# Broker data port
brokerServicePort=6650

# Broker data port for TLS - By default TLS is disabled
brokerServicePortTls=

# Port to use to server HTTP request
webServicePort=8080

# Port to use to server HTTPS request - By default TLS is disabled
webServicePortTls=

# Specify the tls protocols the broker's web service will use to negotiate during TLS handshake
# (a comma-separated list of protocol names).
# Examples:- [TLSv1.3, TLSv1.2]
webServiceTlsProtocols=

# Specify the tls cipher the broker will use to negotiate during TLS Handshake
# (a comma-separated list of ciphers).
# Examples:- [TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256]
webServiceTlsCiphers=

# Hostname or IP address the service binds on, default is 0.0.0.0.
bindAddress=0.0.0.0

# Extra bind addresses for the service: <listener_name>:<scheme>://<host>:<port>,[...]
bindAddresses=

# Hostname or IP address the service advertises to the outside world. If not set, the value of InetAddress.getLocalHost().getHostName() is used.
advertisedAddress=

# Used to specify multiple advertised listeners for the broker.
# The value must format as <listener_name>:pulsar://<host>:<port>,
# multiple listeners should separate with commas.
# Do not use this configuration with advertisedAddress and brokerServicePort.
# The Default value is absent means use advertisedAddress and brokerServicePort.
# advertisedListeners=

# Used to specify the internal listener name for the broker.
# The listener name must contain in the advertisedListeners.
# The Default value is absent, the broker uses the first listener as the internal listener.
# internalListenerName=

# Enable or disable the HAProxy protocol.
haProxyProtocolEnabled=false

# Number of threads to config Netty Acceptor. Default is 1
numAcceptorThreads=

# Number of threads to use for Netty IO. Default is set to 2 * Runtime.getRuntime().availableProcessors()
numIOThreads=

# Number of threads to use for ordered executor. The ordered executor is used to operate with zookeeper,
# such as init zookeeper client, get namespace policies from zookeeper etc. It also used to split bundle. Default is 8
numOrderedExecutorThreads=8

# Number of threads to use for HTTP requests processing. Default is set to 2 * Runtime.getRuntime().availableProcessors()
numHttpServerThreads=

# Number of thread pool size to use for pulsar broker service.
# The executor in thread pool will do basic broker operation like load/unload bundle, update managedLedgerConfig,
# update topic/subscription/replicator message dispatch rate, do leader election etc.
# Default is Runtime.getRuntime().availableProcessors()
numExecutorThreadPoolSize=

# Number of thread pool size to use for pulsar zookeeper callback service
# The cache executor thread pool is used for restarting global zookeeper session.
# Default is 10
numCacheExecutorThreadPoolSize=10

# Option to enable busy-wait settings. Default is false.
# WARNING: This option will enable spin-waiting on executors and IO threads in order to reduce latency during
# context switches. The spinning will consume 100% CPU even when the broker is not doing any work. It is recommended to
# reduce the number of IO threads and BK client threads to only have few CPU cores busy.
enableBusyWait=false

# Flag to control features that are meant to be used when running in standalone mode
isRunningStandalone=

# Name of the cluster to which this broker belongs to
clusterName=

# The maximum number of tenants that each pulsar cluster can create
# This configuration is not precise control, in a concurrent scenario, the threshold will be exceeded
maxTenants=0

# Enable cluster's failure-domain which can distribute brokers into logical region
failureDomainsEnabled=false

# Metadata store session timeout in milliseconds
metadataStoreSessionTimeoutMillis=30000

# Metadata store operation timeout in seconds
metadataStoreOperationTimeoutSeconds=30

# Metadata store cache expiry time in seconds
metadataStoreCacheExpirySeconds=300

# Time to wait for broker graceful shutdown. After this time elapses, the process will be killed
brokerShutdownTimeoutMs=60000

# Flag to skip broker shutdown when broker handles Out of memory error
skipBrokerShutdownOnOOM=false

# Factory class-name to create topic with custom workflow
topicFactoryClassName=

# Enable backlog quota check. Enforces action on topic when the quota is reached
backlogQuotaCheckEnabled=true

# How often to check for topics that have reached the quota
backlogQuotaCheckIntervalInSeconds=60

# Default per-topic backlog quota limit, less than 0 means no limitation. default is -1.
backlogQuotaDefaultLimitBytes=-1

# Default per-topic backlog quota time limit in second, less than 0 means no limitation. default is -1.
backlogQuotaDefaultLimitSecond=-1

# Default backlog quota retention policy. Default is producer_request_hold
# 'producer_request_hold' Policy which holds producer's send request until the resource becomes available (or holding times out)
# 'producer_exception' Policy which throws javax.jms.ResourceAllocationException to the producer
# 'consumer_backlog_eviction' Policy which evicts the oldest message from the slowest consumer's backlog
backlogQuotaDefaultRetentionPolicy=producer_request_hold

# Default ttl for namespaces if ttl is not already configured at namespace policies. (disable default-ttl with value 0)
ttlDurationDefaultInSeconds=0

# Enable topic auto creation if new producer or consumer connected (disable auto creation with value false)
allowAutoTopicCreation=true

# The type of topic that is allowed to be automatically created.(partitioned/non-partitioned)
allowAutoTopicCreationType=partitioned

# Enable subscription auto creation if new consumer connected (disable auto creation with value false)
allowAutoSubscriptionCreation=true

# The number of partitioned topics that is allowed to be automatically created if allowAutoTopicCreationType is partitioned.
defaultNumPartitions=1

# Enable the deletion of inactive topics. This parameter need to cooperate with the allowAutoTopicCreation parameter.
# If brokerDeleteInactiveTopicsEnabled is set to true, we should ensure that allowAutoTopicCreation is also set to true.
brokerDeleteInactiveTopicsEnabled=false

# How often to check for inactive topics
brokerDeleteInactiveTopicsFrequencySeconds=60

# Set the inactive topic delete mode. Default is delete_when_no_subscriptions
# 'delete_when_no_subscriptions' mode only delete the topic which has no subscriptions and no active producers
# 'delete_when_subscriptions_caught_up' mode only delete the topic that all subscriptions has no backlogs(caught up)
# and no active producers/consumers
brokerDeleteInactiveTopicsMode=delete_when_no_subscriptions

# Metadata of inactive partitioned topic will not be cleaned up automatically by default.
# Note: If `allowAutoTopicCreation` and this option are enabled at the same time,
# it may appear that a partitioned topic has just been deleted but is automatically created as a non-partitioned topic.
brokerDeleteInactivePartitionedTopicMetadataEnabled=false

# Max duration of topic inactivity in seconds, default is not present
# If not present, 'brokerDeleteInactiveTopicsFrequencySeconds' will be used
# Topics that are inactive for longer than this value will be deleted
brokerDeleteInactiveTopicsMaxInactiveDurationSeconds=

# Allow you to delete a tenant forcefully.
forceDeleteTenantAllowed=false

# Allow you to delete a namespace forcefully.
forceDeleteNamespaceAllowed=false

# Max pending publish requests per connection to avoid keeping large number of pending
# requests in memory. Default: 1000
maxPendingPublishRequestsPerConnection=1000

# How frequently to proactively check and purge expired messages
messageExpiryCheckIntervalInMinutes=5

# How long to delay rewinding cursor and dispatching messages when active consumer is changed
activeConsumerFailoverDelayTimeMillis=1000

# How long to delete inactive subscriptions from last consuming
# When it is 0, inactive subscriptions are not deleted automatically
subscriptionExpirationTimeMinutes=0

# Enable subscription message redelivery tracker to send redelivery count to consumer (default is enabled)
subscriptionRedeliveryTrackerEnabled=true

# How frequently to proactively check and purge expired subscription
subscriptionExpiryCheckIntervalInMinutes=5

# Enable subscription types (default is all type enabled)
# SubscriptionTypes : Exclusive,Shared,Failover,Key_Shared
# Example : Exclusive,Shared
# Above example will disable Failover and Key_Shared subscription types
subscriptionTypesEnabled=Exclusive,Shared,Failover,Key_Shared

# On KeyShared subscriptions, with default AUTO_SPLIT mode, use splitting ranges or
# consistent hashing to reassign keys to new consumers
subscriptionKeySharedUseConsistentHashing=true

# On KeyShared subscriptions, number of points in the consistent-hashing ring.
# The higher the number, the more equal the assignment of keys to consumers
subscriptionKeySharedConsistentHashingReplicaPoints=100

# Maximum time in ms for a Analise backlog operation to complete
subscriptionBacklogScanMaxTimeMs=120000

# Maximum number of entries to be read within a Analise backlog operation
subscriptionBacklogScanMaxEntries=10000

# Set the default behavior for message deduplication in the broker
# This can be overridden per-namespace. If enabled, broker will reject
# messages that were already stored in the topic
brokerDeduplicationEnabled=false

# Maximum number of producer information that it's going to be
# persisted for deduplication purposes
brokerDeduplicationMaxNumberOfProducers=10000

# How often is the thread pool scheduled to check whether a snapshot needs to be taken.(disable with value 0)
brokerDeduplicationSnapshotFrequencyInSeconds=120
# If this time interval is exceeded, a snapshot will be taken.
# It will run simultaneously with `brokerDeduplicationEntriesInterval`
brokerDeduplicationSnapshotIntervalSeconds=120

# Number of entries after which a dedup info snapshot is taken.
# A larger interval will lead to fewer snapshots being taken, though it would
# increase the topic recovery time when the entries published after the
# snapshot need to be replayed.
brokerDeduplicationEntriesInterval=1000

# Time of inactivity after which the broker will discard the deduplication information
# relative to a disconnected producer. Default is 6 hours.
brokerDeduplicationProducerInactivityTimeoutMinutes=360

# When a namespace is created without specifying the number of bundle, this
# value will be used as the default
defaultNumberOfNamespaceBundles=4

# The maximum number of namespaces that each tenant can create
# This configuration is not precise control, in a concurrent scenario, the threshold will be exceeded
maxNamespacesPerTenant=0

# Max number of topics allowed to be created in the namespace. When the topics reach the max topics of the namespace,
# the broker should reject the new topic request(include topic auto-created by the producer or consumer)
# until the number of connected consumers decrease.
# Using a value of 0, is disabling maxTopicsPerNamespace-limit check.
maxTopicsPerNamespace=0

# The maximum number of connections in the broker. If it exceeds, new connections are rejected.
brokerMaxConnections=0

# The maximum number of connections per IP. If it exceeds, new connections are rejected.
brokerMaxConnectionsPerIp=0

# Allow schema to be auto updated at broker level. User can override this by
# 'is_allow_auto_update_schema' of namespace policy.
isAllowAutoUpdateSchemaEnabled=true

# Enable check for minimum allowed client library version
clientLibraryVersionCheckEnabled=false

# Path for the file used to determine the rotation status for the broker when responding
# to service discovery health checks
statusFilePath=/pulsar/status

# If true, (and ModularLoadManagerImpl is being used), the load manager will attempt to
# use only brokers running the latest software version (to minimize impact to bundles)
preferLaterVersions=false

# Max number of unacknowledged messages allowed to receive messages by a consumer on a shared subscription. Broker will stop sending
# messages to consumer once, this limit reaches until consumer starts acknowledging messages back.
# Using a value of 0, is disabling unackeMessage limit check and consumer can receive messages without any restriction
maxUnackedMessagesPerConsumer=10000

# Max number of unacknowledged messages allowed per shared subscription. Broker will stop dispatching messages to
# all consumers of the subscription once this limit reaches until consumer starts acknowledging messages back and
# unack count reaches to limit/2. Using a value of 0, is disabling unackedMessage-limit
# check and dispatcher can dispatch messages without any restriction
maxUnackedMessagesPerSubscription=50000

# Max number of unacknowledged messages allowed per broker. Once this limit reaches, broker will stop dispatching
# messages to all shared subscription which has higher number of unack messages until subscriptions start
# acknowledging messages back and unack count reaches to limit/2. Using a value of 0, is disabling
# unackedMessage-limit check and broker doesn't block dispatchers
maxUnackedMessagesPerBroker=0

# Once broker reaches maxUnackedMessagesPerBroker limit, it blocks subscriptions which has higher unacked messages
# than this percentage limit and subscription will not receive any new messages until that subscription acks back
# limit/2 messages
maxUnackedMessagesPerSubscriptionOnBrokerBlocked=0.16

# Broker periodically checks if subscription is stuck and unblock if flag is enabled. (Default is disabled)
unblockStuckSubscriptionEnabled=false

# Tick time to schedule task that checks topic publish rate limiting across all topics
# Reducing to lower value can give more accuracy while throttling publish but
# it uses more CPU to perform frequent check. (Disable publish throttling with value 0)
topicPublisherThrottlingTickTimeMillis=10

# Enable precise rate limit for topic publish
preciseTopicPublishRateLimiterEnable=false

# Tick time to schedule task that checks broker publish rate limiting across all topics
# Reducing to lower value can give more accuracy while throttling publish but
# it uses more CPU to perform frequent check. (Disable publish throttling with value 0)
brokerPublisherThrottlingTickTimeMillis=50

# Max Rate(in 1 seconds) of Message allowed to publish for a broker if broker publish rate limiting enabled
# (Disable message rate limit with value 0)
brokerPublisherThrottlingMaxMessageRate=0

# Max Rate(in 1 seconds) of Byte allowed to publish for a broker if broker publish rate limiting enabled.
# (Disable byte rate limit with value 0)
brokerPublisherThrottlingMaxByteRate=0

# Max Rate(in 1 seconds) of Message allowed to publish for a topic if topic publish rate limiting enabled
# (Disable byte rate limit with value 0)
maxPublishRatePerTopicInMessages=0

#Max Rate(in 1 seconds) of Byte allowed to publish for a topic if topic publish rate limiting enabled.
# (Disable byte rate limit with value 0)
maxPublishRatePerTopicInBytes=0

# Too many subscribe requests from a consumer can cause broker rewinding consumer cursors and loading data from bookies,
# hence causing high network bandwidth usage
# When the positive value is set, broker will throttle the subscribe requests for one consumer.
# Otherwise, the throttling will be disabled. The default value of this setting is 0 - throttling is disabled.
subscribeThrottlingRatePerConsumer=0

# Rate period for {subscribeThrottlingRatePerConsumer}. Default is 30s.
subscribeRatePeriodPerConsumerInSecond=30

# Default messages per second dispatch throttling-limit for whole broker. Using a value of 0, is disabling default
# message dispatch-throttling
dispatchThrottlingRateInMsg=0

# Default bytes per second dispatch throttling-limit for whole broker. Using a value of 0, is disabling
# default message-byte dispatch-throttling
dispatchThrottlingRateInByte=0

# Default messages per second dispatch throttling-limit for every topic. Using a value of 0, is disabling default
# message dispatch-throttling
dispatchThrottlingRatePerTopicInMsg=0

# Default bytes per second dispatch throttling-limit for every topic. Using a value of 0, is disabling
# default message-byte dispatch-throttling
dispatchThrottlingRatePerTopicInByte=0

# Apply dispatch rate limiting on batch message instead individual
# messages with in batch message. (Default is disabled)
dispatchThrottlingOnBatchMessageEnabled=false

# Default number of message dispatching throttling-limit for a subscription.
# Using a value of 0, is disabling default message dispatch-throttling.
dispatchThrottlingRatePerSubscriptionInMsg=0

# Default number of message-bytes dispatching throttling-limit for a subscription.
# Using a value of 0, is disabling default message-byte dispatch-throttling.
dispatchThrottlingRatePerSubscriptionInByte=0

# Default messages per second dispatch throttling-limit for every replicator in replication.
# Using a value of 0, is disabling replication message dispatch-throttling
dispatchThrottlingRatePerReplicatorInMsg=0

# Default bytes per second dispatch throttling-limit for every replicator in replication.
# Using a value of 0, is disabling replication message-byte dispatch-throttling
dispatchThrottlingRatePerReplicatorInByte=0

# Dispatch rate-limiting relative to publish rate.
# (Enabling flag will make broker to dynamically update dispatch-rate relatively to publish-rate:
# throttle-dispatch-rate = (publish-rate + configured dispatch-rate).
dispatchThrottlingRateRelativeToPublishRate=false

# By default we enable dispatch-throttling for both caught up consumers as well as consumers who have
# backlog.
dispatchThrottlingOnNonBacklogConsumerEnabled=true

# Max number of entries to read from bookkeeper. By default it is 100 entries.
dispatcherMaxReadBatchSize=1000

# Dispatch messages and execute broker side filters in a per-subscription thread
dispatcherDispatchMessagesInSubscriptionThread=true

# Max size in bytes of entries to read from bookkeeper. By default it is 5MB.
dispatcherMaxReadSizeBytes=5242880

# Min number of entries to read from bookkeeper. By default it is 1 entries.
# When there is an error occurred on reading entries from bookkeeper, the broker
# will backoff the batch size to this minimum number."
dispatcherMinReadBatchSize=1

# Max number of entries to dispatch for a shared subscription. By default it is 20 entries.
dispatcherMaxRoundRobinBatchSize=20

# The read failure backoff initial time in milliseconds. By default it is 15s.
dispatcherReadFailureBackoffInitialTimeInMs=15000

# The read failure backoff max time in milliseconds. By default it is 60s.
dispatcherReadFailureBackoffMaxTimeInMs=60000

# The read failure backoff mandatory stop time in milliseconds. By default it is 0s.
dispatcherReadFailureBackoffMandatoryStopTimeInMs=0

# Precise dispathcer flow control according to history message number of each entry
preciseDispatcherFlowControl=true

# Class name of Pluggable entry filter that can decide whether the entry needs to be filtered
# You can use this class to decide which entries can be sent to consumers.
# Multiple classes need to be separated by commas.
entryFilterNames=

# The directory for all the entry filter implementations
entryFiltersDirectory=

# Max number of concurrent lookup request broker allows to throttle heavy incoming lookup traffic
maxConcurrentLookupRequest=1000

# Max number of concurrent topic loading request broker allows to control number of zk-operations
maxConcurrentTopicLoadRequest=1000

# Max concurrent non-persistent message can be processed per connection
maxConcurrentNonPersistentMessagePerConnection=1000

# Number of worker threads to serve non-persistent topic
numWorkerThreadsForNonPersistentTopic=

# Enable broker to load persistent topics
enablePersistentTopics=true

# Enable broker to load non-persistent topics
enableNonPersistentTopics=true

# Enable to run bookie along with broker
enableRunBookieTogether=false

# Enable to run bookie autorecovery along with broker
enableRunBookieAutoRecoveryTogether=false

# Max number of producers allowed to connect to topic. Once this limit reaches, Broker will reject new producers
# until the number of connected producers decrease.
# Using a value of 0, is disabling maxProducersPerTopic-limit check.
maxProducersPerTopic=10000

# Max number of producers with the same IP address allowed to connect to topic.
# Once this limit reaches, Broker will reject new producers until the number of
# connected producers with the same IP address decrease.
# Using a value of 0, is disabling maxSameAddressProducersPerTopic-limit check.
maxSameAddressProducersPerTopic=0

# Enforce producer to publish encrypted messages.(default disable).
encryptionRequireOnProducer=false

# Max number of consumers allowed to connect to topic. Once this limit reaches, Broker will reject new consumers
# until the number of connected consumers decrease.
# Using a value of 0, is disabling maxConsumersPerTopic-limit check.
maxConsumersPerTopic=5000

# Max number of consumers with the same IP address allowed to connect to topic.
# Once this limit reaches, Broker will reject new consumers until the number of
# connected consumers with the same IP address decrease.
# Using a value of 0, is disabling maxSameAddressConsumersPerTopic-limit check.
maxSameAddressConsumersPerTopic=0

# Max number of subscriptions allowed to subscribe to topic. Once this limit reaches, broker will reject
# new subscription until the number of subscribed subscriptions decrease.
# Using a value of 0, is disabling maxSubscriptionsPerTopic limit check.
maxSubscriptionsPerTopic=0

# Max number of consumers allowed to connect to subscription. Once this limit reaches, Broker will reject new consumers
# until the number of connected consumers decrease.
# Using a value of 0, is disabling maxConsumersPerSubscription-limit check.
maxConsumersPerSubscription=5000

# Max size of messages.
maxMessageSize=5242880

# Interval between checks to see if topics with compaction policies need to be compacted
brokerServiceCompactionMonitorIntervalInSeconds=60

# The estimated backlog size is greater than this threshold, compression will be triggered.
# Using a value of 0, is disabling compression check.
brokerServiceCompactionThresholdInBytes=0

# Timeout for the compaction phase one loop.
# If the execution time of the compaction phase one loop exceeds this time, the compaction will not proceed.
brokerServiceCompactionPhaseOneLoopTimeInSeconds=30

# Whether to enable the delayed delivery for messages.
# If disabled, messages will be immediately delivered and there will
# be no tracking overhead.
delayedDeliveryEnabled=true

# Control the tick time for when retrying on delayed delivery,
# affecting the accuracy of the delivery time compared to the scheduled time.
# Note that this time is used to configure the HashedWheelTimer's tick time for the
# InMemoryDelayedDeliveryTrackerFactory (the default DelayedDeliverTrackerFactory).
# Default is 1 second.
delayedDeliveryTickTimeMillis=1000

# When using the InMemoryDelayedDeliveryTrackerFactory (the default DelayedDeliverTrackerFactory), whether
# the deliverAt time is strictly followed. When false (default), messages may be sent to consumers before the deliverAt
# time by as much as the tickTimeMillis. This can reduce the overhead on the broker of maintaining the delayed index
# for a potentially very short time period. When true, messages will not be sent to consumer until the deliverAt time
# has passed, and they may be as late as the deliverAt time plus the tickTimeMillis for the topic plus the
# delayedDeliveryTickTimeMillis.
isDelayedDeliveryDeliverAtTimeStrict=false

# Size of the lookahead window to use when detecting if all the messages in the topic
# have a fixed delay.
# Default is 50,000. Setting the lookahead window to 0 will disable the logic to handle
# fixed delays in messages in a different way.
delayedDeliveryFixedDelayDetectionLookahead=50000

# Whether to enable acknowledge of batch local index.
acknowledgmentAtBatchIndexLevelEnabled=false

# Enable tracking of replicated subscriptions state across clusters.
enableReplicatedSubscriptions=false

# Frequency of snapshots for replicated subscriptions tracking.
replicatedSubscriptionsSnapshotFrequencyMillis=1000

# Timeout for building a consistent snapshot for tracking replicated subscriptions state.
replicatedSubscriptionsSnapshotTimeoutSeconds=30

# Max number of snapshot to be cached per subscription.
replicatedSubscriptionsSnapshotMaxCachedPerSubscription=10

# Max memory size for broker handling messages sending from producers.
# If the processing message size exceed this value, broker will stop read data
# from the connection. The processing messages means messages are sends to broker
# but broker have not send response to client, usually waiting to write to bookies.
# It's shared across all the topics running in the same broker.
# Use -1 to disable the memory limitation. Default is 1/2 of direct memory.
maxMessagePublishBufferSizeInMB=

# Check between intervals to see if consumed ledgers need to be trimmed
# Use 0 or negative number to disable the check
retentionCheckIntervalInSeconds=120

# Control the frequency of checking the max message size in the topic policy.
# The default interval is 60 seconds.
maxMessageSizeCheckIntervalInSeconds=60

# Max number of partitions per partitioned topic
# Use 0 or negative number to disable the check
maxNumPartitionsPerPartitionedTopic=2048

# There are two policies to apply when broker metadata session expires: session expired happens, "shutdown" or "reconnect".
# With "shutdown", the broker will be restarted.
# With "reconnect", the broker will keep serving the topics, while attempting to recreate a new session.
zookeeperSessionExpiredPolicy=reconnect

# Enable or disable system topic
systemTopicEnabled=true

# Deploy the schema compatibility checker for a specific schema type to enforce schema compatibility check
schemaRegistryCompatibilityCheckers=org.apache.pulsar.broker.service.schema.JsonSchemaCompatibilityCheck,org.apache.pulsar.broker.service.schema.AvroSchemaCompatibilityCheck,org.apache.pulsar.broker.service.schema.ProtobufNativeSchemaCompatibilityCheck

# The schema compatibility strategy is used for system topics.
# Available values: ALWAYS_INCOMPATIBLE, ALWAYS_COMPATIBLE, BACKWARD, FORWARD, FULL, BACKWARD_TRANSITIVE, FORWARD_TRANSITIVE, FULL_TRANSITIVE
systemTopicSchemaCompatibilityStrategy=ALWAYS_COMPATIBLE

# Enable or disable topic level policies, topic level policies depends on the system topic
# Please enable the system topic first.
topicLevelPoliciesEnabled=true

# If a topic remains fenced for this number of seconds, it will be closed forcefully.
# If it is set to 0 or a negative number, the fenced topic will not be closed.
topicFencingTimeoutSeconds=0

### --- Authentication --- ###
# Role names that are treated as "proxy roles". If the broker sees a request with
#role as proxyRoles - it will demand to see a valid original principal.
proxyRoles=

# If this flag is set then the broker authenticates the original Auth data
# else it just accepts the originalPrincipal and authorizes it (if required).
authenticateOriginalAuthData=false

# Tls cert refresh duration in seconds (set 0 to check on every new connection)
tlsCertRefreshCheckDurationSec=300

# Path for the TLS certificate file
tlsCertificateFilePath=

# Path for the TLS private key file
tlsKeyFilePath=

# Path for the trusted TLS certificate file.
# This cert is used to verify that any certs presented by connecting clients
# are signed by a certificate authority. If this verification
# fails, then the certs are untrusted and the connections are dropped.
tlsTrustCertsFilePath=

# Accept untrusted TLS certificate from client.
# If true, a client with a cert which cannot be verified with the
# 'tlsTrustCertsFilePath' cert will allowed to connect to the server,
# though the cert will not be used for client authentication.
tlsAllowInsecureConnection=false

# Whether the hostname is validated when the broker creates a TLS connection with other brokers
tlsHostnameVerificationEnabled=false

# Specify the tls protocols the broker will use to negotiate during TLS handshake
# (a comma-separated list of protocol names).
# Examples:- [TLSv1.3, TLSv1.2]
tlsProtocols=

# Specify the tls cipher the broker will use to negotiate during TLS Handshake
# (a comma-separated list of ciphers).
# Examples:- [TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256]
tlsCiphers=

# Trusted client certificates are required for to connect TLS
# Reject the Connection if the Client Certificate is not trusted.
# In effect, this requires that all connecting clients perform TLS client
# authentication.
tlsRequireTrustedClientCertOnConnect=false

# Specify the TLS provider for the broker service:
# When using TLS authentication with CACert, the valid value is either OPENSSL or JDK.
# When using TLS authentication with KeyStore, available values can be SunJSSE, Conscrypt and etc.
tlsProvider=

# Specify the TLS provider for the web service: SunJSSE, Conscrypt and etc.
webServiceTlsProvider=Conscrypt

### --- KeyStore TLS config variables --- ###
## Note that some of the above TLS configs also apply to the KeyStore TLS configuration.

# Enable TLS with KeyStore type configuration in broker.
tlsEnabledWithKeyStore=false

# TLS KeyStore type configuration in broker: JKS, PKCS12
tlsKeyStoreType=JKS

# TLS KeyStore path in broker
tlsKeyStore=

# TLS KeyStore password for broker
tlsKeyStorePassword=

# TLS TrustStore type configuration in broker: JKS, PKCS12
tlsTrustStoreType=JKS

# TLS TrustStore path in broker
tlsTrustStore=

# TLS TrustStore password in broker, default value is empty password
tlsTrustStorePassword=

# Whether internal client uses TLS to connection the broker
brokerClientTlsEnabled=false

# Whether internal client use KeyStore type to authenticate with Pulsar brokers
brokerClientTlsEnabledWithKeyStore=false

# The TLS Provider used by internal client to authenticate with other Pulsar brokers
brokerClientSslProvider=

# TLS trusted certificate file for internal client,
# used by the internal client to authenticate with Pulsar brokers
brokerClientTrustCertsFilePath=

# TLS private key file for internal client,
# used by the internal client to authenticate with Pulsar brokers
brokerClientKeyFilePath=

# TLS certificate file for internal client,
# used by the internal client to authenticate with Pulsar brokers
brokerClientCertificateFilePath=

# TLS KeyStore type configuration for internal client: JKS, PKCS12
# used by the internal client to authenticate with Pulsar brokers
brokerClientTlsKeyStoreType=JKS

# TLS KeyStore path for internal client
# used by the internal client to authenticate with Pulsar brokers
brokerClientTlsKeyStore=

# TLS KeyStore password for internal client,
# used by the internal client to authenticate with Pulsar brokers
brokerClientTlsKeyStorePassword=

# TLS TrustStore type configuration for internal client: JKS, PKCS12
# used by the internal client to authenticate with Pulsar brokers
brokerClientTlsTrustStoreType=JKS

# TLS TrustStore path for internal client
# used by the internal client to authenticate with Pulsar brokers
brokerClientTlsTrustStore=

# TLS TrustStore password for internal client,
# used by the internal client to authenticate with Pulsar brokers
brokerClientTlsTrustStorePassword=

# Specify the tls cipher the internal client will use to negotiate during TLS Handshake
# (a comma-separated list of ciphers)
# e.g.  [TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256].
# used by the internal client to authenticate with Pulsar brokers
brokerClientTlsCiphers=

# Specify the tls protocols the broker will use to negotiate during TLS handshake
# (a comma-separated list of protocol names).
# e.g.  [TLSv1.3, TLSv1.2]
# used by the internal client to authenticate with Pulsar brokers
brokerClientTlsProtocols=

# You can add extra configuration options for the Pulsar Client and the Pulsar Admin Client
# by prefixing them with "brokerClient_". These configurations are applied after hard coded configuration
# and before the above brokerClient configurations named above.

### --- Metadata Store --- ###

# Whether we should enable metadata operations batching
metadataStoreBatchingEnabled=true

# Maximum delay to impose on batching grouping
metadataStoreBatchingMaxDelayMillis=5

# Maximum number of operations to include in a singular batch
metadataStoreBatchingMaxOperations=1000

# Maximum size of a batch
metadataStoreBatchingMaxSizeKb=128


### --- Authentication --- ###

# Enable authentication
authenticationEnabled=false

# Authentication provider name list, which is comma separated list of class names
authenticationProviders=

# Interval of time for checking for expired authentication credentials
authenticationRefreshCheckSeconds=60

# Enforce authorization
authorizationEnabled=false

# Authorization provider fully qualified class-name
authorizationProvider=org.apache.pulsar.broker.authorization.PulsarAuthorizationProvider

# Allow wildcard matching in authorization
# (wildcard matching only applicable if wildcard-char:
# * presents at first or last position eg: *.pulsar.service, pulsar.service.*)
authorizationAllowWildcardsMatching=false

# Role names that are treated as "super-user", meaning they will be able to do all admin
# operations and publish/consume from all topics
superUserRoles=

# Authentication settings of the broker itself. Used when the broker connects to other brokers,
# either in same or other clusters
brokerClientAuthenticationPlugin=
brokerClientAuthenticationParameters=

# Supported Athenz provider domain names(comma separated) for authentication
athenzDomainNames=

# When this parameter is not empty, unauthenticated users perform as anonymousUserRole
anonymousUserRole=

## Configure the datasource of basic authenticate, supports the file and Base64 format.
# file:
# basicAuthConf=/path/my/.htpasswd
# use Base64 to encode the contents of .htpasswd:
# basicAuthConf=YOUR-BASE64-DATA
basicAuthConf=

### --- Token Authentication Provider --- ###

## Symmetric key
# Configure the secret key to be used to validate auth tokens
# The key can be specified like:
# tokenSecretKey=data:;base64,xxxxxxxxx
# tokenSecretKey=file:///my/secret.key    ( Note: key file must be DER-encoded )
tokenSecretKey=

## Asymmetric public/private key pair
# Configure the public key to be used to validate auth tokens
# The key can be specified like:
# tokenPublicKey=data:;base64,xxxxxxxxx
# tokenPublicKey=file:///my/public.key    ( Note: key file must be DER-encoded )
tokenPublicKey=

# The token "claim" that will be interpreted as the authentication "role" or "principal" by AuthenticationProviderToken (defaults to "sub" if blank)
tokenAuthClaim=

# The token audience "claim" name, e.g. "aud", that will be used to get the audience from token.
# If not set, audience will not be verified.
tokenAudienceClaim=

# The token audience stands for this broker. The field `tokenAudienceClaim` of a valid token, need contains this.
tokenAudience=

### --- SASL Authentication Provider --- ###

# This is a regexp, which limits the range of possible ids which can connect to the Broker using SASL.
# Default value: `SaslConstants.JAAS_CLIENT_ALLOWED_IDS_DEFAULT`, which is ".*pulsar.*",
# so only clients whose id contains 'pulsar' are allowed to connect.
saslJaasClientAllowedIds=.*pulsar.*

# Service Principal, for login context name.
# Default value `SaslConstants.JAAS_DEFAULT_BROKER_SECTION_NAME`, which is "PulsarBroker".
saslJaasServerSectionName=PulsarBroker

# Path to file containing the secret to be used to SaslRoleTokenSigner
# The Path can be specified like:
# saslJaasServerRoleTokenSignerSecretPath=file:///my/saslRoleTokenSignerSecret.key
saslJaasServerRoleTokenSignerSecretPath=

### --- HTTP Server config --- ###

# If >0, it will reject all HTTP requests with bodies larged than the configured limit
httpMaxRequestSize=-1

# If true, the broker will reject all HTTP requests using the TRACE and TRACK verbs.
# This setting may be necessary if the broker is deployed into an environment that uses http port
# scanning and flags web servers allowing the TRACE method as insecure.
disableHttpDebugMethods=false

# Enable the enforcement of limits on the incoming HTTP requests
httpRequestsLimitEnabled=false

# Max HTTP requests per seconds allowed. The excess of requests will be rejected with HTTP code 429 (Too many requests)
httpRequestsMaxPerSecond=100.0

# Capacity for thread pool queue in the HTTP server
httpServerThreadPoolQueueSize=8192

# Capacity for accept queue in the HTTP server
httpServerAcceptQueueSize=8192

# Maximum number of inbound http connections. (0 to disable limiting)
maxHttpServerConnections=2048

# Max concurrent web requests
maxConcurrentHttpRequests=1024

### --- BookKeeper Client --- ###

# Metadata service uri that bookkeeper is used for loading corresponding metadata driver
# and resolving its metadata service location.
# This value can be fetched using `bookkeeper shell whatisinstanceid` command in BookKeeper cluster.
# For example: zk+hierarchical://localhost:2181/ledgers
# The metadata service uri list can also be semicolon separated values like below:
# zk+hierarchical://zk1:2181;zk2:2181;zk3:2181/ledgers
bookkeeperMetadataServiceUri=

# Authentication plugin to use when connecting to bookies
bookkeeperClientAuthenticationPlugin=

# BookKeeper auth plugin implementatation specifics parameters name and values
bookkeeperClientAuthenticationParametersName=
bookkeeperClientAuthenticationParameters=

# Timeout for BK add / read operations
bookkeeperClientTimeoutInSeconds=30

# Number of BookKeeper client worker threads
# Default is Runtime.getRuntime().availableProcessors()
bookkeeperClientNumWorkerThreads=

# Number of BookKeeper client IO threads
# Default is Runtime.getRuntime().availableProcessors() * 2
bookkeeperClientNumIoThreads=

# Use separated IO threads for BookKeeper client
# Default is false, which will use Pulsar IO threads
bookkeeperClientSeparatedIoThreadsEnabled=false

# Speculative reads are initiated if a read request doesn't complete within a certain time
# Using a value of 0, is disabling the speculative reads
bookkeeperClientSpeculativeReadTimeoutInMillis=0

# Number of channels per bookie
bookkeeperNumberOfChannelsPerBookie=16

# Use older Bookkeeper wire protocol with bookie
bookkeeperUseV2WireProtocol=true

# Enable bookies health check. Bookies that have more than the configured number of failure within
# the interval will be quarantined for some time. During this period, new ledgers won't be created
# on these bookies
bookkeeperClientHealthCheckEnabled=true
bookkeeperClientHealthCheckIntervalSeconds=60
bookkeeperClientHealthCheckErrorThresholdPerInterval=5
bookkeeperClientHealthCheckQuarantineTimeInSeconds=1800

#bookie quarantine ratio to avoid all clients quarantine the high pressure bookie servers at the same time
bookkeeperClientQuarantineRatio=1.0

# Specify options for the GetBookieInfo check. These settings can be useful
# to help ensure the list of bookies is up to date on the brokers.
bookkeeperClientGetBookieInfoIntervalSeconds=86400
bookkeeperClientGetBookieInfoRetryIntervalSeconds=60

# Enable rack-aware bookie selection policy. BK will chose bookies from different racks when
# forming a new bookie ensemble
# This parameter related to ensemblePlacementPolicy in conf/bookkeeper.conf, if enabled, ensemblePlacementPolicy
# should be set to org.apache.bookkeeper.client.RackawareEnsemblePlacementPolicy
bookkeeperClientRackawarePolicyEnabled=true

# Enable region-aware bookie selection policy. BK will chose bookies from
# different regions and racks when forming a new bookie ensemble
# If enabled, the value of bookkeeperClientRackawarePolicyEnabled is ignored
# This parameter related to ensemblePlacementPolicy in conf/bookkeeper.conf, if enabled, ensemblePlacementPolicy
# should be set to org.apache.bookkeeper.client.RegionAwareEnsemblePlacementPolicy
bookkeeperClientRegionawarePolicyEnabled=false

# Minimum number of racks per write quorum. BK rack-aware bookie selection policy will try to
# get bookies from at least 'bookkeeperClientMinNumRacksPerWriteQuorum' racks for a write quorum.
bookkeeperClientMinNumRacksPerWriteQuorum=2

# Enforces rack-aware bookie selection policy to pick bookies from 'bookkeeperClientMinNumRacksPerWriteQuorum'
# racks for a writeQuorum.
# If BK can't find bookie then it would throw BKNotEnoughBookiesException instead of picking random one.
bookkeeperClientEnforceMinNumRacksPerWriteQuorum=false

# Enable/disable reordering read sequence on reading entries.
bookkeeperClientReorderReadSequenceEnabled=true

# Enable bookie isolation by specifying a list of bookie groups to choose from. Any bookie
# outside the specified groups will not be used by the broker
bookkeeperClientIsolationGroups=

# Enable bookie secondary-isolation group if bookkeeperClientIsolationGroups doesn't
# have enough bookie available.
bookkeeperClientSecondaryIsolationGroups=

# Minimum bookies that should be available as part of bookkeeperClientIsolationGroups
# else broker will include bookkeeperClientSecondaryIsolationGroups bookies in isolated list.
bookkeeperClientMinAvailableBookiesInIsolationGroups=

# Enable/disable having read operations for a ledger to be sticky to a single bookie.
# If this flag is enabled, the client will use one single bookie (by preference) to read
# all entries for a ledger.
bookkeeperEnableStickyReads=true

# Set the client security provider factory class name.
# Default: org.apache.bookkeeper.tls.TLSContextFactory
bookkeeperTLSProviderFactoryClass=org.apache.bookkeeper.tls.TLSContextFactory

# Enable tls authentication with bookie
bookkeeperTLSClientAuthentication=false

# Supported type: PEM, JKS, PKCS12. Default value: PEM
bookkeeperTLSKeyFileType=PEM

#Supported type: PEM, JKS, PKCS12. Default value: PEM
bookkeeperTLSTrustCertTypes=PEM

# Path to file containing keystore password, if the client keystore is password protected.
bookkeeperTLSKeyStorePasswordPath=

# Path to file containing truststore password, if the client truststore is password protected.
bookkeeperTLSTrustStorePasswordPath=

# Path for the TLS private key file
bookkeeperTLSKeyFilePath=

# Path for the TLS certificate file
bookkeeperTLSCertificateFilePath=

# Path for the trusted TLS certificate file
bookkeeperTLSTrustCertsFilePath=

# Tls cert refresh duration at bookKeeper-client in seconds (0 to disable check)
bookkeeperTlsCertFilesRefreshDurationSeconds=300

# Whether the hostname is validated when the broker creates a TLS connection to a bookkeeper
bookkeeper_tlsHostnameVerificationEnabled=false

# Enable/disable disk weight based placement. Default is false
bookkeeperDiskWeightBasedPlacementEnabled=true

# Set the interval to check the need for sending an explicit LAC
# A value of '0' disables sending any explicit LACs. Default is 0.
bookkeeperExplicitLacIntervalInMills=0

# Expose bookkeeper client managed ledger stats to prometheus. default is false
# bookkeeperClientExposeStatsToPrometheus=false

### --- Managed Ledger --- ###

# Number of bookies to use when creating a ledger
managedLedgerDefaultEnsembleSize=2

# Number of copies to store for each message
managedLedgerDefaultWriteQuorum=2

# Number of guaranteed copies (acks to wait before write is complete)
managedLedgerDefaultAckQuorum=2

# with OpportunisticStriping=true the ensembleSize is adapted automatically to writeQuorum
# in case of lack of enough bookies
#bookkeeper_opportunisticStriping=false

# You can add other configuration options for the BookKeeper client
# by prefixing them with "bookkeeper_". These configurations are applied
# to all bookkeeper clients started by the broker (including the managed ledger bookkeeper clients as well as
# the BookkeeperPackagesStorage bookkeeper client), except the distributed log bookkeeper client.
# The dlog bookkeeper client is configured in the functions worker configuration file.

# How frequently to flush the cursor positions that were accumulated due to rate limiting. (seconds).
# Default is 60 seconds
managedLedgerCursorPositionFlushSeconds=60

# How frequently to refresh the stats. (seconds). Default is 60 seconds
managedLedgerStatsPeriodSeconds=60

# Default type of checksum to use when writing to BookKeeper. Default is "CRC32C"
# Other possible options are "CRC32", "MAC" or "DUMMY" (no checksum).
managedLedgerDigestType=CRC32C

# Number of threads to be used for managed ledger scheduled tasks
managedLedgerNumSchedulerThreads=2

# Amount of memory to use for caching data payload in managed ledger. This memory
# is allocated from JVM direct memory and it's shared across all the topics
# running  in the same broker. By default, uses 1/5th of available direct memory
managedLedgerCacheSizeMB=

# Whether we should make a copy of the entry payloads when inserting in cache
managedLedgerCacheCopyEntries=false

# Threshold to which bring down the cache level when eviction is triggered
managedLedgerCacheEvictionWatermark=0.9

# Configure the cache eviction interval in milliseconds for the managed ledger cache
managedLedgerCacheEvictionIntervalMs=10

# All entries that have stayed in cache for more than the configured time, will be evicted
managedLedgerCacheEvictionTimeThresholdMillis=1000

# Configure the threshold (in number of entries) from where a cursor should be considered 'backlogged'
# and thus should be set as inactive.
managedLedgerCursorBackloggedThreshold=1000

# Rate limit the amount of writes per second generated by consumer acking the messages
managedLedgerDefaultMarkDeleteRateLimit=1.0

# Max number of entries to append to a ledger before triggering a rollover
# A ledger rollover is triggered after the min rollover time has passed
# and one of the following conditions is true:
#  * The max rollover time has been reached
#  * The max entries have been written to the ledger
#  * The max ledger size has been written to the ledger
managedLedgerMaxEntriesPerLedger=50000

# Minimum time between ledger rollover for a topic
managedLedgerMinLedgerRolloverTimeMinutes=10

# Maximum time before forcing a ledger rollover for a topic
managedLedgerMaxLedgerRolloverTimeMinutes=240

# Time to rollover ledger for inactive topic (duration without any publish on that topic)
# Disable rollover with value 0 (Default value 0)
managedLedgerInactiveLedgerRolloverTimeSeconds=0

# Maximum ledger size before triggering a rollover for a topic (MB)
managedLedgerMaxSizePerLedgerMbytes=2048

# Delay between a ledger being successfully offloaded to long term storage
# and the ledger being deleted from bookkeeper (default is 4 hours)
managedLedgerOffloadDeletionLagMs=14400000

# The number of bytes before triggering automatic offload to long term storage
# (default is -1, which is disabled)
managedLedgerOffloadAutoTriggerSizeThresholdBytes=-1

# Max number of entries to append to a cursor ledger
managedLedgerCursorMaxEntriesPerLedger=50000

# Max time before triggering a rollover on a cursor ledger
managedLedgerCursorRolloverTimeInSeconds=14400

# Max number of "acknowledgment holes" that are going to be persistently stored.
# When acknowledging out of order, a consumer will leave holes that are supposed
# to be quickly filled by acking all the messages. The information of which
# messages are acknowledged is persisted by compressing in "ranges" of messages
# that were acknowledged. After the max number of ranges is reached, the information
# will only be tracked in memory and messages will be redelivered in case of
# crashes.
managedLedgerMaxUnackedRangesToPersist=10000

# Maximum amount of memory used hold data read from storage (or from the cache).
# This mechanism prevents the broker to have too many concurrent
# reads from storage and fall into Out of Memory errors in case
# of multiple concurrent reads to multiple concurrent consumers.
# Set 0 in order to disable the feature.
#
managedLedgerMaxReadsInFlightSizeInMB=0

# Max number of "acknowledgment holes" that can be stored in MetadataStore. If number of unack message range is higher
# than this limit then broker will persist unacked ranges into bookkeeper to avoid additional data overhead into
# MetadataStore.
managedLedgerMaxUnackedRangesToPersistInMetadataStore=1000

# Skip reading non-recoverable/unreadable data-ledger under managed-ledger's list. It helps when data-ledgers gets
# corrupted at bookkeeper and managed-cursor is stuck at that ledger.
autoSkipNonRecoverableData=false

# Whether to recover cursors lazily when trying to recover a managed ledger backing a persistent topic.
# It can improve write availability of topics.
# The caveat is now when recovered ledger is ready to write we're not sure if all old consumers last mark
# delete position can be recovered or not.
lazyCursorRecovery=false

# operation timeout while updating managed-ledger metadata.
managedLedgerMetadataOperationsTimeoutSeconds=60

# Read entries timeout when broker tries to read messages from bookkeeper.
managedLedgerReadEntryTimeoutSeconds=0

# Add entry timeout when broker tries to publish message to bookkeeper (0 to disable it).
managedLedgerAddEntryTimeoutSeconds=0

# Managed ledger prometheus stats latency rollover seconds (default: 60s)
managedLedgerPrometheusStatsLatencyRolloverSeconds=60

# Whether trace managed ledger task execution time
managedLedgerTraceTaskExecution=true

# New entries check delay for the cursor under the managed ledger.
# If no new messages in the topic, the cursor will try to check again after the delay time.
# For consumption latency sensitive scenario, can set to a smaller value or set to 0.
# Of course, use a smaller value may degrade consumption throughput. Default is 10ms.
managedLedgerNewEntriesCheckDelayInMillis=10

### --- Load balancer --- ###

# Enable load balancer
loadBalancerEnabled=true

# Percentage of change to trigger load report update
loadBalancerReportUpdateThresholdPercentage=10

# minimum interval to update load report
loadBalancerReportUpdateMinIntervalMillis=5000

# maximum interval to update load report
loadBalancerReportUpdateMaxIntervalMinutes=15

# Frequency of report to collect
loadBalancerHostUsageCheckIntervalMinutes=1

# Enable/disable automatic bundle unloading for load-shedding
loadBalancerSheddingEnabled=true

# Load shedding interval. Broker periodically checks whether some traffic should be offload from
# some over-loaded broker to other under-loaded brokers
loadBalancerSheddingIntervalMinutes=1

# Prevent the same topics to be shed and moved to other broker more than once within this timeframe
loadBalancerSheddingGracePeriodMinutes=30

# Usage threshold to allocate max number of topics to broker
loadBalancerBrokerMaxTopics=50000

# Usage threshold to determine a broker as over-loaded
loadBalancerBrokerOverloadedThresholdPercentage=85

# Interval to flush dynamic resource quota to ZooKeeper
loadBalancerResourceQuotaUpdateIntervalMinutes=15

# enable/disable namespace bundle auto split
loadBalancerAutoBundleSplitEnabled=true

# enable/disable automatic unloading of split bundles
loadBalancerAutoUnloadSplitBundlesEnabled=true

# maximum topics in a bundle, otherwise bundle split will be triggered
loadBalancerNamespaceBundleMaxTopics=1000

# maximum sessions (producers + consumers) in a bundle, otherwise bundle split will be triggered
# (disable threshold check with value -1)
loadBalancerNamespaceBundleMaxSessions=1000

# maximum msgRate (in + out) in a bundle, otherwise bundle split will be triggered
loadBalancerNamespaceBundleMaxMsgRate=30000

# maximum bandwidth (in + out) in a bundle, otherwise bundle split will be triggered
loadBalancerNamespaceBundleMaxBandwidthMbytes=100

# maximum number of bundles in a namespace
loadBalancerNamespaceMaximumBundles=128

# Override the auto-detection of the network interfaces max speed.
# This option is useful in some environments (eg: EC2 VMs) where the max speed
# reported by Linux is not reflecting the real bandwidth available to the broker.
# Since the network usage is employed by the load manager to decide when a broker
# is overloaded, it is important to make sure the info is correct or override it
# with the right value here. The configured value can be a double (eg: 0.8) and that
# can be used to trigger load-shedding even before hitting on NIC limits.
loadBalancerOverrideBrokerNicSpeedGbps=

# Name of load manager to use
loadManagerClassName=org.apache.pulsar.broker.loadbalance.impl.ModularLoadManagerImpl

# Supported algorithms name for namespace bundle split.
# "range_equally_divide" divides the bundle into two parts with the same hash range size.
# "topic_count_equally_divide" divides the bundle into two parts with the same topics count.
supportedNamespaceBundleSplitAlgorithms=range_equally_divide,topic_count_equally_divide,specified_positions_divide

# Default algorithm name for namespace bundle split
defaultNamespaceBundleSplitAlgorithm=topic_count_equally_divide

# load shedding strategy, support OverloadShedder and ThresholdShedder, default is ThresholdShedder since 2.10.0
loadBalancerLoadSheddingStrategy=org.apache.pulsar.broker.loadbalance.impl.ThresholdShedder

# load balance placement strategy, support LeastLongTermMessageRate and LeastResourceUsageWithWeight
loadBalancerLoadPlacementStrategy=org.apache.pulsar.broker.loadbalance.impl.LeastLongTermMessageRate

# The broker resource usage threshold.
# When the broker resource usage is greater than the pulsar cluster average resource usage,
# the threshold shedder will be triggered to offload bundles from the broker.
# It only takes effect in the ThresholdShedder strategy.
loadBalancerBrokerThresholdShedderPercentage=10

# The broker average resource usage difference threshold.
# Average resource usage difference threshold to determine a broker whether to be a best candidate in LeastResourceUsageWithWeight.
# (eg: broker1 with 10% resource usage with weight and broker2 with 30% and broker3 with 80% will have 40% average resource usage.
# The placement strategy can select broker1 and broker2 as best candidates.)
# It only takes effect in the LeastResourceUsageWithWeight strategy.
loadBalancerAverageResourceUsageDifferenceThresholdPercentage=10

# Message-rate percentage threshold between highest and least loaded brokers for
# uniform load shedding. (eg: broker1 with 50K msgRate and broker2 with 30K msgRate
# will have 66% msgRate difference and load balancer can unload bundles from broker-1
# to broker-2)
loadBalancerMsgRateDifferenceShedderThreshold=50

# Message-throughput threshold between highest and least loaded brokers for
# uniform load shedding. (eg: broker1 with 450MB msgRate and broker2 with 100MB msgRate
# will have 4.5 times msgThroughout difference and load balancer can unload bundles
# from broker-1 to broker-2)
loadBalancerMsgThroughputMultiplierDifferenceShedderThreshold=4

# When calculating new resource usage, the history usage accounts for.
# It only takes effect in the ThresholdShedder strategy.
loadBalancerHistoryResourcePercentage=0.9

# The BandWithIn usage weight when calculating new resource usage.
# It only takes effect in the ThresholdShedder strategy.
loadBalancerBandwithInResourceWeight=1.0

# The BandWithOut usage weight when calculating new resource usage.
# It only takes effect in the ThresholdShedder strategy.
loadBalancerBandwithOutResourceWeight=1.0

# The CPU usage weight when calculating new resource usage.
# It only takes effect in the ThresholdShedder strategy.
loadBalancerCPUResourceWeight=1.0

# The heap memory usage weight when calculating new resource usage.
# It only takes effect in the ThresholdShedder strategy.
loadBalancerMemoryResourceWeight=1.0

# The direct memory usage weight when calculating new resource usage.
# It only takes effect in the ThresholdShedder strategy.
loadBalancerDirectMemoryResourceWeight=1.0

# Bundle unload minimum throughput threshold (MB), avoiding bundle unload frequently.
# It only takes effect in the ThresholdShedder strategy.
loadBalancerBundleUnloadMinThroughputThreshold=10

# Time to wait for the unloading of a namespace bundle
namespaceBundleUnloadingTimeoutMs=60000

### --- Replication --- ###

# Enable replication metrics
replicationMetricsEnabled=true

# Max number of connections to open for each broker in a remote cluster
# More connections host-to-host lead to better throughput over high-latency
# links.
replicationConnectionsPerBroker=16

# Replicator producer queue size
replicationProducerQueueSize=1000

# Replicator prefix used for replicator producer name and cursor name
replicatorPrefix=pulsar.repl

# Duration to check replication policy to avoid replicator inconsistency
# due to missing ZooKeeper watch (disable with value 0)
replicationPolicyCheckDurationSeconds=600

# Default message retention time.
# The default value is 0, which means the data is removed after all the subscriptions are consumed.
# Value less than 0 means messages never expire.
defaultRetentionTimeInMinutes=0

# Default retention size.
# The default value is 0, which means the data is removed after all the subscriptions are consumed.
# Value less than 0 means no infinite size quota.
defaultRetentionSizeInMB=0

# How often to check whether the connections are still alive
keepAliveIntervalSeconds=30

# bootstrap namespaces
bootstrapNamespaces=

### --- WebSocket --- ###

# Enable the WebSocket API service in broker
webSocketServiceEnabled=false

# Number of IO threads in Pulsar Client used in WebSocket proxy
webSocketNumIoThreads=

# Number of threads used by Websocket service
webSocketNumServiceThreads=

# Number of connections per Broker in Pulsar Client used in WebSocket proxy
webSocketConnectionsPerBroker=

# Time in milliseconds that idle WebSocket session times out
webSocketSessionIdleTimeoutMillis=300000

# The maximum size of a text message during parsing in WebSocket proxy
webSocketMaxTextFrameSize=1048576

### --- Metrics --- ###

# Whether the '/metrics' endpoint requires authentication. Defaults to false
authenticateMetricsEndpoint=false

# Enable topic level metrics
exposeTopicLevelMetricsInPrometheus=true

# Enable consumer level metrics. default is false
exposeConsumerLevelMetricsInPrometheus=false

# Enable cache metrics data, default value is false
metricsBufferResponse=false

# Enable producer level metrics. default is false
exposeProducerLevelMetricsInPrometheus=false

# Enable managed ledger metrics (aggregated by namespace). default is true
exposeManagedLedgerMetricsInPrometheus=true

# Enable cursor level metrics. default is false
exposeManagedCursorMetricsInPrometheus=false

# Classname of Pluggable JVM GC metrics logger that can log GC specific metrics
# jvmGCMetricsLoggerClassName=

# Time in milliseconds that metrics endpoint would time out. Default is 30s.
# Increase it if there are a lot of topics to expose topic-level metrics.
# Set it to 0 to disable timeout.
metricsServletTimeoutMs=30000

# Enable or disable broker bundles metrics. The default value is false.
exposeBundlesMetricsInPrometheus=false

### --- Functions --- ###

# Enable Functions Worker Service in Broker
functionsWorkerEnabled=false

#Enable Functions Worker to use packageManagement Service to store package.
functionsWorkerEnablePackageManagement=false

### --- Broker Web Stats --- ###

# Enable topic level metrics
exposePublisherStats=true
statsUpdateFrequencyInSecs=60
statsUpdateInitialDelayInSecs=60

# Enable expose the precise backlog stats.
# Set false to use published counter and consumed counter to calculate, this would be more efficient but may be inaccurate.
# Default is false.
exposePreciseBacklogInPrometheus=false

# Enable splitting topic and partition label in Prometheus.
# If enabled, a topic name will split into 2 parts, one is topic name without partition index,
# another one is partition index, e.g. (topic=xxx, partition=0).
# If the topic is a non-partitioned topic, -1 will be used for the partition index.
# If disabled, one label to represent the topic and partition, e.g. (topic=xxx-partition-0)
# Default is false.

splitTopicAndPartitionLabelInPrometheus=false

# If true and the client supports partial producer, aggregate publisher stats of PartitionedTopicStats by producerName.
# Otherwise, aggregate it by list index.
aggregatePublisherStatsByProducerName=false

### --- Schema storage --- ###
# The schema storage implementation used by this broker
schemaRegistryStorageClassName=org.apache.pulsar.broker.service.schema.BookkeeperSchemaStorageFactory

# Whether to enable schema validation.
# When schema validation is enabled, if a producer without a schema attempts to produce the message to a topic with schema, the producer is rejected and disconnected.
isSchemaValidationEnforced=false

# The schema compatibility strategy at broker level.
# Available values: ALWAYS_INCOMPATIBLE, ALWAYS_COMPATIBLE, BACKWARD, FORWARD, FULL, BACKWARD_TRANSITIVE, FORWARD_TRANSITIVE, FULL_TRANSITIVE
schemaCompatibilityStrategy=FULL

### --- Ledger Offloading --- ###

# The directory for all the offloader implementations
offloadersDirectory=./offloaders

# Driver to use to offload old data to long term storage (Possible values: aws-s3, google-cloud-storage, azureblob, aliyun-oss, filesystem)
# When using google-cloud-storage, Make sure both Google Cloud Storage and Google Cloud Storage JSON API are enabled for
# the project (check from Developers Console -> Api&auth -> APIs).
managedLedgerOffloadDriver=

# Maximum number of thread pool threads for ledger offloading
managedLedgerOffloadMaxThreads=2

# The extraction directory of the nar package.
# Available for Protocol Handler, Additional Servlets, Entry Filter, Offloaders, Broker Interceptor.
# Default is System.getProperty("java.io.tmpdir").
narExtractionDirectory=

# Maximum prefetch rounds for ledger reading for offloading
managedLedgerOffloadPrefetchRounds=1

# Use Open Range-Set to cache unacked messages
managedLedgerUnackedRangesOpenCacheSetEnabled=true

# For Amazon S3 ledger offload, AWS region
s3ManagedLedgerOffloadRegion=

# For Amazon S3 ledger offload, Bucket to place offloaded ledger into
s3ManagedLedgerOffloadBucket=

# For Amazon S3 ledger offload, Alternative endpoint to connect to (useful for testing)
s3ManagedLedgerOffloadServiceEndpoint=

# For Amazon S3 ledger offload, Max block size in bytes. (64MB by default, 5MB minimum)
s3ManagedLedgerOffloadMaxBlockSizeInBytes=67108864

# For Amazon S3 ledger offload, Read buffer size in bytes (1MB by default)
s3ManagedLedgerOffloadReadBufferSizeInBytes=1048576

# For Google Cloud Storage ledger offload, region where offload bucket is located.
# reference this page for more details: https://cloud.google.com/storage/docs/bucket-locations
gcsManagedLedgerOffloadRegion=

# For Google Cloud Storage ledger offload, Bucket to place offloaded ledger into
gcsManagedLedgerOffloadBucket=

# For Google Cloud Storage ledger offload, Max block size in bytes. (64MB by default, 5MB minimum)
gcsManagedLedgerOffloadMaxBlockSizeInBytes=67108864

# For Google Cloud Storage ledger offload, Read buffer size in bytes (1MB by default)
gcsManagedLedgerOffloadReadBufferSizeInBytes=1048576

# For Google Cloud Storage, path to json file containing service account credentials.
# For more details, see the "Service Accounts" section of https://support.google.com/googleapi/answer/6158849
gcsManagedLedgerOffloadServiceAccountKeyFile=

# For Azure BlobStore ledger offload, Blob Container (Bucket) to place offloaded ledger into
managedLedgerOffloadBucket=

#For File System Storage, file system profile path
fileSystemProfilePath=../conf/filesystem_offload_core_site.xml

#For File System Storage, file system uri
fileSystemURI=

### --- Transaction config variables --- ###

# Enable transaction coordinator in broker
transactionCoordinatorEnabled=false
transactionMetadataStoreProviderClassName=org.apache.pulsar.transaction.coordinator.impl.MLTransactionMetadataStoreProvider

# Transaction buffer takes a snapshot after the number of transaction operations reaches this value.
transactionBufferSnapshotMaxTransactionCount=1000

# Transaction buffer take snapshot interval time
# Unit : millisecond
transactionBufferSnapshotMinTimeInMillis=5000

# The max concurrent requests for transaction buffer client, default is 1000
transactionBufferClientMaxConcurrentRequests=1000

# The max active transactions per transaction coordinator, default value 0 indicates no limit.
maxActiveTransactionsPerCoordinator=0

# MLPendingAckStore maintains a ConcurrentSkipListMap pendingAckLogIndex,
# It stores the position in pendingAckStore as its value and saves a position used to determine
# whether the previous data can be cleaned up as a key.
# transactionPendingAckLogIndexMinLag is used to configure the minimum lag between indexes
transactionPendingAckLogIndexMinLag=500

# The transaction buffer client's operation timeout in milliseconds.
transactionBufferClientOperationTimeoutInMills=3000

# Provide a mechanism allowing the Transaction Log Store to aggregate multiple records into a batched record and
# persist into a single BK entry. This will make Pulsar transactions work more  efficiently, aka batched log.
# see: https://github.com/apache/pulsar/issues/15370
transactionLogBatchedWriteEnabled=false

# If enabled the feature that transaction log batch, this attribute means maximum log records count in a batch.
transactionLogBatchedWriteMaxRecords=512

# If enabled the feature that transaction log batch, this attribute means bytes size in a batch，default 4m.
transactionLogBatchedWriteMaxSize=4194304

# If enabled the feature that transaction log batch, this attribute means maximum wait time(in millis) for the first
# record in a batch
transactionLogBatchedWriteMaxDelayInMillis=1

# Provide a mechanism allowing the Pending Ack Store to aggregate multiple records into a batched record and persist
# into a single BK entry. This will make Pulsar transactions work more efficiently, aka batched log.
# see: https://github.com/apache/pulsar/issues/15370
transactionPendingAckBatchedWriteEnabled=false

# If enabled the feature that transaction pending ack log batch, this attribute means maximum log records count in a
# batch.
transactionPendingAckBatchedWriteMaxRecords=512

# If enabled the feature that transaction pending ack log batch, this attribute means bytes size in a batch, default:4m.
transactionPendingAckBatchedWriteMaxSize=4194304

# If enabled the feature that transaction pending ack log batch, this attribute means maximum wait time(in millis) for
# the first record in a batch.
transactionPendingAckBatchedWriteMaxDelayInMillis=1

### --- Packages management service configuration variables (begin) --- ###

# Enable the packages management service or not
enablePackagesManagement=false

# The packages management service storage service provide
packagesManagementStorageProvider=org.apache.pulsar.packages.management.storage.bookkeeper.BookKeeperPackagesStorageProvider

# When the packages storage provider is bookkeeper, you can use this configuration to
# control the number of replicas for storing the package
packagesReplicas=1

# The bookkeeper ledger root path
packagesManagementLedgerRootPath=/ledgers

# When using BookKeeperPackagesStorageProvider, you can configure the
# bookkeeper client by prefixing configurations with "bookkeeper_".
# This config applies to managed ledger bookkeeper clients, as well.

### --- Packages management service configuration variables (end) --- ###

#enable or disable strict bookie affinity
strictBookieAffinityEnabled=false

### --- Deprecated settings --- ###

# These settings are left here for compatibility

# Zookeeper quorum connection string
# Deprecated: use metadataStoreUrl instead
zookeeperServers=

# Configuration Store connection string
# Deprecated. Use configurationMetadataStoreUrl
globalZookeeperServers=
configurationStoreServers=

# Zookeeper session timeout in milliseconds
# Deprecated: use metadataStoreSessionTimeoutMillis
zooKeeperSessionTimeoutMillis=-1

# ZooKeeper operation timeout in seconds
# Deprecated: use metadataStoreOperationTimeoutSeconds
zooKeeperOperationTimeoutSeconds=-1

# ZooKeeper cache expiry time in seconds
# Deprecated: use metadataStoreCacheExpirySeconds
zooKeeperCacheExpirySeconds=-1

# Deprecated - Enable TLS when talking with other clusters to replicate messages
replicationTlsEnabled=false

# Deprecated. Use brokerDeleteInactiveTopicsFrequencySeconds
brokerServicePurgeInactiveFrequencyInSeconds=60

# Deprecated - Use backlogQuotaDefaultLimitByte instead.
backlogQuotaDefaultLimitGB=-1

# Deprecated - Use webServicePortTls and brokerServicePortTls instead
tlsEnabled=false

# Enable Key_Shared subscription (default is enabled)
# @deprecated since 2.8.0 subscriptionTypesEnabled is preferred over subscriptionKeySharedEnable.
subscriptionKeySharedEnable=true

# Max number of "acknowledgment holes" that can be stored in Zookeeper. If number of unack message range is higher
# than this limit then broker will persist unacked ranges into bookkeeper to avoid additional data overhead into
# zookeeper.
# Deprecated: use managedLedgerMaxUnackedRangesToPersistInMetadataStore
managedLedgerMaxUnackedRangesToPersistInZooKeeper=-1

# If enabled, the maximum "acknowledgment holes" will not be limited and "acknowledgment holes" are stored in
# multiple entries.
persistentUnackedRangesWithMultipleEntriesEnabled=false

# Deprecated - Use managedLedgerCacheEvictionIntervalMs instead
managedLedgerCacheEvictionFrequency=0