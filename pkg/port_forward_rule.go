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
	return "<" + p.RuleName + ">" + p.TargetPort + ">" + p.TargetIP + ">" + p.SourcePort + ">" + p.Protocol + ">"
}

func NewPortForwardRulesFromK8sService(service corev1.Service, targetAddress string) map[string]PortForwardRule {
	rules := make(map[string]PortForwardRule)
	annotations := service.GetAnnotations()
	for _, port := range service.Spec.Ports {
		if annotations["autoportforward/"+port.Name+"."+"ignoreport"] == "true" {
			continue
		}
		targetPort := strconv.Itoa(int(port.Port))
		if annotations["autoportforward/"+port.Name+"."+"usenodeport"] == "true" {
			targetPort = strconv.Itoa(int(port.NodePort))
		}

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

func GetServiceNamePrefix(serviceName string, namespace string) string {
	return "KAPF@" + getSubstringUpToCharacterLimit(namespace, 5) + "-" + getSubstringUpToCharacterLimit(serviceName, 10) + "-"
}

func GetPortNameFromK8sPort(serviceName string, namespace string, port corev1.ServicePort) string {
	targetPort := strconv.Itoa(int(port.Port))
	proto := string(port.Protocol)
	return GetServiceNamePrefix(serviceName, namespace) + getSubstringUpToCharacterLimit(targetPort, 5) + "-" + getSubstringUpToCharacterLimit(proto, 4)
}
