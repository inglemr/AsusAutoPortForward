package pkg

import (
	"strconv"

	corev1 "k8s.io/api/core/v1"
)

type PortForwardRule struct {
	RuleName   string
	TargetPort string
	SourcePort string
	Protocol   string
	TargetIP   string
}

func (p *PortForwardRule) RouterString() string {
	return "<" + p.RuleName + ">" + p.SourcePort + ">" + p.TargetIP + ">" + p.TargetPort + ">" + p.Protocol + ">"
}

func NewPortForwardRulesFromK8sService(service corev1.Service, targetAddress string) map[string]PortForwardRule {
	rules := make(map[string]PortForwardRule)
	for _, port := range service.Spec.Ports {
		targetPort := strconv.Itoa(int(port.Port))
		sourcePort := strconv.Itoa(int(port.NodePort))
		proto := string(port.Protocol)
		rule := PortForwardRule{
			RuleName:   GetPortNameFromK8sPort(service.Name, service.Namespace, port),
			TargetPort: targetPort,
			SourcePort: sourcePort,
			Protocol:   proto,
			TargetIP:   targetAddress,
		}
		rules[rule.RuleName] = rule
	}
	return rules
}

func getSubstringUpToCharacterLimit(input string, limit int) string {
	if len(input) <= limit {
		return input
	}
	return input[:limit]
}

func GetPortNameFromK8sPort(serviceName string, namespace string, port corev1.ServicePort) string {
	targetPort := strconv.Itoa(int(port.Port))
	proto := string(port.Protocol)
	return getSubstringUpToCharacterLimit(namespace, 8) + "-" + getSubstringUpToCharacterLimit(serviceName, 12) + "-" + getSubstringUpToCharacterLimit(targetPort, 5) + "-" + getSubstringUpToCharacterLimit(proto, 4)
}
