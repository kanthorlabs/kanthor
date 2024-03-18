package constants

import "fmt"

var TopicCore = "core"

var MessageTopic = fmt.Sprintf("%s.message", TopicCore)
var MessageCreateSubject = fmt.Sprintf("%s.create", MessageTopic)
var MessagePublisher = "publisher_message"

func MessageSubscriber(name string) string {
	return fmt.Sprintf("subscriber_message_%s", name)
}

var RequestTopic = fmt.Sprintf("%s.request", TopicCore)
var RequestScheduleSubject = fmt.Sprintf("%s.schedule", RequestTopic)
var RequestPublisher = "publisher_request"

func RequestSubscriber(name string) string {
	return fmt.Sprintf("subscriber_request_%s", name)
}

var ResponseTopic = fmt.Sprintf("%s.response", TopicCore)
var ResponseNotifySubject = fmt.Sprintf("%s.notify", ResponseTopic)
var ResponsePublisher = "publisher_response"

func ResponseSubscriber(name string) string {
	return fmt.Sprintf("subscriber_response_%s", name)
}

func Subscriber(name string) string {
	return fmt.Sprintf("subscriber_%s", name)
}
