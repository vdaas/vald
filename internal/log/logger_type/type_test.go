package loggertype

import "testing"

func TestString(t *testing.T) {
	type test struct {
		name       string
		loggerType LoggerType
		want       string
	}

	tests := []test{
		{
			name:       "returns glg",
			loggerType: GLG,
			want:       "glg",
		},

		{
			name:       "returns zap",
			loggerType: ZAP,
			want:       "zap",
		},

		{
			name:       "returns zerolog",
			loggerType: ZEROLOG,
			want:       "zerolog",
		},

		{
			name:       "returns logrus",
			loggerType: LOGRUS,
			want:       "logrus",
		},

		{
			name:       "returns klog",
			loggerType: KLOG,
			want:       "klog",
		},

		{
			name:       "returns unknown",
			loggerType: LoggerType(100),
			want:       "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.loggerType.String()
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}

func TestAtot(t *testing.T) {
	type test struct {
		name string
		str  string
		want LoggerType
	}

	tests := []test{
		{
			name: "returns GLG when str is glg",
			str:  "glg",
			want: GLG,
		},

		{
			name: "returns GLG when str is GLg",
			str:  "GLg",
			want: GLG,
		},

		{
			name: "returns ZAP when str is zap",
			str:  "zap",
			want: ZAP,
		},

		{
			name: "returns ZAP when str is ZAp",
			str:  "ZAp",
			want: ZAP,
		},

		{
			name: "returns ZEROLOG when str is zerolog",
			str:  "zerolog",
			want: ZEROLOG,
		},

		{
			name: "returns ZEROLOG when str is ZEROLOg",
			str:  "ZEROLOg",
			want: ZEROLOG,
		},

		{
			name: "returns LOGRUS when str is logrus",
			str:  "logrus",
			want: LOGRUS,
		},

		{
			name: "returns LOGRUS when str is LOGRUs",
			str:  "LOGRUs",
			want: LOGRUS,
		},

		{
			name: "returns KLOG when str is klog",
			str:  "klog",
			want: KLOG,
		},

		{
			name: "returns KLOG when str is KLOg",
			str:  "KLog",
			want: KLOG,
		},

		{
			name: "returns unknown when str is Vald",
			str:  "Vald",
			want: Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Atot(tt.str)
			if got != tt.want {
				t.Errorf("not equals. want: %v, but got: %v", tt.want, got)
			}
		})
	}
}
