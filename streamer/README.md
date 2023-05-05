# `streamer`

Provides a data stream service.

## Kafka Broker

This incarnation of `streamer` employs a [Kafka Broker](https://developer.confluent.io/quickstart/kafka-docker/)
for implementation of the stream.

### The Stream

The stream is created with the following command:

```shell
$ docker exec broker kafka-topics \
  --bootstrap-server broker:9092 \
  --create \
  --topic stream
```

### Useful Commands

With the `broker` container running in `docker`:

```shell
$ : list the known topics
$ docker exec broker kafka-topics --bootstrap-server broker:9092 --list
$ : create a topics
$ docker exec broker kafka-topics --bootstrap-server broker:9092 --create --topic stream
$ docker exec --interactive --tty broker \
kafka-console-producer --bootstrap-server broker:9092 \
                       --topic stream
Hello, stream.
^D
$ docker exec --interactive --tty broker \
kafka-console-consumer --bootstrap-server broker:9092 \
                       --topic stream \
                       --from-beginning
Hello, stream.
^CProcessed a total of 1 messages
```

## 3/28/23 Notes

### Kubernetization of Kafka

Looks daunting.  See:
- [with Helm](https://dev.to/thegroo/running-kafka-on-kubernetes-for-local-development-with-helm-2ne8)
- [from docker-compose](https://kubernetes.io/docs/tasks/configure-pod-container/translate-compose-kubernetes/)
  - Probably naiive to think it relevant in this case
- [Bitnami #9406](https://github.com/bitnami/charts/issues/9406)
- [... in 5 minutes](https://tsuyoshiushio.medium.com/local-kafka-cluster-on-kubernetes-on-your-pc-in-5-minutes-651a2ff4dcde)
  - Possibly quite relevant, but likely outdated (from 2019)
- [Setting Up Kafka on K8s the Easy Way](https://blog.datumo.io/setting-up-kafka-on-kubernetes-an-easy-way-26ae150b9ca8)
  - This one might be exactly what we need!
- [How to Connect to Apache Kafka running in Docker](https://www.youtube.com/watch?v=L--VuzFiYrM)
- [Kubernetes Kafka Manifests](https://kow3ns.github.io/kubernetes-kafka/manifests/)
- [How To](https://levelup.gitconnected.com/how-to-deploy-apache-kafka-with-kubernetes-9bd5caf7694f)
- [Another How To](https://dzone.com/articles/how-to-deploy-apache-kafka-with-kubernetes)
- [... on Docker Desktop](https://collabnix.com/deploy-apache-kafka-on-kubernetes-kafka-tutorials/)
- [... Best Practices](https://www.weave.works/blog/kafka-on-kubernetes-and-deploying-best-practice)

Conclusion for today: just live with it outside of k8s.

### Deployment of Kafka

```shell
$ cd ~/repos/noodnik2/incubator20/k8s/kafkloud/streamer
$ docker-compose up zookeeper -d
$ sleep 30 # this is needed to force "zookeeper" to "time out" the previous session
$ docker-compose up broker -d
$ docker exec broker kafka-topics --bootstrap-server broker:9092 --list
__consumer_offsets
stream
```

Since we see the topic that was created previously (i.e., `stream`), we're
good to go for now.

## 4/7/23 Notes

### Monitoring

Success using [kafkaui](https://github.com/provectus/kafka-ui) discovered via
[this article](https://needablackcoffee.medium.com/my-top-5-tools-to-manage-develop-with-apache-kafka-2e6790a88ef2).

```shell
$ docker run -it -p 8060:8080 -e DYNAMIC_CONFIG_ENABLED=true provectuslabs/kafka-ui
```

Browsing [https://localhost:8060](https://localhost:8060) revealed
its ability to browse the status of brokers and topics, and create
or delete topics.

#### Also See
- [Kafka Open Source Monitoring Tools](https://sematext.com/blog/kafka-open-source-monitoring-tools/)

## 3/17 Updates, Towards Pushing to GCP
Added `namespace` qualifier and `component` label to `k8s` specs, enabling grouping of `streamer`
resources, e.g. simplifying cleanup:

```shell
$ kubectl delete all -n kafkloud -l component=streamer
```

First attempt to push `streamer` to GCP:

```shell
$ gcloud container clusters create-auto kafkloud-cluster --region=us-west1
$ kubectl create namespace kafkloud
$ kubectl config set-context --current --namespace=kafkloud
$ kubectl apply -f k8s
$ : wait a long time for the auto-scaling to enable creation of the containers...
$ kubectl port-forward pod/kafkaui-dpl-678b9bb985-n9wb4 8080 9092
$ : browse to the KafkaUI app at "localhost:8080", enter "kafka-svc" at port "9092" \
> Use the "Topics" tab to create a topic and send a message
$ gcloud container clusters delete kafkloud-cluster --zone us-west1
```

NOTES:
- it took over 5 minutes to create the 3-node cluster; consider specifying 1 node next time? 

