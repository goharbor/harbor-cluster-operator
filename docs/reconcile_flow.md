# Harbor Cluster Reconciler

## Control Loop

 ![Control Loop](./assets/reconcile.png)

### The Criterion of READY

   The important point of Control Loop is how to evaluate the CRs(Redis, PostgreSql, Minio) status is READY.
   When all of the three dependencies is READY, operator can create the Harbor CR.

 - Redis
    
    When the available replicas of redis is equal to `.spec.redis.replicas` in Redis CR. 
    And the available replicas of redis sentinel is equal to `.spec.sentinel.replicas` is Redis CR.
    
 - PostgreSQL
    
    In Postgres CR's status, there is string field named `PostgresClusterStatus`, which indicates whether the PostgreSQL services is Healthy.
    In current stage, we directly use it to evaluate the PostgreSQL is ready when this field is equal to `running`.
    
 - Minio
 
    In Minio CR's status, there is int32 field named `AvailableReplicas`, which indicates the available replicas of MinIO instance.
    So, when `.status.availableReplicas` is equal to `.spec.replicas`, it is ready. 
    