# Piped Piper

Experiments with golang channel and patters:

### Inspirations:
- Fork Join model
- [GenStage](https://hexdocs.pm/gen_stage/ConsumerSupervisor.html) and Actor model
- [Akka](https://doc.akka.io/docs/akka/current/typed/dispatchers.html)
- [Reactive Stream](https://gist.github.com/staltz/868e7e9bc2a7b8c1f754) and [BackPressure](https://github.com/ReactiveX/RxJava/wiki/Backpressure)
- [Docs](http://conal.net/talks/)


### Running
`go run -race main.go`


### Thoughts.

- Producer would be responsible for streaming a batch of data as requested by Consumer
- Consumer would consume
- Dispatcher is a wrapper around consumer. It takes a list and sends data to consumer for consumption 
- Supervisor, is what connects the producer and consumer/dispatcher
- ProducerConsumer, consumes events and then send it to channel


##### Producer

```
type Result[E any] struct {
    Data E
    Err error
}

type Producer interface {
	Next(ctx context.Context, opts Opts) chan Result
	Stop() error
}
```

##### Dispatcher

```
type Dispatcher interface {
	Dispatch(ctx context.Context, result Result) chan Result
}

```

We have two types of `Dispatcher`
- `ConsumerGroupDispatcher`, takes consumers and loadbalancer and balances amongst consumers
- `BroadcastDispatcher`, takes consumers and sends events to all consumers


##### Consumer

```
// Either

type Consumer interface {
    Handle(ctx context.Context, result Result) error
}

// Or

type Consumer interface {
    Handle(ctx context.Context, result Result) chan Result
}
```


##### ProducerConsumer

```
type ProducerConsumer struct {
    Consumer Consumer
}

func (pc *ProducerConsumer) Thinking(ctx context.Context, opts Opts) chan Result {
}
```
