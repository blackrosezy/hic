import socket


class Parser:
    def split_url(self, url):
        elements = url.split(':')
        if len(elements) != 2:
            return True, '', ''

        if 'http' != elements[0] and 'https' != elements[0]:
            return True, '', ''

        domain = elements[1][2:]
        return False, elements[0], domain

    def split_ipdata(self, ipdata):
        elements = ipdata.split(':')
        if len(elements) != 3:
            return True, '', '', 0

        if 'http' != elements[0] and 'https' != elements[0]:
            return True, '', '', 0

        if len(elements[1]) < 10:
            return True, '', '', 0

        ip = elements[1][2:]
        if not self.__is_valid_ipv4_address(ip):
            return True, '', '', 0

        try:
            port = int(elements[2])
            if port <= 0:
                return True, '', '', 0
        except ValueError:
            return True, '', '', 0

        return False, elements[0], ip, port


    def __is_valid_ipv4_address(self, address):
        try:
            socket.inet_pton(socket.AF_INET, address)
        except AttributeError:  # no inet_pton here, sorry
            try:
                socket.inet_aton(address)
            except socket.error:
                return False
            return address.count('.') == 3
        except socket.error:  # not a valid address
            return False

        return True