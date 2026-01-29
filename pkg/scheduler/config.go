package scheduler

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type TerminusArgs struct {
	metav1.TypeMeta       `json:",inline"`
	OversubscriptionRatio float64 `json:"oversubscriptionRatio"`
}

// 默认配置
func (args *TerminusArgs) SetDefaults() {
	if args.OversubscriptionRatio == 0 {
		args.OversubscriptionRatio = 1.0
	}
}
