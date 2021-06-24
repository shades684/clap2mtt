package detection

import "time"

const LEAD_OUT_DETECTION_MS = 80

type Detection struct {
	start            time.Time
	lastNonDetection time.Time
	lastDetection    time.Time
}

func NewDetection() *Detection {
	now := time.Now()
	d := Detection{now, now, now}

	return &d
}

func (d Detection) Duration() time.Duration {
	return d.lastDetection.Sub(d.start)
}

func (d Detection) GetStart() time.Time {
	return d.start
}

func (d Detection) GetEnd() time.Time {
	return d.lastDetection
}

func (d Detection) HasStopped() bool {
	return d.lastDetection.Sub(d.lastNonDetection).Milliseconds() > LEAD_OUT_DETECTION_MS
}

func (d Detection) Update(signal bool) {
	now := time.Now()

	if signal {
		d.lastDetection = now
	} else {
		d.lastNonDetection = now
	}
}
