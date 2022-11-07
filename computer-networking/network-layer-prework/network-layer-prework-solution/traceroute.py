import argparse
import socket
import sys
import time


def is_icmp(data):
    return data[9] == 1


def icmp_ttl_expired(data):
    ihl = (data[0] & 0b1111) * 4
    return data[ihl] == 11  # ICMP type "Time-to-live exceeded"


def icmp_destination_unreachable(data):
    ihl = (data[0] & 0b1111) * 4
    return data[ihl] == 3 and data[ihl + 1] == 3


def trace(host, options, out=sys.stdout):
    """
    Trace a basic path to the host
    """
    sender = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    receiver = socket.socket(socket.AF_INET, socket.SOCK_RAW,
                             socket.IPPROTO_ICMP)
    receiver.settimeout(options.waittime)
    target_port = options.port

    for ttl in range(options.first_ttl, options.max_ttl):
        out.write(f'{ttl:>2}  ')
        prior_ip = None
        for _ in range(options.nqueries):
            sender.setsockopt(socket.SOL_IP, socket.IP_TTL, ttl)
            t = time.time()
            sender.sendto(b'\x00' * 24, (host, target_port))
            out.flush()
            try:
                while True:
                    data, (responding_ip, responding_port) \
                            = receiver.recvfrom(4096)
                    name, _ = socket.getnameinfo(
                            (responding_ip, responding_port), 0)
                    if not is_icmp(data) and (
                            icmp_ttl_expired(data) or
                            icmp_destination_unreachable(data)):
                        continue
                    rtt = (time.time() - t) * 1000
                    if responding_ip != prior_ip:
                        if prior_ip is not None:
                            out.write('\n    ')
                        out.write(f'{name} ({responding_ip})')
                        prior_ip = responding_ip
                    out.write(f'  {rtt:.3f} ms')
                    break
            except socket.timeout:
                if prior_ip is None:
                    out.write('* ')
                else:
                    out.write('\n    *')
        out.write('\n')
        if icmp_destination_unreachable(data):
            break
        target_port += 1


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='A simple traceroute clone')
    parser.add_argument('host', help='A host IP address')
    parser.add_argument('-f', dest='first_ttl', default=1, type=int,
                        help='Set the initial time to live, default 1')
    parser.add_argument('-m', dest='max_ttl', default=64, type=int,
                        help='Set the max time to live, default 64')
    parser.add_argument('-p', dest='port', default=33434, type=int,
                        help='Set the base port number used in probes, '
                        'default 33434')
    parser.add_argument('-q', dest='nqueries', default=3, type=int,
                        help='Set the number of probes per ttl, default 3')
    parser.add_argument('-w', dest='waittime', default=5, type=int,
                        help='Set the time to wait for a response to a probe, '
                        'default 5 seconds')
    args = parser.parse_args()

    trace(args.host, args)
