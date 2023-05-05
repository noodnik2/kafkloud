# Kafkloud Design / Implementation

## First Kafkloud Use Case
![First Use Case Diagram](kafkloud-firstusecase-1xbg.png)
_(diagram created with [excalidraw](https://excalidraw.com/))_

Keeping it simple, the first implementation milestone for Kafkloud is to demonstrate
the flow of information across all components.

The following endpoints are suggested to facilitate demonstration of having met this goal:

### Endpoints
#### Portal
- `DoIt` - UI button triggering a call to `producer` resulting in a message being sent to `streamer`
- `IsDone` - UI control or display informing user of the status of processing of the message

#### Producer
- `Produce` - HTTP `POST` endpoint triggering creation and delivery of message(s) to `streamer`

#### Streamer
- Standard Kafka APIs are used the inbound and outbound endpoint(s) for `streamer`

#### Consumer
- `Status` - HTTP `GET` endpoint retrieving the evidence of `producer` message(s)
