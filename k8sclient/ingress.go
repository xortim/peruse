package k8sclient

import (
	"net/url"
	"strings"

	v1beta1 "k8s.io/api/extensions/v1beta1"
)

// IngressURLs returns a slice of URLs that are captured by the Ingress
// supports external-dns annotation external-dns.alpha.kubernetes.io/hostname=fqdn.
func IngressURLs(ing v1beta1.Ingress) []url.URL {
	urls := []url.URL{}
	externalHost := IngressExternalDNSName(&ing)
	statusName := IngressStatusName(&ing)
	for _, rule := range ing.Spec.Rules {
		url := url.URL{Scheme: "http"}
		if IngressHostTLS(rule.Host, ing.Spec.TLS) {
			url.Scheme = "https"
		}

		// default to information reported by the Ingress status
		// items like external-dns and ingress controllers use this for backend info
		url.Host = statusName
		// if the rule has a host set, then use it
		if len(rule.Host) != 0 {
			url.Host = rule.Host
		}

		// if the externalHost is set, then use that
		if len(externalHost) != 0 {
			url.Host = externalHost
		}

		for _, path := range rule.HTTP.Paths {
			url.Path = path.Path
		}
		urls = append(urls, url)
	}
	return urls
}

// IngressExternalDNSName returns the value of the external-dns annotation
func IngressExternalDNSName(ing *v1beta1.Ingress) string {
	// trim the trailing `.` - assumes external-dns is not configured for default domain appending
	// as such, if this is the case we're also going to assume that that said domain is in the search
	// configuration for hosts that would have access to this information.
	return strings.Trim(ing.GetObjectMeta().GetAnnotations()[ExternalDNSHostnameAnnotation], ".")
}

// IngressStatusName returns the first available hostname or IP reported in the status field
func IngressStatusName(ing *v1beta1.Ingress) string {
	if len(ing.Status.LoadBalancer.Ingress) == 0 {
		return ""
	}
	if len(ing.Status.LoadBalancer.Ingress[0].Hostname) != 0 {
		return ing.Status.LoadBalancer.Ingress[0].Hostname
	}
	return ing.Status.LoadBalancer.Ingress[0].IP
}

// IngressHostTLS returns true if the path has a corresponding TLS host entry
func IngressHostTLS(needle string, ingTLSs []v1beta1.IngressTLS) bool {
	for _, ingTLS := range ingTLSs {
		for _, host := range ingTLS.Hosts {
			if host == needle {
				return true
			}
		}
	}
	return false
}
