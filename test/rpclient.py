import json
import socket
import time
import random

class RPCClient(object):

    def __init__(self, addr, codec=json):
        self._socket = socket.create_connection(addr)
        self._codec = codec

    def _message(self, name, *params):
        return dict(id=random.randint(1,1000),
                    params=list(params),
                    method=name)

    def call(self, name, *params):
        req = self._message(name, *params)
        id = req.get('id')

        mesg = self._codec.dumps(req)
        self._socket.sendall(bytes(mesg,encoding='utf8'))

        # This will actually have to loop if resp is bigger
        resp = self._socket.recv(4096)
        resp = self._codec.loads(resp)

        if resp.get('id') != id:
            raise Exception("expected id=%s, received id=%s: %s"
                            %(id, resp.get('id'), resp.get('error')))

        if resp.get('error') is not None:
            raise Exception(resp.get('error'))

        return resp.get('result')

    def close(self):
        self._socket.close()


if __name__ == '__main__':
    # modify the connect info for your rpc server listen
    rpc = RPCClient(("127.0.0.1", 9876))
    
    print ("AccRPC.JsonTest",rpc.call("AccRPC.JsonTest", {"name":"Trias","age":24}))
    print ("AccRPC.Test",rpc.call("AccRPC.Test","The test is ok"))
    print("==========================")
    
    
    print ("AccRPC.CreateAcc",rpc.call("AccRPC.CreateAcc", {"path":"/tmp/testacc.json","pass":"1234qwer"}))
    print ("AccRPC.ImportAcc",rpc.call("AccRPC.ImportAcc", {"path":"/tmp/testacc.json","pass":"1234qwer"}))
    acclist=rpc.call("AccRPC.GetAcclist", "GetAcclist")
    print ("AccRPC.GetAcclist",acclist)
    if not acclist:
        print("Have Not account in Account Server!")
        exit(1)

    print ("AccRPC.ExportAcc",rpc.call("AccRPC.ExportAcc", {"addr":acclist[0],"path":"/tmp","pass1":"1234qwer","pass2":"1q2w3e4r"}))
    
    tsign=rpc.call("AccRPC.Sign", {"addr":acclist[0],"hash":"hashtext","pass":"1234qwer"})
    print ("AccRPC.Sign",tsign)
    print ("AccRPC.Verify",rpc.call("AccRPC.Verify", {"pubkey":tsign['Pubkey'],"hash":"hashtext","stext":tsign['Sigadata']}))
    
    cs=rpc.call("AccRPC.CreateShieldAddr", acclist[0])
    print ("AccRPC.CreateShieldAddr",cs)
    print ("AccRPC.Verify_shield",rpc.call("AccRPC.Verify_shield", {"addr":acclist[0],"shieldaddr":cs['Shieldaddr'],"shieldpkey":cs['ShieldpKey']}))
    stsign=rpc.call("AccRPC.Shield_Sign", {"addr":acclist[0],"pass":"1234qwer","hash":"hashtext","shieldpkey":cs['ShieldpKey']})
    print ("AccRPC.Shield_Sign",stsign)
    print ("AccRPC.Verify",rpc.call("AccRPC.Verify", {"pubkey":stsign['Pubkey'],"hash":"hashtext","stext":stsign['Sigadata']}))
    
    print("==========================")
    print ("AccRPC.Pubkey2Addr",rpc.call("AccRPC.Pubkey2Addr", tsign['Pubkey']))
    print ("AccRPC.Shield_Pubkey2Addr",rpc.call("AccRPC.Shield_Pubkey2Addr", stsign['Pubkey']))
    