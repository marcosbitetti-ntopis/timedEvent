package tests

import (
	"fmt"
	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
	"strconv"
	"time"
)

const (
	DATE_FORMAT        = "2006-01-02 15:04:05Z"
	TEST_ENDPOINT      = "http://localhost:9010/v1/event"
	CONTENT_TYPE       = "Content-Type"
	CONTENT_TYPE_CE    = "application/cloudevents"
	PUBLISH_DATE       = "publishDate"
	PUBLISH_QUEUE      = "publishQueue"
	PUBLISH_TYPE       = "publishtype"
	TEST_PUBLISH_QUEUE = "throwAt"
	TEST_PUBLISH_TYPE  = "data_only"
)

func CreateEventRequest() {
	ginkgo.It("Valid msg", func() {
		for i := 1; i < 300; i++ {
			strIvalue := strconv.Itoa(i)
			fmt.Println("Trying to create an event " + strIvalue)
			mockReader, err := GetMockReader(getMockEvent(time.Now(), strIvalue))
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			h := make(map[string]string)
			h[CONTENT_TYPE] = CONTENT_TYPE_CE
			h[PUBLISH_DATE] = time.Now().Add(time.Duration(i) * 100 * time.Millisecond).UTC().Format(DATE_FORMAT)
			h[PUBLISH_QUEUE] = TEST_PUBLISH_QUEUE
			h[PUBLISH_TYPE] = TEST_PUBLISH_TYPE
			resp, err := SendPostRequestWithHeaders(TEST_ENDPOINT, mockReader, h)
			gomega.Expect(err).ShouldNot(gomega.HaveOccurred())
			gomega.Expect(resp.StatusCode).To(gomega.Equal(201))
		}
	})
}

func getMockEvent(publihsDate time.Time, ref string) interface{} {
	return struct {
		SpecVersion  string `json:"specversion"`
		Type         string `json:"type"`
		Source       string `json:"source"`
		ID           string `json:"id"`
		PublishDate  string `json:"publishdate"`
		PublishQueue string `json:"publishqueue"`
		PublishType  string `json:"publishtype"`
		Data         string `json:"data"`
	}{
		SpecVersion:  "0.2",
		Type:         "TestEvent",
		Source:       "sourceEvent",
		ID:           fmt.Sprintf("mockEvent%s", ref),
		PublishDate:  publihsDate.Format(DATE_FORMAT),
		PublishQueue: TEST_PUBLISH_QUEUE,
		PublishType:  TEST_PUBLISH_TYPE,
		Data:         fmt.Sprintf("Mock event ref: %s, generated at %s", ref, publihsDate.Format(DATE_FORMAT)),
	}
}
