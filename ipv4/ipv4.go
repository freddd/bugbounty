package ipv4

import (
	"fmt"
	"net"
)

// Borrowed from: https://github.com/nathanlong85/iptools/tree/master/ipv4range

// IPv4Range represents an IPv4 subnet
type IPv4Range struct {
	ips           []net.IP
	availableIPs  []net.IP
	unavailbleIPs []net.IP
	mask          net.IPMask
	network       net.IP
	broadcast     net.IP
}

// New IPv4 network range. Results include network and broadcase IPs.
// cidrNet format should be "192.168.0.0/23"
func New(cidrNet string) (IPv4Range, error) {
	var r IPv4Range

	ip, ipnet, err := net.ParseCIDR(cidrNet)
	if err != nil {
		return r, err
	}

	// Calculate the range of IPs in the subnet
	var ips []net.IP
	for ip.Mask(ipnet.Mask); ipnet.Contains(ip); increment(ip) {
		copiedIP := make([]byte, len(ip))
		copy(copiedIP, ip)
		ips = append(ips, copiedIP)
	}

	// Set all initial properties on IPV4Range
	r.ips = ips
	r.mask = ipnet.Mask
	r.network = ipnet.IP
	r.broadcast = Broadcast(ipnet.IP, ipnet.Mask)

	// Set availableIPs, remove broadcast/network addresses
	r.availableIPs = make([]net.IP, len(r.ips))
	copy(r.availableIPs, r.ips)
	r.Remove(r.broadcast)
	r.Remove(r.network)
	return r, nil
}

// All returns a slice of all IPs that were calculated in the subnet.
// This includes broadcast, network, and any removed addresses.
func (r *IPv4Range) All() []net.IP {
	return r.ips
}

// Mask returns a net.IPMask netmask for the subnet
func (r *IPv4Range) Mask() net.IPMask {
	return r.mask
}

// Network returns a net.IP network address for the subnet
func (r *IPv4Range) Network() net.IP {
	return r.network
}

// Broadcast returns a net.IP broadcast address for the subnet
func (r *IPv4Range) Broadcast() net.IP {
	return r.broadcast
}

// Available returns only usable IPs. Broadcast/Network are filtered out by
// default. Any other IPs that have been removed with Remove() will not
// be returned as well.
func (r *IPv4Range) Available() []net.IP {
	return r.availableIPs
}

// NextAvailable returns the next available IP(s) in the IP Range. The number of
// available IPs returned should be specified as a parameter.
func (r *IPv4Range) NextAvailable(num int) ([]net.IP, error) {
	if len(r.availableIPs) < num {
		return nil, fmt.Errorf("Requested %d IPs, only %d available", num, len(r.availableIPs))
	}

	return r.availableIPs[:num], nil
}

// Unavailable returns whichever IPs are missing from Available().
func (r *IPv4Range) Unavailable() []net.IP {
	return r.unavailbleIPs
}

// Remove manually removes an IP from an IPv4Range
func (r *IPv4Range) Remove(ip net.IP) bool {
	found := false
	for idx, val := range r.availableIPs {
		if val.Equal(ip) {
			found = true
			r.availableIPs = append(r.availableIPs[:idx], r.availableIPs[idx+1:]...)
			r.unavailbleIPs = append(r.unavailbleIPs, ip)
		}
	}

	return found
}

// Broadcast calculates a subnet's broadcast address based on an IP and Mask
func Broadcast(ip net.IP, mask net.IPMask) net.IP {
	return net.IPv4(
		ip[0]|^mask[0],
		ip[1]|^mask[1],
		ip[2]|^mask[2],
		ip[3]|^mask[3])
}

// Helper function for IPV4Range.New
func increment(ip net.IP) {
	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}
}