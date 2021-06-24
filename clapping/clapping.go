package clapping

import (
	"time"

	"clap2mqtt/detection"
)

const LEAD_OUT_CLAP_MS = 80
const LEAD_IN_CLAPPING_MS = 500
const LEAD_OUT_CLAPPING_MS = 500

type Clapping struct {
	start      time.Time
	detections []detection.Detection
}

func NewClapping() Clapping {
	return Clapping{time.Now(), nil}
}

func (c Clapping) isClap(d detection.Detection) bool {
	return d.Duration().Milliseconds() >= 0 && d.Duration().Milliseconds() < LEAD_OUT_CLAP_MS
}

func (c Clapping) hasLeadIn() bool {
	return c.detections[0].GetStart().Sub(c.start).Milliseconds() > LEAD_IN_CLAPPING_MS
}

func (c Clapping) hasLeadOut() bool {
	return time.Now().Sub(c.detections[len(c.detections)-1].GetEnd()).Milliseconds() > LEAD_OUT_CLAPPING_MS
}

func (c Clapping) isValid() bool {
	if len(c.detections) > 0 && !c.hasLeadIn() {
		return false
	}

	if len(c.detections) > 1 {
		for i := 0; i < len(c.detections); i++ {
			//do more checks to validate clap array;
		}
	}

	return true
}

func (c Clapping) AddDetection(d detection.Detection) {
	if c.isValid() && c.isClap(d) {
		c.detections = append(c.detections, d)
	} else {
		c.Reset()
	}
}

func (c Clapping) HasStopped() bool {
	return len(c.detections) > 0 && c.hasLeadIn() && c.hasLeadOut()
}

func (c Clapping) Count() int {
	if c.isValid() {
		return len(c.detections)
	}

	return 0
}

func (c Clapping) Reset() {
	c.detections = make([]detection.Detection, 0)
	c.start = time.Now()
}
